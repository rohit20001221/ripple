package network

import (
	"log"
	"net"
	"os"
	"time"

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

		log.Println("establishing connection:", peer.String())
		conn, err := net.DialTimeout("tcp", peer.String(), time.Duration(time.Second*2))
		if err != nil {
			log.Println(err)
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
		log.Println("enqueing piece:", i)
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

	outFile, err := os.Create(n.Torrent.OutPath)
	if err != nil {
		log.Fatalln(err)
	}

	defer outFile.Close()

	for range len(n.Torrent.PieceHashes) {
		piece := <-n.PieceResult

		log.Println("writing piece at:", piece.start)
		outFile.Seek(int64(piece.start), 0)
		outFile.Write(piece.piece)
	}

	close(n.PieceResult)
	close(n.TaskQueue)
}
