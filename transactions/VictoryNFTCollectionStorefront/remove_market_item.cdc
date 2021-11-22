import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import VictoryNFTCollectionItem from "../../contracts/VictoryNFTCollectionItem.cdc"
import VictoryNFTCollectionStorefront from "../../contracts/VictoryNFTCollectionStorefront.cdc"

transaction(bundleID: UInt64) {
    let marketCollection: &VictoryNFTCollectionStorefront.Collection
    let victoryCollection: Capability<&VictoryNFTCollectionItem.Collection{NonFungibleToken.Provider, VictoryNFTCollectionItem.VictoryNFTCollectionItemCollectionPublic}>

    prepare(signer: AuthAccount) {
        self.marketCollection = signer.borrow<&VictoryNFTCollectionStorefront.Collection>(from: VictoryNFTCollectionStorefront.CollectionStoragePath)
            ?? panic("Missing or mis-typed VictoryNFTCollectionStorefront Collection")

        // we need a provider capability, but one is not provided by default so we create one.
        let VictoryNFTCollectionItemCollectionProviderPrivatePath = /private/VictoryNFTCollectionItemCollectionProvider

        if !signer.getCapability<&VictoryNFTCollectionItem.Collection{NonFungibleToken.Provider, VictoryNFTCollectionItem.VictoryNFTCollectionItemCollectionPublic}>(VictoryNFTCollectionItemCollectionProviderPrivatePath)!.check() {
            signer.link<&VictoryNFTCollectionItem.Collection{NonFungibleToken.Provider, VictoryNFTCollectionItem.VictoryNFTCollectionItemCollectionPublic}>(VictoryNFTCollectionItemCollectionProviderPrivatePath, target: VictoryNFTCollectionItem.CollectionStoragePath)
        }

        self.victoryCollection = signer.getCapability<&VictoryNFTCollectionItem.Collection{NonFungibleToken.Provider, VictoryNFTCollectionItem.VictoryNFTCollectionItemCollectionPublic}>(VictoryNFTCollectionItemCollectionProviderPrivatePath)!
        assert(self.victoryCollection.borrow() != nil, message: "Missing or mis-typed VictoryNFTCollectionItemCollection provider")
    }

    execute {
        let offer <-self.marketCollection.remove(bundleID: bundleID)
        destroy offer

        self.victoryCollection.borrow()!.removeBundle(bundleID: bundleID)
    }
}
