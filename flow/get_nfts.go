// 0x4ac0ee1c903bf362
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/access/grpc"
)

const getNFTsScript = `
import FooBarV4 from 0x4ac0ee1c903bf362
import NonFungibleToken from 0x631e88ae7f1d7c20
import MetadataViews from 0x631e88ae7f1d7c20

access(all) fun main(address: Address): [MetadataViews.Display] {
    let account = getAccount(address)

    let collectionRef = account.capabilities.borrow<&{NonFungibleToken.Collection}>(
            FooBarV4.CollectionPublicPath
        ) ?? panic("Could not borrow capability from collection at specified path")

    let ids = collectionRef.getIDs()

    let result: [MetadataViews.Display] = []

    for id in ids {
        let nftRef = collectionRef.borrowNFT(id)
            ?? panic("NFT not found in collection")

        let display = nftRef.resolveView(Type<MetadataViews.Display>())
            ?? panic("Display view not available")

        result.append(display as! MetadataViews.Display)
    }

    return result
}
`

func main() {
	ctx := context.Background()

	client, err := grpc.NewClient("access.devnet.nodes.onflow.org:9000")
	if err != nil {
		log.Fatalf("Failed to connect to Flow: %v", err)
	}
	defer client.Close()

	userAddress := flow.HexToAddress("0x4ac0ee1c903bf362") // Replace with target address

	addressArg := cadence.NewAddress(userAddress)

	value, err := client.ExecuteScriptAtLatestBlock(ctx, []byte(getNFTsScript), []cadence.Value{addressArg})
	if err != nil {
		log.Fatalf("Failed to execute script: %v", err)
	}

	// Просто выводим полученные данные как есть
	fmt.Println("Raw NFT data:")
	fmt.Println(value.String())
}
