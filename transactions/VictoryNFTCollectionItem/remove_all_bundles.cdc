import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import VictoryNFTCollectionItem from "../../contracts/VictoryNFTCollectionItem.cdc"

transaction() {
    let victoryBundler: Capability<&VictoryNFTCollectionItem.Collection{NonFungibleToken.Provider, VictoryNFTCollectionItem.VictoryNFTCollectionItemBundle}>

    prepare(signer: AuthAccount) {
        // we need a provider capability, but one is not provided by default so we create one.
        let bundlePath = /private/VictoryNFTCollectionItemBundle

        if !signer.getCapability<&VictoryNFTCollectionItem.Collection{NonFungibleToken.Provider, VictoryNFTCollectionItem.VictoryNFTCollectionItemBundle}>(bundlePath)!.check() {
            signer.link<&VictoryNFTCollectionItem.Collection{NonFungibleToken.Provider, VictoryNFTCollectionItem.VictoryNFTCollectionItemBundle}>(bundlePath, target: VictoryNFTCollectionItem.CollectionStoragePath)
        }

        self.victoryBundler = signer.getCapability<&VictoryNFTCollectionItem.Collection{NonFungibleToken.Provider, VictoryNFTCollectionItem.VictoryNFTCollectionItemBundle}>(bundlePath)!
        assert(self.victoryBundler.borrow() != nil, message: "Missing or mis-typed VictoryNFTCollectionItemBundle provider")
    }

    execute {
        self.victoryBundler.borrow()!.removeAllBundles()
    }
}
