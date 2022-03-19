package gen3

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
)

const gameSectionDataSize = 0xF80
const sectionRawDataSize = 0xFF4
const sectionSize = 0x1000
const saveFileSize = sectionSize * 32
const numGameSaveSections = 14
const numHallOfFameSections = 2
const securityValue = 0x8012025

// GameSaveSection represents one of the logical sections of the save data structure.
type GameSaveSection struct {
	data     []byte
	id       uint16
	checksum uint16
	security uint32
	counter  uint32
	empty    bool
}

// HallOfFameSection represents one of the logical hall of fame sections of the save data structure.
type HallOfFameSection struct {
	data     []byte
	checksum uint16
	security uint32
	empty    bool
}

// SaveData is full representation of the save data.
type SaveData struct {
	gameSaveA      [numGameSaveSections]GameSaveSection
	gameSaveB      [numGameSaveSections]GameSaveSection
	hallOfFame     [numHallOfFameSections]HallOfFameSection
	trainerHill    []byte
	recordedBattle []byte
	activeGameSave *[numGameSaveSections]GameSaveSection
}

// LoadSaveFile reads a game save file.
func LoadSaveFile(filename string) (SaveData, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return SaveData{}, err
	}
	if len(bytes) != saveFileSize {
		return SaveData{}, fmt.Errorf("Expected save file size is %d bytes, but %s is %d bytes", saveFileSize, filename, len(bytes))
	}

	return LoadSaveFileFromBytes(bytes)
}

// LoadSaveFileFromBytes reads a game save file from slice of bytes.
func LoadSaveFileFromBytes(bytes []byte) (SaveData, error) {
	var err error
	saveData := SaveData{
		activeGameSave: nil,
	}

	// First, load the two sets of 14 game save sections.
	for i := 0; i < numGameSaveSections; i++ {
		offset := i * sectionSize
		saveData.gameSaveA[i], err = loadGameSaveSection(bytes[offset : offset+sectionSize])
		if err != nil {
			return saveData, err
		}
		offset = (i + numGameSaveSections) * sectionSize
		saveData.gameSaveB[i], err = loadGameSaveSection(bytes[offset : offset+sectionSize])
		if err != nil {
			return saveData, err
		}
	}

	// Load Hall of Fame save sections.
	for i := 0; i < numHallOfFameSections; i++ {
		offset := (i + numGameSaveSections*2) * sectionSize
		saveData.hallOfFame[i], err = loadHallOfFameSection(bytes[offset : offset+sectionSize])
		if err != nil {
			return saveData, err
		}
	}

	// Load Trainer Hill save section.
	offset := (numGameSaveSections*2 + numHallOfFameSections) * sectionSize
	saveData.trainerHill = loadTrainerHillSection(bytes[offset : offset+sectionSize])

	// Load Recorded Battle save section.
	offset = (numGameSaveSections*2 + numHallOfFameSections + 1) * sectionSize
	saveData.recordedBattle = loadRecordedBattleSection(bytes[offset : offset+sectionSize])

	return saveData, saveData.CheckCorruption()
}

func isSectionEmpty(sectionBytes []byte) bool {
	for _, b := range sectionBytes[:gameSectionDataSize] {
		if b != 0xFF {
			return false
		}
	}
	return true
}

func loadGameSaveSection(sectionBytes []byte) (GameSaveSection, error) {
	section := GameSaveSection{
		data: make([]byte, sectionRawDataSize),
	}
	if len(sectionBytes) != sectionSize {
		return section, fmt.Errorf("Failed to load save section because only %d bytes were provided", len(sectionBytes))
	}
	n := copy(section.data, sectionBytes)
	if n != sectionRawDataSize {
		return section, fmt.Errorf("Failed to load save section because only %d bytes could be copied", n)
	}
	section.id = binary.LittleEndian.Uint16(sectionBytes[sectionRawDataSize : sectionRawDataSize+2])
	section.checksum = binary.LittleEndian.Uint16(sectionBytes[sectionRawDataSize+2 : sectionRawDataSize+4])
	section.security = binary.LittleEndian.Uint32(sectionBytes[sectionRawDataSize+4 : sectionRawDataSize+8])
	section.counter = binary.LittleEndian.Uint32(sectionBytes[sectionRawDataSize+8 : sectionRawDataSize+12])
	section.empty = isSectionEmpty(sectionBytes)
	return section, nil
}

