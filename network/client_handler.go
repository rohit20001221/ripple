package network

import (
	"bytes"
	"log"

	"github.com/rohit20001221/ripple/messages"
)

func (c *Client) Close() {
	c.Conn.Close()
}

func (c *Client) startDownloadHandler() {
	defer c.Close()

	handshake := messages.NewHandshake(c.Torrent.InfoHash, c.Torrent.PeerID)

	var buf bytes.Buffer
	handshake.Encode(&buf)

	log.Println("sending handshake message to:", c.Peer.String())
	_, err := c.Conn.Write(buf.Bytes())
	if err != nil {
		log.Println(err)
		return
	}

	for task := range c.TaskQueue {
		log.Println(task)
	}
}
