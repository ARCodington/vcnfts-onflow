import FungibleToken from "../../contracts/FungibleToken.cdc"
import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import FUSD from "../../contracts/FUSD.cdc"
import VictoryNFTCollectionItem from "../../contracts/VictoryNFTCollectionItem.cdc"
import VictoryNFTCollectionStorefront from "../../contracts/VictoryNFTCollectionStorefront.cdc"

transaction(bundleID: UInt64, marketCollectionAddress: Address) {
    let paymentVault: @FungibleToken.Vault
    let victoryCollection: &VictoryNFTCollectionItem.Collection{NonFungibleToken.Receiver}
    let marketCollection: &VictoryNFTCollectionStorefront.Collection{VictoryNFTCollectionStorefront.CollectionPublic}
    let ownerVault: Capability<&{FungibleToken.Receiver}>

    prepare(signer: AuthAccount) {
        self.marketCollection = getAccount(marketCollectionAddress)
            .getCapability<&VictoryNFTCollectionStorefront.Collection{VictoryNFTCollectionStorefront.CollectionPublic}>(
                VictoryNFTCollectionStorefront.CollectionPublicPath
            )!
            .borrow()
            ?? panic("Could not borrow market collection from market address")

        let saleItem = self.marketCollection.borrowSaleItem(bundleID: bundleID)
                    ?? panic("No item with that ID")
        let price = saleItem.price

        let ownerAccount = getAccount(saleItem.originalOwner)
        self.ownerVault = ownerAccount.getCapability<&{FungibleToken.Receiver}>(/public/fusdReceiver)

        let mainFUSDVault = signer.borrow<&FUSD.Vault>(from: /storage/fusdVault)
            ?? panic("Cannot borrow FUSD vault from acct storage")
        self.paymentVault <- mainFUSDVault.withdraw(amount: price)
 
        // we need a provider capability, but one is not provided by default so we create one.
        let VictoryNFTCollectionItemCollectionProviderPrivatePath = /private/VictoryNFTCollectionItemCollectionProvider

        if !signer.getCapability<&VictoryNFTCollectionItem.Collection{NonFungibleToken.Provider, VictoryNFTCollectionItem.VictoryNFTCollectionItemCollectionPublic}>(VictoryNFTCollectionItemCollectionProviderPrivatePath)!.check() {
            signer.link<&VictoryNFTCollectionItem.Collection{NonFungibleToken.Provider, VictoryNFTCollectionItem.VictoryNFTCollectionItemCollectionPublic}>(VictoryNFTCollectionItemCollectionProviderPrivatePath, target: VictoryNFTCollectionItem.CollectionStoragePath)
        }
        self.victoryCollection = signer.borrow<&VictoryNFTCollectionItem.Collection{NonFungibleToken.Receiver}>(
            from: VictoryNFTCollectionItem.CollectionStoragePath
        ) ?? panic("Cannot borrow VictoryNFTCollectionItem collection receiver from acct")
    }

    execute {
        self.marketCollection.purchase(
            bundleID: bundleID,
            buyerCollection: self.victoryCollection,
            buyerPayment: <- self.paymentVault,
            ownerPaymentReceiver: self.ownerVault
        )
    }
}