func loadHallOfFameSection(sectionBytes []byte) (HallOfFameSection, error) {
	section := HallOfFameSection{
		data: make([]byte, sectionRawDataSize),
	}
	if len(sectionBytes) != sectionSize {
		return section, fmt.Errorf("Failed to load save section because only %d bytes were provided", len(sectionBytes))
	}
	n := copy(section.data, sectionBytes)
	if n != sectionRawDataSize {
		return section, fmt.Errorf("Failed to load save section because only %d bytes could be copied", n)
	}
	section.checksum = binary.LittleEndian.Uint16(sectionBytes[sectionRawDataSize : sectionRawDataSize+2])
	section.security = binary.LittleEndian.Uint32(sectionBytes[sectionRawDataSize+4 : sectionRawDataSize+8])
	section.empty = isSectionEmpty(sectionBytes)
	return section, nil
}

func loadTrainerHillSection(sectionBytes []byte) []byte {
	return sectionBytes
}

func loadRecordedBattleSection(sectionBytes []byte) []byte {
	return sectionBytes
}

func (s *SaveData) getLatestGameSaveSection() *[numGameSaveSections]GameSaveSection {
	counterA := s.gameSaveA[len(s.gameSaveA)-1].counter
	if s.gameSaveA[len(s.gameSaveA)-1].empty {
		counterA = 0
	}
	counterB := s.gameSaveB[len(s.gameSaveB)-1].counter
	if s.gameSaveB[len(s.gameSaveB)-1].empty {
		counterB = 0
	}
	if counterA > counterB {
		s.activeGameSave = &s.gameSaveA
	} else {
		s.activeGameSave = &s.gameSaveB
	}
	return s.activeGameSave
}

// CheckCorruption inspects the save data to see if any of the data is corrupted.
func (s *SaveData) CheckCorruption() error {
	gameSaveSection := s.getLatestGameSaveSection()
	for i := 0; i < numGameSaveSections; i++ {
		if err := gameSaveSection[i].checkCorruption(); err != nil {
			return err
		}
	}
	for i := 0; i < numHallOfFameSections; i++ {
		if err := s.hallOfFame[i].checkCorruption(i); err != nil {
			return err
		}
	}
	return nil
}

// FixChecksums recalculates all of the save data checksums.
func (s *SaveData) FixChecksums() {
	gameSaveSections := s.getLatestGameSaveSection()
	for i := 0; i < numGameSaveSections; i++ {
		gameSaveSections[i].checksum = gameSaveSections[i].calculateChecksum()
	}
	for i := 0; i < numHallOfFameSections; i++ {
		s.hallOfFame[i].empty = isSectionEmpty(s.hallOfFame[i].data)
		if !s.hallOfFame[i].empty {
			s.hallOfFame[i].checksum = s.hallOfFame[i].calculateChecksum()
		}
	}
}

