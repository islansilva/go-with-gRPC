package main

import (
	"net"
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
	"example/gRPC/internal/database"
	"example/gRPC/internal/service"
	"example/gRPC/internal/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)


func main() {
	db, err := sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	categoryDB := database.NewCategory(db)
	categoryService:= service.NewCategoryService(*categoryDB)

	grpcServer := grpc.NewServer()
	pb.RegisterCategoryServiceServer(grpcServer, categoryService)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":50051")

	if err != nil {
		panic(err)
	}

	if err:= grpcServer.Serve(lis); err != nil {
		panic(err)
	}

}