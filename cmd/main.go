package main

import (
	"log"
	"os"

	"github.com/rohit20001221/ripple/network"
	"github.com/rohit20001221/ripple/torrent"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("error: provide the path to torrent")
	}

	filePath := os.Args[1]
	outPath := os.Args[2]

	torrent, err := torrent.New(filePath, outPath)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Announce:", torrent.Announce)
	log.Println("File Length:", torrent.Length)
	log.Println("Piece Length:", torrent.PieceLength)

	for _, peer := range torrent.Peers {
		log.Println(peer)
	}

	network := network.NewPeerNetwork(torrent)

	// start the torrent network
	network.Start()
}
