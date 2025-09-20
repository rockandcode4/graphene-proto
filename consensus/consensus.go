package consensus

import (
	"log"
	"sync"
	"time"

	"github.com/yourorg/graphene-proto/p2p"
	"github.com/yourorg/graphene-proto/state"
)

type Block struct {
	Number   uint64
	Prev     []byte
	Time     int64
	Txns     [][]byte
	Proposer string
	Hash     []byte
}

type Consensus struct {
	state *state.StateDB
	p2p   *p2p.P2P

	mu      sync.Mutex
	running bool

	// in-memory chain
	chain []*Block

	validators []string
}

func NewConsensus(st *state.StateDB, p *p2p.P2P) *Consensus {
	genesis := &Block{Number: 0, Prev: nil, Time: time.Now().Unix(), Proposer: "genesis"}
	return &Consensus{
		state:      st,
		p2p:        p,
		chain:      []*Block{genesis},
		validators: []string{},
	}
}

func (c *Consensus) Start() {
	c.mu.Lock()
	if c.running {
		c.mu.Unlock()
		return
	}
	c.running = true
	c.mu.Unlock()
	log.Println("Consensus started")
	go c.loop()
}

func (c *Consensus) Stop() {
	c.mu.Lock()
	c.running = false
	c.mu.Unlock()
}

func (c *Consensus) loop() {
	ticker := time.NewTicker(3 * time.Second)
	for range ticker.C {
		c.mu.Lock()
		if !c.running {
			c.mu.Unlock()
			break
		}
		proposer := "local-proposer"
		if len(c.validators) > 0 {
			proposer = c.validators[len(c.chain)%uint64(len(c.validators))]
		}
		b := &Block{
			Number:   uint64(len(c.chain)),
			Prev:     []byte("prevhash"),
			Time:     time.Now().Unix(),
			Txns:     [][]byte{},
			Proposer: proposer,
		}
		c.chain = append(c.chain, b)
		log.Printf("Proposed block %d by %s", b.Number, proposer)
		_ = c.finalizeBlock(b)
		c.mu.Unlock()
	}
}

func (c *Consensus) finalizeBlock(b *Block) error {
	log.Printf("Finalized block %d", b.Number)
	return nil
}

func (c *Consensus) SubmitTx(from, to string, amount uint64) error {
	return c.state.Transfer(from, to, amount)
}

func (c *Consensus) GetBalance(addr string) (uint64, error) {
	a, err := c.state.GetAccount(addr)
	if err != nil {
		return 0, err
	}
	return a.Balance, nil
}

func (c *Consensus) SetValidators(vals []string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.validators = vals
}
