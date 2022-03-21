import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import VictoryCollectible from "../../contracts/VictoryCollectible.cdc"

// This script returns whether a particular item is already for sale.

pub fun main(address: Address, itemID: UInt64): Bool {
    // get the public account object for the token owner
    let owner = getAccount(address)

    let collection = owner.getCapability(VictoryCollectible.CollectionPublicPath)!
        .borrow<&{VictoryCollectible.VictoryCollectibleCollectionPublic}>()
        ?? panic("Could not borrow VictoryCollectibleCollectionPublic")

    return collection.isNFTForSale(id: itemID)
}
