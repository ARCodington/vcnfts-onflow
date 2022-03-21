import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import VictoryCollectible from "../../contracts/VictoryCollectible.cdc"

transaction(itemIDs: [UInt64]) {
    let bundleProvider: Capability<&VictoryCollectible.Collection{NonFungibleToken.Provider, VictoryCollectible.VictoryCollectibleBundle}>

    prepare(signer: AuthAccount) {
        // we need a provider capability, but one is not provided by default so we create one.
        let bundlePath = /private/VictoryCollectibleBundle

        if !signer.getCapability<&VictoryCollectible.Collection{NonFungibleToken.Provider, VictoryCollectible.VictoryCollectibleBundle}>(bundlePath)!.check() {
            signer.link<&VictoryCollectible.Collection{NonFungibleToken.Provider, VictoryCollectible.VictoryCollectibleBundle}>(bundlePath, target: VictoryCollectible.CollectionStoragePath)
        }

        self.bundleProvider = signer.getCapability<&VictoryCollectible.Collection{NonFungibleToken.Provider, VictoryCollectible.VictoryCollectibleBundle}>(bundlePath)!
        assert(self.bundleProvider.borrow() != nil, message: "Missing or mis-typed VictoryCollectibleCollection provider")
    }

    execute {
        self.bundleProvider.borrow()!.createBundle(itemIDs: itemIDs)
    }
}
