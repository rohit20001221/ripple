package torrent

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"os"

	"github.com/jackpal/bencode-go"
)

type torrentFileInfo struct {
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
	PieceLength int    `bencode:"piece length"`
	Pieces      string `bencode:"pieces"`
}

type torrentFile struct {
	Announce string          `bencode:"announce"`
	Info     torrentFileInfo `bencode:"info"`
}

func (t *torrentFile) InfoHash() [20]byte {
	var buf bytes.Buffer

	bencode.Marshal(&buf, t.Info)

	h := sha1.New()
	h.Write(buf.Bytes())
	infoHash := h.Sum(nil)

	return [20]byte(infoHash)
}

func (t *torrentFile) PieceHashes() [][20]byte {
	pieces := make([][20]byte, 0)

	offset := 0
	limit := 20

	piceBytes := []byte(t.Info.Pieces)

	for offset < len(piceBytes) {
		piece := piceBytes[offset : offset+limit]

		pieces = append(pieces, [20]byte(piece))

		offset += limit
	}

	return pieces
}

func Open(path string) (*TorrentFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	var torrentFile torrentFile

	err = bencode.Unmarshal(file, &torrentFile)
	if err != nil {
		return nil, err
	}

	infoHash := torrentFile.InfoHash()
	pieceHashes := torrentFile.PieceHashes()

	var peerID [20]byte
	rand.Read(peerID[:])

	return &TorrentFile{
		Announce:    torrentFile.Announce,
		Length:      torrentFile.Info.Length,
		Name:        torrentFile.Info.Name,
		PieceLength: torrentFile.Info.PieceLength,
		PieceHashes: pieceHashes,
		InfoHash:    infoHash,
		PeerID:      peerID,
	}, nil
}
