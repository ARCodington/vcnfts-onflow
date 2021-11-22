import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import VictoryNFTCollectionItem from "../../contracts/VictoryNFTCollectionItem.cdc"

// This script returns whether a particular item is already for sale.

pub fun main(address: Address, itemID: UInt64): Bool {
    // get the public account object for the token owner
    let owner = getAccount(address)

    let collection = owner.getCapability(VictoryNFTCollectionItem.CollectionPublicPath)!
        .borrow<&{VictoryNFTCollectionItem.VictoryNFTCollectionItemCollectionPublic}>()
        ?? panic("Could not borrow VictoryNFTCollectionItemCollectionPublic")

    return collection.isNFTForSale(id: itemID)
}
