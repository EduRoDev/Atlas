package http

import "github.com/EduRoDev/Atlas/internal/auth/domain"

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,min=3,max=30"`
	Password string `json:"password" binding:"required"`
	FullName string `json:"full_name" binding:"required,min=8"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type userResponse struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	Username   string `json:"username"`
	FullName   string `json:"full_name"`
	IsVerified bool   `json:"is_verified"`
}

func toUserResponse(u *domain.User) userResponse {
	return userResponse{
		ID: u.ID.String(),
		Email: u.Email,
		Username: u.Username,
		FullName: u.FullName,
		IsVerified: u.IsVerified,
	}
}