import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import VictoryNFTCollectionItem from "../../contracts/VictoryNFTCollectionItem.cdc"

// This transaction transfers a Victory Collection Item from one account to another.

transaction(recipient: Address, withdrawID: UInt64) {
    prepare(signer: AuthAccount) {
        
        // get the recipients public account object
        let recipient = getAccount(recipient)

        // borrow a reference to the signer's NFT collection
        let collectionRef = signer.borrow<&VictoryNFTCollectionItem.Collection>(from: VictoryNFTCollectionItem.CollectionStoragePath)
            ?? panic("Could not borrow a reference to the owner's collection")

        if collectionRef.isNFTForSale(id: withdrawID) {
            panic("Can't transfer an item that is already for sale")
        }

        // borrow a public reference to the receivers collection
        let depositRef = recipient.getCapability(VictoryNFTCollectionItem.CollectionPublicPath)!.borrow<&{NonFungibleToken.CollectionPublic}>()!

        // withdraw the NFT from the owner's collection
        let nft <- collectionRef.withdraw(withdrawID: withdrawID)

        // Deposit the NFT in the recipient's collection
        depositRef.deposit(token: <-nft)
    }
}
