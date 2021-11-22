import FungibleToken from "../../contracts/FungibleToken.cdc"
import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import FUSD from "../../contracts/FUSD.cdc"
import VictoryNFTCollectionItem from "../../contracts/VictoryNFTCollectionItem.cdc"
import VictoryNFTCollectionStorefront from "../../contracts/VictoryNFTCollectionStorefront.cdc"

transaction {
    prepare(signer: AuthAccount) {

        // if the account doesn't already have a vaule
        if signer.borrow<&FUSD.Vault>(from: /storage/fusdVault) == nil {
            // Create a new FUSD Vault and put it in storage
            signer.save(<-FUSD.createEmptyVault(), to: /storage/fusdVault)

            // Create a public capability to the Vault that only exposes
            // the deposit function through the Receiver interface
            signer.link<&FUSD.Vault{FungibleToken.Receiver}>(
            /public/fusdReceiver,
            target: /storage/fusdVault
            )

            // Create a public capability to the Vault that only exposes
            // the balance field through the Balance interface
            signer.link<&FUSD.Vault{FungibleToken.Balance}>(
            /public/fusdBalance,
            target: /storage/fusdVault
            )
        }

        // if the account doesn't already have a collection
        if signer.borrow<&VictoryNFTCollectionItem.Collection>(from: VictoryNFTCollectionItem.CollectionStoragePath) == nil {

            // create a new empty collection
            let collection <- VictoryNFTCollectionItem.createEmptyCollection()
            
            // save it to the account
            signer.save(<-collection, to: VictoryNFTCollectionItem.CollectionStoragePath)

            // create a public capability for the collection
            signer.link<&VictoryNFTCollectionItem.Collection{NonFungibleToken.CollectionPublic, VictoryNFTCollectionItem.VictoryNFTCollectionItemCollectionPublic}>(VictoryNFTCollectionItem.CollectionPublicPath, target: VictoryNFTCollectionItem.CollectionStoragePath)
        }

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
