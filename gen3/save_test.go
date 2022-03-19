package gen3

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadSaveFileFromBytes(t *testing.T) {
	t.Parallel()

	t.Run("loading a fresh save with only save game B", func(t *testing.T) {
		t.Parallel()
		sd, err := LoadSaveFileFromBytes(eFreshSaveB)
		require.Nil(t, err, "should not return an error")
		assert.Equal(t, &sd.gameSaveB, sd.activeGameSave)
		for i, s := range sd.gameSaveA {
			assert.Truef(t, s.empty, "section %d of save A is not empty", i)
		}
		for i, s := range sd.gameSaveB {
			assert.Falsef(t, s.empty, "section %d of save B is empty", i)
			assert.Equalf(t, s.counter, uint32(1), "section %d of save B has wrong counter", i)
		}
		for i, s := range sd.hallOfFame {
			assert.Truef(t, s.empty, "section %d of HoF is not empty", i)
		}
		assert.Equal(t, sd.GetPlayerGender(), Male)
		assert.Equal(t, sd.GetPlayerName(), "ACID")
	})

	t.Run("loading a fresh save with two save slots", func(t *testing.T) {
		t.Parallel()
		sd, err := LoadSaveFileFromBytes(eFreshSaveAB)
		require.Nil(t, err, "should not return an error")
		assert.Equal(t, &sd.gameSaveA, sd.activeGameSave)
		for i, s := range sd.gameSaveA {
			assert.Falsef(t, s.empty, "section %d of save A is empty", i)
			assert.Equalf(t, s.counter, uint32(2), "section %d of save A has wrong counter", i)
		}
		for i, s := range sd.gameSaveB {
			assert.Falsef(t, s.empty, "section %d of save B is empty", i)
			assert.Equalf(t, s.counter, uint32(1), "section %d of save B has wrong counter", i)
		}
		for i, s := range sd.hallOfFame {
			assert.Truef(t, s.empty, "section %d of HoF is not empty", i)
		}
		assert.Equal(t, sd.GetPlayerGender(), Male)
		assert.Equal(t, sd.GetPlayerName(), "ACID")
	})

	t.Run("loading a PKHeX file w/ empty checksummed HoF", func(t *testing.T) {
		t.Parallel()
		sd, err := LoadSaveFileFromBytes(ePKHeXAfterStarter)
		require.Nil(t, err, "should not return an error")
		assert.Equal(t, &sd.gameSaveA, sd.activeGameSave)
		for i, s := range sd.gameSaveA {
			assert.Falsef(t, s.empty, "section %d of save A is empty", i)
			assert.Equalf(t, s.counter, uint32(8), "section %d of save A has wrong counter", i)
		}
		for i, s := range sd.gameSaveB {
			assert.Falsef(t, s.empty, "section %d of save B is empty", i)
			assert.Equalf(t, s.counter, uint32(7), "section %d of save B has wrong counter", i)
		}
		for i, s := range sd.hallOfFame {
			assert.Truef(t, s.empty, "section %d of HoF is not empty", i)
		}
	})

	testHoFnotEmpty := func(t *testing.T, bytes []byte, counterA, counterB uint32) {
		t.Parallel()
		sd, err := LoadSaveFileFromBytes(bytes)
		require.Nil(t, err, "should not return an error")
		assert.Equal(t, &sd.gameSaveB, sd.activeGameSave)
		for i, s := range sd.gameSaveA {
			assert.Falsef(t, s.empty, "section %d of save A is empty", i)
			assert.Equalf(t, s.counter, counterA, "section %d of save A has wrong counter", i)
		}
		for i, s := range sd.gameSaveB {
			assert.Falsef(t, s.empty, "section %d of save B is empty", i)
			assert.Equalf(t, s.counter, counterB, "section %d of save B has wrong counter", i)
		}
		assert.False(t, sd.hallOfFame[0].empty, "section 0 of HoF should not be empty")
		assert.False(t, sd.hallOfFame[1].empty, "section 1 of HoF should not be empty")
		assert.NotZero(t, sd.hallOfFame[0].checksum)
		assert.Zero(t, sd.hallOfFame[1].checksum)
	}

	t.Run("loading a PKHeX file w/ HoF sector 0 data", func(t *testing.T) {
		testHoFnotEmpty(t, ePKHeXChamp, 0x2C, 0x2D)
	})

	t.Run("loading a file w/ HoF sector 0 data", func(t *testing.T) {
		testHoFnotEmpty(t, eHoF0Save, 0x4F2, 0x4F3)
	})
}

func TestWrite(t *testing.T) {
	t.Parallel()

	write := func(t *testing.T, sd SaveData, bs []byte) *bytes.Buffer {
		t.Helper()
		buf := bytes.NewBuffer(nil)
		sd.Write(buf)
		return buf
	}

	testUnmodified := func(t *testing.T, bs []byte) {
		t.Parallel()
		sd, err := LoadSaveFileFromBytes(bs)
		require.Nil(t, err, "should not return an error")
		buf := write(t, sd, bs)
		assert.Equal(t, bs, buf.Bytes())
	}

	t.Run("write w/o modifications", func(t *testing.T) {
		t.Run("fresh save B", func(t *testing.T) { testUnmodified(t, eFreshSaveB) })
		t.Run("fresh save AB", func(t *testing.T) { testUnmodified(t, eFreshSaveAB) })
		t.Run("HoF sector 0", func(t *testing.T) { testUnmodified(t, eHoF0Save) })
		t.Run("PKHeX after starter", func(t *testing.T) { testUnmodified(t, ePKHeXAfterStarter) })
		t.Run("PKHeX after champ", func(t *testing.T) { testUnmodified(t, ePKHeXChamp) })
	})

	testModify := func(t *testing.T, bs []byte) {
		t.Parallel()

		const (
			gender = Female
			name   = "POKE"
		)

		sd, err := LoadSaveFileFromBytes(bs)
		require.Nil(t, err)

		sd.SetPlayerGender(gender)
		require.Equal(t, sd.GetPlayerGender(), gender)
		sd.SetPlayerName(name)
		require.Equal(t, sd.GetPlayerName(), name)

		buf := write(t, sd, eFreshSaveB)
		sdMod, err := LoadSaveFileFromBytes(buf.Bytes())
		require.Nil(t, err)
		assert.Equal(t, sdMod.GetPlayerGender(), sd.GetPlayerGender())
		assert.Equal(t, sdMod.GetPlayerName(), sd.GetPlayerName())
	}

	t.Run("write after simple modifications", func(t *testing.T) {
		t.Run("fresh save B", func(t *testing.T) { testModify(t, eFreshSaveB) })
		t.Run("fresh save AB", func(t *testing.T) { testModify(t, eFreshSaveAB) })
		t.Run("HoF sector 0", func(t *testing.T) { testModify(t, eHoF0Save) })
		t.Run("PKHeX after starter", func(t *testing.T) { testModify(t, ePKHeXAfterStarter) })
		t.Run("PKHeX after champ", func(t *testing.T) { testModify(t, ePKHeXChamp) })
	})
}
