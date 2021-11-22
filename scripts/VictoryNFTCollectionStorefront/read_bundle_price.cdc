import VictoryNFTCollectionStorefront from "../../contracts/VictoryNFTCollectionStorefront.cdc"

// This script returns an array of all the NFT IDs for sale 
// in an account's SaleOffer collection.

pub fun main(address: Address, bundleID: UInt64): UFix64 {
    let marketCollectionRef = getAccount(address)
        .getCapability<&VictoryNFTCollectionStorefront.Collection{VictoryNFTCollectionStorefront.CollectionPublic}>(
            VictoryNFTCollectionStorefront.CollectionPublicPath
        )
        .borrow()
        ?? panic("Could not borrow market collection from market address")
    
    let saleOffer = marketCollectionRef.borrowSaleItem(bundleID: bundleID)
                    ?? panic("Could not borrow sale item from market collection")
    return saleOffer!.price
}
