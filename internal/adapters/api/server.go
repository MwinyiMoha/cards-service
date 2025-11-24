package api

import (
	"cards-service/internal/core/ports"
	"context"

	"github.com/mwinyimoha/protos/gen/go/pb"
)

type Server struct {
	pb.UnimplementedCardsServiceServer
	service ports.AppService
}

func NewServer(svc ports.AppService) *Server {
	return &Server{service: svc}
}

func (srv *Server) ValidateCardNumber(ctx context.Context, req *pb.ValidateCardNumberRequest) (*pb.ValidateCardNumberResponse, error) {
	cardInfo, err := srv.service.ValidateCardNumber(req.GetCardNumber())
	if err != nil {
		return nil, err
	}

	return &pb.ValidateCardNumberResponse{
		CardNumber:    cardInfo.CardNumber,
		ProviderName:  cardInfo.CardProvider,
		ProviderBadge: cardInfo.ProviderBadge,
	}, nil
}
