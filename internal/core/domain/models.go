package domain

import (
	"fmt"
	"strings"

	"github.com/mwinyimoha/commons/pkg/errors"
)

type CardNumberPayload struct {
	CardNumber string `validate:"required,valid_card_number"`
}

type CardInfo struct {
	CardNumber    string
	CardProvider  string
	ProviderBadge string
}

func NewCardInfo(cardNumber string) (*CardInfo, error) {
	providers := map[string]string{
		"3": "AMEX",
		"4": "VISA",
		"5": "MASTERCARD",
		"6": "DISCOVER",
	}
	providerName, exists := providers[cardNumber[:1]]
	if !exists {
		return nil, errors.NewErrorf(errors.InvalidArgument, "unknown card provider")
	}

	return &CardInfo{
		CardNumber:    cardNumber,
		CardProvider:  providerName,
		ProviderBadge: fmt.Sprintf("https://dummy.com/card-provider-icons/%s.png", strings.ToLower(providerName)),
	}, nil
}
