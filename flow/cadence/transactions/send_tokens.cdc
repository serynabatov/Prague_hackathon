import "FungibleToken"
import "FlowToken"

transaction(recipient: Address, amount: UFix64) {

    // Reference to the sender's Vault
    let senderVault: &FlowToken.Vault

    prepare(signer: AuthAccount) {
        // Borrow a reference to the signer's FlowToken Vault
        self.senderVault = signer.borrow<&FlowToken.Vault>(from: /storage/flowTokenVault)
            ?? panic("Could not borrow reference to the sender's Vault")
    }

    execute {
        // Get the recipient's public account
        let recipientAccount = getAccount(recipient)

        // Borrow a reference to the recipient's Receiver capability
        let receiver = recipientAccount
            .getCapability(/public/flowTokenReceiver)
            .borrow<&{FungibleToken.Receiver}>()
            ?? panic("Could not borrow receiver reference")

        // Withdraw tokens from sender vault
        let sentVault <- self.senderVault.withdraw(amount: amount)

        // Deposit tokens to the recipient
        receiver.deposit(from: <- sentVault)
        log("Tokens sent successfully!")
    }
}