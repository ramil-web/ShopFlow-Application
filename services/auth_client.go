package services

import (
	"context"
	"time"

	"shopflow/application/proto/authpb"

	"google.golang.org/grpc"
)

type AuthClient interface {
	VerifyToken(userID uint32, token string) (bool, string, error)
}

type authClientGRPC struct {
	client authpb.AuthServiceClient
}

func NewAuthClient(grpcAddr string) (AuthClient, error) {
	conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := authpb.NewAuthServiceClient(conn)
	return &authClientGRPC{client: client}, nil
}

func (a *authClientGRPC) VerifyToken(userID uint32, token string) (bool, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resp, err := a.client.VerifyToken(ctx, &authpb.AuthRequest{
		UserId: userID,
		Token:  token,
	})
	if err != nil {
		return false, "", err
	}
	return resp.Valid, resp.Email, nil
}
