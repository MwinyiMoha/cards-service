package app

import (
	"cards-service/internal/core/domain"

	"github.com/go-playground/validator/v10"
	"github.com/mwinyimoha/commons/pkg/errors"
)

type Service struct {
	validation *validator.Validate
}

func NewService() (*Service, error) {
	validation, err := createValidator()
	if err != nil {
		return nil, err
	}

	return &Service{validation: validation}, nil
}

func (svc *Service) ValidateCardNumber(cardNumber string) (*domain.CardInfo, error) {
	payload := &domain.CardNumberPayload{CardNumber: cardNumber}

	if err := svc.validation.Struct(payload); err != nil {
		if verr, ok := err.(validator.ValidationErrors); ok {
			violations := errors.BuildViolations(verr)
			return nil, errors.NewValidationError(violations)
		}

		return nil, errors.WrapError(err, errors.Internal, "validation failed")
	}

	return domain.NewCardInfo(cardNumber), nil
}
