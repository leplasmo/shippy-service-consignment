// shippy-service-consignment/main.go

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	// import the generated protobuf code

	pb "github.com/leplasmo/shippy-service-consignment/proto/consignment"
	vesselProto "github.com/leplasmo/shippy-service-vessel/proto/vessel"
	"github.com/micro/go-micro/v2"
)

const (
	defaultHost = "datastore:27017"
)

func main() {

	// Service Discovery
	// Create a new service. Optionally include some options here.
	service := micro.NewService(
		// This name must match the package name given in the
		// protobuf definition
		micro.Name("shippy.service.consignment"),
	)

	// Init will parse the command line flags.
	service.Init()

	uri := os.Getenv("DB_HOST")
	if uri == "" {
		uri = defaultHost
	}

	client, err := CreateClient(context.Background(), uri, 0)
	if err != nil {
		log.Panic(err)
	}
	defer client.Disconnect(context.Background())

	consignmentCollection := client.Database("shippy").Collection("consignments")

	repository := &MongoRepository{consignmentCollection}

	vesselClient := vesselProto.NewVesselService("shippy.service.client", service.Client())
	h := &handler{repository, vesselClient}

	// Register handlers
	pb.RegisterShippingServiceHandler(service.Server(), h)

	// Run the server
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
