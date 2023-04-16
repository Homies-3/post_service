package main

import (
	"log"
	"net"
	"os"
	pb "post_service/pb"
	"post_service/utils"

	repo "post_service/repositories"
	services "post_service/services"

	"google.golang.org/grpc"
)

func main() {
	logger := log.New(os.Stdout, "post-service", log.LstdFlags)
	env := utils.LoadEnv(logger)
	db := env.ConnectToDB()

	lis, err := net.Listen("tcp", env.ServerPort)

	if err != nil {
		log.Fatalln("Failed to listing:", err)
	}

	grpcServer := grpc.NewServer()
	pR := repo.NewPostRepository(logger, *db.Collection("posts"))
	pb.RegisterPostServiceServer(grpcServer, services.NewService(logger, pR))
	log.Println("started server")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}

}
