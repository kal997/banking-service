package service

import (
	"github.com/kal997/banking/domain"
	"github.com/kal997/banking/dto"
	"github.com/kal997/banking/errs"
)

type DefaultAuthService struct {
	repo domain.AuthRepository
}

func (asd DefaultAuthService) Login(request dto.LoginRequest) (*dto.LoginResponse, *errs.AppError) {
	// do service side validation
	login, err := asd.repo.FindBy(request.Username, request.Password)
	if err != nil {
		return nil, err
	}

	// build claims from login domain obj
	claims := login.ClaimsForAccessToken()
	// generate new auth token
	authToken := domain.NewAuthToken(claims)
	// sign the auth token
	accessToken, err := authToken.NewAccessToken()
	if err != nil {
		return nil, err
	}

	// convert it into dto.LoginResponce and return it
	return &dto.LoginResponse{AccessToken: accessToken}, nil

}

func NewLoginService(repo domain.AuthRepository) DefaultAuthService {
	return DefaultAuthService{repo: repo}

}
