package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/dgrijalva/jwt-go"
	"github.com/kal997/banking/domain"
	"github.com/kal997/banking/dto"
	"github.com/kal997/banking/errs"
	"github.com/kal997/banking/logger"
)

type DefaultAuthService struct {
	repo            domain.AuthRepository
	rolePermissions domain.RolePermissions
}

func (asd DefaultAuthService) IsAuthorized(token string, routeName string, vars map[string]string) bool {

	u := buildVerifyURL(token, routeName, vars)

	response, err := http.Get(u)
	if err != nil {
		logger.Error("error while sending Verify route " + err.Error())
		return false
	}

	m := map[string]bool{}
	err = json.NewDecoder(response.Body).Decode(&m)
	if err != nil {
		var errorStr string
		err = json.NewDecoder(response.Body).Decode(&errorStr)
		if err != nil {
			return false
		} else {
			logger.Info(errorStr)
			return false
		}
	}

	return m["isAuthorized"]

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

func NewAuthService(repo domain.AuthRepository, rolePermissions domain.RolePermissions) DefaultAuthService {
	return DefaultAuthService{repo: repo, rolePermissions: rolePermissions}

}

func (asd DefaultAuthService) Verify(urlParams map[string]string) *errs.AppError {
	// convert the string token claims into JWT struct
	jwtToken, err := jwtTokenFromString(urlParams["token"])
	// if failed, return unexpected app error
	if err != nil {
		return errs.NewAuthorizationError(err.Error())
	}
	// conversion is successfull
	// check if the token is valid in terms on time and signture
	// if passed
	if jwtToken.Valid {

		// cast the JWT claims into our User defined claims struct
		claims := jwtToken.Claims.(*domain.AccessTokenClaims)

		// if role is user, check if customer_id in the claims and the one in the request are matching
		if claims.IsUserRole() {
			if !claims.IsRequestedVerifiedWithTokenClaims(urlParams) {
				return errs.NewAuthorizationError("request not verified with the token claims")
			}

		}
		// if true,  verify the if the route can be used by this role
		// else return false, unauthorized route
		if !asd.rolePermissions.IsAuthorizedFor(claims.Role, urlParams["routeName"]) {
			return errs.NewAuthorizationError(fmt.Sprintf("%s role is not authorized", claims.Role))

		}
		// if true return true
		return nil

	} else {
		// if failed, token is expired or invalid signture (if expired, refresh)
		return errs.NewAuthorizationError("Invalid token")

	}

}

func jwtTokenFromString(tokenStr string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &domain.AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(domain.HMAC_SAMPLE_SECRET), nil
	})
	if err != nil {
		logger.Error("Error while parsing token: " + err.Error())
		return nil, err
	}
	return token, nil
}

/*
This will generate a url for token verification in the below format

/auth/verify?token={token string}

	&routeName={current route name}
	&customer_id={customer id from the current route}
	&account_id={account id from current route if available}

Sample: /auth/verify?token=aaaa.bbbb.cccc&routeName=MakeTransaction&customer_id=2000&account_id=95470
*/
func buildVerifyURL(token string, routeName string, vars map[string]string) string {
	u := url.URL{Host: "127.0.0.1:8005", Path: "/auth/verify", Scheme: "http"}
	q := u.Query()
	q.Add("token", token)
	q.Add("routeName", routeName)

	for k, v := range vars {
		q.Add(k, v)
	}

	u.RawQuery = q.Encode()
	return u.String()
}
