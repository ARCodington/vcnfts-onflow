import VictoryNFTCollectionStorefront from "../../contracts/VictoryNFTCollectionStorefront.cdc"

transaction(start: Int, end: Int) {
    let marketCollection: &VictoryNFTCollectionStorefront.Collection

    prepare(signer: AuthAccount) {
        self.marketCollection = signer.borrow<&VictoryNFTCollectionStorefront.Collection>(from: VictoryNFTCollectionStorefront.CollectionStoragePath)
            ?? panic("Missing or mis-typed VictoryNFTCollectionStorefront Collection")
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
