package network

type pieceTask struct {
	start int
	end   int
	size  int
	hash  [20]byte
}
