package custom_flow

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/access/grpc"
	"github.com/onflow/flow-go-sdk/crypto"
)

const sendTokensScript = `
import FungibleToken from 0x9a0766d93b6608b7;
import FlowToken from 0x7e60df042a9c0868;
import FungibleTokenMetadataViews from 0x9a0766d93b6608b7;

transaction(amount: UFix64, to: Address) {
    let vaultData: FungibleTokenMetadataViews.FTVaultData
    let sentVault: @{FungibleToken.Vault}

    prepare(signer: auth(BorrowValue) &Account) {
        self.vaultData = FlowToken.resolveContractView(
            resourceType: nil,
            viewType: Type<FungibleTokenMetadataViews.FTVaultData>()
        ) as! FungibleTokenMetadataViews.FTVaultData?
            ?? panic("Could not resolve FTVaultData view")

        let vaultRef = signer.storage.borrow<auth(FungibleToken.Withdraw) &FlowToken.Vault>(
            from: self.vaultData.storagePath
        ) ?? panic("Could not borrow Vault reference from storage path: ".concat(self.vaultData.storagePath.toString()))

        self.sentVault <- vaultRef.withdraw(amount: amount)
    }

    execute {
        let recipient = getAccount(to)
        let receiverRef = recipient.capabilities.borrow<&{FungibleToken.Receiver}>(
            self.vaultData.receiverPath
        ) ?? panic("Could not borrow Receiver reference to path: ".concat(self.vaultData.receiverPath.toString()))

        receiverRef.deposit(from: <-self.sentVault)
    }
}
`

func SendTokens(
	ctx context.Context,
	client *grpc.Client,
	signerAddress flow.Address,
	signerPrivateKey crypto.PrivateKey,
	keyIndex int,
	amount float64,
	recipientAddress flow.Address,
) (*flow.TransactionResult, error) {
	// 1. Получаем последний блок для reference block
	latestBlock, err := client.GetLatestBlockHeader(ctx, true)
	signerAccount, err := client.GetAccount(ctx, signerAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest block: %w", err)
	}

	// Получаем текущий sequence number
	sequenceNumber := signerAccount.Keys[keyIndex].SequenceNumber

	// 2. Получаем последний блок для reference block
	latestBlock, err = client.GetLatestBlockHeader(ctx, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest block: %w", err)
	}

	// 2. Создаем транзакцию
	tx := flow.NewTransaction().
		SetScript([]byte(sendTokensScript)).
		SetGasLimit(100).
		SetProposalKey(
			signerAddress,
			uint32(keyIndex),
			sequenceNumber,
		).
		SetReferenceBlockID(latestBlock.ID).
		SetPayer(signerAddress).
		AddAuthorizer(signerAddress)

	// 3. Добавляем аргументы в том же порядке, что и в Cadence
	// Аргумент 1: amount (UFix64)
	amountUFix64, err := cadence.NewUFix64(fmt.Sprintf("%.8f", amount))
	if err != nil {
		return nil, fmt.Errorf("failed to create amount UFix64: %w", err)
	}
	if err := tx.AddArgument(amountUFix64); err != nil {
		return nil, fmt.Errorf("failed to add amount argument: %w", err)
	}

	// Аргумент 2: recipient address (Address)
	recipientCadence := cadence.NewAddress(recipientAddress)
	if err := tx.AddArgument(recipientCadence); err != nil {
		return nil, fmt.Errorf("failed to add recipient argument: %w", err)
	}

	// 4. Подписываем транзакцию
	signer, err := crypto.NewInMemorySigner(signerPrivateKey, crypto.SHA3_256)
	if err != nil {
		return nil, fmt.Errorf("failed to create signer: %w", err)
	}

	if err := tx.SignEnvelope(signerAddress, uint32(keyIndex), signer); err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	// 5. Отправляем транзакцию
	if err := client.SendTransaction(ctx, *tx); err != nil {
		return nil, fmt.Errorf("failed to send transaction: %w", err)
	}

	// 6. Ждем завершения транзакции
	return WaitForTransactionSeal(ctx, client, tx.ID())
}

func WaitForTransactionSeal(ctx context.Context, client *grpc.Client, txID flow.Identifier) (*flow.TransactionResult, error) {
	fmt.Printf("Waiting for transaction %s to be sealed...\n", txID)

	for {
		result, err := client.GetTransactionResult(ctx, txID)
		if err != nil {
			return nil, fmt.Errorf("failed to get transaction result: %w", err)
		}

		switch result.Status {
		case flow.TransactionStatusSealed:
			return result, nil
		case flow.TransactionStatusExpired:
			return nil, fmt.Errorf("transaction expired")
		}

		select {
		case <-time.After(time.Second):
			// Продолжаем ожидание
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func main() {
	ctx := context.Background()

	// 1. Подключение к Flow (testnet)
	client, err := grpc.NewClient("access.devnet.nodes.onflow.org:9000")
	if err != nil {
		log.Fatalf("Failed to connect to Flow: %v", err)
	}

	// 2. Конфигурация отправителя (SIGNER_ADDRESS)
	signerPrivateKeyHex := "b29ef36d10870513effbd9721de7d6fa2a41450c8d6ba3710c007f399403eb28"
	signerAddress := flow.HexToAddress("0x534ee6a91a6def16")
	keyIndex := 0

	// 3. Конфигурация получателя (TARGET_ADDRESS)
	recipientAddress := flow.HexToAddress("0x51846a0f69492bba")
	amount := 10.0

	// 4. Декодирование приватного ключа
	privateKey, err := crypto.DecodePrivateKeyHex(crypto.ECDSA_P256, signerPrivateKeyHex)
	if err != nil {
		log.Fatalf("Failed to decode private key: %v", err)
	}

	// 5. Выполнение транзакции
	result, err := SendTokens(
		ctx,
		client,
		signerAddress,
		privateKey,
		keyIndex,
		amount,
		recipientAddress,
	)
	if err != nil {
		log.Fatalf("Failed to send tokens: %v", err)
	}

	if result.Error != nil {
		log.Fatalf("Transaction failed: %v", result.Error)
	}

	log.Printf("Transaction successful! ID: %s", result.TransactionID)
}
