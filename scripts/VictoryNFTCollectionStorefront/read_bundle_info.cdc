import VictoryCollectibleSaleOffer from "../../contracts/VictoryCollectibleSaleOffer.cdc"
import VictoryCollectible from "../../contracts/VictoryCollectible.cdc"

// This script returns all the public properties of a sale offer.

pub struct BundleInfo {
    pub var saleCompleted: Bool
    pub var bundleID: UInt64
    pub var itemIDs: [UInt64]
    pub var price: UFix64
    pub var saleType: UInt8
    pub var startTime: UInt32
    pub var endTime: UInt32
    pub var targetPrice: UFix64
    pub let royaltyFactor: UFix64
    pub let originalOwner: Address
    pub var seller: Address
    pub var winner: Address

    init(
        saleCompleted: Bool,
        bundleID: UInt64,
        itemIDs: [UInt64],
        price: UFix64,
        saleType: UInt8,
        startTime: UInt32,
        endTime: UInt32,
        targetPrice: UFix64,
        royaltyFactor: UFix64,
        originalOwner: Address,
        seller: Address,
        winner: Address
        ) 
    {
            self.saleCompleted = saleCompleted
            self.bundleID = bundleID
            self.itemIDs = itemIDs
            self.price = price
            self.saleType = saleType
            self.startTime = startTime
            self.endTime = endTime
            self.targetPrice = targetPrice
            self.royaltyFactor = royaltyFactor
            self.originalOwner = originalOwner
            self.seller = seller
            self.winner = winner
    }
}

pub fun main(address: Address, bundleID: UInt64): BundleInfo {
    let owner = getAccount(address)

    let marketCollectionRef = owner
            .getCapability<&VictoryCollectibleSaleOffer.Collection{VictoryCollectibleSaleOffer.CollectionPublic}>(VictoryCollectibleSaleOffer.CollectionPublicPath)
            .borrow()
        ?? panic("Could not borrow market collection from market address")
    
    let itemCollectionRef = owner.getCapability(VictoryCollectible.CollectionPublicPath)!
        .borrow<&{VictoryCollectible.VictoryCollectibleCollectionPublic}>()
        ?? panic("Could not borrow VictoryCollectibleCollectionPublic")

    let saleOffer = marketCollectionRef.borrowSaleItem(bundleID: bundleID)
                    ?? panic("Could not borrow sale item from market collection")

    return BundleInfo(
            saleCompleted: saleOffer!.saleCompleted,
            bundleID: saleOffer!.bundleID,
            itemIDs: itemCollectionRef.getBundleIDs(bundleID: bundleID),
            price: saleOffer!.price,
            saleType: saleOffer!.saleType,
            startTime: saleOffer!.startTime,
            endTime: saleOffer!.endTime,
            targetPrice: saleOffer!.targetPrice,
            royaltyFactor: saleOffer!.royaltyFactor,
            originalOwner: saleOffer!.originalOwner,
            seller: saleOffer!.seller,
            winner: saleOffer!.winner
    )
}
