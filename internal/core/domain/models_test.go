package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCardInfo(t *testing.T) {

	t.Run("Valid Provider Prefix", func(t *testing.T) {
		testNumber := "4111111111111111"
		cardInfo, err := NewCardInfo(testNumber)

		assert.NoError(t, err)
		assert.Equal(t, testNumber, cardInfo.CardNumber)
		assert.Equal(t, "VISA", cardInfo.CardProvider)
		assert.Equal(t, "https://dummy.com/card-provider-icons/visa.png", cardInfo.ProviderBadge)
	})

	t.Run("Invalid Provider Prefix", func(t *testing.T) {
		testNumber := "7111111111111111"
		cardInfo, err := NewCardInfo(testNumber)

		assert.Nil(t, cardInfo)
		assert.Error(t, err)
	})
}
