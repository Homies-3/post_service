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
	redis := env.InitRedis()

	queue := env.InitQueue()
	ch, err := queue.Channel()
	if err != nil {
		log.Fatalln(err)
	}

	defer func() {
		queue.Close()
		ch.Close()
	}()

	// create queue
	cQ, err := ch.QueueDeclare("create-queue",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalln(err)
	}

	// update queue
	uQ, err := ch.QueueDeclare("test-name",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalln(err)
	}

	lis, err := net.Listen("tcp", env.ServerPort)

	if err != nil {
		log.Fatalln("Failed to listing:", err)
	}

	grpcServer := grpc.NewServer()
	pR := repo.NewPostRepository(logger, db)
	pb.RegisterPostServiceServer(grpcServer, services.NewService(logger, pR, redis, &cQ, &uQ, ch))
	log.Println("server started")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}

}
