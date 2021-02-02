package gen3

// GetLocalTimeOffset gets the save's local time offset, which is used
// with the RTC to determine the game's local time.
func (s *SaveData) GetLocalTimeOffset() (int16, int8, int8, int8) {
	days := s.readS16(0x98)
	hours := s.readS8(0x9A)
	minutes := s.readS8(0x9B)
	seconds := s.readS8(0x9C)
	return days, hours, minutes, seconds
}

// SetLocalTimeOffset sets the save's local time offset, which is used
// with the RTC to determine the game's local time.
func (s *SaveData) SetLocalTimeOffset(days int16, hours, minutes, seconds int8) {
	s.writeS16(days, 0x98)
	s.writeS8(hours, 0x9A)
	s.writeS8(minutes, 0x9B)
	s.writeS8(seconds, 0x9C)
}

// GetLastBerryTreeUpdate gets the save's last berry tree update time offset,
// which is compared to the game's local time to determine if the berry trees
// need to be updated.
func (s *SaveData) GetLastBerryTreeUpdate() (int16, int8, int8, int8) {
	days := s.readS16(0xA0)
	hours := s.readS8(0xA2)
	minutes := s.readS8(0xA3)
	seconds := s.readS8(0xA4)
	return days, hours, minutes, seconds
}

// SetLastBerryTreeUpdate sets the save's last berry tree update time offset,
// which is compared to the game's local time to determine if the berry trees
// need to be updated.
func (s *SaveData) SetLastBerryTreeUpdate(days int16, hours, minutes, seconds int8) {
	s.writeS16(days, 0xA0)
	s.writeS8(hours, 0xA2)
	s.writeS8(minutes, 0xA3)
	s.writeS8(seconds, 0xA4)
}
