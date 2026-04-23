package network

import (
	"net"
	"sync"

	"github.com/rohit20001221/ripple/peer"
	"github.com/rohit20001221/ripple/torrent"
)

type Client struct {
	Peer    peer.Peer
	Torrent *torrent.Torrent
	Conn    net.Conn
}

type PeerNetwork struct {
	Clients []*Client
	Torrent *torrent.Torrent
	wg      sync.WaitGroup
}

func NewPeerNetwork(torrent *torrent.Torrent) *PeerNetwork {
	clients := make([]*Client, 0)

	for _, peer := range torrent.Peers {
		// establish a tcp connection
		conn, err := net.Dial("tcp", peer.String())
		if err != nil {
			continue
		}

		clients = append(clients, &Client{
			Peer:    peer,
			Torrent: torrent,
			Conn:    conn,
		})
	}

	return &PeerNetwork{
		Clients: clients,
		Torrent: torrent,
		wg:      sync.WaitGroup{},
	}
}

func (n *PeerNetwork) Start() {
	for _, client := range n.Clients {
		n.wg.Go(client.startDownloadHandler)
	}

	n.wg.Wait()
}
