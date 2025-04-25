package service

import (
	"github.com/kal997/banking-lib/errs"
	"github.com/kal997/banking/dto"
)

type AuthService interface {
	Login(dto.LoginRequest) (*dto.LoginResponse, *errs.AppError)
	Verify(urlParams map[string]string) *errs.AppError
	IsAuthorized(token string, routeName string, vars map[string]string) bool
}
