import VictoryNFTCollectionStorefront from "../../contracts/VictoryNFTCollectionStorefront.cdc"

// This transaction configures an account to hold SaleOffer items.

transaction {
    prepare(signer: AuthAccount) {

        // if the account doesn't already have a collection
        if signer.borrow<&VictoryNFTCollectionStorefront.Collection>(from: VictoryNFTCollectionStorefront.CollectionStoragePath) == nil {

            // create a new empty collection
            let collection <- VictoryNFTCollectionStorefront.createEmptyCollection() as! @VictoryNFTCollectionStorefront.Collection
            
            // save it to the account
            signer.save(<-collection, to: VictoryNFTCollectionStorefront.CollectionStoragePath)

            // create a public capability for the collection
            signer.link<&VictoryNFTCollectionStorefront.Collection{VictoryNFTCollectionStorefront.CollectionPublic}>(VictoryNFTCollectionStorefront.CollectionPublicPath, target: VictoryNFTCollectionStorefront.CollectionStoragePath)
        }
    }
}
