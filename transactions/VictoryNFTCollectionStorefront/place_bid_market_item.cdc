import VictoryNFTCollectionStorefront from "../../contracts/VictoryNFTCollectionStorefront.cdc"

transaction(bundleID: UInt64, bidPrice: UFix64, bidder: Address, marketCollectionAddress: Address) {
    let marketCollection: &VictoryNFTCollectionStorefront.Collection{VictoryNFTCollectionStorefront.CollectionPublic}

    prepare(signer: AuthAccount) {
        self.marketCollection = getAccount(marketCollectionAddress)
            .getCapability<&VictoryNFTCollectionStorefront.Collection{VictoryNFTCollectionStorefront.CollectionPublic}>(
                VictoryNFTCollectionStorefront.CollectionPublicPath
            )!
            .borrow()
            ?? panic("Could not borrow market collection from market address")
    }

    execute {
        self.marketCollection.placeBid(
            bundleID: bundleID,
            bidPrice: bidPrice,
            bidder: bidder
        )
    }
}