func (s *SaveData) Write(w io.Writer) error {
	s.FixChecksums()
	var err error
	tryWrite := func(bytes []byte) {
		if err != nil {
			return
		}
		_, err = w.Write(bytes)
	}

	for i := 0; i < len(s.gameSaveA); i++ {
		section := s.gameSaveA[i]
		tryWrite(section.data)
		idBytes := make([]byte, 2)
		binary.LittleEndian.PutUint16(idBytes, section.id)
		tryWrite(idBytes)
		checksumBytes := make([]byte, 2)
		binary.LittleEndian.PutUint16(checksumBytes, section.checksum)
		tryWrite(checksumBytes)
		securityBytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(securityBytes, section.security)
		tryWrite(securityBytes)
		counterBytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(counterBytes, section.counter)
		tryWrite(counterBytes)
	}
	for i := 0; i < len(s.gameSaveB); i++ {
		section := s.gameSaveB[i]
		tryWrite(section.data)
		idBytes := make([]byte, 2)
		binary.LittleEndian.PutUint16(idBytes, section.id)
		tryWrite(idBytes)
		checksumBytes := make([]byte, 2)
		binary.LittleEndian.PutUint16(checksumBytes, section.checksum)
		tryWrite(checksumBytes)
		securityBytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(securityBytes, section.security)
		tryWrite(securityBytes)
		counterBytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(counterBytes, section.counter)
		tryWrite(counterBytes)
	}
	for i := 0; i < len(s.hallOfFame); i++ {
		section := s.hallOfFame[i]
		unusedVal := 0
		if section.empty {
			unusedVal = 0xFFFFFFFF
		}
		tryWrite(section.data)
		checksumBytes := make([]byte, 2)
		binary.LittleEndian.PutUint16(checksumBytes, section.checksum)
		tryWrite(checksumBytes)
		unusedBytes := make([]byte, 2)
		binary.LittleEndian.PutUint16(unusedBytes, uint16(unusedVal))
		tryWrite(unusedBytes)
		securityBytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(securityBytes, section.security)
		tryWrite(securityBytes)
		unusedBytes = make([]byte, 4)
		binary.LittleEndian.PutUint32(unusedBytes, uint32(unusedVal))
		tryWrite(unusedBytes)
	}
	tryWrite(s.trainerHill)
	tryWrite(s.recordedBattle)

	return err
}

func (s *GameSaveSection) calculateChecksum() uint16 {
	var sum uint32
	for i := 0; i < gameSectionDataSize/4; i++ {
		offset := i * 4
		sum += binary.LittleEndian.Uint32(s.data[offset : offset+4])
	}
	return uint16(sum>>16) + uint16(sum&0xFFFF)
}

func (s *HallOfFameSection) calculateChecksum() uint16 {
	var sum uint32
	for i := 0; i < gameSectionDataSize/4; i++ {
		offset := i * 4
		sum += binary.LittleEndian.Uint32(s.data[offset : offset+4])
	}
	return uint16(sum>>16) + uint16(sum&0xFFFF)
}

func (s *GameSaveSection) checkCorruption() error {
	checksum := s.calculateChecksum()
	if checksum != s.checksum {
		return fmt.Errorf("Game save section %d has incorrect checksum %#x. Expected %#x", s.id, s.checksum, checksum)
	}
	if s.security != securityValue {
		return fmt.Errorf("Game save section %d has incorrect security value %#x. Expected %#x", s.id, s.security, securityValue)
	}
	return nil
}

func (s *HallOfFameSection) checkCorruption(id int) error {
	if s.empty {
		return nil
	}
	checksum := s.calculateChecksum()
	if checksum != s.checksum {
		return fmt.Errorf("Hall of Fame save section %d has incorrect checksum %#x. Expected %#x", id, s.checksum, checksum)
	}
	if s.security != securityValue {
		return fmt.Errorf("Hall of Fame save section %d has incorrect security value %#x. Expected %#x", id, s.security, securityValue)
	}
	return nil
}

func (s *SaveData) getGameSaveSection(id uint) *GameSaveSection {
	gameSaveSection := s.getLatestGameSaveSection()
	for i := 0; i < numGameSaveSections; i++ {
		if gameSaveSection[i].id == uint16(id) {
			return &gameSaveSection[i]
		}
	}
	return nil
}

func (s *SaveData) readU8(address uint) uint8 {
	sectionID := address / gameSectionDataSize
	section := s.getGameSaveSection(sectionID)
	offset := address % gameSectionDataSize
	return section.data[offset]
}

func (s *SaveData) readS8(address uint) int8 {
	return int8(s.readU8(address))
}

func (s *SaveData) readU16(address uint) uint16 {
	lo := s.readU8(address)
	hi := s.readU8(address + 1)
	return (uint16(hi) << 8) | uint16(lo)
}

func (s *SaveData) readS16(address uint) int16 {
	return int16(s.readU16(address))
}

