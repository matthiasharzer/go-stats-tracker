package sheetutils_test

import (
	"testing"

	"github.com/matthiasharzer/go-stats-tracker/utils/sheetutils"
	"github.com/stretchr/testify/assert"
)

func TestIncrementColumn(t *testing.T) {
	t.Run("increments single letters", func(t *testing.T) {
		t.Run("increments A to B", func(t *testing.T) {
			nextColumn := sheetutils.IncrementColumn("A")
			assert.Equal(t, "B", nextColumn)
		})
		t.Run("increments G to H", func(t *testing.T) {
			nextColumn := sheetutils.IncrementColumn("G")
			assert.Equal(t, "H", nextColumn)
		})
		t.Run("increments Z to AA", func(t *testing.T) {
			nextColumn := sheetutils.IncrementColumn("Z")
			assert.Equal(t, "AA", nextColumn)
		})
	})
	t.Run("increments multiple letters", func(t *testing.T) {
		t.Run("increments AA to AB", func(t *testing.T) {
			nextColumn := sheetutils.IncrementColumn("AA")
			assert.Equal(t, "AB", nextColumn)
		})
		t.Run("increments AZ to BA", func(t *testing.T) {
			nextColumn := sheetutils.IncrementColumn("AZ")
			assert.Equal(t, "BA", nextColumn)
		})
		t.Run("increments ZZ to AAA", func(t *testing.T) {
			nextColumn := sheetutils.IncrementColumn("ZZ")
			assert.Equal(t, "AAA", nextColumn)
		})
	})
}

func TestIncrementColumnN(t *testing.T) {
	t.Run("increments A to C", func(t *testing.T) {
		nextColumn := sheetutils.IncrementColumnN("A", 2)
		assert.Equal(t, "C", nextColumn)
	})
	t.Run("increments G to Z", func(t *testing.T) {
		nextColumn := sheetutils.IncrementColumnN("G", 19)
		assert.Equal(t, "Z", nextColumn)
	})
	t.Run("increments Z to AC", func(t *testing.T) {
		nextColumn := sheetutils.IncrementColumnN("Z", 3)
		assert.Equal(t, "AC", nextColumn)
	})
	t.Run("increments AZ to BC", func(t *testing.T) {
		nextColumn := sheetutils.IncrementColumnN("AZ", 3)
		assert.Equal(t, "BC", nextColumn)
	})
	t.Run("increments Z to AAA", func(t *testing.T) {
		nextColumn := sheetutils.IncrementColumnN("Z", 26*26+1)
		assert.Equal(t, "AAA", nextColumn)
	})
}
