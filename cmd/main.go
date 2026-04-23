package main

import (
	"log"
	"os"

	"github.com/rohit20001221/ripple/network"
	"github.com/rohit20001221/ripple/torrent"
)

func main() {
	if len(os.Args) < 1 {
		log.Fatalln("error: provide the path to torrent")
	}

	filePath := os.Args[1]

	torrent, err := torrent.New(filePath)
	if err != nil {
		log.Fatalln(err)
	}

	for _, peer := range torrent.Peers {
		log.Println(peer)
	}

	network := network.NewPeerNetwork(torrent)

	// start the torrent network
	network.Start()
}
