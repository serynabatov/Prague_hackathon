package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/access/grpc"
	"github.com/onflow/flow-go-sdk/crypto"
)

const setupCollectionScript = `
import FooBarV4 from 0x4ac0ee1c903bf362
import NonFungibleToken from 0x631e88ae7f1d7c20

transaction {

    prepare(signer: auth(BorrowValue, IssueStorageCapabilityController, PublishCapability, SaveValue, UnpublishCapability) &Account) {

        // Return early if the account already has a collection
        if signer.storage.borrow<&FooBarV4.Collection>(from: FooBarV4.CollectionStoragePath) != nil {
            log("Collection already exists, skipping creation.")
            return
        }

        // Create a new empty collection
        let collection <- FooBarV4.createEmptyCollection(nftType: Type<@FooBarV4.NFT>())

        // save it to the account
        signer.storage.save(<-collection, to: FooBarV4.CollectionStoragePath)
        log("Saved collection to storage.")

        // issue a capability and publish it
        let collectionCap = signer.capabilities.storage.issue<&FooBarV4.Collection>(FooBarV4.CollectionStoragePath)
        signer.capabilities.publish(collectionCap, at: FooBarV4.CollectionPublicPath)
        log("Published collection capability.")
    }

    execute {
        log("Transaction executed successfully.")
    }
}
`

// SetupAccountForNFTs prepares an account to receive FooBarV4 NFTs by creating a collection.
func SetupAccountForNFTs(
	ctx context.Context,
	client *grpc.Client,
	signerAddress flow.Address,
	signerPrivateKey crypto.PrivateKey,
	keyIndex int,
) (*flow.TransactionResult, error) {
	// 1. Get the signer account to retrieve the latest sequence number
	signerAccount, err := client.GetAccount(ctx, signerAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get signer account: %w", err)
	}
	sequenceNumber := signerAccount.Keys[keyIndex].SequenceNumber

	// 2. Get the latest sealed block as the reference block
	latestBlock, err := client.GetLatestBlockHeader(ctx, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest block: %w", err)
	}

	// 3. Create the transaction
	tx := flow.NewTransaction().
		SetScript([]byte(setupCollectionScript)).
		SetGasLimit(200). // Increased gas limit slightly, adjust if needed
		SetProposalKey(
			signerAddress,
			uint32(keyIndex),
			sequenceNumber,
		).
		SetReferenceBlockID(latestBlock.ID).
		SetPayer(signerAddress).
		AddAuthorizer(signerAddress) // The signer is the only authorizer needed

	// This transaction takes no arguments, so no tx.AddArgument() calls are needed.

	// 4. Sign the transaction
	signer, err := crypto.NewInMemorySigner(signerPrivateKey, crypto.SHA3_256)
	if err != nil {
		return nil, fmt.Errorf("failed to create signer: %w", err)
	}

	if err := tx.SignEnvelope(signerAddress, uint32(keyIndex), signer); err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	// 5. Send the transaction
	if err := client.SendTransaction(ctx, *tx); err != nil {
		return nil, fmt.Errorf("failed to send transaction: %w", err)
	}
	log.Printf("Transaction sent with ID: %s", tx.ID())

	// 6. Wait for the transaction to be sealed
	return WaitForTransactionSeal(ctx, client, tx.ID())
}

// WaitForTransactionSeal waits for a transaction to be sealed.
func WaitForTransactionSeal(ctx context.Context, client *grpc.Client, txID flow.Identifier) (*flow.TransactionResult, error) {
	log.Printf("Waiting for transaction %s to be sealed...", txID)

	for {
		result, err := client.GetTransactionResult(ctx, txID)
		if err != nil {
			return nil, fmt.Errorf("failed to get transaction result: %w", err)
		}

		if result.Error != nil {
			// Log Cadence execution errors
			log.Printf("Transaction %s execution error: %v", txID, result.Error)
			// It's often useful to see the events even if there's an error
			for _, event := range result.Events {
				log.Printf("Event from errored transaction: %s", event.String())
			}
		} else {
			// Log events for successful transactions
			for _, event := range result.Events {
				log.Printf("Event: %s", event.String())
			}
		}


		switch result.Status {
		case flow.TransactionStatusSealed:
			log.Printf("Transaction %s sealed.", txID)
			return result, result.Error // Return the Cadence error if present
		case flow.TransactionStatusExpired:
			return nil, fmt.Errorf("transaction %s expired", txID)
		default:
			log.Printf("Transaction %s in unknown status: %s", txID, result.Status)
		}


		select {
		case <-time.After(2 * time.Second): // Check every 2 seconds
			// Continue polling
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func main() {
	ctx := context.Background()

	// 1. Connection to Flow (Devnet in this example)
	// Ensure this matches the network where FooBarV4 and NonFungibleToken are deployed
	client, err := grpc.NewClient("access.devnet.nodes.onflow.org:9000")
	// For local emulator:
	// client, err := grpc.NewClient("127.0.0.1:3569")
	if err != nil {
		log.Fatalf("Failed to connect to Flow: %v", err)
	}
	defer client.Close()
	log.Println("Connected to Flow network.")

	// 2. Signer configuration
	// Replace with your actual private key and address for the account that will have the collection
	signerPrivateKeyHex := "71faf2f567d303ef4ae83cd1765a3fe125f8e5ae117056f600a6dd8970f81fff" // IMPORTANT: Replace with your private key
	signerAddress := flow.HexToAddress("0x51846a0f69492bba")                                  // IMPORTANT: Replace with your address
	keyIndex := 0

	// 3. Decode the private key
	privateKey, err := crypto.DecodePrivateKeyHex(crypto.ECDSA_P256, signerPrivateKeyHex)
	if err != nil {
		log.Fatalf("Failed to decode private key: %v", err)
	}

	// 4. Execute the transaction to setup the NFT collection
	log.Printf("Attempting to setup NFT collection for account %s...", signerAddress.String())
	result, err := SetupAccountForNFTs(
		ctx,
		client,
		signerAddress,
		privateKey,
		keyIndex,
	)
	if err != nil {
		// This catches Go-level errors (network, signing, etc.) or transaction expiry
		log.Fatalf("Failed to setup NFT collection (Go error or expiry): %v", err)
	}

	// The result.Error from WaitForTransactionSeal will contain Cadence execution errors
	if result.Error != nil {
		log.Fatalf("Transaction %s failed (Cadence error): %v", result.TransactionID, result.Error)
	}

	log.Printf("Transaction to setup NFT collection successful! ID: %s, Status: %s", result.TransactionID, result.Status)
}
