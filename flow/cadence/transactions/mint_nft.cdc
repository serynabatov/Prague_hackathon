import "NonFungibleToken"
import "FooBar"

transaction(
    recipient: Address
) {

    /// local variable for storing the minter reference
    let minter: &FooBar.NFTMinter

    /// Reference to the receiver's collection
    let recipientCollectionRef: &{NonFungibleToken.Receiver}

    prepare(signer: auth(BorrowValue) &Account) {
        
        // borrow a reference to the NFTMinter resource in storage
        self.minter = signer.storage.borrow<&FooBar.NFTMinter>(from: FooBar.MinterStoragePath)
            ?? panic("Account does not store an object at the specified path")

        // Borrow the recipient's public NFT collection reference
        self.recipientCollectionRef = getAccount(recipient).capabilities.borrow<&{NonFungibleToken.Receiver}>(
                FooBar.CollectionPublicPath
            ) ?? panic("Could not get receiver reference to the NFT Collection")
    }

    execute {
        // Mint the NFT and deposit it to the recipient's collection
        let mintedNFT <- self.minter.createNFT(name: "FooBar NFT DIVERSO", description: "A unique FooBar NFT",url: "https://example.com/foobar-nft")
        self.recipientCollectionRef.deposit(token: <-mintedNFT)
    }
}