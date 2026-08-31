package stringutils_test

import (
	"testing"

	"github.com/matthiasharzer/go-stats-tracker/utils/stringutils"
	"github.com/stretchr/testify/assert"
)

func TestRandomString(t *testing.T) {
	t.Run("generates a random string of the specified length", func(t *testing.T) {
		for i := 1; i <= 100; i++ {
			randomStr := stringutils.RandomString(i)
			assert.Equal(t, i, len(randomStr))
		}
	})
	t.Run("generates different strings on subsequent calls", func(t *testing.T) {
		randomStr1 := stringutils.RandomString(255)
		randomStr2 := stringutils.RandomString(255)

		// It's highly unlikely that two random strings of length 255 will be the same
		assert.NotEqual(t, randomStr1, randomStr2)
	})
	t.Run("generates a string with only alphanumeric characters", func(t *testing.T) {
		randomStr := stringutils.RandomString(255)
		for _, c := range randomStr {
			if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')) {
				assert.Fail(t, "Random string contains non-alphanumeric character", "Character: %c", c)
			}
		}
	})
}
