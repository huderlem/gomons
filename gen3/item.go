package gen3

import (
	"encoding/binary"
	"fmt"
)

type itemPocket byte

const (
	itemPocketPc itemPocket = iota
	itemPocketItems
	itemPocketKeyItems
	itemPocketBalls
	itemPocketTmHm
	itemPocketBerries
)

type itemPocketInfo struct {
	offset            int
	capacity          int
	encryptedQuantity bool
}

var itemPockets = map[itemPocket]itemPocketInfo{
	itemPocketPc:       itemPocketInfo{0x498, 50, false},
	itemPocketItems:    itemPocketInfo{0x560, 30, true},
	itemPocketKeyItems: itemPocketInfo{0x5D8, 30, true},
	itemPocketBalls:    itemPocketInfo{0x650, 16, true},
	itemPocketTmHm:     itemPocketInfo{0x690, 64, true},
	itemPocketBerries:  itemPocketInfo{0x790, 46, true},
}

// GetRegisteredItem gets the player's currently-registered item.
func (s *SaveData) GetRegisteredItem() (uint16, error) {
	section, err := s.getGameSaveSection(1)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint16(section.data[0x496:0x498]), nil
}

// SetRegisteredItem sets the player's currently-registered item.
func (s *SaveData) SetRegisteredItem(item uint16) error {
	section, err := s.getGameSaveSection(1)
	if err != nil {
		return err
	}
	binary.LittleEndian.PutUint16(section.data[0x496:0x498], item)
	return nil
}

// GetNumPcItems gets the number of items stored in the PC.
func (s *SaveData) GetNumPcItems() (int, error) {
	return s.getNumItemsInPocket(itemPocketPc)
}

// GetPcItem gets the item id and quantity at the given PC slot index.
func (s *SaveData) GetPcItem(index int) (uint16, uint16, error) {
	return s.getItemInPocket(itemPocketPc, index)
}

// SetPcItem sets the item id and quantity at the given PC slot index.
func (s *SaveData) SetPcItem(itemID, quantity uint16, index int) error {
	return s.setItemInPocket(itemPocketPc, itemID, quantity, index)
}

// AddPcItem add the given item and quantity to the end of the player's PC item storage.
func (s *SaveData) AddPcItem(itemID, quantity uint16) error {
	return s.addItemToPocket(itemPocketPc, itemID, quantity)
}

// RemovePcItem removes the item at the given index in the player's PC item storage, and
// then shifts any following items up in the list.
func (s *SaveData) RemovePcItem(index int) error {
	return s.removeItemFromPocket(itemPocketPc, index)
}

// GetNumItems gets the number of items stored in main items pocket.
func (s *SaveData) GetNumItems() (int, error) {
	return s.getNumItemsInPocket(itemPocketItems)
}

// GetItem gets the item id and quantity at the given index in the main items pocket.
func (s *SaveData) GetItem(index int) (uint16, uint16, error) {
	return s.getItemInPocket(itemPocketItems, index)
}

// SetItem sets the item id and quantity at the given main items pocket index.
func (s *SaveData) SetItem(itemID, quantity uint16, index int) error {
	return s.setItemInPocket(itemPocketItems, itemID, quantity, index)
}

// AddItem add the given item and quantity to the end of the player's main items pocket.
func (s *SaveData) AddItem(itemID, quantity uint16) error {
	return s.addItemToPocket(itemPocketItems, itemID, quantity)
}

// RemoveItem removes the item at the given index in the player's main items pocket, and
// then shifts any following items up in the list.
func (s *SaveData) RemoveItem(index int) error {
	return s.removeItemFromPocket(itemPocketItems, index)
}

// GetNumKeyItems gets the number of items stored in key items pocket.
func (s *SaveData) GetNumKeyItems() (int, error) {
	return s.getNumItemsInPocket(itemPocketKeyItems)
}

// GetKeyItem gets the item id at the given index in the key items pocket.
func (s *SaveData) GetKeyItem(index int) (uint16, error) {
	itemID, _, err := s.getItemInPocket(itemPocketKeyItems, index)
	return itemID, err
}

