package gen3

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/huderlem/gomons/util"
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
}

// HallOfFameSection represents one of the logical hall of fame sections of the save data structure.
type HallOfFameSection struct {
	data     []byte
	checksum uint16
	security uint32
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
	counterB := s.gameSaveB[len(s.gameSaveB)-1].counter
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
		s.hallOfFame[i].checksum = s.hallOfFame[i].calculateChecksum()
	}
}

func (s *SaveData) Write(w io.Writer) error {
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
		tryWrite(section.data)
		checksumBytes := make([]byte, 2)
		binary.LittleEndian.PutUint16(checksumBytes, section.checksum)
		tryWrite(checksumBytes)
		unusedBytes := make([]byte, 2)
		binary.LittleEndian.PutUint16(unusedBytes, 0)
		tryWrite(unusedBytes)
		securityBytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(securityBytes, section.security)
		tryWrite(securityBytes)
		unusedBytes = make([]byte, 4)
		binary.LittleEndian.PutUint32(unusedBytes, 0)
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
	checksum := s.calculateChecksum()
	if checksum != s.checksum {
		return fmt.Errorf("Hall of Fame save section %d has incorrect checksum %#x. Expected %#x", id, s.checksum, checksum)
	}
	if s.security != securityValue {
		return fmt.Errorf("Hall of Fame save section %d has incorrect security value %#x. Expected %#x", id, s.security, securityValue)
	}
	return nil
}

func (s *SaveData) getGameSaveSection(id int) (*GameSaveSection, error) {
	if id < 0 || id >= numGameSaveSections {
		return nil, fmt.Errorf("Invalid game save section id %d", id)
	}
	gameSaveSection := s.getLatestGameSaveSection()
	for i := 0; i < numGameSaveSections; i++ {
		if gameSaveSection[i].id == uint16(id) {
			return &gameSaveSection[i], nil
		}
	}
	return nil, fmt.Errorf("Missing game save section id %d", id)
}

// GetPlayerGender gets the player's gender.
func (s *SaveData) GetPlayerGender() (PlayerGender, error) {
	section, err := s.getGameSaveSection(0)
	if err != nil {
		return Invalid, err
	}
	gender := PlayerGender(section.data[8])
	if gender != Male && gender != Female {
		return Invalid, fmt.Errorf("Save data has invalid player gender %#x", gender)
	}
	return gender, nil
}

// SetPlayerGender sets the player's gender.
func (s *SaveData) SetPlayerGender(gender PlayerGender) error {
	section, err := s.getGameSaveSection(0)
	if err != nil {
		return err
	}
	section.data[8] = byte(gender)
	return nil
}

// GetPlayerName gets the player's OT name.
func (s *SaveData) GetPlayerName() (string, error) {
	section, err := s.getGameSaveSection(0)
	if err != nil {
		return "", err
	}
	var sb strings.Builder
	for _, b := range section.data[0:7] {
		if b == endOfString {
			break
		}
		if letter, ok := reverseCharmap[b]; ok {
			sb.WriteRune(letter)
		} else {
			sb.WriteByte(b)
		}
	}
	return sb.String(), nil
}

// SetPlayerName gets the player's OT name.
func (s *SaveData) SetPlayerName(name string) error {
	if len(name) < 1 || len(name) > 7 {
		return fmt.Errorf("Player name must be between 1 and 7 characters long")
	}
	buffer := make([]byte, 7)
	pos := 0
	for _, letter := range name {
		if b, ok := charmap[letter]; ok {
			buffer[pos] = b
			pos++
		} else {
			return fmt.Errorf("Cannot set player name to %s because the character '%c' is unsupported", name, letter)
		}
	}
	// Pad with null-terminating characters
	for pos < 7 {
		buffer[pos] = endOfString
		pos++
	}
	section, err := s.getGameSaveSection(0)
	if err != nil {
		return err
	}
	copy(section.data[0:7], buffer)
	return nil
}

// GetPlayerTrainerID gets the player's raw trainer id. The trainer id is composed of two parts--public and secret.
// The public id is viewable in-game, while the secret id is not. The trainer id is always used by the game
// engine as its full 4-byte value.
func (s *SaveData) GetPlayerTrainerID() (uint32, error) {
	section, err := s.getGameSaveSection(0)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(section.data[0xA:0xE]), nil
}

// GetPlayerPublicID gets the player's public trainer id part. This is the id viewable in-game.
func (s *SaveData) GetPlayerPublicID() (uint16, error) {
	otid, err := s.GetPlayerTrainerID()
	if err != nil {
		return 0, err
	}
	return uint16(otid), nil
}

// GetPlayerSecretID gets the player's secret trainer id part. This id is not viewable in-game, but is
// still a part of the player's trainer identity. It is used, for example, when determining a mon's shininess.
func (s *SaveData) GetPlayerSecretID() (uint16, error) {
	otid, err := s.GetPlayerTrainerID()
	if err != nil {
		return 0, err
	}
	return uint16(otid >> 16), nil
}

// SetPlayerTrainerID sets the player's raw trainer id. The trainer id is composed of two parts--public and secret.
// The public id is viewable in-game, while the secret id is not. The trainer id is always used by the game
// engine as its full 4-byte value.
func (s *SaveData) SetPlayerTrainerID(id uint32) error {
	section, err := s.getGameSaveSection(0)
	if err != nil {
		return err
	}
	binary.LittleEndian.PutUint32(section.data[0xA:0xE], id)
	return nil
}

