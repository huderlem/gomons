package util

import (
	"math"
)

// BoolToByte returns 1 if the given bool is true, and 0 otherwise.
func BoolToByte(b bool) byte {
	if b {
		return 1
	}
	return 0
}

// BoolToU32 returns 1 if the given bool is true, and 0 otherwise.
func BoolToU32(b bool) uint32 {
	if b {
		return 1
	}
	return 0
}

// CheckBit checks if the bit is set at the given position.
func CheckBit(value byte, index int) bool {
	return (value & (1 << index)) != 0
}

// WriteBit set or clears the bit at the given position.
func WriteBit(value byte, index int, set bool) byte {
	if set {
		return value | (1 << index)
	}
	return value &^ (1 << index)
}

// WriteBits set or clears the bits at the given position.
func WriteBits(originalValue, value byte, index int, width int) byte {
	mask := byte((int(math.Pow(2, float64(width))) - 1) << index)
	return (originalValue &^ mask) | (value << index)
}

// WriteBitsU16 set or clears the bits at the given position.
func WriteBitsU16(originalValue, value uint16, index int, width int) uint16 {
	mask := uint16((int(math.Pow(2, float64(width))) - 1) << index)
	return (originalValue &^ mask) | (value << index)
}

// WriteBitsU32 set or clears the bits at the given position.
func WriteBitsU32(originalValue, value uint32, index int, width int) uint32 {
	mask := uint32((int(math.Pow(2, float64(width))) - 1) << index)
	return (originalValue &^ mask) | (value << index)
}

// CheckBitU32 checks if the bit is set at the given position.
func CheckBitU32(value uint32, index int) bool {
	return (value & (1 << index)) != 0
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
