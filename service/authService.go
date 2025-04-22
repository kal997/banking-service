package service

import (
	"github.com/kal997/banking/dto"
	"github.com/kal997/banking/errs"
)

type AuthService interface {
	Login(dto.LoginRequest) (*dto.LoginResponse, *errs.AppError)
	Verify(urlParams map[string]string) *errs.AppError
	IsAuthorized(token string, routeName string, vars map[string]string) bool
}
