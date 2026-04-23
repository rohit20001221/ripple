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

func NewHandshake() *HandshakeMessage {
	infoHash := make([]byte, 20)
	peerID := make([]byte, 20)

	return &HandshakeMessage{
		Length:         19,
		ProtocolString: []byte("BitTorrent protocol"),
		Reserved:       [8]byte{0, 0, 0, 0, 0, 0, 0, 0},
		InfoHash:       [20]byte(infoHash),
		PeerID:         [20]byte(peerID),
	}
}

func (h *HandshakeMessage) SetInfoHash(hash [20]byte) {
	h.InfoHash = hash
}

func (h *HandshakeMessage) SetPeerID(peerID [20]byte) {
	h.PeerID = peerID
}

func (h *HandshakeMessage) Encode(buf io.Writer) {
	binary.Write(buf, binary.BigEndian, h.Length)

	protocolString := make([]byte, 19)
	copy(protocolString, h.ProtocolString)

	binary.Write(buf, binary.BigEndian, protocolString)
	binary.Write(buf, binary.BigEndian, h.Reserved)
	binary.Write(buf, binary.BigEndian, h.InfoHash)
	binary.Write(buf, binary.BigEndian, h.PeerID)
}

func (h *HandshakeMessage) Decode(r io.Reader) *HandshakeMessage {
	h.ProtocolString = make([]byte, 19)

	binary.Read(r, binary.BigEndian, &h.Length)
	binary.Read(r, binary.BigEndian, &h.ProtocolString)
	binary.Read(r, binary.BigEndian, &h.Reserved)
	binary.Read(r, binary.BigEndian, &h.InfoHash)
	binary.Read(r, binary.BigEndian, &h.PeerID)

	return h
}
