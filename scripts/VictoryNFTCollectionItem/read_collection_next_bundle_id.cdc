import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import VictoryCollectible from "../../contracts/VictoryCollectible.cdc"

// This script returns an array of all the NFT IDs in an account's collection.

pub fun main(address: Address): UInt64 {
    let account = getAccount(address)

    let collectionRef = account.getCapability(VictoryCollectible.CollectionPublicPath)!
        .borrow<&{VictoryCollectible.VictoryCollectibleCollectionPublic}>()
        ?? panic("Could not borrow VictoryCollectibleCollectionPublic")
    
    return collectionRef.getNextBundleID()
}
