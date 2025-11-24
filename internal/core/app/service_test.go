package app

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/mwinyimoha/commons/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService(t *testing.T) {
	val := validator.New()
	svc := NewService(val)

	t.Run("Valid Cards", func(t *testing.T) {
		tests := []string{
			"4111111111111111", // Visa 16 digits
			"4012888888881881", // Visa
			"378282246310005",  // Amex 15 digits
			"5555555555554444", // MasterCard
			"6011111111111117", // Discover
		}

		for _, card := range tests {
			t.Run(card, func(t *testing.T) {
				info, err := svc.ValidateCardNumber(card)
				require.NoError(t, err)
				require.NotNil(t, info)
				assert.Equal(t, card, info.CardNumber)
			})
		}
	})

	t.Run("Invalid Cards", func(t *testing.T) {
		invalidCards := []string{
			"123456789012345",  // invalid prefix
			"4111111111111",    // too short
			"37828224631000",   // Amex too short
			"4111111111111112", // Luhn fail
			"6011abcd11111117", // invalid chars
			"9111111111111111", // unsupported prefix
			"3245678901234561", // Amex longer than 15
		}

		for _, card := range invalidCards {
			t.Run(card, func(t *testing.T) {
				info, err := svc.ValidateCardNumber(card)
				require.Error(t, err)
				assert.Nil(t, info)

				_, ok := err.(*errors.Error)
				assert.True(t, ok, "expected error of type *errors.Error")
			})
		}
	})
}
