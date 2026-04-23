package bitfield

type BitField []byte

func (b BitField) HasIndex(index int) bool {
	idx, offset := index/8, index%8

	return (b[idx] & (1 << uint(offset))) != 0
}

func (b BitField) SetIndex(index int) {
	idx, offset := index/8, index%8

	b[idx] |= (1 << uint(offset))
}
