import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import VictoryNFTCollectionItem from "../../contracts/VictoryNFTCollectionItem.cdc"

// This script returns all the public properties of an NFT in an account's collection.

pub struct ItemInfo {
    pub var originalOwner: Address
    pub var id: UInt64
    pub var typeID: UInt64
    pub var brandID: UInt64
    pub var dropID: UInt64
    pub var issueNum: UInt32
    pub var maxIssueNum: UInt32
    pub var contentHash: String
    pub var metaURL: String
    pub var geoURL: String

    init(
        originalOwner: Address,
        id: UInt64,
        typeID: UInt64,
        brandID: UInt64,
        dropID: UInt64,
        issueNum: UInt32,
        maxIssueNum: UInt32,
        contentHash: String,
        metaURL: String,
        geoURL: String,
        ) 
    {
            self.originalOwner = originalOwner
            self.id = id
            self.typeID = typeID
            self.brandID = brandID
            self.dropID = dropID
            self.issueNum = issueNum
            self.maxIssueNum = maxIssueNum
            self.contentHash = contentHash
            self.metaURL = metaURL
            self.geoURL = geoURL
    }
}

pub fun main(address: Address, itemID: UInt64): ItemInfo {

    // get the public account object for the token owner
    let owner = getAccount(address)

    let collectionBorrow = owner.getCapability(VictoryNFTCollectionItem.CollectionPublicPath)!
        .borrow<&{VictoryNFTCollectionItem.VictoryNFTCollectionItemCollectionPublic}>()
        ?? panic("Could not borrow VictoryNFTCollectionItemCollectionPublic")

    // borrow a reference to a specific NFT in the collection
    let victoryItem = collectionBorrow.borrowVictoryItem(id: itemID)
        ?? panic("No such itemID in that collection")

    let hash = victoryItem.contentHash
    let hashBytes = hash.toBigEndianBytes()
    let hashStr = String.encodeHex(hashBytes)

    return ItemInfo(
            originalOwner: victoryItem!.originalOwner,
            id: victoryItem!.id,
            typeID: victoryItem!.typeID,
            brandID: victoryItem!.brandID,
            dropID: victoryItem!.dropID,
            issueNum: victoryItem!.issueNum,
            maxIssueNum: victoryItem!.maxIssueNum,
            contentHash: hashStr,
            metaURL: victoryItem!.metaURL,
            geoURL: victoryItem!.geoURL
        )
}
