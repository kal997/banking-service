package service

import (
	"github.com/kal997/banking/dto"
	"github.com/kal997/banking/errs"
)

type AuthService interface {
	Login(dto.LoginRequest) (*dto.LoginResponse, *errs.AppError)
}
