import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import VictoryNFTCollectionItem from "../../contracts/VictoryNFTCollectionItem.cdc"

// This transaction configures an account to hold Victory Collection Items.

transaction {
    prepare(signer: AuthAccount) {
        // if the account doesn't already have a collection
        if signer.borrow<&VictoryNFTCollectionItem.Collection>(from: VictoryNFTCollectionItem.CollectionStoragePath) == nil {

            // create a new empty collection
            let collection <- VictoryNFTCollectionItem.createEmptyCollection()
            
            // save it to the account
            signer.save(<-collection, to: VictoryNFTCollectionItem.CollectionStoragePath)

            // create a public capability for the collection
            signer.link<&VictoryNFTCollectionItem.Collection{NonFungibleToken.CollectionPublic, VictoryNFTCollectionItem.VictoryNFTCollectionItemCollectionPublic}>(VictoryNFTCollectionItem.CollectionPublicPath, target: VictoryNFTCollectionItem.CollectionStoragePath)
        }
    }
}
