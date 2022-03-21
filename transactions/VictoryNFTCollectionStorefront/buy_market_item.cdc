import FungibleToken from "../../contracts/FungibleToken.cdc"
import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import FUSD from "../../contracts/FUSD.cdc"
import VictoryCollectible from "../../contracts/VictoryCollectible.cdc"
import VictoryCollectibleSaleOffer from "../../contracts/VictoryCollectibleSaleOffer.cdc"

transaction(bundleID: UInt64, marketCollectionAddress: Address) {
    let paymentVault: @FungibleToken.Vault
    let victoryCollection: &VictoryCollectible.Collection{NonFungibleToken.Receiver}
    let marketCollection: &VictoryCollectibleSaleOffer.Collection{VictoryCollectibleSaleOffer.CollectionPublic}
    let ownerVault: Capability<&{FungibleToken.Receiver}>

    prepare(signer: AuthAccount) {
        self.marketCollection = getAccount(marketCollectionAddress)
            .getCapability<&VictoryCollectibleSaleOffer.Collection{VictoryCollectibleSaleOffer.CollectionPublic}>(
                VictoryCollectibleSaleOffer.CollectionPublicPath
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
        let VictoryCollectibleCollectionProviderPrivatePath = /private/VictoryCollectibleCollectionProvider

        if !signer.getCapability<&VictoryCollectible.Collection{NonFungibleToken.Provider, VictoryCollectible.VictoryCollectibleCollectionPublic}>(VictoryCollectibleCollectionProviderPrivatePath)!.check() {
            signer.link<&VictoryCollectible.Collection{NonFungibleToken.Provider, VictoryCollectible.VictoryCollectibleCollectionPublic}>(VictoryCollectibleCollectionProviderPrivatePath, target: VictoryCollectible.CollectionStoragePath)
        }
        self.victoryCollection = signer.borrow<&VictoryCollectible.Collection{NonFungibleToken.Receiver}>(
            from: VictoryCollectible.CollectionStoragePath
        ) ?? panic("Cannot borrow VictoryCollectible collection receiver from acct")
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
