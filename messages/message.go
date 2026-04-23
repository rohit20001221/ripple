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
	binary.Write(w, binary.BigEndian, msg.Payload)
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
