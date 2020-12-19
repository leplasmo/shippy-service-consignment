// shippy-service-consignment/main.go

package main

import (
	"context"
	"log"

	// import the generated protobuf code

	pb "github.com/leplasmo/shippy-service-consignment/proto/consignment"
	"github.com/micro/go-micro/v2"
)

// const (
// 	port = ":30051"
// )

type repository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
	GetAll() []*pb.Consignment
}

// Repository - Dummy repository, this simulates the use of a data store
// of some kind. We'll replace this with a real implementation later on.
type Repository struct {
	// mu           sync.RWMutex
	consignments []*pb.Consignment
}

// Create a new consignment
func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	// repo.mu.Lock()
	updated := append(repo.consignments, consignment)
	repo.consignments = updated
	// repo.mu.Unlock()
	return consignment, nil
}

// Get all consignments
func (repo *Repository) GetAll() []*pb.Consignment {
	return repo.consignments
}

// Service should implement all of the methods to satisfy the service
// we defined in our protobuf definition. You can check the interface
// in the generated code itself for the exact method signatures etc
// to give you a better idea.
type consignmentService struct {
	repo repository
}

// CreateConsignment - we created just one method on our service,
// which is a create method, which takes a context and a request as an
// argument, these are handled by the gRPC server.
func (s *consignmentService) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {
	// Save our consignment
	consignment, err := s.repo.Create(req)
	if err != nil {
		return err
	}

	// Return matching the `Response` message we created in our
	// protobuf definition.
	// return &pb.Response{Created: true, Consignment: consignment}, nil
	res.Created = true
	res.Consignment = consignment
	return nil
}

// GetConsignments
func (s *consignmentService) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	consignments := s.repo.GetAll()
	// return &pb.Response{Consignments: consignments}, nil
	res.Consignments = consignments
	return nil
}

func main() {
	repo := &Repository{}

	// Service Discovery
	// Create a new service. Optionally include some options here.
	service := micro.NewService(
		// This name must match the package name given in the
		// protobuf definition
		micro.Name("shippy.service.consignment"),
	)

	// Init will parse the command line flags.
	service.Init()

	// Register the service
	if err := pb.RegisterShippingServiceHandler(service.Server(), &consignmentService{repo}); err != nil {
		log.Panic(err)
	}

	// Run the server
	if err := service.Run(); err != nil {
		log.Panic(err)
	}

	// // Set-up our gRPC server.
	// lis, err := net.Listen("tcp", port)
	// if err != nil {
	// 	log.Fatalf("failed to listen: %v", err)
	// }
	// s := grpc.NewServer()

	// // Register our service with the gRPC server, this will tie our
	// // implementation into the auto-generated interface code for our
	// // protobuf definition.
	// pb.RegisterShippingServiceServer(s, &service{repo})

	// // Register reflection service on gRPC server.
	// reflection.Register(s)

	// log.Println("Running on port:", port)
	// if err := s.Serve(lis); err != nil {
	// 	log.Fatalf("failed to serve: %v", err)
	// }
}