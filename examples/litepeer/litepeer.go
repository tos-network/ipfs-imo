package main

// This example launches an IPFS-Imo peer and fetches a hello-world
// hash from the IPFS network.

import (
	"context"
	"fmt"
	"io"

	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/multiformats/go-multiaddr"
	ipfsimo "github.com/tos-network/ipfs-imo"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ds := ipfsimo.NewInMemoryDatastore()
	priv, _, err := crypto.GenerateKeyPair(crypto.RSA, 2048)
	if err != nil {
		panic(err)
	}

	listen, _ := multiaddr.NewMultiaddr("/ip4/0.0.0.0/tcp/4005")

	h, dht, err := ipfsimo.SetupLibp2p(
		ctx,
		priv,
		nil,
		[]multiaddr.Multiaddr{listen},
		ds,
		ipfsimo.Libp2pOptionsExtra...,
	)

	if err != nil {
		panic(err)
	}

	lite, err := ipfsimo.New(ctx, ds, nil, h, dht, nil)
	if err != nil {
		panic(err)
	}

	lite.Bootstrap(ipfsimo.DefaultBootstrapPeers())

	c, _ := cid.Decode("QmWATWQ7fVPP2EFGu71UkfnqhYXDYH566qy47CnJDgvs8u")
	rsc, err := lite.GetFile(ctx, c)
	if err != nil {
		panic(err)
	}
	defer rsc.Close()
	content, err := io.ReadAll(rsc)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(content))
}
