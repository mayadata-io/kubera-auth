package generates

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	errs "errors"

	"github.com/dgrijalva/jwt-go"
	log "github.com/golang/glog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/mayadata-io/kubera-auth/pkg/errors"
	"github.com/mayadata-io/kubera-auth/pkg/k8s"
	"github.com/mayadata-io/kubera-auth/pkg/models"
	"github.com/mayadata-io/kubera-auth/pkg/types"
	"github.com/mayadata-io/kubera-auth/pkg/utils/random"
)

// JWTAccessClaims jwt claims
type JWTAccessClaims struct {
	UUID     string      `json:"uuid,omitempty"`
	Role     models.Role `json:"role,omitempty"`
	UserName string      `json:"username,omitempty"`
	Email    *string     `json:"email,omitempty"`
	Name     string      `json:"name,omitempty"`
	jwt.StandardClaims
}

func init() {
	if os.Getenv("CONFIGMAP_NAME") == "" {
		log.Fatal("Environment variable CONFIGMAP_NAME is not set")
	}
}

// NewJWTAccessGenerate create to generate the jwt access token instance
func NewJWTAccessGenerate(method jwt.SigningMethod) *JWTAccessGenerate {

	key, err := initializeSecret()
	if err != nil {
		log.Fatal(err)
	}
	return &JWTAccessGenerate{
		SignedKey:    []byte(key),
		SignedMethod: method,
	}
}

// GenerateBasic provide the basis of the generated token data
type GenerateBasic struct {
	UserInfo  *models.PublicUserInfo
	CreateAt  *time.Time
	TokenInfo *models.Token
}

// JWTAccessGenerate generate the jwt access token
type JWTAccessGenerate struct {
	SignedKey    []byte
	SignedMethod jwt.SigningMethod
}

func initializeSecret() (string, error) {

	if k8s.ClientSet == nil {
		return "", errors.New("ClienSet not found")
	}

	cm, err := k8s.ClientSet.CoreV1().ConfigMaps(types.DefaultNamespace).Get(context.TODO(), types.DefaultConfigMap, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	if cm.Data[types.JWTSecretString] != "" {
		return cm.Data[types.JWTSecretString], nil
	}

	secret := random.GetRandomString(10)
	cm.Data[types.JWTSecretString] = secret
	_, err = k8s.ClientSet.CoreV1().ConfigMaps(types.DefaultNamespace).Update(context.TODO(), cm, metav1.UpdateOptions{})
	return secret, err
}

// Token based on the UUID generated token
func (a *JWTAccessGenerate) Token(data *GenerateBasic) (string, error) {
	claims := &JWTAccessClaims{
		UUID:     data.UserInfo.GetID(),
		Role:     data.UserInfo.GetRole(),
		UserName: data.UserInfo.GetUserName(),
		Email:    data.UserInfo.GetEmail(),
		Name:     data.UserInfo.GetName(),
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: data.TokenInfo.GetAccessCreateAt().Add(data.TokenInfo.GetAccessExpiresIn()).Unix(),
		},
	}
	token := jwt.NewWithClaims(a.SignedMethod, claims)
	var key interface{}
	if a.isEs() {
		v, err := jwt.ParseECPrivateKeyFromPEM(a.SignedKey)
		if err != nil {
			return "", err
		}
		key = v
	} else if a.isRsOrPS() {
		v, err := jwt.ParseRSAPrivateKeyFromPEM(a.SignedKey)
		if err != nil {
			return "", err
		}
		key = v
	} else if a.isHs() {
		key = a.SignedKey
	} else {
		return "", errs.New("unsupported sign method")
	}

	access, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	return access, nil
}

func (a *JWTAccessGenerate) isEs() bool {
	return strings.HasPrefix(a.SignedMethod.Alg(), "ES")
}

func (a *JWTAccessGenerate) isRsOrPS() bool {
	isRs := strings.HasPrefix(a.SignedMethod.Alg(), "RS")
	isPs := strings.HasPrefix(a.SignedMethod.Alg(), "PS")
	return isRs || isPs
}

func (a *JWTAccessGenerate) isHs() bool {
	return strings.HasPrefix(a.SignedMethod.Alg(), "HS")
}

// Validate validates  the token
func (a *JWTAccessGenerate) Validate(tokenString string) (bool, error) {

	token, err := a.parseToken(tokenString)
	if err != nil {
		return false, err
	}

	return token.Valid, nil
}

// Parse parses a UserName from a token
func (a *JWTAccessGenerate) Parse(tokenString string) (*models.PublicUserInfo, error) {

	token, err := a.parseToken(tokenString)
	if err != nil {
		return nil, err
	}

	var userInfo *models.PublicUserInfo = new(models.PublicUserInfo)

	if claims, ok := token.Claims.(*JWTAccessClaims); ok && token.Valid {
		userInfo.Email = claims.Email
		userInfo.Name = claims.Name
		userInfo.UserName = claims.UserName
		userInfo.Role = claims.Role
		userInfo.ID = claims.UUID
		return userInfo, nil
	}
	return nil, errors.ErrInvalidAccessToken
}

func (a *JWTAccessGenerate) parseToken(tokenString string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, &JWTAccessClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validating the signing method
		if ok := token.Method.Alg() == a.SignedMethod.Alg(); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(a.SignedKey), nil
	})
}
