// 0x4ac0ee1c903bf362
package custom_flow

import (
	"context"
	"fmt"
	"log"
	"regexp"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/access/grpc"
)

const getNFTs = `
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

type NFTDisplay struct {
	Name        string
	Description string
	Thumbnail   string
}

func parseNFTDisplays(value cadence.Value) ([]NFTDisplay, error) {
	// Получаем строковое представление
	str := value.String()
	fmt.Println("Parsing string:", str)

	// Regex для извлечения данных из строкового представления
	displayRegex := regexp.MustCompile(`MetadataViews\.Display\(name: "([^"]*)", description: "([^"]*)", thumbnail: [^(]*\(url: "([^"]*)"\)\)`)
	matches := displayRegex.FindAllStringSubmatch(str, -1)

	var nfts []NFTDisplay

	for _, match := range matches {
		if len(match) >= 4 {
			nft := NFTDisplay{
				Name:        match[1],
				Description: match[2],
				Thumbnail:   match[3],
			}
			nfts = append(nfts, nft)
		}
	}

	return nfts, nil
}

func GetNfts() ([]NFTDisplay, error) {
	ctx := context.Background()

	client, err := grpc.NewClient("access.devnet.nodes.onflow.org:9000")
	defer client.Close()

	if err != nil {
		log.Fatalf("Failed to connect to Flow: %v", err)
		return nil, err
	}

	userAddress := flow.HexToAddress("0x4ac0ee1c903bf362") // Replace with target address

	addressArg := cadence.NewAddress(userAddress)

	value, err := client.ExecuteScriptAtLatestBlock(ctx, []byte(getNFTs), []cadence.Value{addressArg})
	if err != nil {
		log.Fatalf("Failed to execute script: %v", err)
		return nil, err
	}

	nfts, err := parseNFTDisplays(value)
	if err != nil {
		log.Fatalf("Failed to parse NFT data: %v", err)
		return nil, err
	}

	return nfts, nil
}
