package ports

import "cards-service/internal/core/domain"

type AppService interface {
	ValidateCardNumber(cardNumber string) (*domain.CardInfo, error)
}
