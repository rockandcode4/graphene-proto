package main

import (
    "crypto/ecdsa"
    "crypto/elliptic"
    "crypto/rand"
    "crypto/x509"
    "encoding/hex"
    "fmt"
)

func main() {
    priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
    if err != nil {
        panic(err)
    }

    privBytes, _ := x509.MarshalECPrivateKey(priv)
    pubBytes, _ := x509.MarshalPKIXPublicKey(&priv.PublicKey)

    fmt.Println("Private Key:", hex.EncodeToString(privBytes))
    fmt.Println("Public Key:", hex.EncodeToString(pubBytes))
}
