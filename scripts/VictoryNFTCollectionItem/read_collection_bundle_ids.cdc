import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import VictoryNFTCollectionItem from "../../contracts/VictoryNFTCollectionItem.cdc"

// This script returns an array of all the NFT IDs in a particular bundle in an account's collection.

pub fun main(address: Address, bundleID: UInt64): [UInt64] {
    // get the public account object for the token owner
    let owner = getAccount(address)

    let collectionBorrow = owner.getCapability(VictoryNFTCollectionItem.CollectionPublicPath)!
        .borrow<&{VictoryNFTCollectionItem.VictoryNFTCollectionItemCollectionPublic}>()
        ?? panic("Could not borrow VictoryNFTCollectionItemCollectionPublic")

    return collectionBorrow.getBundleIDs(bundleID: bundleID)
}
