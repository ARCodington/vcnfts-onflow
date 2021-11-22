import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import VictoryNFTCollectionItem from "../../contracts/VictoryNFTCollectionItem.cdc"

// This script returns the geo URL property of an NFT in an account's collection.

pub fun main(address: Address, itemID: UInt64): String {

    // get the public account object for the token owner
    let owner = getAccount(address)

    let collectionBorrow = owner.getCapability(VictoryNFTCollectionItem.CollectionPublicPath)!
        .borrow<&{VictoryNFTCollectionItem.VictoryNFTCollectionItemCollectionPublic}>()
        ?? panic("Could not borrow VictoryNFTCollectionItemCollectionPublic")

    // borrow a reference to a specific NFT in the collection
    let VictoryItem = collectionBorrow.borrowVictoryItem(id: itemID)
        ?? panic("No such itemID in that collection")

    return VictoryItem.geoURL
}