// SetKeyItem sets the item id at the given key items pocket index.
func (s *SaveData) SetKeyItem(itemID uint16, index int) error {
	return s.setItemInPocket(itemPocketKeyItems, itemID, 1, index)
}

// AddKeyItem add the given item and quantity to the end of the player's key items pocket.
func (s *SaveData) AddKeyItem(itemID uint16) error {
	return s.addItemToPocket(itemPocketKeyItems, itemID, 1)
}

// RemoveKeyItem removes the item at the given index in the player's key items pocket, and
// then shifts any following items up in the list.
func (s *SaveData) RemoveKeyItem(index int) error {
	return s.removeItemFromPocket(itemPocketKeyItems, index)
}

// GetNumBallItems gets the number of items stored in pokeball items pocket.
func (s *SaveData) GetNumBallItems() (int, error) {
	return s.getNumItemsInPocket(itemPocketBalls)
}

// GetBallItem gets the item id and quantity at the given index in the pokeball items pocket.
func (s *SaveData) GetBallItem(index int) (uint16, uint16, error) {
	return s.getItemInPocket(itemPocketBalls, index)
}

// SetBallItem sets the item id and quantity at the given pokeball items pocket index.
func (s *SaveData) SetBallItem(itemID, quantity uint16, index int) error {
	return s.setItemInPocket(itemPocketBalls, itemID, quantity, index)
}

// AddBallItem add the given item and quantity to the end of the player's pokeball items pocket.
func (s *SaveData) AddBallItem(itemID, quantity uint16) error {
	return s.addItemToPocket(itemPocketBalls, itemID, quantity)
}

// RemoveBallItem removes the item at the given index in the player's pokeball items pocket, and
// then shifts any following items up in the list.
func (s *SaveData) RemoveBallItem(index int) error {
	return s.removeItemFromPocket(itemPocketBalls, index)
}

// GetNumTmHmItems gets the number of items stored in TM/HM items pocket.
func (s *SaveData) GetNumTmHmItems() (int, error) {
	return s.getNumItemsInPocket(itemPocketTmHm)
}

// GetTmHmItem gets the item id and quantity at the given index in the TM/HM items pocket.
func (s *SaveData) GetTmHmItem(index int) (uint16, uint16, error) {
	return s.getItemInPocket(itemPocketTmHm, index)
}

// SetTmHmItem sets the item id and quantity at the given TM/HM items pocket index.
func (s *SaveData) SetTmHmItem(itemID, quantity uint16, index int) error {
	return s.setItemInPocket(itemPocketTmHm, itemID, quantity, index)
}

// AddTmHmItem add the given item and quantity to the end of the player's TM/HM items pocket.
func (s *SaveData) AddTmHmItem(itemID, quantity uint16) error {
	return s.addItemToPocket(itemPocketTmHm, itemID, quantity)
}

// RemoveTmHmItem removes the item at the given index in the player's TM/HM items pocket, and
// then shifts any following items up in the list.
func (s *SaveData) RemoveTmHmItem(index int) error {
	return s.removeItemFromPocket(itemPocketTmHm, index)
}

// GetNumBerryItems gets the number of items stored in berry items pocket.
func (s *SaveData) GetNumBerryItems() (int, error) {
	return s.getNumItemsInPocket(itemPocketBerries)
}

// GetBerryItem gets the item id and quantity at the given index in the berry items pocket.
func (s *SaveData) GetBerryItem(index int) (uint16, uint16, error) {
	return s.getItemInPocket(itemPocketBerries, index)
}

// SetBerryItem sets the item id and quantity at the given berry items pocket index.
func (s *SaveData) SetBerryItem(itemID, quantity uint16, index int) error {
	return s.setItemInPocket(itemPocketBerries, itemID, quantity, index)
}

