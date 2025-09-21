package main

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"

    "gfn/node"
    "gfn/store"
)

func main() {
    fmt.Println("Starting GFN Blockchain...")

    // Init node
    if err := node.InitNode(); err != nil {
        fmt.Println("Init error:", err)
        return
    }

    // Handle shutdown signals (CTRL+C, kill, etc.)
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        <-quit
        fmt.Println("\nShutting down node...")
        store.CloseDB()
        os.Exit(0)
    }()

    // Start node consensus
    if err := node.StartNode(); err != nil {
        fmt.Println("Node error:", err)
    }
}
