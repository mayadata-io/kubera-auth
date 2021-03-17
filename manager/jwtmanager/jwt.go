package jwtmanager

import (
	"time"

	"github.com/mayadata-io/kubera-auth/manager/usermanager"
	"github.com/mayadata-io/kubera-auth/pkg/generates"
	"github.com/mayadata-io/kubera-auth/pkg/models"
	"github.com/mayadata-io/kubera-auth/pkg/store"
)

// TokenGenerateRequest provide to generate the token request parameters
type TokenGenerateRequest struct {
	UserInfo       *models.PublicUserInfo
	AccessTokenExp time.Duration
}

// Config authorization configuration parameters
type Config struct {
	// access token expiration time, 0 means it doesn't expire
	AccessTokenExp time.Duration
}

// default configs
var (
	DefaultTokenCfg = &Config{AccessTokenExp: time.Hour * 24}
)

// ParseToken validates the token
func ParseToken(userStore *store.UserStore, accessGenerate *generates.JWTAccessGenerate, tokenString string) (*models.UserCredentials, error) {
	claimedUser, err := accessGenerate.Parse(tokenString)
	if err != nil {
		return nil, err
	}
	user, err := usermanager.GetUserByUID(userStore, claimedUser.UID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GenerateAuthToken generate the authorization token(code)
func GenerateAuthToken(accessGenerate *generates.JWTAccessGenerate, tgr *TokenGenerateRequest, jwtType models.TokenType) (*models.Token, error) {
	ti := models.NewToken(jwtType)

	createAt := time.Now()
	td := &generates.GenerateBasic{
		UserInfo:  tgr.UserInfo,
		CreateAt:  &createAt,
		TokenInfo: ti,
	}

	cfg := DefaultTokenCfg
	aexp := cfg.AccessTokenExp
	if exp := tgr.AccessTokenExp; exp > 0 {
		aexp = exp
	}
	ti.SetAccessCreateAt(createAt)
	ti.SetAccessExpiresIn(aexp)

	tv, err := accessGenerate.Token(td)
	if err != nil {
		return nil, err
	}
	ti.SetAccess(tv)
	return ti, nil
}
