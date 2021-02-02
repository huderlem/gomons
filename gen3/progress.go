package gen3

// GetFlag gets the specified event flag.
func (s *SaveData) GetFlag(flag int) bool {
	offset := uint(flag) / 8
	bit := uint(flag) % 8
	return s.readBit(0x21F0+offset, bit)
}

// SetFlag sets the specified event flag.
func (s *SaveData) SetFlag(value bool, flag int) {
	offset := uint(flag) / 8
	bit := uint(flag) % 8
	s.writeBit(value, 0x21F0+offset, bit)
}

// GetVar gets the specified event variable.
func (s *SaveData) GetVar(variable int) uint16 {
	offset := uint(variable-0x4000) * 2
	return s.readU16(0x231C + offset)
}

// SetVar sets the specified event variable.
func (s *SaveData) SetVar(value uint16, variable int) {
	offset := uint(variable-0x4000) * 2
	s.writeU16(value, 0x231C+offset)
}

// GetGameStat gets the specified game progress stat.
func (s *SaveData) GetGameStat(gameStat int) uint32 {
	key := s.GetEncryptionKey()
	offset := uint(gameStat) * 4
	return s.readU32(0x251C+offset) ^ key
}

// SetGameStat sets the specified game progress stat.
func (s *SaveData) SetGameStat(value uint32, gameStat int) {
	key := s.GetEncryptionKey()
	offset := uint(gameStat) * 4
	s.writeU32(value^key, 0x251C+offset)
}
