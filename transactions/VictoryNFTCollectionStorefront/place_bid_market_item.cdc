import VictoryNFTCollectionStorefront from "../../contracts/VictoryNFTCollectionStorefront.cdc"
import FungibleToken from "../../contracts/FungibleToken.cdc"
import FUSD from "../../contracts/FUSD.cdc"

transaction(bundleID: UInt64, bidPrice: UFix64, bidder: Address, marketCollectionAddress: Address) {
    let marketCollection: &VictoryNFTCollectionStorefront.Collection{VictoryNFTCollectionStorefront.CollectionPublic}
    let bidderReceiver: Capability<&{FungibleToken.Receiver}>
    let bidVault: @FungibleToken.Vault

    prepare(signer: AuthAccount) {
        self.marketCollection = getAccount(marketCollectionAddress)
            .getCapability<&VictoryNFTCollectionStorefront.Collection{VictoryNFTCollectionStorefront.CollectionPublic}>(
                VictoryNFTCollectionStorefront.CollectionPublicPath
            )!
            .borrow()
            ?? panic("Could not borrow market collection from market address")

        // temporarily borrow the bid amount
        let mainFUSDVault = signer.borrow<&FUSD.Vault>(from: /storage/fusdVault)
             ?? panic("Cannot borrow FUSD vault from acct storage")
        self.bidVault <- mainFUSDVault.withdraw(amount: bidPrice)

        // pass in the receiver so we can deposit the bid amount back
        self.bidderReceiver = signer.getCapability<&{FungibleToken.Receiver}>(/public/fusdReceiver)
    }

    execute {
        self.marketCollection.placeBid(
            bundleID: bundleID,
            bidPrice: bidPrice,
            bidder: bidder,
            bidderReceiver: self.bidderReceiver,
            bidVault: <- self.bidVault
        )
    }
}
