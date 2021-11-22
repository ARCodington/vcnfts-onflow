import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import VictoryNFTCollectionItem from "../../contracts/VictoryNFTCollectionItem.cdc"

transaction(itemIDs: [UInt64]) {
    let VictoryNFTCollectionItemProvider: Capability<&VictoryNFTCollectionItem.Collection{NonFungibleToken.Provider, VictoryNFTCollectionItem.VictoryNFTCollectionItemCollectionPublic}>

    prepare(signer: AuthAccount) {
        // we need a provider capability, but one is not provided by default so we create one.
        let VictoryNFTCollectionItemCollectionProviderPrivatePath = /private/VictoryNFTCollectionItemCollectionProvider

        if !signer.getCapability<&VictoryNFTCollectionItem.Collection{NonFungibleToken.Provider, VictoryNFTCollectionItem.VictoryNFTCollectionItemCollectionPublic}>(VictoryNFTCollectionItemCollectionProviderPrivatePath)!.check() {
            signer.link<&VictoryNFTCollectionItem.Collection{NonFungibleToken.Provider, VictoryNFTCollectionItem.VictoryNFTCollectionItemCollectionPublic}>(VictoryNFTCollectionItemCollectionProviderPrivatePath, target: VictoryNFTCollectionItem.CollectionStoragePath)
        }

        self.VictoryNFTCollectionItemProvider = signer.getCapability<&VictoryNFTCollectionItem.Collection{NonFungibleToken.Provider, VictoryNFTCollectionItem.VictoryNFTCollectionItemCollectionPublic}>(VictoryNFTCollectionItemCollectionProviderPrivatePath)!
        assert(self.VictoryNFTCollectionItemProvider.borrow() != nil, message: "Missing or mis-typed VictoryNFTCollectionItemCollection provider")
    }

    execute {
        self.VictoryNFTCollectionItemProvider.borrow()!.createBundle(itemIDs: itemIDs)
    }
}
