import VictoryCollectible from "../../contracts/VictoryCollectible.cdc"

// This scripts returns the number of VictoryCollectible currently in existence.

pub fun main(): UInt64 {    
    return VictoryCollectible.totalSupply
}
