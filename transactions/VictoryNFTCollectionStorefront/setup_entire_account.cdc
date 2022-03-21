import FungibleToken from "../../contracts/FungibleToken.cdc"
import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import FUSD from "../../contracts/FUSD.cdc"
import VictoryCollectible from "../../contracts/VictoryCollectible.cdc"
import VictoryCollectibleSaleOffer from "../../contracts/VictoryCollectibleSaleOffer.cdc"

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
        if signer.borrow<&VictoryCollectible.Collection>(from: VictoryCollectible.CollectionStoragePath) == nil {

            // create a new empty collection
            let collection <- VictoryCollectible.createEmptyCollection()
            
            // save it to the account
            signer.save(<-collection, to: VictoryCollectible.CollectionStoragePath)

            // create a public capability for the collection
            signer.link<&VictoryCollectible.Collection{NonFungibleToken.CollectionPublic, VictoryCollectible.VictoryCollectibleCollectionPublic}>(VictoryCollectible.CollectionPublicPath, target: VictoryCollectible.CollectionStoragePath)
        }

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
