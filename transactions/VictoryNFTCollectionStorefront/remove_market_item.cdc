import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import VictoryCollectible from "../../contracts/VictoryCollectible.cdc"
import VictoryCollectibleSaleOffer from "../../contracts/VictoryCollectibleSaleOffer.cdc"

transaction(bundleID: UInt64) {
    let marketCollection: &VictoryCollectibleSaleOffer.Collection
    let victoryBundler: Capability<&VictoryCollectible.Collection{NonFungibleToken.Provider, VictoryCollectible.VictoryCollectibleBundle}>

    prepare(signer: AuthAccount) {
        self.marketCollection = signer.borrow<&VictoryCollectibleSaleOffer.Collection>(from: VictoryCollectibleSaleOffer.CollectionStoragePath)
            ?? panic("Missing or mis-typed VictoryCollectibleSaleOffer Collection")

        // we need a provider capability, but one is not provided by default so we create one.
        let bundlePath = /private/VictoryCollectibleBundle

        if !signer.getCapability<&VictoryCollectible.Collection{NonFungibleToken.Provider, VictoryCollectible.VictoryCollectibleBundle}>(bundlePath)!.check() {
            signer.link<&VictoryCollectible.Collection{NonFungibleToken.Provider, VictoryCollectible.VictoryCollectibleBundle}>(bundlePath, target: VictoryCollectible.CollectionStoragePath)
        }

        self.victoryBundler = signer.getCapability<&VictoryCollectible.Collection{NonFungibleToken.Provider, VictoryCollectible.VictoryCollectibleBundle}>(bundlePath)!
        assert(self.victoryBundler.borrow() != nil, message: "Missing or mis-typed VictoryCollectibleBundle provider")
    }

    execute {
        let offer <-self.marketCollection.remove(bundleID: bundleID)
        destroy offer

        self.victoryBundler.borrow()!.removeBundle(bundleID: bundleID)
    }
}
