package network

import (
	"log"
	"net"

	"github.com/rohit20001221/ripple/peer"
	"github.com/rohit20001221/ripple/torrent"
)

type Client struct {
	Peer        peer.Peer
	Torrent     *torrent.Torrent
	Conn        net.Conn
	TaskQueue   chan *pieceTask
	PieceResult chan *pieceResponse
}

type PeerNetwork struct {
	Clients     []*Client
	Torrent     *torrent.Torrent
	TaskQueue   chan *pieceTask
	PieceResult chan *pieceResponse
}

func NewPeerNetwork(torrent *torrent.Torrent) *PeerNetwork {
	clients := make([]*Client, 0)

	taskQueue := make(chan *pieceTask, len(torrent.PieceHashes))
	pieceResult := make(chan *pieceResponse)

	for _, peer := range torrent.Peers {
		// establish a tcp connection
		conn, err := net.Dial("tcp", peer.String())
		if err != nil {
			continue
		}

		clients = append(clients, &Client{
			Peer:        peer,
			Torrent:     torrent,
			Conn:        conn,
			TaskQueue:   taskQueue,
			PieceResult: pieceResult,
		})
	}

	return &PeerNetwork{
		Clients:     clients,
		Torrent:     torrent,
		TaskQueue:   taskQueue,
		PieceResult: pieceResult,
	}
}

func (n *PeerNetwork) Start() {
	for i, piece := range n.Torrent.PieceHashes {
		start, end, size := n.Torrent.GetPeicePosition(i)

		n.TaskQueue <- &pieceTask{
			start: start,
			end:   end,
			size:  size,
			hash:  piece,
			index: i,
		}
	}

	for _, client := range n.Clients {
		go client.startDownloadHandler()
	}

	for range len(n.Torrent.PieceHashes) {
		piece := <-n.PieceResult

		// collect indivudial pieces
		log.Println(string(piece.piece))
	}

	close(n.PieceResult)
	close(n.TaskQueue)
}
