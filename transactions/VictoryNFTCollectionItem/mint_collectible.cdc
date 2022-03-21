import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import VictoryCollectible from "../../contracts/VictoryCollectible.cdc"

// This transction uses the NFTMinter resource to mint a new NFT.
//
// It must be run with the account that has the minter resource
// stored at path /storage/NFTMinter.

transaction(recipient: Address, typeID: UInt64, brandID: UInt64, seriesID: UInt64, dropID: UInt64, contentHash: String, startIssueNum: UInt32, maxIssueNum: UInt32, totalIssueNum: UInt32) {
    
    // local variable for storing the minter reference
    let minter: &VictoryCollectible.NFTMinter

    prepare(signer: AuthAccount) {
        // borrow a reference to the NFTMinter resource in storage
        self.minter = signer.borrow<&VictoryCollectible.NFTMinter>(from: VictoryCollectible.MinterStoragePath)
            ?? panic("Could not borrow a reference to the NFT minter")
    }

    execute {
        // get the public account object for the recipient
        let recipientAcct = getAccount(recipient)

        // borrow the recipient's public NFT collection reference
        let receiver = recipientAcct
            .getCapability(VictoryCollectible.CollectionPublicPath)!
            .borrow<&{NonFungibleToken.CollectionPublic}>()
            ?? panic("Could not get receiver reference to the NFT Collection")

        // convert hash hex string to UInt256
        var hashInt: UInt256 = 0
        var decodeHash: [UInt8] = contentHash.decodeHex()
        var i = 0
        while i < decodeHash.length {
            hashInt = hashInt * (256 as UInt256)
            let chunkValue: UInt256 = UInt256(decodeHash[i])
            hashInt = hashInt + chunkValue
            i = i + 1
        }

        // mint the NFT and deposit it to the recipient's collection
        self.minter.mintNFT(recipient: receiver, owner: recipient, 
                            typeID: typeID, 
                            brandID: brandID, seriesID: seriesID, dropID: dropID,
                            contentHash:hashInt, 
                            startIssueNum: startIssueNum,
                            maxIssueNum: maxIssueNum, 
                            totalIssueNum: totalIssueNum)
    }
}
