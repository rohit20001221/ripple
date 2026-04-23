package network

type pieceTask struct {
	start int
	end   int
	size  int
	index int
	hash  [20]byte
}
