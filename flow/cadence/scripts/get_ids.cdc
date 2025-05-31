import "NonFungibleToken"
import "FooBar"

access(all) fun main(address: Address): [UInt64] {
    let account = getAccount(address)

    let collectionRef = account.capabilities.borrow<&{NonFungibleToken.Collection}>(
            FooBar.CollectionPublicPath
        ) ?? panic("Could not borrow capability from collection at specified path")

    return collectionRef.getIDs()
}