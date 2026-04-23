package messages

type messageID uint8

const (
	MSG_CHOKE      messageID = 0
	MSG_UNCHOKE    messageID = 1
	MSG_INTERESTED messageID = 2
	MSG_BITFIELD   messageID = 5
	MSG_REQUEST    messageID = 6
	MSG_PIECE      messageID = 7
)
