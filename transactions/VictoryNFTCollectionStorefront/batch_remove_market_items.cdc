import VictoryCollectibleSaleOffer from "../../contracts/VictoryCollectibleSaleOffer.cdc"

transaction(start: Int, end: Int) {
    let marketCollection: &VictoryCollectibleSaleOffer.Collection

    prepare(signer: AuthAccount) {
        self.marketCollection = signer.borrow<&VictoryCollectibleSaleOffer.Collection>(from: VictoryCollectibleSaleOffer.CollectionStoragePath)
            ?? panic("Missing or mis-typed VictoryCollectibleSaleOffer Collection")
    }

    execute {
        let ids = self.marketCollection.getSaleOfferIDs()
        var i: Int = start;
        while i <= end {
            let offer <-self.marketCollection.remove(bundleID: ids[i])
            destroy offer
            i = i + 1
        }
    }
}
