import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import VictoryCollectible from "../../contracts/VictoryCollectible.cdc"

// This transaction configures an account to hold Victory Collection Items.

transaction {
    prepare(signer: AuthAccount) {
        // if the account doesn't already have a collection
        if signer.borrow<&VictoryCollectible.Collection>(from: VictoryCollectible.CollectionStoragePath) == nil {

            // create a new empty collection
            let collection <- VictoryCollectible.createEmptyCollection()
            
            // save it to the account
            signer.save(<-collection, to: VictoryCollectible.CollectionStoragePath)

            // create a public capability for the collection
            signer.link<&VictoryCollectible.Collection{NonFungibleToken.CollectionPublic, VictoryCollectible.VictoryCollectibleCollectionPublic}>(VictoryCollectible.CollectionPublicPath, target: VictoryCollectible.CollectionStoragePath)
        }
    }
}
