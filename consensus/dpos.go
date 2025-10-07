package consensus

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "math/rand"
    "sync"
    "time"

    "gfn/state"
    "gfn/store"
)

// ------------------- Types -------------------

// Validator represents a staking node
type Validator struct {
    Address string
    Stake   uint64
    Active  bool
}

// Delegation from a user to a validator
type Delegation struct {
    Delegator string
    Validator string
    Amount    uint64
}

// Block structure
type Block struct {
    Index     int
    Timestamp int64
    PrevHash  string
    Hash      string
    Validator string
    Txns      []Transaction
}

// Transaction for transfer, stake, or delegation
type Transaction struct {
    From      string
    To        string
    Amount    uint64
    Type      string // "transfer", "stake", "delegate"
    Validator string // used for "delegate"
}

var (
    Blockchain []Block
    Validators []Validator
    Delegations []Delegation
    mu          sync.Mutex
)

// ------------------- Core -------------------

// Initialize genesis block
func InitGenesis() {
    genesis := Block{
        Index:     0,
        Timestamp: time.Now().Unix(),
        PrevHash:  "",
        Validator: "genesis",
        Txns:      []Transaction{},
    }
    genesis.Hash = calculateHash(genesis)
    Blockchain = append(Blockchain, genesis)
    store.SaveBlock(genesis)
    fmt.Println("Genesis block created.")
}

// Load blockchain from store
func LoadBlockchain() error {
    blocks, err := store.LoadBlocks()
    if err != nil {
        return err
    }
    Blockchain = blocks
    return nil
}

// RunConsensus — simple DPoS loop
func RunConsensus() {
    for {
        time.Sleep(5 * time.Second)
        mu.Lock()
        proposer := electValidator()
        if proposer == "" {
            mu.Unlock()
            continue
        }
        block := generateBlock(proposer)
        Blockchain = append(Blockchain, block)
        store.SaveBlock(block)
        fmt.Printf("⛓️ Block %d produced by %s (stake=%d)\n", block.Index, proposer, getValidatorStake(proposer))
        mu.Unlock()
    }
}

// ------------------- Staking / Delegation -------------------

// Handle staking transaction
func Stake(address string, amount uint64) error {
    if state.Balances[address] < amount {
        return fmt.Errorf("insufficient balance")
    }
    state.Debit(address, amount)
    found := false
    for i := range Validators {
        if Validators[i].Address == address {
            Validators[i].Stake += amount
            found = true
            break
        }
    }
    if !found {
        Validators = append(Validators, Validator{Address: address, Stake: amount, Active: true})
    }
    fmt.Printf("✅ %s staked %d GFN\n", address, amount)
    return nil
}

// Handle delegation transaction
func Delegate(delegator, validator string, amount uint64) error {
    if state.Balances[delegator] < amount {
        return fmt.Errorf("insufficient balance")
    }
    state.Debit(delegator, amount)
    Delegations = append(Delegations, Delegation{Delegator: delegator, Validator: validator, Amount: amount})

    // Increase validator's effective stake
    for i := range Validators {
        if Validators[i].Address == validator {
            Validators[i].Stake += amount
            fmt.Printf("🤝 %s delegated %d GFN to %s\n", delegator, amount, validator)
            return nil
        }
    }

    // If validator not found, create it
    Validators = append(Validators, Validator{Address: validator, Stake: amount, Active: true})
    fmt.Printf("🤝 %s delegated %d GFN to new validator %s\n", delegator, amount, validator)
    return nil
}

// Elect validator weighted by stake
func electValidator() string {
    var totalStake uint64
    for _, v := range Validators {
        if v.Active {
            totalStake += v.Stake
        }
    }
    if totalStake == 0 {
        return ""
    }
    r := rand.Uint64() % totalStake
    var cumulative uint64
    for _, v := range Validators {
        if v.Active {
            cumulative += v.Stake
            if r < cumulative {
                return v.Address
            }
        }
    }
    return ""
}

// ------------------- Block Generation -------------------

func generateBlock(validator string) Block {
    prev := Blockchain[len(Blockchain)-1]
    block := Block{
        Index:     prev.Index + 1,
        Timestamp: time.Now().Unix(),
        PrevHash:  prev.Hash,
        Validator: validator,
        Txns:      []Transaction{},
    }
    block.Hash = calculateHash(block)
    return block
}

func calculateHash(b Block) string {
    record := fmt.Sprintf("%d%d%s%s", b.Index, b.Timestamp, b.PrevHash, b.Validator)
    hash := sha256.Sum256([]byte(record))
    return hex.EncodeToString(hash[:])
}

func getValidatorStake(addr string) uint64 {
    for _, v := range Validators {
        if v.Address == addr {
            return v.Stake
        }
    }
    return 0
}
