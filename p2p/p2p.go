package p2p

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	libp2p "github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	ma "github.com/multiformats/go-multiaddr"
	peerstore "github.com/libp2p/go-libp2p/core/peer"
	host "github.com/libp2p/go-libp2p/core/host"
	crypto "github.com/libp2p/go-libp2p/core/crypto"
)

const (
	BlocksTopic = "graphene-blocks"
	TxTopic     = "graphene-tx"
)

// P2P wraps libp2p host + pubsub
type P2P struct {
	ctx    context.Context
	host   host.Host
	ps     *pubsub.PubSub
	blocks *pubsub.Topic
	tx     *pubsub.Topic
	sub    *pubsub.Subscription
}

// NewP2P creates a new libp2p host and gossip pubsub instance.
// listenAddr is a multiaddr string like "/ip4/0.0.0.0/tcp/0"
func NewP2P(ctx context.Context, listenAddr string) (*P2P, error) {
	// generate ephemeral keypair for this host
	priv, _, err := crypto.GenerateKeyPair(crypto.Ed25519, -1)
	if err != nil {
		return nil, err
	}

	// create host
	h, err := libp2p.New(libp2p.ListenAddrStrings(listenAddr), libp2p.Identity(priv))
	if err != nil {
		return nil, err
	}

	// create pubsub
	ps, err := pubsub.NewGossipSub(ctx, h)
	if err != nil {
		_ = h.Close()
		return nil, err
	}

	blocksTopic, err := ps.Join(BlocksTopic)
	if err != nil {
		_ = h.Close()
		return nil, err
	}
	txTopic, err := ps.Join(TxTopic)
	if err != nil {
		_ = h.Close()
		return nil, err
	}
	sub, err := blocksTopic.Subscribe()
	if err != nil {
		_ = h.Close()
		return nil, err
	}

	p := &P2P{
		ctx:    ctx,
		host:   h,
		ps:     ps,
		blocks: blocksTopic,
		tx:     txTopic,
		sub:    sub,
	}

	log.Printf("libp2p host started: id=%s addrs=%v\n", h.ID().Pretty(), h.Addrs())
	return p, nil
}

// HostID returns peer id string
func (p *P2P) HostID() string {
	if p == nil || p.host == nil {
		return ""
	}
	return p.host.ID().Pretty()
}

// ConnectToPeer dials a peer by multiaddr string like "/ip4/127.0.0.1/tcp/4001/p2p/Qm..."
func (p *P2P) ConnectToPeer(maddr string) error {
	if p == nil || p.host == nil {
		return fmt.Errorf("p2p host not initialized")
	}
	maAddr, err := ma.NewMultiaddr(maddr)
	if err != nil {
		return err
	}
	pi, err := peerstore.AddrInfoFromP2pAddr(maAddr)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(p.ctx, 5*time.Second)
	defer cancel()
	if err := p.host.Connect(ctx, *pi); err != nil {
		return err
	}
	log.Printf("Connected to peer %s\n", pi.ID.Pretty())
	return nil
}

// ConnectToPeers tries to connect to a list of multiaddr strings (best-effort)
func (p *P2P) ConnectToPeers(addrs []string) {
	for _, a := range addrs {
		go func(maS string) {
			if err := p.ConnectToPeer(maS); err != nil {
				log.Printf("bootstrap connect %s failed: %v\n", maS, err)
			}
		}(a)
	}
}

// PublishBlock publishes a JSON-encoded block to the blocks topic.
func (p *P2P) PublishBlock(v interface{}) error {
	if p == nil || p.blocks == nil {
		return fmt.Errorf("blocks topic not ready")
	}
	bz, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return p.blocks.Publish(p.ctx, bz)
}

// SubscribeBlocks starts a goroutine that reads blocks from pubsub and calls the handler for each message.
// Handler should unmarshal the message and process it.
func (p *P2P) SubscribeBlocks(handler func(msg []byte)) {
	go func() {
		for {
			msg, err := p.sub.Next(p.ctx)
			if err != nil {
				// context closed or subscription error
				log.Printf("blocks sub next err: %v\n", err)
				return
			}
			// ignore our own published messages
			if msg.ReceivedFrom == p.host.ID() {
				continue
			}
			// deliver payload to handler
			handler(msg.Data)
		}
	}()
}

// Stop closes host and pubsub
func (p *P2P) Stop() error {
	if p.sub != nil {
		_ = p.sub.Cancel()
	}
	if p.host != nil {
		return p.host.Close()
	}
	return nil
}
