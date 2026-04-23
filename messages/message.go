package messages

import (
	"encoding/binary"
	"io"

	"github.com/rohit20001221/ripple/bitfield"
)

type PeerMessage struct {
	Length    uint32
	MessageID messageID
	Payload   []byte
}

func NewPeerMessage(messageID messageID, payload []byte) PeerMessage {
	return PeerMessage{
		Length:    uint32(1 + len(payload)),
		MessageID: messageID,
		Payload:   payload,
	}
}

func ReadPeerMessage(r io.Reader) PeerMessage {
	var message PeerMessage

	binary.Read(r, binary.BigEndian, &message.Length)
	binary.Read(r, binary.BigEndian, &message.MessageID)

	message.Payload = make([]byte, message.Length-1)
	binary.Read(r, binary.BigEndian, &message.Payload)

	return message
}

func (msg PeerMessage) Write(w io.Writer) {
	binary.Write(w, binary.BigEndian, msg.Length)
	binary.Write(w, binary.BigEndian, msg.MessageID)

	if len(msg.Payload) > 0 {
		binary.Write(w, binary.BigEndian, msg.Payload)
	}
}

func ReadBitField(r io.Reader) bitfield.BitField {
	message := ReadPeerMessage(r)

	return bitfield.BitField(message.Payload)
}

// write interested message
func SendInterested(w io.Writer) {
	message := PeerMessage{
		Length:    1,
		MessageID: MSG_INTERESTED,
		Payload:   make([]byte, 0),
	}

	message.Write(w)
}

// get unchoke message
func IsUnchoke(r io.Reader) bool {
	msg := ReadPeerMessage(r)

	return msg.MessageID == MSG_UNCHOKE
}

// request for a piece
func RequestPiece(index, begin, length uint32, w io.Writer) {
	// one uint32 => 4 bytes
	message := PeerMessage{
		Length:    13,
		MessageID: MSG_REQUEST,
	}

	// write the message header
	message.Write(w)

	// write the payload
	binary.Write(w, binary.BigEndian, index)
	binary.Write(w, binary.BigEndian, begin)
	binary.Write(w, binary.BigEndian, length)
}

func ReadBlock(message PeerMessage) (int, int, []byte) {
	var index uint32
	var begin uint32
	var block []byte

	index = binary.BigEndian.Uint32(message.Payload[0:4])
	begin = binary.BigEndian.Uint32(message.Payload[4:8])
	block = message.Payload[8:]

	return int(index), int(begin), block
}
