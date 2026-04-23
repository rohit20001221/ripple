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

func ReadPeerMessage(r io.Reader) (PeerMessage, error) {
	var message PeerMessage

	binary.Read(r, binary.BigEndian, &message.Length)
	binary.Read(r, binary.BigEndian, &message.MessageID)

	message.Payload = make([]byte, message.Length-1)
	err := binary.Read(r, binary.BigEndian, &message.Payload)

	return message, err
}

func (msg PeerMessage) Write(w io.Writer) error {
	binary.Write(w, binary.BigEndian, msg.Length)
	err := binary.Write(w, binary.BigEndian, msg.MessageID)
	if err != nil {
		return err
	}

	if len(msg.Payload) > 0 {
		err := binary.Write(w, binary.BigEndian, msg.Payload)
		if err != nil {
			return err
		}
	}

	return nil
}

func ReadBitField(r io.Reader) (bitfield.BitField, error) {
	message, err := ReadPeerMessage(r)

	return bitfield.BitField(message.Payload), err
}

// write interested message
func SendInterested(w io.Writer) error {
	message := PeerMessage{
		Length:    1,
		MessageID: MSG_INTERESTED,
		Payload:   make([]byte, 0),
	}

	err := message.Write(w)
	return err
}

// request for a piece
func RequestPiece(index, begin, length uint32, w io.Writer) error {
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
	err := binary.Write(w, binary.BigEndian, length)

	return err
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
