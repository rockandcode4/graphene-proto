package consensus

import (
    "crypto/sha256"
    "encoding/hex"
    "time"
)

type Block struct {
    Height    int64
    Timestamp int64
    PrevHash  string
    Data      []byte
    Validator string
    Hash      string
}

func NewBlock(height int64, prevHash, validator string, data []byte) *Block {
    b := &Block{
        Height:    height,
        Timestamp: time.Now().Unix(),
        PrevHash:  prevHash,
        Data:      data,
        Validator: validator,
    }
    b.Hash = b.CalculateHash()
    return b
}

func (b *Block) CalculateHash() string {
    record := string(b.Height) + b.PrevHash + b.Validator + string(b.Data) + string(b.Timestamp)
    hash := sha256.Sum256([]byte(record))
    return hex.EncodeToString(hash[:])
}
