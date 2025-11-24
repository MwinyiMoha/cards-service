package api

import (
	"cards-service/internal/core/domain"
	"cards-service/internal/core/ports"
	"context"
	"net"
	"testing"

	"github.com/mwinyimoha/protos/gen/go/pb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

type mockAppService struct {
	cardInfo *domain.CardInfo
	err      error
}

func (m *mockAppService) ValidateCardNumber(cardNumber string) (*domain.CardInfo, error) {
	if m.err != nil {
		return nil, m.err
	}

	return m.cardInfo, nil
}

func setupGRPCServer(t *testing.T, svc ports.AppService) (*grpc.ClientConn, func()) {
	listener := bufconn.Listen(bufSize)

	server := grpc.NewServer()
	pb.RegisterCardsServiceServer(server, NewServer(svc))

	go func() {
		if err := server.Serve(listener); err != nil {
			t.Logf("gRPC server stopped: %v", err)
		}
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}

	conn, err := grpc.DialContext(
		context.Background(),
		"bufnet",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)

	cleanup := func() {
		server.Stop()
		conn.Close()
	}

	return conn, cleanup
}

func TestValidateCardNumberRPC(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		expected := &domain.CardInfo{
			CardNumber:    "4111111111111111",
			CardProvider:  "Visa",
			ProviderBadge: "visa",
		}
		mockSvc := &mockAppService{cardInfo: expected}

		conn, cleanup := setupGRPCServer(t, mockSvc)
		defer cleanup()

		client := pb.NewCardsServiceClient(conn)

		req := &pb.ValidateCardNumberRequest{CardNumber: expected.CardNumber}
		resp, err := client.ValidateCardNumber(context.Background(), req)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, expected.CardNumber, resp.CardNumber)
		assert.Equal(t, expected.CardProvider, resp.ProviderName)
		assert.Equal(t, expected.ProviderBadge, resp.ProviderBadge)
	})

	t.Run("Error", func(t *testing.T) {
		mockErr := assert.AnError
		mockSvc := &mockAppService{err: mockErr}

		conn, cleanup := setupGRPCServer(t, mockSvc)
		defer cleanup()

		client := pb.NewCardsServiceClient(conn)

		req := &pb.ValidateCardNumberRequest{CardNumber: "4111111111111111"}
		resp, err := client.ValidateCardNumber(context.Background(), req)

		assert.Nil(t, resp)
		require.Error(t, err)
		assert.Contains(t, err.Error(), mockErr.Error())
	})
}
