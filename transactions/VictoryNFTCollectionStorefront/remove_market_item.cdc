import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import VictoryNFTCollectionItem from "../../contracts/VictoryNFTCollectionItem.cdc"
import VictoryNFTCollectionStorefront from "../../contracts/VictoryNFTCollectionStorefront.cdc"

transaction(bundleID: UInt64) {
    let marketCollection: &VictoryNFTCollectionStorefront.Collection
    let victoryBundler: Capability<&VictoryNFTCollectionItem.Collection{NonFungibleToken.Provider, VictoryNFTCollectionItem.VictoryNFTCollectionItemBundle}>

    prepare(signer: AuthAccount) {
        self.marketCollection = signer.borrow<&VictoryNFTCollectionStorefront.Collection>(from: VictoryNFTCollectionStorefront.CollectionStoragePath)
            ?? panic("Missing or mis-typed VictoryNFTCollectionStorefront Collection")

        // we need a provider capability, but one is not provided by default so we create one.
        let bundlePath = /private/VictoryNFTCollectionItemBundle

        if !signer.getCapability<&VictoryNFTCollectionItem.Collection{NonFungibleToken.Provider, VictoryNFTCollectionItem.VictoryNFTCollectionItemBundle}>(bundlePath)!.check() {
            signer.link<&VictoryNFTCollectionItem.Collection{NonFungibleToken.Provider, VictoryNFTCollectionItem.VictoryNFTCollectionItemBundle}>(bundlePath, target: VictoryNFTCollectionItem.CollectionStoragePath)
        }

        self.victoryBundler = signer.getCapability<&VictoryNFTCollectionItem.Collection{NonFungibleToken.Provider, VictoryNFTCollectionItem.VictoryNFTCollectionItemBundle}>(bundlePath)!
        assert(self.victoryBundler.borrow() != nil, message: "Missing or mis-typed VictoryNFTCollectionItemBundle provider")
    }

    execute {
        let offer <-self.marketCollection.remove(bundleID: bundleID)
        destroy offer

        self.victoryBundler.borrow()!.removeBundle(bundleID: bundleID)
    }
}
