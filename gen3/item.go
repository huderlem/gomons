package gen3

import (
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
	offset            uint
	capacity          uint
	encryptedQuantity bool
}

var itemPockets = map[itemPocket]itemPocketInfo{
	itemPocketPc:       {0x1418, 50, false},
	itemPocketItems:    {0x14E0, 30, true},
	itemPocketKeyItems: {0x1558, 30, true},
	itemPocketBalls:    {0x15D0, 16, true},
	itemPocketTmHm:     {0x1610, 64, true},
	itemPocketBerries:  {0x1710, 46, true},
}

// GetRegisteredItem gets the player's currently-registered item.
func (s *SaveData) GetRegisteredItem() uint16 {
	return s.readU16(0x1416)
}

// SetRegisteredItem sets the player's currently-registered item.
func (s *SaveData) SetRegisteredItem(item uint16) {
	s.writeU16(item, 0x1416)
}

// GetNumPcItems gets the number of items stored in the PC.
func (s *SaveData) GetNumPcItems() int {
	return s.getNumItemsInPocket(itemPocketPc)
}

// GetPcItem gets the item id and quantity at the given PC slot index.
func (s *SaveData) GetPcItem(index int) (uint16, uint16) {
	return s.getItemInPocket(itemPocketPc, index)
}

// SetPcItem sets the item id and quantity at the given PC slot index.
func (s *SaveData) SetPcItem(itemID, quantity uint16, index int) {
	s.setItemInPocket(itemPocketPc, itemID, quantity, index)
}

// AddPcItem add the given item and quantity to the end of the player's PC item storage.
func (s *SaveData) AddPcItem(itemID, quantity uint16) error {
	return s.addItemToPocket(itemPocketPc, itemID, quantity)
}

// RemovePcItem removes the item at the given index in the player's PC item storage, and
// then shifts any following items up in the list.
func (s *SaveData) RemovePcItem(index int) {
	s.removeItemFromPocket(itemPocketPc, index)
}

// GetNumItems gets the number of items stored in main items pocket.
func (s *SaveData) GetNumItems() int {
	return s.getNumItemsInPocket(itemPocketItems)
}

// GetItem gets the item id and quantity at the given index in the main items pocket.
func (s *SaveData) GetItem(index int) (uint16, uint16) {
	return s.getItemInPocket(itemPocketItems, index)
}

// SetItem sets the item id and quantity at the given main items pocket index.
func (s *SaveData) SetItem(itemID, quantity uint16, index int) {
	s.setItemInPocket(itemPocketItems, itemID, quantity, index)
}

// AddItem add the given item and quantity to the end of the player's main items pocket.
func (s *SaveData) AddItem(itemID, quantity uint16) error {
	return s.addItemToPocket(itemPocketItems, itemID, quantity)
}

// RemoveItem removes the item at the given index in the player's main items pocket, and
// then shifts any following items up in the list.
func (s *SaveData) RemoveItem(index int) {
	s.removeItemFromPocket(itemPocketItems, index)
}

// GetNumKeyItems gets the number of items stored in key items pocket.
func (s *SaveData) GetNumKeyItems() int {
	return s.getNumItemsInPocket(itemPocketKeyItems)
}

// GetKeyItem gets the item id at the given index in the key items pocket.
func (s *SaveData) GetKeyItem(index int) uint16 {
	itemID, _ := s.getItemInPocket(itemPocketKeyItems, index)
	return itemID
}

// SetKeyItem sets the item id at the given key items pocket index.
func (s *SaveData) SetKeyItem(itemID uint16, index int) {
	s.setItemInPocket(itemPocketKeyItems, itemID, 1, index)
}

// AddKeyItem add the given item and quantity to the end of the player's key items pocket.
func (s *SaveData) AddKeyItem(itemID uint16) error {
	return s.addItemToPocket(itemPocketKeyItems, itemID, 1)
}

// RemoveKeyItem removes the item at the given index in the player's key items pocket, and
// then shifts any following items up in the list.
func (s *SaveData) RemoveKeyItem(index int) {
	s.removeItemFromPocket(itemPocketKeyItems, index)
}

// GetNumBallItems gets the number of items stored in pokeball items pocket.
func (s *SaveData) GetNumBallItems() int {
	return s.getNumItemsInPocket(itemPocketBalls)
}

// GetBallItem gets the item id and quantity at the given index in the pokeball items pocket.
func (s *SaveData) GetBallItem(index int) (uint16, uint16) {
	return s.getItemInPocket(itemPocketBalls, index)
}

// SetBallItem sets the item id and quantity at the given pokeball items pocket index.
func (s *SaveData) SetBallItem(itemID, quantity uint16, index int) {
	s.setItemInPocket(itemPocketBalls, itemID, quantity, index)
}

// AddBallItem add the given item and quantity to the end of the player's pokeball items pocket.
func (s *SaveData) AddBallItem(itemID, quantity uint16) error {
	return s.addItemToPocket(itemPocketBalls, itemID, quantity)
}

