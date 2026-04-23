package messages

import (
	"encoding/binary"
	"io"
)

type HandshakeMessage struct {
	Length         byte
	ProtocolString []byte // always BitTorrent protocol [size: always 19 bytes]
	Reserved       [8]byte
	InfoHash       [20]byte
	PeerID         [20]byte
}

func NewHandshake(infoHash, peerID [20]byte) *HandshakeMessage {
	return &HandshakeMessage{
		Length:         19,
		ProtocolString: []byte("BitTorrent protocol"),
		Reserved:       [8]byte{0, 0, 0, 0, 0, 0, 0, 0},
		InfoHash:       infoHash,
		PeerID:         peerID,
	}
}

func (h *HandshakeMessage) Encode(buf io.Writer) *HandshakeMessage {
	binary.Write(buf, binary.BigEndian, h.Length)

	protocolString := make([]byte, 19)
	copy(protocolString, h.ProtocolString)

	binary.Write(buf, binary.BigEndian, protocolString)
	binary.Write(buf, binary.BigEndian, h.Reserved)
	binary.Write(buf, binary.BigEndian, h.InfoHash)
	binary.Write(buf, binary.BigEndian, h.PeerID)

	return h
}

func (h *HandshakeMessage) Decode(r io.Reader) *HandshakeMessage {
	h.ProtocolString = make([]byte, 19)

	binary.Read(r, binary.BigEndian, h.Length)
	binary.Read(r, binary.BigEndian, h.ProtocolString)
	binary.Read(r, binary.BigEndian, h.Reserved)
	binary.Read(r, binary.BigEndian, h.InfoHash)
	binary.Read(r, binary.BigEndian, h.PeerID)

	return h
}
