import "FooBar" 
import "NonFungibleToken" 
import "MetadataViews"

access(all) fun main(address: Address): [MetadataViews.Display] {
    let account = getAccount(address)

    let collectionRef = account.capabilities.borrow<&{NonFungibleToken.Collection}>(
            FooBar.CollectionPublicPath
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