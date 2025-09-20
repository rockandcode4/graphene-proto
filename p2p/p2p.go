package p2p

import (
	"context"
	"log"

	libp2p "github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	host "github.com/libp2p/go-libp2p/core/host"
	golog "github.com/ipfs/go-log/v2"
)

var logger = golog.Logger("p2p")

type P2P struct {
	ctx  context.Context
	host host.Host
	ps   *pubsub.PubSub
}

func NewP2P(ctx context.Context, bind string) (*P2P, error) {
	h, err := libp2p.New()
	if err != nil {
		return nil, err
	}
	ps, err := pubsub.NewGossipSub(ctx, h)
	if err != nil {
		return nil, err
	}
	p := &P2P{ctx: ctx, host: h, ps: ps}
	return p, nil
}

func (p *P2P) Start() {
	addrs := p.host.Addrs()
	pid := p.host.ID()
	log.Printf("libp2p started: id=%s addrs=%v\n", pid.Pretty(), addrs)
	// create default topics for blocks/tx
	_, err := p.ps.Join("graphene-blocks")
	if err != nil {
		logger.Error(err)
	}
	_, err = p.ps.Join("graphene-tx")
	if err != nil {
		logger.Error(err)
	}
}

func (p *P2P) Stop() {
	_ = p.host.Close()
}

func (p *P2P) HostID() string {
	return p.host.ID().Pretty()
}
