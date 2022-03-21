import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import VictoryCollectible from "../../contracts/VictoryCollectible.cdc"

// This script returns an array of all the NFT IDs in a particular bundle in an account's collection.

pub fun main(address: Address, bundleID: UInt64): [UInt64] {
    // get the public account object for the token owner
    let owner = getAccount(address)

    let collectionBorrow = owner.getCapability(VictoryCollectible.CollectionPublicPath)!
        .borrow<&{VictoryCollectible.VictoryCollectibleCollectionPublic}>()
        ?? panic("Could not borrow VictoryCollectibleCollectionPublic")

    return collectionBorrow.getBundleIDs(bundleID: bundleID)
}
