package consensus

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"gfn/state"
	"gfn/store"
	"gfn/p2p"
)

// ------------------- Types -------------------

type Validator struct {
	Address string
	Stake   uint64
	Active  bool
}

type Delegation struct {
	Delegator string
	Validator string
	Amount    uint64
}

type Transaction struct {
	From      string
	To        string
	Amount    uint64
	Type      string // "transfer", "stake", "delegate"
	Validator string // used for "delegate"
}

type Block struct {
	Index     int
	Timestamp int64
	PrevHash  string
	Hash      string
	Validator string
	Txns      []Transaction
}

var (
	Blockchain  []Block
	Validators  []Validator
	Delegations []Delegation
	mu          sync.Mutex
	p2pNet      *p2p.P2P
)

// ------------------- Genesis -------------------

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
	fmt.Println("✅ Genesis block created.")
}

func LoadBlockchain() error {
	blocks, err := store.LoadBlocks()
	if err != nil {
		return err
	}
	Blockchain = blocks
	return nil
}

// ------------------- P2P Integration -------------------

func RegisterP2PNetwork(p *p2p.P2P) {
	p2pNet = p
	p2pNet.SubscribeBlocks(func(msg []byte) {
		handleIncomingBlock(msg)
	})
}

// ------------------- Core Consensus -------------------

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
		fmt.Printf("⛓️  Block %d produced by %s (stake=%d)\n", block.Index, proposer, getValidatorStake(proposer))

		// publish to peers
		if p2pNet != nil {
			if err := p2pNet.PublishBlock(block); err != nil {
				log.Printf("publish error: %v\n", err)
			}
		}
		mu.Unlock()
	}
}

// ------------------- Staking & Delegation -------------------

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

func Delegate(delegator, validator string, amount uint64) error {
	if state.Balances[delegator] < amount {
		return fmt.Errorf("insufficient balance")
	}
	state.Debit(delegator, amount)
	Delegations = append(Delegations, Delegation{Delegator: delegator, Validator: validator, Amount: amount})
	for i := range Validators {
		if Validators[i].Address == validator {
			Validators[i].Stake += amount
			fmt.Printf("🤝 %s delegated %d GFN to %s\n", delegator, amount, validator)
			return nil
		}
	}
	Validators = append(Validators, Validator{Address: validator, Stake: amount, Active: true})
	fmt.Printf("🤝 %s delegated %d GFN to new validator %s\n", delegator, amount, validator)
	return nil
}

// ------------------- Block Logic -------------------

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

// ------------------- Election -------------------

func electValidator() string {
	var total uint64
	for _, v := range Validators {
		if v.Active {
			total += v.Stake
		}
	}
	if total == 0 {
		return ""
	}
	r := rand.Uint64() % total
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

// ------------------- Block Sync -------------------

func handleIncomingBlock(bz []byte) {
	var incoming Block
	if err := json.Unmarshal(bz, &incoming); err != nil {
		log.Printf("❌ Invalid block json: %v", err)
		return
	}
	if blockExists(incoming.Hash) {
		return
	}
	headHash := ""
	if len(Blockchain) > 0 {
		headHash = Blockchain[len(Blockchain)-1].Hash
	}
	if incoming.PrevHash != headHash {
		log.Printf("⚠️  Incoming block prev mismatch: have=%s want=%s", headHash, incoming.PrevHash)
		return
	}
	store.SaveBlock(incoming)
	Blockchain = append(Blockchain, incoming)
	log.Printf("📦 Imported block %d from peer %s", incoming.Index, incoming.Validator)
}

func blockExists(hash string) bool {
	for _, b := range Blockchain {
		if b.Hash == hash {
			return true
		}
	}
	return false
}
