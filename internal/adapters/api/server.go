package api

import (
	"cards-service/internal/core/ports"
	"context"

	"github.com/mwinyimoha/commons/pkg/errors"
	"github.com/mwinyimoha/protos/gen/go/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		if verr, ok := err.(*errors.ValidationError); ok {
			return nil, status.Convert(verr).Err()
		}

		if cerr, ok := err.(*errors.Error); ok {
			return nil, status.Convert(cerr).Err()
		}

		return nil, status.Errorf(codes.Internal, "internal error: %v", err)
	}

	return &pb.ValidateCardNumberResponse{
		CardNumber:    cardInfo.CardNumber,
		ProviderName:  cardInfo.CardProvider,
		ProviderBadge: cardInfo.ProviderBadge,
	}, nil
}
