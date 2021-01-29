package util

// BoolToByte returns 1 if the given bool is true, and 0 otherwise.
func BoolToByte(b bool) byte {
	if b {
		return 1
	}
	return 0
}

// CheckBitSetInArray checks if the bit is set at the given position.
func CheckBitSetInArray(bitArray []byte, index int) bool {
	slot := index / 8
	bit := index % 8
	return (bitArray[slot] & (1 << bit)) != 0
}

// SetBitInArray sets the bit at the given position.
func SetBitInArray(bitArray []byte, index int) {
	slot := index / 8
	bit := index % 8
	bitArray[slot] |= (1 << bit)
}

// ClearBitInArray clears the bit at the given position.
func ClearBitInArray(bitArray []byte, index int) {
	slot := index / 8
	bit := index % 8
	bitArray[slot] &^= (1 << bit)
}
