package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

var ErrInvalidToken = errors.New("invalid token")

const GuestType = "guest"
const RealType = "real"

type UserInfo struct {
	UserID   string `json:"uid"`
	DeviceID string `json:"did"`
	GameID   string `json:"gid"`
	AppID    string `json:"aid"`

	// guest or real users
	// real means the user signed up using social network
	Type string `json:"utp"`
}

// Claims are typical jwt claims shared by all platforms.
type Claims struct {
	UserInfo
	jwt.StandardClaims
}

type Signer struct {
	hmacSecret []byte
}

func NewSigner(hmacSecret string) *Signer {
	return &Signer{hmacSecret: []byte(hmacSecret)}
}

func (s *Signer) Encode(claims Claims) (string, error) {
	claims.StandardClaims = jwt.StandardClaims{
		IssuedAt: time.Now().UTC().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.hmacSecret)
}

func (s *Signer) Decode(ts string) (*Claims, error) {
	claimsOut := &Claims{}

	token, err := jwt.ParseWithClaims(ts, claimsOut, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return s.hmacSecret, nil
	})
	if err != nil {
		return claimsOut, errors.Wrapf(err, ErrInvalidToken.Error())
	}

	if !token.Valid {
		return claimsOut, ErrInvalidToken
	}

	return claimsOut, nil
}
