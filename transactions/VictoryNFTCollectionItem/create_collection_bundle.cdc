import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import VictoryNFTCollectionItem from "../../contracts/VictoryNFTCollectionItem.cdc"

transaction(itemIDs: [UInt64]) {
    let bundleProvider: Capability<&VictoryNFTCollectionItem.Collection{NonFungibleToken.Provider, VictoryNFTCollectionItem.VictoryNFTCollectionItemBundle}>

    prepare(signer: AuthAccount) {
        // we need a provider capability, but one is not provided by default so we create one.
        let bundlePath = /private/VictoryNFTCollectionItemBundle

        if !signer.getCapability<&VictoryNFTCollectionItem.Collection{NonFungibleToken.Provider, VictoryNFTCollectionItem.VictoryNFTCollectionItemBundle}>(bundlePath)!.check() {
            signer.link<&VictoryNFTCollectionItem.Collection{NonFungibleToken.Provider, VictoryNFTCollectionItem.VictoryNFTCollectionItemBundle}>(bundlePath, target: VictoryNFTCollectionItem.CollectionStoragePath)
        }

        self.bundleProvider = signer.getCapability<&VictoryNFTCollectionItem.Collection{NonFungibleToken.Provider, VictoryNFTCollectionItem.VictoryNFTCollectionItemBundle}>(bundlePath)!
        assert(self.bundleProvider.borrow() != nil, message: "Missing or mis-typed VictoryNFTCollectionItemCollection provider")
    }

    execute {
        self.bundleProvider.borrow()!.createBundle(itemIDs: itemIDs)
    }
}
