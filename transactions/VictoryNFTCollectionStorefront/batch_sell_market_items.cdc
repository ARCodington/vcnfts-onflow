import FungibleToken from "../../contracts/FungibleToken.cdc"
import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import FUSD from "../../contracts/FUSD.cdc"
import VictoryCollectible from "../../contracts/VictoryCollectible.cdc"
import VictoryCollectibleSaleOffer from "../../contracts/VictoryCollectibleSaleOffer.cdc"

transaction(sellerAddress: Address, royaltyAddress: Address, startID: UInt64, endID: UInt64, 
            type: UInt8, price: UFix64, startTime: UInt32, endTime: UInt32, targetPrice: UFix64, 
            royaltyFactor: UFix64) 
{
    let victoryBundler: Capability<&VictoryCollectible.Collection{NonFungibleToken.Provider, VictoryCollectible.VictoryCollectibleBundle}>
    let victoryCollection: Capability<&VictoryCollectible.Collection{NonFungibleToken.Provider, VictoryCollectible.VictoryCollectibleCollectionPublic}>
    let marketCollection: &VictoryCollectibleSaleOffer.Collection
    let sellerPaymentReceiver: Capability<&{FungibleToken.Receiver}>
    let royaltyPaymentReceiver: Capability<&{FungibleToken.Receiver}>

    prepare(signer: AuthAccount) {
        // we need a provider capability, but one is not provided by default so we create one.
        let bundlePath = /private/VictoryCollectibleBundle

        if !signer.getCapability<&VictoryCollectible.Collection{NonFungibleToken.Provider, VictoryCollectible.VictoryCollectibleBundle}>(bundlePath)!.check() {
            signer.link<&VictoryCollectible.Collection{NonFungibleToken.Provider, VictoryCollectible.VictoryCollectibleBundle}>(bundlePath, target: VictoryCollectible.CollectionStoragePath)
        }

        self.victoryBundler = signer.getCapability<&VictoryCollectible.Collection{NonFungibleToken.Provider, VictoryCollectible.VictoryCollectibleBundle}>(bundlePath)!
        assert(self.victoryBundler.borrow() != nil, message: "Missing or mis-typed VictoryCollectibleBundle provider")

        // we need a provider capability, but one is not provided by default so we create one.
        let providerPath = /private/VictoryCollectibleCollectionProvider

        if !signer.getCapability<&VictoryCollectible.Collection{NonFungibleToken.Provider, VictoryCollectible.VictoryCollectibleCollectionPublic}>(providerPath)!.check() {
            signer.link<&VictoryCollectible.Collection{NonFungibleToken.Provider, VictoryCollectible.VictoryCollectibleCollectionPublic}>(providerPath, target: VictoryCollectible.CollectionStoragePath)
        }

        self.victoryCollection = signer.getCapability<&VictoryCollectible.Collection{NonFungibleToken.Provider, VictoryCollectible.VictoryCollectibleCollectionPublic}>(providerPath)!
        assert(self.victoryCollection.borrow() != nil, message: "Missing or mis-typed VictoryCollectibleCollection provider")

        self.marketCollection = signer.borrow<&VictoryCollectibleSaleOffer.Collection>(from: VictoryCollectibleSaleOffer.CollectionStoragePath)
            ?? panic("Missing or mis-typed VictoryCollectibleSaleOffer Collection")

        let sellerAccount = getAccount(sellerAddress)
        self.sellerPaymentReceiver = sellerAccount.getCapability<&FUSD.Vault{FungibleToken.Receiver}>(/public/fusdReceiver)!
        assert(self.sellerPaymentReceiver.borrow() != nil, message: "Missing or mis-typed FUSD receiver for seller")

        // link to royalty receiver
        let royaltyAccount = getAccount(royaltyAddress)
        self.royaltyPaymentReceiver = royaltyAccount.getCapability<&FUSD.Vault{FungibleToken.Receiver}>(/public/fusdReceiver)!
        assert(self.royaltyPaymentReceiver.borrow() != nil, message: "Missing or mis-typed FUSD receiver for royalty")
    }

    execute {
        var itemID: UInt64 = startID;

        // create a series of identical sale offers, one per item - used to post "drops" of items at the same starting point
        while itemID <= endID {
            let bundleID = self.victoryBundler.borrow()!.createBundle(itemIDs: [itemID])

            let offer <- self.marketCollection.createSaleOffer (
                sellerItemProvider: self.victoryCollection,
                bundleID: bundleID,
                sellerPaymentReceiver: self.sellerPaymentReceiver,
                price: price,
                saleType: type,
                startTime: startTime,
                endTime: endTime,
                targetPrice: targetPrice,
                royaltyPaymentReceiver: self.royaltyPaymentReceiver,
                royaltyFactor: royaltyFactor
            )
            self.marketCollection.insert(offer: <-offer)

            itemID = itemID + (1 as UInt64)
        }
    }
}