func (s *SaveData) readU32(address uint) uint32 {
	lo := s.readU16(address)
	hi := s.readU16(address + 2)
	return (uint32(hi) << 16) | uint32(lo)
}

func (s *SaveData) readS32(address uint) int32 {
	return int32(s.readU32(address))
}

func (s *SaveData) writeU8(value uint8, address uint) {
	sectionID := address / gameSectionDataSize
	section := s.getGameSaveSection(sectionID)
	offset := address % gameSectionDataSize
	section.data[offset] = value
}

func (s *SaveData) writeS8(value int8, address uint) {
	s.writeU8(uint8(value), address)
}

func (s *SaveData) writeU16(value uint16, address uint) {
	lo := uint8(value & 0xFF)
	hi := uint8(value >> 8)
	s.writeU8(lo, address)
	s.writeU8(hi, address+1)
}

func (s *SaveData) writeS16(value int16, address uint) {
	s.writeU16(uint16(value), address)
}

func (s *SaveData) writeU32(value uint32, address uint) {
	lo := uint16(value & 0xFFFF)
	hi := uint16(value >> 16)
	s.writeU16(lo, address)
	s.writeU16(hi, address+2)
}

func (s *SaveData) writeS32(value int32, address uint) {
	s.writeU32(uint32(value), address)
}

func (s *SaveData) readU8Slice(outBuffer []uint8, address uint) {
	for i := 0; i < len(outBuffer); i++ {
		v := s.readU8(address + uint(i))
		outBuffer[i] = v
	}
}

func (s *SaveData) readU16Slice(outBuffer []uint16, address uint) {
	for i := 0; i < len(outBuffer); i++ {
		outBuffer[i] = s.readU16(address + 2*uint(i))
	}
}

func (s *SaveData) readU32Slice(outBuffer []uint32, address uint) {
	for i := 0; i < len(outBuffer); i++ {
		outBuffer[i] = s.readU32(address + 4*uint(i))
	}
}

func (s *SaveData) writeU8Slice(buffer []uint8, address uint) {
	for i := 0; i < len(buffer); i++ {
		s.writeU8(buffer[i], address+uint(i))
	}
}

func (s *SaveData) writeU16Slice(buffer []uint16, address uint) {
	for i := 0; i < len(buffer); i++ {
		s.writeU16(buffer[i], address+2*uint(i))
	}
}

func (s *SaveData) writeU32Slice(buffer []uint32, address uint) {
	for i := 0; i < len(buffer); i++ {
		s.writeU32(buffer[i], address+4*uint(i))
	}
}

func (s *SaveData) readBit(address uint, bit uint) bool {
	v := s.readU8(address + bit/8)
	mask := uint8(1 << (bit % 8))
	return (v & mask) != 0
}

func (s *SaveData) readBitsU8(address uint, bit, width uint) uint8 {
	return uint8(s.readBitsU32(address, bit, width))
}

func (s *SaveData) readBitsU16(address uint, bit, width uint) uint16 {
	return uint16(s.readBitsU32(address, bit, width))
}

func (s *SaveData) readBitsU32(address uint, bit, width uint) uint32 {
	var value uint32
	for i := 0; uint(i) < width; i++ {
		if s.readBit(address+bit/8, bit%8) {
			value |= (1 << i)
		}
		bit++
	}
	return value
}

func (s *SaveData) writeBit(set bool, address uint, bit uint) {
	v := s.readU8(address + bit/8)
	var mask uint8 = 1 << (bit % 8)
	if set {
		v |= mask
	} else {
		v &^= mask
	}
	s.writeU8(v, address+bit/8)
}

func (s *SaveData) writeBitsU8(value uint8, address uint, bit, width uint) {
	s.writeBitsU32(uint32(value), address, bit, width)
}

func (s *SaveData) writeBitsU16(value uint16, address uint, bit, width uint) {
	s.writeBitsU32(uint32(value), address, bit, width)
}

func (s *SaveData) writeBitsU32(value uint32, address uint, bit, width uint) {
	for width > 0 {
		set := (value & 1) != 0
		s.writeBit(set, address+bit/8, bit%8)
		value >>= 1
		width--
		bit++
	}
}
