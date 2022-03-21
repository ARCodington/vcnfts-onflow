import VictoryCollectibleSaleOffer from "../../contracts/VictoryCollectibleSaleOffer.cdc"

// This transaction configures an account to hold SaleOffer items.

transaction {
    prepare(signer: AuthAccount) {

        // if the account doesn't already have a collection
        if signer.borrow<&VictoryCollectibleSaleOffer.Collection>(from: VictoryCollectibleSaleOffer.CollectionStoragePath) == nil {

            // create a new empty collection
            let collection <- VictoryCollectibleSaleOffer.createEmptyCollection() as! @VictoryCollectibleSaleOffer.Collection
            
            // save it to the account
            signer.save(<-collection, to: VictoryCollectibleSaleOffer.CollectionStoragePath)

            // create a public capability for the collection
            signer.link<&VictoryCollectibleSaleOffer.Collection{VictoryCollectibleSaleOffer.CollectionPublic}>(VictoryCollectibleSaleOffer.CollectionPublicPath, target: VictoryCollectibleSaleOffer.CollectionStoragePath)
        }
    }
}
