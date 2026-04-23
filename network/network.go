package network

import (
	"net"
	"sync"

	"github.com/rohit20001221/ripple/peer"
	"github.com/rohit20001221/ripple/torrent"
)

type Client struct {
	Peer      peer.Peer
	Torrent   *torrent.Torrent
	Conn      net.Conn
	TaskQueue chan *pieceTask
}

type PeerNetwork struct {
	Clients   []*Client
	Torrent   *torrent.Torrent
	wg        sync.WaitGroup
	TaskQueue chan *pieceTask
}

func NewPeerNetwork(torrent *torrent.Torrent) *PeerNetwork {
	clients := make([]*Client, 0)

	taskQueue := make(chan *pieceTask, len(torrent.PieceHashes))

	for _, peer := range torrent.Peers {
		// establish a tcp connection
		conn, err := net.Dial("tcp", peer.String())
		if err != nil {
			continue
		}

		clients = append(clients, &Client{
			Peer:      peer,
			Torrent:   torrent,
			Conn:      conn,
			TaskQueue: taskQueue,
		})
	}

	return &PeerNetwork{
		Clients:   clients,
		Torrent:   torrent,
		wg:        sync.WaitGroup{},
		TaskQueue: taskQueue,
	}
}

func (n *PeerNetwork) Start() {
	for i, piece := range n.Torrent.PieceHashes {
		start, end := n.Torrent.GetPeicePosition(i)

		n.TaskQueue <- &pieceTask{
			start: start,
			end:   end,
			hash:  piece,
		}
	}

	for _, client := range n.Clients {
		n.wg.Go(client.startDownloadHandler)
	}

	n.wg.Wait()
}