// AddBerryItem add the given item and quantity to the end of the player's berry items pocket.
func (s *SaveData) AddBerryItem(itemID, quantity uint16) error {
	return s.addItemToPocket(itemPocketBerries, itemID, quantity)
}

// RemoveBerryItem removes the item at the given index in the player's berry items pocket, and
// then shifts any following items up in the list.
func (s *SaveData) RemoveBerryItem(index int) error {
	return s.removeItemFromPocket(itemPocketBerries, index)
}

func (s *SaveData) getNumItemsInPocket(pocket itemPocket) (int, error) {
	section, err := s.getGameSaveSection(1)
	if err != nil {
		return 0, err
	}
	count := 0
	for count < itemPockets[pocket].capacity {
		offset := itemPockets[pocket].offset + count*4
		itemID := binary.LittleEndian.Uint16(section.data[offset : offset+2])
		if itemID == 0 {
			break
		}
		count++
	}
	return count, nil
}

func (s *SaveData) getItemInPocket(pocket itemPocket, index int) (uint16, uint16, error) {
	section, err := s.getGameSaveSection(1)
	if err != nil {
		return 0, 0, err
	}
	offset := itemPockets[pocket].offset + index*4
	itemID := binary.LittleEndian.Uint16(section.data[offset : offset+2])
	quantity := binary.LittleEndian.Uint16(section.data[offset+2 : offset+4])
	if itemPockets[pocket].encryptedQuantity {
		key, err := s.GetEncryptionKey()
		if err != nil {
			return 0, 0, err
		}
		quantity ^= uint16(key)
	}
	return itemID, quantity, nil
}

func (s *SaveData) setItemInPocket(pocket itemPocket, itemID, quantity uint16, index int) error {
	section, err := s.getGameSaveSection(1)
	if err != nil {
		return err
	}
	offset := itemPockets[pocket].offset + index*4
	binary.LittleEndian.PutUint16(section.data[offset:offset+2], itemID)
	if itemPockets[pocket].encryptedQuantity {
		key, err := s.GetEncryptionKey()
		if err != nil {
			return err
		}
		quantity ^= uint16(key)
	}
	binary.LittleEndian.PutUint16(section.data[offset+2:offset+4], quantity)
	return nil
}

func (s *SaveData) addItemToPocket(pocket itemPocket, itemID, quantity uint16) error {
	section, err := s.getGameSaveSection(1)
	if err != nil {
		return err
	}
	count, err := s.getNumItemsInPocket(pocket)
	if err != nil {
		return err
	}
	if count == itemPockets[pocket].capacity {
		return fmt.Errorf("Cannot add item because the pocket is already full")
	}
	offset := itemPockets[pocket].offset + (count * 4)
	binary.LittleEndian.PutUint16(section.data[offset:offset+2], itemID)
	if itemPockets[pocket].encryptedQuantity {
		key, err := s.GetEncryptionKey()
		if err != nil {
			return err
		}
		quantity ^= uint16(key)
	}
	binary.LittleEndian.PutUint16(section.data[offset+2:offset+4], quantity)
	return nil
}

func (s *SaveData) removeItemFromPocket(pocket itemPocket, index int) error {
	section, err := s.getGameSaveSection(1)
	if err != nil {
		return err
	}
	count, err := s.getNumItemsInPocket(pocket)
	if err != nil {
		return err
	}
	if index >= count {
		return nil
	}
	for index < count-1 {
		offset := itemPockets[pocket].offset + (index * 4)
		itemID := binary.LittleEndian.Uint16(section.data[offset+4 : offset+6])
		quantity := binary.LittleEndian.Uint16(section.data[offset+6 : offset+8])
		binary.LittleEndian.PutUint16(section.data[offset:offset+2], itemID)
		binary.LittleEndian.PutUint16(section.data[offset+2:offset+4], quantity)
		index++
	}
	offset := itemPockets[pocket].offset + (index * 4)
	binary.LittleEndian.PutUint16(section.data[offset:offset+2], 0)
	binary.LittleEndian.PutUint16(section.data[offset+2:offset+4], 0)
	return nil
}
