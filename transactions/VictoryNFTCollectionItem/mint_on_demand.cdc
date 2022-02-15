import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import VictoryNFTCollectionItem from "../../contracts/VictoryNFTCollectionItem.cdc"

// this script mints a single NFT based on an existing reference NFT (a.k.a. "mint on demand")
// the metadata is consistent across all NFTs minted in this way, so that they appear to be a cohesive set
// even though they are minted individually

transaction(recipient: Address, referenceItemID: UInt64, issueNum: UInt32) {
    
    // local variable for storing the minter reference
    let minter: &VictoryNFTCollectionItem.NFTMinter

    prepare(signer: AuthAccount) {
        // borrow a reference to the NFTMinter resource in storage
        self.minter = signer.borrow<&VictoryNFTCollectionItem.NFTMinter>(from: VictoryNFTCollectionItem.MinterStoragePath)
            ?? panic("Could not borrow a reference to the NFT minter")
    }

    execute {
        // get the public account object for the recipient
        let recipientAcct = getAccount(recipient)

        // borrow the recipient's public NFT collection reference
        let receiver = recipientAcct
            .getCapability(VictoryNFTCollectionItem.CollectionPublicPath)!
            .borrow<&{NonFungibleToken.CollectionPublic}>()
            ?? panic("Could not get receiver reference to the NFT Collection")

        let collection = recipientAcct.getCapability(VictoryNFTCollectionItem.CollectionPublicPath)!
            .borrow<&{VictoryNFTCollectionItem.VictoryNFTCollectionItemCollectionPublic}>()
            ?? panic("Could not borrow VictoryNFTCollectionItemCollectionPublic")

        // borrow a reference to a specific NFT in the collection
        let victoryItem = collection.borrowVictoryItem(id: referenceItemID)
            ?? panic("No such itemID in that collection")

        // mint the NFT and deposit it to the recipient's collection
        // populate most of the values from the reference NFT
        self.minter.mintNFT(recipient: receiver, 
                            owner: recipient, 
                            typeID: victoryItem.typeID, 
                            brandID: victoryItem.brandID, 
                            dropID: victoryItem.dropID,
                            contentHash: victoryItem.contentHash, 
                            startIssueNum: issueNum,
                            maxIssueNum: 1, 
                            totalIssueNum: victoryItem.maxIssueNum,
                            metaURL: victoryItem.metaURL, 
                            geoURL: victoryItem.geoURL)
    }
}