package domain

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"

	"github.com/kal997/banking-lib/logger"
)

type RemoteAuthRepository struct {
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
	authPort := os.Getenv("AUTH_SERVER_PORT")
	authIp := os.Getenv("SERVER_ADDRESS")
	host := authIp + ":" + authPort

	u := url.URL{Host: host, Path: "/auth/verify", Scheme: "http"}
	q := u.Query()
	q.Add("token", token)
	q.Add("routeName", routeName)

	for k, v := range vars {
		q.Add(k, v)
	}

	u.RawQuery = q.Encode()
	return u.String()
}

func (asd RemoteAuthRepository) IsAuthorized(token string, routeName string, vars map[string]string) bool {

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

func NewRemoteAuthRepository() RemoteAuthRepository {
	return RemoteAuthRepository{}
}
