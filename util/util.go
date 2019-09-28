package util

// BoolToByte returns 1 if the given bool is true, and 0 otherwise.
func BoolToByte(b bool) byte {
	if b {
		return 1
	}
	return 0
}
