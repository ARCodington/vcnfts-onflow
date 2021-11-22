import FungibleToken from "../../contracts/FungibleToken.cdc"
import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import FUSD from "../../contracts/FUSD.cdc"
import VictoryNFTCollectionItem from "../../contracts/VictoryNFTCollectionItem.cdc"
import VictoryNFTCollectionStorefront from "../../contracts/VictoryNFTCollectionStorefront.cdc"

transaction(royaltyAddress: Address, startID: UInt64, endID: UInt64, type: UInt8, price: UFix64, startTime: UInt32, endTime: UInt32, targetPrice: UFix64, royaltyFactor: UFix64) {
    let fusdReceiver: Capability<&FUSD.Vault{FungibleToken.Receiver}>
    let victoryCollection: Capability<&VictoryNFTCollectionItem.Collection{NonFungibleToken.Provider, VictoryNFTCollectionItem.VictoryNFTCollectionItemCollectionPublic}>
    let marketCollection: &VictoryNFTCollectionStorefront.Collection
    let royaltyPaymentReceiver: Capability<&FUSD.Vault{FungibleToken.Receiver}>

    prepare(signer: AuthAccount) {
        // we need a provider capability, but one is not provided by default so we create one.
        let VictoryNFTCollectionItemCollectionProviderPrivatePath = /private/VictoryNFTCollectionItemCollectionProvider

        if !signer.getCapability<&VictoryNFTCollectionItem.Collection{NonFungibleToken.Provider, VictoryNFTCollectionItem.VictoryNFTCollectionItemCollectionPublic}>(VictoryNFTCollectionItemCollectionProviderPrivatePath)!.check() {
            signer.link<&VictoryNFTCollectionItem.Collection{NonFungibleToken.Provider, VictoryNFTCollectionItem.VictoryNFTCollectionItemCollectionPublic}>(VictoryNFTCollectionItemCollectionProviderPrivatePath, target: VictoryNFTCollectionItem.CollectionStoragePath)
        }

        self.victoryCollection = signer.getCapability<&VictoryNFTCollectionItem.Collection{NonFungibleToken.Provider, VictoryNFTCollectionItem.VictoryNFTCollectionItemCollectionPublic}>(VictoryNFTCollectionItemCollectionProviderPrivatePath)!
        assert(self.victoryCollection.borrow() != nil, message: "Missing or mis-typed VictoryNFTCollectionItemCollection provider")

        self.marketCollection = signer.borrow<&VictoryNFTCollectionStorefront.Collection>(from: VictoryNFTCollectionStorefront.CollectionStoragePath)
            ?? panic("Missing or mis-typed VictoryNFTCollectionStorefront Collection")

        self.fusdReceiver = signer.getCapability<&FUSD.Vault{FungibleToken.Receiver}>(/public/fusdReceiver)!
        assert(self.fusdReceiver.borrow() != nil, message: "Missing or mis-typed FUSD receiver")

        // link to royalty receiver
        let royaltyAccount = getAccount(royaltyAddress)
        self.royaltyPaymentReceiver = royaltyAccount.getCapability<&FUSD.Vault{FungibleToken.Receiver}>(/public/fusdReceiver)!
    }

    execute {
        var itemID: UInt64 = startID;

        // create a series of identical sale offers, one per item - used to post "drops" of items at the same starting point
        while itemID <= endID {
            let bundleID = self.victoryCollection.borrow()!.createBundle(itemIDs: [itemID])

            let offer <- self.marketCollection.createSaleOffer (
                sellerItemProvider: self.victoryCollection,
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

            itemID = itemID + (1 as UInt64)
        }
    }
}
