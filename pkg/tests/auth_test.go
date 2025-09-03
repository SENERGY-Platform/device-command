package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"time"

	"github.com/SENERGY-Platform/device-command/pkg/auth"
	"github.com/SENERGY-Platform/device-command/pkg/configuration"
	"github.com/golang-jwt/jwt"
)

func generateUserTokenById(userid string) (token string, err error) {
	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour).Unix(),
		Issuer:    "mgw-device-command",
		Subject:   userid,
	}
	jwtoken := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	unsignedTokenString, err := jwtoken.SigningString()
	if err != nil {
		log.Println("ERROR: GenerateUserTokenById::SigningString()", err, userid)
		return token, err
	}
	tokenString := strings.Join([]string{unsignedTokenString, ""}, ".")
	return tokenString, nil
}

func authMock(config configuration.Config, ctx context.Context, wg *sync.WaitGroup) (configuration.Config, error) {
	token, err := generateUserTokenById("testOwner")
	if err != nil {
		return config, err
	}
	handler := func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println("auth mock:", json.NewEncoder(writer).Encode(auth.OpenidToken{
			AccessToken: token,
		}))
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	wg.Add(1)
	go func() {
		<-ctx.Done()
		server.Close()
		wg.Done()
	}()
	config.AuthEndpoint = server.URL
	return config, nil
}
