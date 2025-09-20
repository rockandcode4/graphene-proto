package consensus

import (
    "context"
    "fmt"
    "log"
    "sync"
    "time"

    "github.com/rockandcode4/graphene-proto/state"
    "github.com/rockandcode4/graphene-proto/p2p"
)

// Simple block + vote model for prototype.
// Not production-ready.

type Block struct {
    Number uint64
    Prev   []byte
    Time   int64
    Txns   [][]byte
    Proposer string
    Hash   []byte
}

type Vote struct {
    BlockHash []byte
    Voter     string
    Signature []byte // prototype: not real sig
}

type Consensus struct {
    state *state.StateDB
    p2p   *p2p.P2P

    mu sync.Mutex
    running bool

    // in-memory chain
    chain []*Block

    validators []string
}

func NewConsensus(st *state.StateDB, p *p2p.P2P) *Consensus {
    // load genesis chain from state or create one
    genesis := &Block{Number:0, Prev:nil, Time: time.Now().Unix(), Proposer:"genesis"}
    return &Consensus{
        state: st, p2p: p, chain: []*Block{genesis}, validators: []string{},
    }
}

func (c *Consensus) Start() {
    c.mu.Lock()
    if c.running { c.mu.Unlock(); return }
    c.running = true
    c.mu.Unlock()
    log.Println("Consensus started")
    go c.loop()
}

func (c *Consensus) Stop() {
    c.mu.Lock(); c.running = false; c.mu.Unlock()
}

func (c *Consensus) loop() {
    ticker := time.NewTicker(3 * time.Second)
    for range ticker.C {
        c.mu.Lock()
        if !c.running { c.mu.Unlock(); break }
        // for prototype: proposer is deterministic round-robin among validators
        proposer := "local-proposer"
        // create empty block
        b := &Block{
            Number: uint64(len(c.chain)),
            Prev:   []byte("prevhash"),
            Time:   time.Now().Unix(),
            Txns:   [][]byte{},
            Proposer: proposer,
        }
        c.chain = append(c.chain, b)
        log.Printf("Proposed block %d by %s", b.Number, proposer)
        // finalize instantly for prototype:
        _ = c.finalizeBlock(b)
        c.mu.Unlock()
    }
}

func (c *Consensus) finalizeBlock(b *Block) error {
    // in Casper you'd gather attester votes and then finalize when >2/3.
    // Prototype: accept block immediately and apply state changes if any.
    log.Printf("Finalized block %d", b.Number)
    // persist or notify p2p
    return nil
}

func (c *Consensus) SubmitTx(from, to string, amount uint64) error {
    // simple immediate state application; in real chain txs go to mempool and included in blocks.
    return c.state.Transfer(from, to, amount)
}

func (c *Consensus) GetBalance(addr string) (uint64, error) {
    a, err := c.state.GetAccount(addr)
    if err != nil { return 0, err }
    return a.Balance, nil
}
