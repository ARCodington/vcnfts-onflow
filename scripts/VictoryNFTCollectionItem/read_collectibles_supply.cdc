import VictoryNFTCollectionItem from "../../contracts/VictoryNFTCollectionItem.cdc"

// This scripts returns the number of VictoryNFTCollectionItem currently in existence.

pub fun main(): UInt64 {    
    return VictoryNFTCollectionItem.totalSupply
}
