package jwtauth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Signer struct {
	parser     *jwt.Parser
	key        []byte
	expiration time.Duration
	TokenName  string
}

func NewSigner(key []byte, expiration time.Duration) *Signer {
	return &Signer{
		parser:     jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name})),
		key:        key,
		expiration: expiration,
		TokenName:  "_token",
	}
}

func (s *Signer) GetFromRequest(req *http.Request) (*jwt.Token, error) {
	var jwtString string
	if jwtString = s.getJwtString(req); jwtString == "" {
		return nil, fmt.Errorf("jwt token not found")
	}
	return s.parser.Parse(jwtString, func(token *jwt.Token) (interface{}, error) {
		return s.key, nil
	})
}

func (s *Signer) ForgeToken(subject string) (string, time.Time, error) {
	expiration := time.Now().Add(s.expiration)
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expiration),
		Subject:   subject,
	}
	tkn, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.key)
	return tkn, expiration, err
}

func (s *Signer) CreateCookie(token string, expiration time.Time) *http.Cookie {
	return &http.Cookie{
		Path:     "/",
		Name:     s.TokenName,
		Value:    token,
		HttpOnly: true,
		Expires:  expiration,
		SameSite: http.SameSiteStrictMode,
	}
}

func (s *Signer) ExpireCookie() *http.Cookie {
	expiration := time.Now().Add(-1 * time.Hour)
	return &http.Cookie{
		Path:     "/",
		Name:     s.TokenName,
		Value:    "",
		HttpOnly: true,
		Expires:  expiration,
		SameSite: http.SameSiteStrictMode,
	}
}

func (s *Signer) getJwtString(req *http.Request) string {
	if cook, err := req.Cookie(s.TokenName); err == nil {
		return cook.Value
	}
	header := req.Header.Get("Authentication")
	if header == "" {
		return ""
	}
	idx := strings.Index(header, "Bearer ")
	if idx != 0 {
		return ""
	}
	return header[7:]
}
