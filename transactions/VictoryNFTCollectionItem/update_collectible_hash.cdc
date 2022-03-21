import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import VictoryCollectible from "../../contracts/VictoryCollectible.cdc"

// This transaction updates a Victory Collection Item's metadata URL.

transaction(updateID: UInt64, newHash: String) {
    prepare(signer: AuthAccount) {
        
        // borrow a reference to the signer's NFT collection
        let collectionRef = signer.borrow<&VictoryCollectible.Collection>(from: VictoryCollectible.CollectionStoragePath)
            ?? panic("Could not borrow a reference to the owner's collection")

        // convert hash hex string to UInt256
        var hashInt: UInt256 = 0
        var decodeHash: [UInt8] = newHash.decodeHex()
        var i = 0
        while i < decodeHash.length {
            hashInt = hashInt * (256 as UInt256)
            let chunkValue: UInt256 = UInt256(decodeHash[i])
            hashInt = hashInt + chunkValue
            i = i + 1
        }

        // update the NFT's metadata URL
        collectionRef.updateHash(id: updateID, contentHash: hashInt)
    }
}
