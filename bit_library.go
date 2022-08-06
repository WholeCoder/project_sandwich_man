package main

type BitsetByte []byte

func InitNewByteset(bray []byte) BitsetByte {
	return BitsetByte(bray)
}

func (b BitsetByte) GetBit(index int) bool {
	pos := index / 8
	j := uint(7 - index%8)
	fmt.Printf("\nbit at element (GetBit):  %v", b[pos])
	return (b[pos] & (byte(1) << j)) != 0
}

func (b BitsetByte) SetBit(index int, value bool) {
	pos := index / 8
	j := uint(7 - index%8)
	fmt.Println("j =", j)
	if value {
		b[pos] |= (byte(1) << j)
	} else {
		b[pos] &= ^(byte(1) << j)
	}
	fmt.Printf("\nBitset element (in SetBit):  %8b\n", b[pos])
}

func (b BitsetByte) Len() int {
	return 8 * len(b)
}
