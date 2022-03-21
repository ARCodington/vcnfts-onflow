import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import VictoryCollectible from "../../contracts/VictoryCollectible.cdc"

// This transaction updates a Victory Collection Item's geo location URL.

transaction(updateID: UInt64, newURL: String) {
    prepare(signer: AuthAccount) {
        
        // borrow a reference to the signer's NFT collection
        let collectionRef = signer.borrow<&VictoryCollectible.Collection>(from: VictoryCollectible.CollectionStoragePath)
            ?? panic("Could not borrow a reference to the owner's collection")

         // update the NFT's metadata URL
        collectionRef.geoLocate(id: updateID, locationRef: newURL)
    }
}
