package domain

import (
	"fmt"
	"strings"
)

type CardNumberPayload struct {
	CardNumber string `validate:"required,valid_card_number"`
}

type CardInfo struct {
	CardNumber    string
	CardProvider  string
	ProviderBadge string
}

func NewCardInfo(cardNumber string) *CardInfo {
	providers := map[string]string{
		"3": "AMEX",
		"4": "VISA",
		"5": "MASTERCARD",
		"6": "DISCOVER",
	}
	providerName := providers[cardNumber[:1]]
	return &CardInfo{
		CardNumber:    cardNumber,
		CardProvider:  providerName,
		ProviderBadge: fmt.Sprintf("https://dummy.com/card-provider-icons/%s.png", strings.ToLower(providerName)),
	}
}
