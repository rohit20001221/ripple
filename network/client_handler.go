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

	handshakeResponse, err := messages.NewHandshake().Decode(c.Conn)
	if err != nil {
		log.Println(err)
		return
	}

	// if the info hash don't match close the peer connection
	if !bytes.Equal(c.Torrent.InfoHash[:], handshakeResponse.InfoHash[:]) {
		c.Conn.Close()

		log.Println("info hash didn't match")
		return
	}

	// get the bitfield message
	bitfield, err := messages.ReadBitField(c.Conn)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("bitfield: %b", bitfield)
	choked := true

	// tell the peer that you are interested to connect with it
	messages.SendInterested(c.Conn)

	// start exchanging the messages
	for task := range c.TaskQueue {
		if !bitfield.HasIndex(task.index) {
			// if the current peer is not having the piece
			// enque the task and continue
			log.Println("peer dosen't have the block")
			c.TaskQueue <- task
			continue
		}

		message, err := messages.ReadPeerMessage(c.Conn)
		if err != nil {
			log.Println(err)
			return
		}

		switch message.MessageID {
		case messages.MSG_CHOKE:
			c.TaskQueue <- task
			choked = true
		case messages.MSG_UNCHOKE:
			choked = false
		}

		if !choked {
			totalBlocks := 0
			piece := make([]byte, task.size)

			// process the piece

			// break the piece into blocks of 16KiB (16 * 1024 bytes) and send a request message for each block
			// payload for request message is as follows
			// index: 0 based index of the piece
			// begin: 0 based byte offset within the piece
			// length of block in bytes
			// length of blocks is 16 * 1024 or 2^14 for each block except last one

			// messages.RequestPiece(index, begin, length)
			for _, block := range task.Blocks() {
				totalBlocks++
				messages.RequestPiece(block.pieceIndex, block.begin, block.length, c.Conn)
			}

			// wait until we receive all the blocks of the pieces
			for totalBlocks > 0 {
				message, err := messages.ReadPeerMessage(c.Conn)
				if err != nil {
					log.Println(err)
					return
				}

				if message.MessageID == messages.MSG_PIECE {
					// parse the message payload to piece block message
					_, begin, block := messages.ReadBlock(message)

					start := begin
					end := start + len(block)

					copy(piece[start:end], block)

					totalBlocks--
				}
			}

			// if the piece received is invalid then assign it back to queue
			if !task.CheckIntegrity(piece) {
				log.Println("invalid integrity")

				c.TaskQueue <- task
				continue
			}

			c.PieceResult <- &pieceResponse{
				start: task.start,
				piece: piece,
			}
		}

	}
}
