package node

import (
    "context"
    "fmt"
    "log"
    "path/filepath"

    "github.com/rockandcode4/graphene-proto/p2p"
    "github.com/rockandcode4/graphene-proto/state"
    "github.com/rockandcode4/graphene-proto/consensus"
    "github.com/rockandcode4/graphene-proto/rpc"
)

type Node struct {
    ctx context.Context
    cfg *Config
    p2p *p2p.P2P
    state *state.StateDB
    consensus *consensus.Consensus
    rpcServer *rpc.Server
}

func NewNode(ctx context.Context, cfg *Config) (*Node, error) {
    dbPath := filepath.Join(cfg.DataDir, "leveldb")
    st, err := state.NewStateDB(dbPath)
    if err != nil { return nil, err }

    p, err := p2p.NewP2P(ctx, cfg.BindAddr)
    if err != nil { return nil, err }

    cons := consensus.NewConsensus(st, p)
    rpcS, err := rpc.NewServer(cons, cfg.RPCPort)
    if err != nil { return nil, err }

    n := &Node{
        ctx: ctx, cfg: cfg, p2p: p, state: st, consensus: cons, rpcServer: rpcS,
    }
    // start services
    go n.p2p.Start()
    go n.consensus.Start()
    go n.rpcServer.Start()

    return n, nil
}

func (n *Node) HostID() string { return n.p2p.HostID() }

func (n *Node) Stop() {
    n.rpcServer.Stop()
    n.consensus.Stop()
    n.p2p.Stop()
    if err := n.state.Close(); err != nil { log.Println("state close:", err) }
    fmt.Println("node stopped")
}
