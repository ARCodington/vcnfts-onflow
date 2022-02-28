import FungibleToken from "../../contracts/FungibleToken.cdc"
import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import FUSD from "../../contracts/FUSD.cdc"
import VictoryNFTCollectionItem from "../../contracts/VictoryNFTCollectionItem.cdc"
import VictoryNFTCollectionStorefront from "../../contracts/VictoryNFTCollectionStorefront.cdc"

transaction(sellerAddress: Address, royaltyAddress: Address, 
            itemIDs: [UInt64], type: UInt8, price: UFix64, startTime: UInt32, endTime: UInt32, targetPrice: UFix64, 
            royaltyFactor: UFix64) 
{
    let victoryBundler: Capability<&VictoryNFTCollectionItem.Collection{NonFungibleToken.Provider, VictoryNFTCollectionItem.VictoryNFTCollectionItemBundle}>
    let victoryCollection: Capability<&VictoryNFTCollectionItem.Collection{NonFungibleToken.Provider, VictoryNFTCollectionItem.VictoryNFTCollectionItemCollectionPublic}>
    let marketCollection: &VictoryNFTCollectionStorefront.Collection
    let sellerPaymentReceiver: Capability<&{FungibleToken.Receiver}>
    let royaltyPaymentReceiver: Capability<&{FungibleToken.Receiver}>

    prepare(signer: AuthAccount) {
        // we need a provider capability, but one is not provided by default so we create one.
        let bundlePath = /private/VictoryNFTCollectionItemBundle

        if !signer.getCapability<&VictoryNFTCollectionItem.Collection{NonFungibleToken.Provider, VictoryNFTCollectionItem.VictoryNFTCollectionItemBundle}>(bundlePath)!.check() {
            signer.link<&VictoryNFTCollectionItem.Collection{NonFungibleToken.Provider, VictoryNFTCollectionItem.VictoryNFTCollectionItemBundle}>(bundlePath, target: VictoryNFTCollectionItem.CollectionStoragePath)
        }

        self.victoryBundler = signer.getCapability<&VictoryNFTCollectionItem.Collection{NonFungibleToken.Provider, VictoryNFTCollectionItem.VictoryNFTCollectionItemBundle}>(bundlePath)!
        assert(self.victoryBundler.borrow() != nil, message: "Missing or mis-typed VictoryNFTCollectionItemBundle provider")

        // we need a provider capability, but one is not provided by default so we create one.
        let VictoryNFTCollectionItemCollectionProviderPrivatePath = /private/VictoryNFTCollectionItemCollectionProvider

        if !signer.getCapability<&VictoryNFTCollectionItem.Collection{NonFungibleToken.Provider, VictoryNFTCollectionItem.VictoryNFTCollectionItemCollectionPublic}>(VictoryNFTCollectionItemCollectionProviderPrivatePath)!.check() {
            signer.link<&VictoryNFTCollectionItem.Collection{NonFungibleToken.Provider, VictoryNFTCollectionItem.VictoryNFTCollectionItemCollectionPublic}>(VictoryNFTCollectionItemCollectionProviderPrivatePath, target: VictoryNFTCollectionItem.CollectionStoragePath)
        }

        self.victoryCollection = signer.getCapability<&VictoryNFTCollectionItem.Collection{NonFungibleToken.Provider, VictoryNFTCollectionItem.VictoryNFTCollectionItemCollectionPublic}>(VictoryNFTCollectionItemCollectionProviderPrivatePath)!
        assert(self.victoryCollection.borrow() != nil, message: "Missing or mis-typed VictoryNFTCollectionItemCollection provider")

        self.marketCollection = signer.borrow<&VictoryNFTCollectionStorefront.Collection>(from: VictoryNFTCollectionStorefront.CollectionStoragePath)
            ?? panic("Missing or mis-typed VictoryNFTCollectionStorefront Collection")

        let sellerAccount = getAccount(sellerAddress)
        self.sellerPaymentReceiver = sellerAccount.getCapability<&FUSD.Vault{FungibleToken.Receiver}>(/public/fusdReceiver)!
        assert(self.sellerPaymentReceiver.borrow() != nil, message: "Missing or mis-typed FUSD receiver for seller")

        // link to royalty receiver
        let royaltyAccount = getAccount(royaltyAddress)
        self.royaltyPaymentReceiver = royaltyAccount.getCapability<&FUSD.Vault{FungibleToken.Receiver}>(/public/fusdReceiver)!
        assert(self.royaltyPaymentReceiver.borrow() != nil, message: "Missing or mis-typed FUSD receiver for royalty")
    }

    execute {
        let bundleID = self.victoryBundler.borrow()!.createBundle(itemIDs: itemIDs)
        
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
    }
}
