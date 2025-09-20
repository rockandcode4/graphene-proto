package main

import (
    "context"
    "flag"
    "fmt"
    "log"
    "os"
    "time"

    "github.com/rockandcode4/graphene-proto/node"
)

func main() {
    cfgFile := flag.String("config", "", "node config json")
    datadir := flag.String("datadir", "./data", "data directory")
    bind := flag.String("bind", "/ip4/127.0.0.1/tcp/0", "libp2p bind multiaddr")
    port := flag.Int("rpc", 8545, "rpc port")
    flag.Parse()

    cfg := node.DefaultConfig()
    cfg.DataDir = *datadir
    cfg.BindAddr = *bind
    cfg.RPCPort = *port
    if *cfgFile != "" {
        if err := node.LoadConfigFromFile(*cfgFile, cfg); err != nil {
            log.Println("warning: failed to load config:", err)
        }
    }

    ctx := context.Background()
    n, err := node.NewNode(ctx, cfg)
    if err != nil {
        fmt.Println("failed to start node:", err)
        os.Exit(1)
    }
    defer n.Stop()

    log.Printf("Node started. RPC on :%d  PeerID=%s", cfg.RPCPort, n.HostID())

    // simple run loop
    for {
        time.Sleep(10 * time.Second)
    }
}
