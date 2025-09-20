package main


import (
"crypto/ecdsa"
"crypto/elliptic"
"crypto/rand"
"encoding/hex"
"fmt"
"log"
)


func main() {
priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
if err != nil { log.Fatal(err) }
privBytes := priv.D.Bytes()
pubX := priv.PublicKey.X.Bytes()
pubY := priv.PublicKey.Y.Bytes()
fmt.Println("PRIVATE:", hex.EncodeToString(privBytes))
fmt.Println("PUB_X:", hex.EncodeToString(pubX))
fmt.Println("PUB_Y:", hex.EncodeToString(pubY))
}
