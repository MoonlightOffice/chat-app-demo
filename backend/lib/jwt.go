package lib

import (
	"github.com/moonlightoffice/jwt-go"
)

const jwtSecret string = "sample-jwt-key"

func ToJWT(payload []byte) (string, error) {
	return jwt.Encode(
		jwt.JoseHeader{Alg: jwt.AlgHS512},
		payload,
		[]byte(jwtSecret),
	)
}

func FromJWT(token string) ([]byte, error) {
	return jwt.Decode(token, []byte(jwtSecret))
}
