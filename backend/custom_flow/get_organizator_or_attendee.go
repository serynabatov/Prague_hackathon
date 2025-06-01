package custom_flow

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/access/grpc"
)

type NFTRecord struct {
	ID            uint64
	Name          string
	Description   string
	Thumbnail     string
	OrganizerBool bool
}

const getNFTsScript = `
import FooBarV4 from 0x4ac0ee1c903bf362
import NonFungibleToken from 0x631e88ae7f1d7c20
import MetadataViews from 0x631e88ae7f1d7c20

access(all) fun main(address: Address): [[AnyStruct]] {
    let account = getAccount(address)

    let collectionRef = account.capabilities.borrow<&{NonFungibleToken.Collection}>(
            FooBarV4.CollectionPublicPath
        ) ?? panic("Could not borrow capability from collection at specified path")

    let ids = collectionRef.getIDs()
    let result: [[AnyStruct]] = []

    for id in ids {
        let nftRef = collectionRef.borrowNFT(id)
            ?? panic("NFT not found in collection")

        // Получаем Display данные
        let display = nftRef.resolveView(Type<MetadataViews.Display>())
            ?? panic("Display view not available")

        // Получаем organizerBool из NFT
        let organizerBool = (nftRef as! &FooBarV4.NFT).organizerBool

        // Создаем запись для этого NFT
        let nftRecord: [AnyStruct] = [
            id,
            display,
            organizerBool
        ]

        result.append(nftRecord)
    }

    return result
}
`

func parseNFTRecords(value cadence.Value) ([]NFTRecord, error) {
	str := value.String()
	fmt.Println("Parsing string:", str)

	// Улучшенное регулярное выражение для обработки вывода
	recordRegex := regexp.MustCompile(`\[(\d+),\s*A\.631e88ae7f1d7c20\.MetadataViews\.Display\(name:\s*"([^"]*)",\s*description:\s*"([^"]*)",\s*thumbnail:\s*A\.631e88ae7f1d7c20\.MetadataViews\.HTTPFile\(url:\s*"([^"]*)"\)\),\s*(true|false)\]`)
	matches := recordRegex.FindAllStringSubmatch(str, -1)

	var records []NFTRecord

	for _, match := range matches {
		if len(match) >= 6 {
			id, err := strconv.ParseUint(match[1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse NFT ID: %v", err)
			}

			organizerBool, err := strconv.ParseBool(match[5])
			if err != nil {
				return nil, fmt.Errorf("failed to parse organizerBool: %v", err)
			}

			record := NFTRecord{
				ID:            id,
				Name:          match[2],
				Description:   match[3],
				Thumbnail:     match[4],
				OrganizerBool: organizerBool,
			}
			records = append(records, record)
		}
	}

	return records, nil
}

func GetOrganizatorOrAttendee() ([]NFTRecord, error) {
	ctx := context.Background()

	client, err := grpc.NewClient("access.devnet.nodes.onflow.org:9000")
	if err != nil {
		log.Fatalf("Failed to connect to Flow: %v", err)
		return nil, err
	}
	defer client.Close()

	userAddress := flow.HexToAddress("0x4ac0ee1c903bf362") // Replace with target address

	addressArg := cadence.NewAddress(userAddress)

	value, err := client.ExecuteScriptAtLatestBlock(ctx, []byte(getNFTsScript), []cadence.Value{addressArg})
	if err != nil {
		log.Fatalf("Failed to execute script: %v", err)
		return nil, err
	}

	records, err := parseNFTRecords(value)
	if err != nil {
		log.Fatalf("Failed to parse NFT records: %v", err)
		return nil, err
	}

	return records, err
}
