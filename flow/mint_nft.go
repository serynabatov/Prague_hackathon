package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/onflow/cadence" // Needed again for constructing arguments
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/access/grpc"
	"github.com/onflow/flow-go-sdk/crypto"
)


const mintNFTScript = `
import NonFungibleToken from 0x631e88ae7f1d7c20
import FooBarV4 from 0x4ac0ee1c903bf362

transaction(
    recipient: Address,
    name: String,
    description: String,
    url: String,
    role: Bool
) {
    let minter: &FooBarV4.NFTMinter
    let recipientCollectionRef: &{NonFungibleToken.Receiver}

    prepare(signer: auth(BorrowValue) &Account) {
        self.minter = signer.storage.borrow<&FooBarV4.NFTMinter>(from: FooBarV4.MinterStoragePath)
            ?? panic("Account does not store an NFTMinter object at the specified path. Make sure the signer account is configured to mint.")

        self.recipientCollectionRef = getAccount(recipient).capabilities.borrow<&{NonFungibleToken.Receiver}>(
                FooBarV4.CollectionPublicPath
            ) ?? panic("Could not get receiver reference to the NFT Collection. Make sure the recipient has a collection setup.")
    }

    execute {
        let mintedNFT <- self.minter.createNFT(name: name, description: description, url: url, organizerBool: role)
        self.recipientCollectionRef.deposit(token: <-mintedNFT)
        log("NFT minted and deposited to recipient's collection.")
    }
}
`


// MintNFT mints an NFT and sends it to the recipient.
// The signerAddress is the account that owns the NFTMinter resource.
func MintNFT(
	ctx context.Context,
	client *grpc.Client,
	signerAddress flow.Address, // Address of the account that will mint (minter owner)
	signerPrivateKey crypto.PrivateKey,
	keyIndex int,
	recipientAddress flow.Address, // Address of the account to receive the NFT
	nftName string,
	nftDescription string,
	nftURL string,
	nftRole bool,
) (*flow.TransactionResult, error) {
	minterAccount, err := client.GetAccount(ctx, signerAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get minter account %s: %w", signerAddress.String(), err)
	}
	sequenceNumber := minterAccount.Keys[keyIndex].SequenceNumber

	latestBlock, err := client.GetLatestBlockHeader(ctx, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest block: %w", err)
	}

	tx := flow.NewTransaction().
		SetScript([]byte(mintNFTScript)).
		SetGasLimit(300). // Minting might require slightly more gas
		SetProposalKey(signerAddress, uint32(keyIndex), sequenceNumber).
		SetReferenceBlockID(latestBlock.ID).
		SetPayer(signerAddress).
		AddAuthorizer(signerAddress)

	// Add arguments in the order they are defined in the Cadence script
	// 1. recipient: Address
	cadenceRecipient := cadence.NewAddress(recipientAddress)
	if err := tx.AddArgument(cadenceRecipient); err != nil {
		return nil, fmt.Errorf("failed to add recipient argument: %w", err)
	}
	// 2. name: String
	cadenceName, err := cadence.NewString(nftName)
	if err != nil {
		return nil, fmt.Errorf("failed to create 'name' argument: %w", err)
	}
	if err := tx.AddArgument(cadenceName); err != nil {
		return nil, fmt.Errorf("failed to add 'name' argument: %w", err)
	}
	// 3. description: String
	cadenceDescription, err := cadence.NewString(nftDescription)
	if err != nil {
		return nil, fmt.Errorf("failed to create 'description' argument: %w", err)
	}
	if err := tx.AddArgument(cadenceDescription); err != nil {
		return nil, fmt.Errorf("failed to add 'description' argument: %w", err)
	}
	// 4. url: String
	cadenceURL, err := cadence.NewString(nftURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create 'url' argument: %w", err)
	}
	if err := tx.AddArgument(cadenceURL); err != nil {
		return nil, fmt.Errorf("failed to add 'url' argument: %w", err)
	}
    // 5. role: Bool
	cadenceRole := cadence.NewBool(nftRole)
	if err := tx.AddArgument(cadenceRole); err != nil {
		return nil, fmt.Errorf("failed to add 'role' argument: %w", err)
	}

	signer, err := crypto.NewInMemorySigner(signerPrivateKey, crypto.SHA3_256)
	if err != nil {
		return nil, fmt.Errorf("failed to create signer: %w", err)
	}
	if err := tx.SignEnvelope(signerAddress, uint32(keyIndex), signer); err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	if err := client.SendTransaction(ctx, *tx); err != nil {
		return nil, fmt.Errorf("failed to send transaction: %w", err)
	}
	log.Printf("Mint NFT Transaction sent (ID: %s) by %s for recipient %s", tx.ID(), signerAddress.String(), recipientAddress.String())
	return WaitForTransactionSeal(ctx, client, tx.ID())
}

