import FungibleToken from "../../contracts/FungibleToken.cdc"
import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import FUSD from "../../contracts/FUSD.cdc"
import VictoryCollectible from "../../contracts/VictoryCollectible.cdc"
import VictoryCollectibleSaleOffer from "../../contracts/VictoryCollectibleSaleOffer.cdc"

transaction(royaltyAddress: Address, bundleID: UInt64, type: UInt8, price: UFix64, startTime: UInt32, endTime: UInt32, targetPrice: UFix64, royaltyFactor: UFix64) {
    let fusdReceiver: Capability<&FUSD.Vault{FungibleToken.Receiver}>
    let VictoryCollectibleProvider: Capability<&VictoryCollectible.Collection{NonFungibleToken.Provider, VictoryCollectible.VictoryCollectibleCollectionPublic}>
    let marketCollection: &VictoryCollectibleSaleOffer.Collection
    let royaltyPaymentReceiver: Capability<&FUSD.Vault{FungibleToken.Receiver}>

    prepare(signer: AuthAccount) {
        // we need a provider capability, but one is not provided by default so we create one.
        let VictoryCollectibleCollectionProviderPrivatePath = /private/VictoryCollectibleCollectionProvider

        if !signer.getCapability<&VictoryCollectible.Collection{NonFungibleToken.Provider, VictoryCollectible.VictoryCollectibleCollectionPublic}>(VictoryCollectibleCollectionProviderPrivatePath)!.check() {
            signer.link<&VictoryCollectible.Collection{NonFungibleToken.Provider, VictoryCollectible.VictoryCollectibleCollectionPublic}>(VictoryCollectibleCollectionProviderPrivatePath, target: VictoryCollectible.CollectionStoragePath)
        }

        self.VictoryCollectibleProvider = signer.getCapability<&VictoryCollectible.Collection{NonFungibleToken.Provider, VictoryCollectible.VictoryCollectibleCollectionPublic}>(VictoryCollectibleCollectionProviderPrivatePath)!
        assert(self.VictoryCollectibleProvider.borrow() != nil, message: "Missing or mis-typed VictoryCollectibleCollection provider")

        self.marketCollection = signer.borrow<&VictoryCollectibleSaleOffer.Collection>(from: VictoryCollectibleSaleOffer.CollectionStoragePath)
            ?? panic("Missing or mis-typed VictoryCollectibleSaleOffer Collection")

        self.fusdReceiver = signer.getCapability<&FUSD.Vault{FungibleToken.Receiver}>(/public/fusdReceiver)!
        assert(self.fusdReceiver.borrow() != nil, message: "Missing or mis-typed FUSD receiver")

        // link to royalty receiver
        let royaltyAccount = getAccount(royaltyAddress)
        self.royaltyPaymentReceiver = royaltyAccount.getCapability<&FUSD.Vault{FungibleToken.Receiver}>(/public/fusdReceiver)!
    }

    execute {
        let offer <- self.marketCollection.createSaleOffer (
            sellerItemProvider: self.VictoryCollectibleProvider,
            bundleID: bundleID,
            sellerPaymentReceiver: self.fusdReceiver,
            price: price,
            saleType: type,
            startTime: startTime,
            endTime: endTime,
            targetPrice: targetPrice,
            royaltyPaymentReceiver: self.royaltyPaymentReceiver,
            royaltyFactor: royaltyFactor
        )
        self.marketCollection.insert(offer: <-offer)
    }
}
