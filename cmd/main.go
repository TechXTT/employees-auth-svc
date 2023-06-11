package main

import (
	"fmt"
	"log"
	"net"

	"github.com/TechXTT/employees-auth-svc/pkg/config"
	"github.com/TechXTT/employees-auth-svc/pkg/db"
	"github.com/TechXTT/employees-auth-svc/pkg/pb"
	"github.com/TechXTT/employees-auth-svc/pkg/services"
	"github.com/TechXTT/employees-auth-svc/pkg/utils"
	"google.golang.org/grpc"
)

func main() {
	c, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	} else {
		//print config
		fmt.Println(c)
	}

	DB := db.Init(c.DatabaseURL)

	jwt := utils.JwtWrapper{
		SecretKey:       c.JWTSecretKey,
		Issuer:          "AuthService",
		ExpirationHours: 24 * 7,
	}

	lis, err := net.Listen("tcp", c.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	fmt.Println("Auth service is running on port", c.Port)

	s := services.NewAdminService(DB, jwt)

	grpcServer := grpc.NewServer()

	pb.RegisterAuthServiceServer(grpcServer, s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
