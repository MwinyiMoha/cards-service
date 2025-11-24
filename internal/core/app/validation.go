package app

import (
	"slices"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

func registerValidators(val *validator.Validate) {
	val.RegisterValidation("valid_card_number", validateCardNumber)
}

func validateCardNumber(fl validator.FieldLevel) bool {
	val := fl.Field().String()

	if !slices.Contains([]string{"3", "4", "5", "6"}, val[:1]) {
		return false
	}

	if len(val) < 15 || len(val) > 16 {
		return false
	}

	if val[:1] == "3" {
		if len(val) != 15 {
			return false
		}
	}

	return luhnValidation(val)
}

func luhnValidation(cardNumber string) bool {
	trimmed := strings.ReplaceAll(cardNumber, " ", "")
	chars := []rune(trimmed)
	slices.Reverse(chars)
	sum := 0

	for i, char := range string(chars) {
		digit, err := strconv.Atoi(string(char))
		if err != nil {
			return false
		}

		if i%2 == 1 {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
	}

	return sum%10 == 0
}