// RemoveBallItem removes the item at the given index in the player's pokeball items pocket, and
// then shifts any following items up in the list.
func (s *SaveData) RemoveBallItem(index int) {
	s.removeItemFromPocket(itemPocketBalls, index)
}

// GetNumTmHmItems gets the number of items stored in TM/HM items pocket.
func (s *SaveData) GetNumTmHmItems() int {
	return s.getNumItemsInPocket(itemPocketTmHm)
}

// GetTmHmItem gets the item id and quantity at the given index in the TM/HM items pocket.
func (s *SaveData) GetTmHmItem(index int) (uint16, uint16) {
	return s.getItemInPocket(itemPocketTmHm, index)
}

// SetTmHmItem sets the item id and quantity at the given TM/HM items pocket index.
func (s *SaveData) SetTmHmItem(itemID, quantity uint16, index int) {
	s.setItemInPocket(itemPocketTmHm, itemID, quantity, index)
}

// AddTmHmItem add the given item and quantity to the end of the player's TM/HM items pocket.
func (s *SaveData) AddTmHmItem(itemID, quantity uint16) error {
	return s.addItemToPocket(itemPocketTmHm, itemID, quantity)
}

// RemoveTmHmItem removes the item at the given index in the player's TM/HM items pocket, and
// then shifts any following items up in the list.
func (s *SaveData) RemoveTmHmItem(index int) {
	s.removeItemFromPocket(itemPocketTmHm, index)
}

// GetNumBerryItems gets the number of items stored in berry items pocket.
func (s *SaveData) GetNumBerryItems() int {
	return s.getNumItemsInPocket(itemPocketBerries)
}

// GetBerryItem gets the item id and quantity at the given index in the berry items pocket.
func (s *SaveData) GetBerryItem(index int) (uint16, uint16) {
	return s.getItemInPocket(itemPocketBerries, index)
}

// SetBerryItem sets the item id and quantity at the given berry items pocket index.
func (s *SaveData) SetBerryItem(itemID, quantity uint16, index int) {
	s.setItemInPocket(itemPocketBerries, itemID, quantity, index)
}

// AddBerryItem add the given item and quantity to the end of the player's berry items pocket.
func (s *SaveData) AddBerryItem(itemID, quantity uint16) error {
	return s.addItemToPocket(itemPocketBerries, itemID, quantity)
}

// RemoveBerryItem removes the item at the given index in the player's berry items pocket, and
// then shifts any following items up in the list.
func (s *SaveData) RemoveBerryItem(index int) {
	s.removeItemFromPocket(itemPocketBerries, index)
}

func (s *SaveData) getNumItemsInPocket(pocket itemPocket) int {
	var count uint
	for count < itemPockets[pocket].capacity {
		offset := itemPockets[pocket].offset + count*4
		itemID := s.readU16(offset)
		if itemID == 0 {
			break
		}
		count++
	}
	return int(count)
}

func (s *SaveData) getItemInPocket(pocket itemPocket, index int) (uint16, uint16) {
	offset := itemPockets[pocket].offset + uint(index)*4
	itemID := s.readU16(offset)
	quantity := s.readU16(offset + 2)
	if itemPockets[pocket].encryptedQuantity {
		key := s.GetEncryptionKey()
		quantity ^= uint16(key)
	}
	return itemID, quantity
}

func (s *SaveData) setItemInPocket(pocket itemPocket, itemID, quantity uint16, index int) {
	offset := itemPockets[pocket].offset + uint(index)*4
	s.writeU16(itemID, offset)
	if itemPockets[pocket].encryptedQuantity {
		key := s.GetEncryptionKey()
		quantity ^= uint16(key)
	}
	s.writeU16(quantity, offset+2)
}

func (s *SaveData) addItemToPocket(pocket itemPocket, itemID, quantity uint16) error {
	count := s.getNumItemsInPocket(pocket)
	if uint(count) == itemPockets[pocket].capacity {
		return fmt.Errorf("Cannot add item because the pocket is already full")
	}
	offset := itemPockets[pocket].offset + uint(count)*4
	s.writeU16(itemID, offset)
	if itemPockets[pocket].encryptedQuantity {
		key := s.GetEncryptionKey()
		quantity ^= uint16(key)
	}
	s.writeU16(quantity, offset+2)
	return nil
}

func (s *SaveData) removeItemFromPocket(pocket itemPocket, index int) {
	count := s.getNumItemsInPocket(pocket)
	if index >= int(count) {
		return
	}
	for index < int(count)-1 {
		offset := itemPockets[pocket].offset + uint(index)*4
		itemID := s.readU16(offset + 4)
		quantity := s.readU16(offset + 6)
		s.writeU16(itemID, offset)
		s.writeU16(quantity, offset+2)
		index++
	}
	offset := itemPockets[pocket].offset + uint(index)*4
	s.writeU16(0, offset)
	s.writeU16(0, offset+2)
}
