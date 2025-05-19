package server

import (
	"context"
	"log"
	"time"

	desc "github.com/Kosfedev/auth/pkg/user/gRPC"
	"github.com/brianvoe/gofakeit"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Server is...
type Server struct {
	desc.UnimplementedAuthV1Server
}

// Create is...
func (s *Server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	return &desc.CreateResponse{Id: gofakeit.Int64()}, nil
}

// Get is...
func (s *Server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	log.Printf("user id: %d\n", req.GetId())

	return &desc.GetResponse{
		Id:        req.GetId(),
		Name:      gofakeit.Name(),
		Email:     gofakeit.Email(),
		Role:      1,
		CreatedAt: timestamppb.New(gofakeit.Date()),
		UpdatedAt: timestamppb.New(gofakeit.Date()),
	}, nil
}

// Update is...
func (s *Server) Update(ctx context.Context, req *desc.UpdateRequest) (*desc.UpdateResponse, error) {
	log.Printf("user id: %d\n", req.GetId())

	updatedUser := &desc.UpdateResponse{
		Id:        req.GetId(),
		Name:      gofakeit.Name(),
		Email:     gofakeit.Email(),
		Role:      gofakeit.Uint32(),
		CreatedAt: timestamppb.New(gofakeit.Date()),
		UpdatedAt: timestamppb.New(time.Now()),
	}

	if name := req.GetName(); name != nil {
		updatedUser.Name = name.GetValue()
	}

	if email := req.GetEmail(); email != nil {
		updatedUser.Email = email.GetValue()
	}

	if role := req.GetRole(); role != nil {
		updatedUser.Role = role.GetValue()
	}

	return updatedUser, nil
}

// Delete is...
func (s *Server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	log.Printf("user id: %d\n", req.GetId())

	return nil, nil
}
