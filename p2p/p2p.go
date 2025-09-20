package p2p

import (
    "context"
    "fmt"
    "log"

    libp2p "github.com/libp2p/go-libp2p"
    host "github.com/libp2p/go-libp2p/core/host"
    peer "github.com/libp2p/go-libp2p/core/peer"
    golog "github.com/ipfs/go-log/v2"
)

var logger = golog.Logger("p2p")

type P2P struct {
    ctx  context.Context
    host host.Host
}

func NewP2P(ctx context.Context, bind string) (*P2P, error) {
    h, err := libp2p.New()
    if err != nil { return nil, err }
    p := &P2P{ctx: ctx, host: h}
    return p, nil
}

func (p *P2P) Start() {
    addrs := p.host.Addrs()
    pid := p.host.ID()
    log.Printf("libp2p started: id=%s addrs=%v\n", pid.Pretty(), addrs)
    // In a real network you'd use DHT/Bootstrap peers and pubsub
}

func (p *P2P) Stop() {
    _ = p.host.Close()
}

func (p *P2P) HostID() string {
    return p.host.ID().Pretty()
}

func (p *P2P) Connect(addr string) error {
    pi, err := peer.Decode(addr)
    if err != nil { return err }
    _, err = p.host.Network().DialPeer(p.ctx, pi)
    if err != nil { return fmt.Errorf("dial peer: %w", err) }
    return nil
}
