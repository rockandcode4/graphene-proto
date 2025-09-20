package main

import (
    "fmt"
    "os"

    "github.com/urfave/cli/v2"
    "gfn/consensus"
    "gfn/core"
    "gfn/node"
)

func main() {
    app := &cli.App{
        Name:  "gfn",
        Usage: "Graphene (GFN) blockchain node",
        Commands: []*cli.Command{
            {
                Name:  "init",
                Usage: "Initialize a new node with genesis",
                Action: func(c *cli.Context) error {
                    return node.InitNode()
                },
            },
            {
                Name:  "start",
                Usage: "Start Graphene node",
                Action: func(c *cli.Context) error {
                    return node.StartNode()
                },
            },
        },
    }

    if err := app.Run(os.Args); err != nil {
        fmt.Println("Error:", err)
    }
}
