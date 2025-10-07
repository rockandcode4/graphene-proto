package node

import (
    "fmt"

    "gfn/consensus"
    "gfn/state"
    "gfn/store"
)

// InitNode sets up the blockchain node
func InitNode() error {
    fmt.Println("Initializing node...")

    // Open LevelDB
    if err := store.OpenDB("data"); err != nil {
        return fmt.Errorf("failed to open database: %v", err)
    }

    // Initialize account state
    state.InitState()

    // Give initial balances (for testing)
    state.Credit("validator1", 1000)
    state.Credit("validator2", 500)
    state.Credit("user1", 0)

    // Try loading existing blockchain
    err := consensus.LoadBlockchain()
    if err != nil {
        fmt.Println("No existing blockchain found, creating genesis...")
        consensus.InitGenesis()
    }

    // Initialize validators
    consensus.Validators = []consensus.Validator{
        {Address: "validator1", Stake: 1000, Active: true},
        {Address: "validator2", Stake: 800, Active: true},
    }

    fmt.Println("Validators initialized.")
    return nil
}

// StartNode begins the consensus loop
func StartNode() error {
    fmt.Println("Starting node...")
    go consensus.RunConsensus()
    select {} // block forever
}
