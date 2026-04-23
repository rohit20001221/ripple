package network

import (
	"bytes"
	"encoding/hex"
	"log"

	"github.com/rohit20001221/ripple/messages"
)

func (c *Client) Close() {
	c.Conn.Close()
}

func (c *Client) startDownloadHandler() {
	defer c.Close()

	handshake := messages.NewHandshake()
	handshake.SetInfoHash(c.Torrent.InfoHash)
	handshake.SetPeerID(c.Torrent.PeerID)

	var buf bytes.Buffer
	handshake.Encode(&buf)

	log.Println("sending handshake message to:", c.Peer.String())
	_, err := c.Conn.Write(buf.Bytes())
	if err != nil {
		log.Println(err)
		return
	}

	handshakeResponse := messages.NewHandshake().Decode(c.Conn)

	// if the info hash don't match close the peer connection
	if !bytes.Equal(c.Torrent.InfoHash[:], handshakeResponse.InfoHash[:]) {
		c.Conn.Close()

		log.Println("info hash didn't match")
		return
	}

	// get the bitfield message
	bitfield := messages.ReadBitField(c.Conn)

	// start exchanging the messages
	for task := range c.TaskQueue {
		if !bitfield.HasIndex(task.index) {
			// if the current peer is not having the piece
			// enque the task and continue
			c.TaskQueue <- task
			continue
		}

		log.Println("Task hash:", hex.EncodeToString(task.hash[:]))
	}
}