// WaitForTransactionSeal waits for a transaction to be sealed.
func WaitForTransactionSeal(ctx context.Context, client *grpc.Client, txID flow.Identifier) (*flow.TransactionResult, error) {
	log.Printf("Waiting for transaction %s to be sealed...", txID)
	for {
		result, err := client.GetTransactionResult(ctx, txID)
		if err != nil {
			return nil, fmt.Errorf("failed to get transaction result for %s: %w", txID, err)
		}

		if result.Error != nil {
			log.Printf("Transaction %s execution error: %v", txID, result.Error)
			for _, event := range result.Events {
				log.Printf("Event from errored transaction %s: Type: %s, Values: %s", txID, event.Type, event.Value.String())
			}
		} else if result.Status == flow.TransactionStatusSealed {
			log.Printf("Transaction %s successfully sealed. Events:", txID)
			for _, event := range result.Events {
				log.Printf("  Event: Type: %s, Values: %s", event.Type, event.Value.String())
			}
		}

		switch result.Status {
		case flow.TransactionStatusSealed:
			log.Printf("Transaction %s sealed.", txID)
			return result, result.Error
		case flow.TransactionStatusExpired:
			log.Printf("Transaction %s expired.", txID)
			return nil, fmt.Errorf("transaction %s expired", txID)
		case flow.TransactionStatusPending, flow.TransactionStatusExecuted:
			log.Printf("Transaction %s status: %s. Waiting...", txID, result.Status)
		default:
			log.Printf("Transaction %s in unhandled status: %s. Waiting...", txID, result.Status)
		}

		select {
		case <-time.After(2 * time.Second):
		case <-ctx.Done():
			log.Printf("Context cancelled while waiting for transaction %s.", txID)
			return nil, ctx.Err()
		}
	}
}

func main() {
	ctx := context.Background()

	client, err := grpc.NewClient("access.devnet.nodes.onflow.org:9000")
	if err != nil {
		log.Fatalf("Failed to connect to Flow: %v", err)
	}
	defer client.Close()
	log.Println("Connected to Flow network.")

	// --- Signer (Minter) Configuration ---
	// This account must OWN the NFTMinter resource.
	// IMPORTANT: Replace with your actual minter private key and address
	minterPrivateKeyHex := "e8c708809c4bf218e3004ae667f77ef43878bccc7eac7b4fde0f71be37c8282d" // Example Minter PK
	minterAddress := flow.HexToAddress("0x4ac0ee1c903bf362")                                  // Example Minter Address
	minterKeyIndex := 0

	minterPrivateKey, err := crypto.DecodePrivateKeyHex(crypto.ECDSA_P256, minterPrivateKeyHex)
	if err != nil {
		log.Fatalf("Failed to decode minter private key: %v", err)
	}

	// --- Recipient Configuration ---
	// This account will RECEIVE the NFT. It must have a collection set up.
	// You might want to use a different account than the minter for the recipient.
	// For testing, you can use the same account if it also has a collection.
	// IMPORTANT: Replace with the recipient's address.
	recipientAddressForNFT := flow.HexToAddress("0x51846a0f69492bba") // Example Recipient Address
    // If the recipient is a different account and needs collection setup,
    // you'll need its private key to run SetupAccountForNFTs for it.
    // For this example, we assume recipient already has a collection or is the same as minter
    // and minter will have a collection setup first.

	// --- NFT Data ---
	nftName := "BOB is the best"
	nftDescription := "all around awesome guy"
	nftURL := "https://upload.wikimedia.org/wikipedia/commons/7/70/Example.png"
	nftRole := true

	// --- Step 2: Mint the NFT ---
	log.Printf("Attempting to mint NFT for recipient %s (minted by %s)...", recipientAddressForNFT.String(), minterAddress.String())
	mintResult, err := MintNFT(
		ctx,
		client,
		minterAddress,
		minterPrivateKey,
		minterKeyIndex,
		recipientAddressForNFT,
		nftName,
		nftDescription,
		nftURL,
		nftRole,
	)
	if err != nil {
		log.Fatalf("Failed to mint NFT (Go error, expiry, or context cancellation): %v", err)
	}
	if mintResult.Error != nil {
		log.Fatalf("NFT Minting Transaction %s failed (Cadence error): %v", mintResult.TransactionID, mintResult.Error)
	}

	log.Printf("NFT Minting Transaction successful! ID: %s, Status: %s", mintResult.TransactionID, mintResult.Status)
}
