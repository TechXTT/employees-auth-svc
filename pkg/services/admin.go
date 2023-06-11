package services

import (
	"context"
	"log"
	"net/http"

	"gorm.io/gorm"

	"github.com/TechXTT/employees-auth-svc/pkg/models"
	"github.com/TechXTT/employees-auth-svc/pkg/pb"
	"github.com/TechXTT/employees-auth-svc/pkg/utils"
)

type AdminService interface {
	Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error)
	Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error)
	ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error)
}

type adminService struct {
	*gorm.DB
	Jwt utils.JwtWrapper
}

func NewAdminService(db *gorm.DB, jwt utils.JwtWrapper) AdminService {
	return &adminService{db, jwt}
}

func (s *adminService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	var admin models.Admin

	// check if user exists in db, where req.EmailOrUsername can be either email or username
	if err := s.DB.WithContext(ctx).Where("email = ?", req.EmailOrUsername).Or("username = ?", req.EmailOrUsername).First(&admin).Error; err != nil {
		return &pb.LoginResponse{
			Error:  "User not found",
			Status: http.StatusNotFound,
		}, nil
	}

	match := utils.CheckPasswordHash(req.Password, admin.Password)

	if !match {
		return &pb.LoginResponse{
			Error:  "Wrong password",
			Status: http.StatusNotFound,
		}, nil
	}

	token, err := s.Jwt.GenerateToken(admin)
	if err != nil {
		return &pb.LoginResponse{
			Error:  err.Error(),
			Status: http.StatusInternalServerError,
		}, nil
	}

	return &pb.LoginResponse{
		Token:  token,
		Status: http.StatusOK,
	}, nil

}

func (s *adminService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	var admin models.Admin

	if err := s.DB.WithContext(ctx).Where("email = ?", req.Email).First(&admin).Error; err == nil {
		return &pb.RegisterResponse{
			Error:  "User already exists",
			Status: http.StatusConflict,
		}, nil
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return &pb.RegisterResponse{
			Error:  err.Error(),
			Status: http.StatusInternalServerError,
		}, nil
	}

	admin = models.Admin{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := s.DB.WithContext(ctx).Create(&admin).Error; err != nil {
		log.Println("Error creating user: ", err)
		return &pb.RegisterResponse{
			Error:  err.Error(),
			Status: http.StatusInternalServerError,
		}, nil
	}

	token, err := s.Jwt.GenerateToken(admin)
	if err != nil {
		return &pb.RegisterResponse{
			Error:  err.Error(),
			Status: http.StatusInternalServerError,
		}, nil
	}

	return &pb.RegisterResponse{
		Status: http.StatusCreated,
		Token:  token,
	}, nil
}

func (s *adminService) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	claims, err := s.Jwt.ValidateToken(req.Token)
	if err != nil {
		return &pb.ValidateTokenResponse{
			Error:  err.Error(),
			Status: http.StatusUnauthorized,
		}, nil
	}

	var admin models.Admin

	if err := s.DB.WithContext(ctx).Where("id = ?", claims.Id).First(&admin).Error; err != nil {
		return &pb.ValidateTokenResponse{
			Error:  "User not found",
			Status: http.StatusNotFound,
		}, nil
	}

	return &pb.ValidateTokenResponse{
		Status: http.StatusOK,
		UserId: admin.ID.String(),
	}, nil
}
