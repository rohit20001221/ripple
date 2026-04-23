package bitfield

type BitField []byte

func (bf BitField) HasIndex(index int) bool {
	byteIndex := index / 8
	offset := index % 8

	if byteIndex < 0 || byteIndex >= len(bf) {
		return false
	}

	return bf[byteIndex]>>uint(7-offset)&1 != 0
}

func (bf BitField) SetIndex(index int) {
	byteIndex := index / 8
	offset := index % 8

	// silently discard invalid bounded index
	if byteIndex < 0 || byteIndex >= len(bf) {
		return
	}

	bf[byteIndex] |= 1 << uint(7-offset)
}
