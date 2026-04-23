package torrent

import (
	"encoding/binary"
	"errors"
	"net"
	"net/http"
	"net/url"
	"strconv"

	"github.com/jackpal/bencode-go"
	"github.com/rohit20001221/ripple/peer"
)

type TorrentFile struct {
	Announce    string
	Length      int
	Name        string
	PieceLength int
	PieceHashes [][20]byte
	InfoHash    [20]byte
	PeerID      [20]byte
}

type Torrent struct {
	Announce    string
	Length      int
	Name        string
	PieceLength int
	PieceHashes [][20]byte
	InfoHash    [20]byte
	PeerID      [20]byte
	Peers       []peer.Peer
}

type trackerResponse struct {
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
}

func (t *TorrentFile) getTrackerURL() string {
	baseURL, _ := url.Parse(t.Announce)

	params := url.Values{}

	params.Add("info_hash", string(t.InfoHash[:]))
	params.Add("peer_id", string(t.PeerID[:]))
	params.Add("port", "6634")
	params.Add("uploaded", "0")
	params.Add("downloaded", "0")
	params.Add("compact", "1")
	params.Add("left", strconv.Itoa(t.Length))

	baseURL.RawQuery = params.Encode()

	return baseURL.String()
}

func parsePeerResponse(peersResponse []byte) ([]peer.Peer, error) {
	if len(peersResponse)%6 != 0 {
		return []peer.Peer{}, errors.New("invalid peers response")
	}

	peers := make([]peer.Peer, 0)

	offset := 0
	limit := 6

	for offset < len(peersResponse) {
		blob := peersResponse[offset : offset+limit]

		ip := net.IP(blob[:4])
		port := binary.BigEndian.Uint16(blob[4:6])

		peers = append(peers, peer.Peer{
			IP:   ip,
			Port: port,
		})

		offset += limit
	}

	return peers, nil
}

func (t *TorrentFile) getPeers() ([]peer.Peer, error) {
	url := t.getTrackerURL()

	resp, err := http.Get(url)
	if err != nil {
		return []peer.Peer{}, err
	}

	defer resp.Body.Close()

	var trackerResponse trackerResponse

	err = bencode.Unmarshal(resp.Body, &trackerResponse)
	if err != nil {
		return []peer.Peer{}, err
	}

	peers, err := parsePeerResponse([]byte(trackerResponse.Peers))
	if err != nil {
		return []peer.Peer{}, err
	}

	return peers, nil
}

func New(path string) (*Torrent, error) {
	torrentFile, err := Open(path)
	if err != nil {
		return nil, err
	}

	peers, err := torrentFile.getPeers()
	if err != nil {
		return nil, err
	}

	return &Torrent{
		Announce:    torrentFile.Announce,
		Length:      torrentFile.Length,
		Name:        torrentFile.Name,
		PieceLength: torrentFile.PieceLength,
		PieceHashes: torrentFile.PieceHashes,
		InfoHash:    torrentFile.InfoHash,
		PeerID:      torrentFile.InfoHash,
		Peers:       peers,
	}, nil
}
