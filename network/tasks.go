package network

import (
	"bytes"
	"crypto/sha1"
)

const (
	MAX_BLOCK_SIZE = 16384
)

type pieceTask struct {
	start int
	end   int
	size  int
	index int
	hash  [20]byte
}

type pieceResponse struct {
	start int
	piece []byte
}

type pieceBlock struct {
	pieceIndex uint32
	begin      uint32
	end        uint32
	length     uint32
}

func (task *pieceTask) Blocks() []pieceBlock {
	blocks := make([]pieceBlock, 0)

	offset := 0
	for offset < task.size {
		end := min(offset+MAX_BLOCK_SIZE, task.size)

		blocks = append(blocks, pieceBlock{
			pieceIndex: uint32(task.index),
			begin:      uint32(offset),
			end:        uint32(end),
			length:     uint32(end - offset),
		})

		offset += MAX_BLOCK_SIZE
	}

	return blocks
}

func (task *pieceTask) CheckIntegrity(piece []byte) bool {
	h := sha1.New()
	h.Write(piece)
	hash := h.Sum(nil)

	return bytes.Equal(task.hash[:], hash)
}