// SetPlayerPublicID sets the player's public trainer id part. This is the id viewable in-game.
func (s *SaveData) SetPlayerPublicID(id uint16) error {
	section, err := s.getGameSaveSection(0)
	if err != nil {
		return err
	}
	binary.LittleEndian.PutUint16(section.data[0xA:0xC], id)
	return nil
}

// SetPlayerSecretID gets the player's secret trainer id part. This id is not viewable in-game, but is
// still a part of the player's trainer identity. It is used, for example, when determining a mon's shininess.
func (s *SaveData) SetPlayerSecretID(id uint16) error {
	section, err := s.getGameSaveSection(0)
	if err != nil {
		return err
	}
	binary.LittleEndian.PutUint16(section.data[0xC:0xE], id)
	return nil
}

// GetPlayTime gets the play time hours, seconds, minutes, and vblanks. One VBlank is roughly 1/60 of a second.
func (s *SaveData) GetPlayTime() (uint16, uint8, uint8, uint8, error) {
	section, err := s.getGameSaveSection(0)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	hours := binary.LittleEndian.Uint16(section.data[0xE:0x10])
	minutes := section.data[0x10]
	seconds := section.data[0x11]
	vblanks := section.data[0x12]
	return hours, minutes, seconds, vblanks, nil
}

// SetPlayTime gets the play time hours, seconds, minutes, and vblanks. One VBlank is roughly 1/60 of a second.
func (s *SaveData) SetPlayTime(hours uint16, minutes uint8, seconds uint8, vblanks uint8) error {
	section, err := s.getGameSaveSection(0)
	if err != nil {
		return err
	}
	if hours > 999 {
		return fmt.Errorf("Hours must be in range 0-999, but got %d", minutes)
	}
	if minutes > 59 {
		return fmt.Errorf("Minutes must be in range 0-59, but got %d", minutes)
	}
	if seconds > 59 {
		return fmt.Errorf("Seconds must be in range 0-59, but got %d", seconds)
	}
	if vblanks > 59 {
		return fmt.Errorf("VBlanks must be in range 0-59, but got %d", vblanks)
	}
	binary.LittleEndian.PutUint16(section.data[0xE:0x10], hours)
	section.data[0x10] = minutes
	section.data[0x11] = seconds
	section.data[0x12] = vblanks
	return nil
}

// GetRegionMapZoomedIn gets whether or not the region map is zoomed in.
func (s *SaveData) GetRegionMapZoomedIn() (bool, error) {
	section, err := s.getGameSaveSection(0)
	if err != nil {
		return false, err
	}
	return (section.data[0x15] & 0x8) != 0, nil
}

// SetRegionMapZoomedIn sets whether or not the region map is zoomed in.
func (s *SaveData) SetRegionMapZoomedIn(isZoomedIn bool) error {
	section, err := s.getGameSaveSection(0)
	if err != nil {
		return err
	}
	section.data[0x15] = (section.data[0x15] &^ 0x8) | (util.BoolToByte(isZoomedIn) << 3)
	return nil
}

// GetEncryptionKey gets the encryption key used to decrypt various save data.
func (s *SaveData) GetEncryptionKey() (uint32, error) {
	section, err := s.getGameSaveSection(0)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(section.data[0xAC:0xB0]), nil
}

// GetMoney gets the player's current money amount.
func (s *SaveData) GetMoney() (uint32, error) {
	section, err := s.getGameSaveSection(1)
	if err != nil {
		return 0, err
	}
	encryptedMoney := binary.LittleEndian.Uint32(section.data[0x490:0x494])
	key, err := s.GetEncryptionKey()
	if err != nil {
		return 0, err
	}
	return encryptedMoney ^ key, nil
}

// SetMoney sets the player's current money amount.
func (s *SaveData) SetMoney(money uint32) error {
	section, err := s.getGameSaveSection(1)
	if err != nil {
		return err
	}
	key, err := s.GetEncryptionKey()
	if err != nil {
		return err
	}
	if money > 999999 {
		money = 999999
	}
	encryptedMoney := money ^ key
	binary.LittleEndian.PutUint32(section.data[0x490:0x494], encryptedMoney)
	return nil
}

// GetPlayerCoordinates gets the player's current x/y coordinates.
func (s *SaveData) GetPlayerCoordinates() (int16, int16, error) {
	section, err := s.getGameSaveSection(1)
	if err != nil {
		return 0, 0, err
	}
	x := int16(binary.LittleEndian.Uint16(section.data[0x0:0x2]))
	y := int16(binary.LittleEndian.Uint16(section.data[0x2:0x4]))
	return x, y, nil
}

// SetPlayerCoordinates sets the player's current x/y coordinates.
func (s *SaveData) SetPlayerCoordinates(x, y int16) error {
	section, err := s.getGameSaveSection(1)
	if err != nil {
		return err
	}
	binary.LittleEndian.PutUint16(section.data[0x0:0x2], uint16(x))
	binary.LittleEndian.PutUint16(section.data[0x2:0x4], uint16(y))
	return nil
}
