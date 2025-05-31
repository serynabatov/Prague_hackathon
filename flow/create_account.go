package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"os"
	"time"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/client"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/onflow/flow-go-sdk/templates"
	"google.golang.org/grpc"
)

type FlowAccount struct {
	Address    string
	PrivateKey string
	PublicKey  string
}

func createFlowAccount(proposerAddress, proposerPrivateKeyHex string) (FlowAccount, error) {
	ctx := context.Background()

	flowClient, err := client.New("access.testnet.nodes.onflow.org:9000", grpc.WithInsecure())
	if err != nil {
		return FlowAccount{}, fmt.Errorf("failed to connect to testnet: %v", err)
	}
	defer flowClient.Close()

	// Parse proposer private key
	proposerPrivateKey, err := crypto.DecodePrivateKeyHex(crypto.ECDSA_P256, proposerPrivateKeyHex)
	if err != nil {
		return FlowAccount{}, fmt.Errorf("failed to decode proposer private key: %v", err)
	}
	signer, err := crypto.NewInMemorySigner(proposerPrivateKey, crypto.SHA3_256)
	if err != nil {
		return FlowAccount{}, fmt.Errorf("failed to create signer: %v", err)
	}

	// Generate new key pair
	seed := make([]byte, 32)
	if _, err := rand.Read(seed); err != nil {
		return FlowAccount{}, fmt.Errorf("failed to generate secure seed: %v", err)
	}
	newPrivateKey, err := crypto.GeneratePrivateKey(crypto.ECDSA_P256, seed)
	if err != nil {
		return FlowAccount{}, fmt.Errorf("failed to generate private key: %v", err)
	}
	newPublicKey := newPrivateKey.PublicKey()

	// Create account key
	accountKey := &flow.AccountKey{
		PublicKey: newPublicKey,
		SigAlgo:   crypto.ECDSA_P256,
		HashAlgo:  crypto.SHA3_256,
		Weight:    flow.AccountKeyWeightThreshold,
	}

	// Get proposer account
	proposerFlowAddr := flow.HexToAddress(proposerAddress)
	proposerAccount, err := flowClient.GetAccount(ctx, proposerFlowAddr)
	if err != nil {
		return FlowAccount{}, fmt.Errorf("failed to get proposer account: %v", err)
	}

	// Get latest block
	latestBlock, err := flowClient.GetLatestBlockHeader(ctx, true)
	if err != nil {
		return FlowAccount{}, fmt.Errorf("failed to get latest block: %v", err)
	}

	// Create transaction
	tx, err := templates.CreateAccount([]*flow.AccountKey{accountKey}, nil, proposerFlowAddr)
	if err != nil {
		return FlowAccount{}, fmt.Errorf("failed to create account transaction: %v", err)
	}
	tx.SetProposalKey(proposerAccount.Address, proposerAccount.Keys[0].Index, proposerAccount.Keys[0].SequenceNumber)
	tx.SetReferenceBlockID(latestBlock.ID)
	tx.SetPayer(proposerFlowAddr)

	// Sign transaction
	err = tx.SignEnvelope(proposerFlowAddr, proposerAccount.Keys[0].Index, signer)
	if err != nil {
		return FlowAccount{}, fmt.Errorf("failed to sign transaction: %v", err)
	}

	// Send transaction
	err = flowClient.SendTransaction(ctx, *tx)
	if err != nil {
		return FlowAccount{}, fmt.Errorf("failed to send transaction: %v", err)
	}

	// Wait for transaction to be sealed (with timeout)
    timeout := time.After(10 * time.Second)
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()

    var result *flow.TransactionResult
    for {
        select {
        case <-timeout:
            return FlowAccount{}, fmt.Errorf("transaction timed out")
        case <-ticker.C:
            result, err = flowClient.GetTransactionResult(ctx, tx.ID())
            if err != nil {
                return FlowAccount{}, fmt.Errorf("failed to get transaction result: %v", err)
            }
            if result.Status == flow.TransactionStatusSealed {
                goto TransactionSealed
            } else if result.Error != nil { // Проверяем наличие ошибки
                return FlowAccount{}, fmt.Errorf("transaction failed: %s", result.Error.Error())
            }
        }
    }

TransactionSealed:
	// Extract new account address
	var newAddress flow.Address
	for _, event := range result.Events {
		if event.Type == flow.EventAccountCreated {
			accountCreatedEvent := flow.AccountCreatedEvent(event)
			newAddress = accountCreatedEvent.Address()
			break
		}
	}
	if newAddress == flow.EmptyAddress {
		return FlowAccount{}, fmt.Errorf("failed to extract new account address from events")
	}

	return FlowAccount{
		Address:    newAddress.Hex(),
		PrivateKey: newPrivateKey.String(),
		PublicKey:  newPublicKey.String(),
	}, nil
}

func main() {
	proposerAddress := "447fc5c9baffdfd1"
	proposerPrivateKey := "45bfe485142dbbe1adc62b8dfc46f3efcaac92508466fb9b042ce75675ade71e"

	account, err := createFlowAccount(proposerAddress, proposerPrivateKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ New Flow Account Created:\n")
	fmt.Printf("Address: 0x%s\n", account.Address)
	fmt.Printf("Private Key: %s\n", account.PrivateKey)
	fmt.Printf("Public Key: %s\n", account.PublicKey)
}
