package node

import (
    "encoding/json"
    "os"
)

type Config struct {
    DataDir    string `json:"data_dir"`
    BindAddr   string `json:"bind_addr"`
    Bootstrap  []string `json:"bootstrap"`
    RPCPort    int    `json:"rpc_port"`
    Genesis    string `json:"genesis_json"`
    NodeKeyHex string `json:"node_key_hex"`
}

func DefaultConfig() *Config {
    return &Config{
        DataDir:  "./data",
        BindAddr: "/ip4/127.0.0.1/tcp/0",
        RPCPort:  8545,
    }
}

func LoadConfigFromFile(path string, cfg *Config) error {
    b, err := os.ReadFile(path)
    if err != nil { return err }
    return json.Unmarshal(b, cfg)
}
