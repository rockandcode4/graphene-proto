package main

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatalf("keygen failed: %v", err)
	}
	privBytes := crypto.FromECDSA(privateKey)
	pubBytes := crypto.FromECDSAPub(&privateKey.PublicKey)
	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	fmt.Println("PRIVATE_KEY_HEX:", hex.EncodeToString(privBytes))
	fmt.Println("PUBKEY_HEX:", hex.EncodeToString(pubBytes))
	fmt.Println("ADDRESS:", address.Hex())
}
