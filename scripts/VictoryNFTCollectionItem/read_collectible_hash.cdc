import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import VictoryNFTCollectionItem from "../../contracts/VictoryNFTCollectionItem.cdc"

// This script returns the content hash of an NFT in an account's collection
// Note: most callers can't handle UInt256, so we return an array of UInt8

pub fun main(address: Address, itemID: UInt64): [UInt8] {

    // get the public account object for the token owner
    let owner = getAccount(address)

    let collectionBorrow = owner.getCapability(VictoryNFTCollectionItem.CollectionPublicPath)!
        .borrow<&{VictoryNFTCollectionItem.VictoryNFTCollectionItemCollectionPublic}>()
        ?? panic("Could not borrow VictoryNFTCollectionItemCollectionPublic")

    // borrow a reference to a specific NFT in the collection
    let victoryItem = collectionBorrow.borrowVictoryItem(id: itemID)
        ?? panic("No such itemID in that collection")

    let hash = victoryItem.contentHash
    var hashBytes: [UInt8] = hash.toBigEndianBytes()
    return hashBytes
}