package server

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	log "github.com/golang/glog"

	"github.com/mayadata-io/kubera-auth/pkg/errors"
	"github.com/mayadata-io/kubera-auth/pkg/generates"
	"github.com/mayadata-io/kubera-auth/pkg/manage"
	"github.com/mayadata-io/kubera-auth/pkg/models"
	"github.com/mayadata-io/kubera-auth/pkg/oauth"
	"github.com/mayadata-io/kubera-auth/pkg/store"
	"github.com/mayadata-io/kubera-auth/pkg/types"
)

func init() {
	if os.Getenv("DB_SERVER") == "" {
		log.Fatal("Environment variables JWT_SECRET or DB_SERVER are not set")
	}
}

// NewServer create authorization server
func NewServer(cfg *Config) *Server {

	manager := manage.NewManager()

	userStoreCfg := store.NewConfig(types.DefaultDBServerURL, types.DefaultAuthDB)

	manager.MustUserStorage(store.NewUserStore(userStoreCfg, store.NewDefaultUserConfig()))

	manager.MapAccessGenerate(generates.NewJWTAccessGenerate(jwt.SigningMethodHS512))

	srv := &Server{
		Config:       cfg,
		Manager:      manager,
		GithubConfig: oauth.NewGithubConfig(),
	}

	return srv
}

// Server Provide authorization server
type Server struct {
	Config       *Config
	Manager      *manage.Manager
	GithubConfig oauth.SocialAuthConfig
}

func (s *Server) redirectError(c *gin.Context, err error) {
	data, code, _ := s.getErrorData(err)
	c.JSON(code, data)
}

func (s *Server) redirect(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

// ValidationAuthenticateRequest the authenticate request validation
func (s *Server) validationAuthenticateRequest(username, password string) (*manage.TokenGenerateRequest, error) {
	if username == "" || password == "" {
		return nil, errors.ErrInvalidRequest
	}

	userInfo, err := s.Manager.VerifyUserPassword(username, password)
	if err != nil {
		return nil, err
	}

	req := &manage.TokenGenerateRequest{
		UserInfo: userInfo,
	}
	return req, nil
}

// LocalLoginRequest the local authentication request handling
func (s *Server) LocalLoginRequest(c *gin.Context, username, password string) {

	tgr, err := s.validationAuthenticateRequest(username, password)
	if err != nil {
		s.redirectError(c, err)
		return
	}

	ti, err := s.Manager.GenerateAuthToken(tgr)
	if err != nil {
		s.redirectError(c, err)
		return
	}

	err = s.Manager.LocalLoginUser(username)
	if err != nil {
		s.redirectError(c, err)
		return
	}

	s.redirect(c, s.getTokenData(ti))
	return
}

//SocialLoginRequest logs in the user with github or gmail
func (s *Server) SocialLoginRequest(c *gin.Context, user *models.UserCredentials) (*models.Token, error) {

	storedUser, err := s.Manager.SocialLoginUser(user)
	if err != nil {
		log.Errorln("Error logging in ", err)
		return nil, err
	}

	tgr := &manage.TokenGenerateRequest{
		UserInfo: storedUser.GetPublicInfo(),
	}

	return s.Manager.GenerateAuthToken(tgr)
}

// LogoutRequest the authorization request handling
func (s *Server) LogoutRequest(c *gin.Context) {

	userInfo, err := s.getUserFromToken(c.Request)
	if err != nil {
		s.redirectError(c, err)
		return
	}

	err = s.Manager.LogoutUser(userInfo.GetID())
	if err != nil {
		s.redirectError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "LoggedOut successfully",
	})
	return
}

// GetTokenData token data
func (s *Server) getTokenData(ti *models.Token) map[string]interface{} {
	data := map[string]interface{}{
		"access_token": ti.GetAccess(),
		"token_type":   s.Config.TokenType,
		"expires_in":   int64(ti.GetAccessExpiresIn() / time.Second),
	}
	return data
}

// GetErrorData get error response data
func (s *Server) getErrorData(err error) (map[string]interface{}, int, http.Header) {
	var re errors.Response
	if v, ok := errors.Descriptions[err]; ok {
		re.Error = err
		re.Description = v
		re.StatusCode = errors.StatusCodes[err]
	} else {
		if fn := s.internalErrorHandler; fn != nil {
			if v := fn(err); v != nil {
				re = *v
			}
		}

		if re.Error == nil {
			re.Error = errors.ErrServerError
			re.Description = errors.Descriptions[errors.ErrServerError]
			re.StatusCode = errors.StatusCodes[errors.ErrServerError]
		}
	}

	if fn := s.responseErrorHandler; fn != nil {
		fn(&re)
	}

	data := make(map[string]interface{})
	if err := re.Error; err != nil {
		data["error"] = err.Error()
	}

	if v := re.ErrorCode; v != 0 {
		data["error_code"] = v
	}

	if v := re.Description; v != "" {
		data["error_description"] = v
	}

	if v := re.URI; v != "" {
		data["error_uri"] = v
	}

	statusCode := http.StatusInternalServerError
	if v := re.StatusCode; v > 0 {
		statusCode = v
	}

	return data, statusCode, re.Header
}

func (s *Server) internalErrorHandler(err error) (re *errors.Response) {
	log.Infoln("Internal Error:", err.Error())
	return
}

func (s *Server) responseErrorHandler(re *errors.Response) {
	log.Infoln("Response Error:", re.Error.Error())
}

func (s *Server) getTokenFromHeader(r *http.Request) (string, error) {
	auth := r.Header.Get("Authorization")
	prefix := "Bearer "
	token := ""

	if auth != "" && strings.HasPrefix(auth, prefix) {
		token = auth[len(prefix):]
	}

	if token == "" {
		return "", errors.ErrInvalidAccessToken
	}

	return token, nil
}

func (s *Server) getUserFromToken(r *http.Request) (*models.PublicUserInfo, error) {
	tokenString, err := s.getTokenFromHeader(r)
	if err != nil {
		return nil, err
	}

	userInfo, err := s.Manager.ParseToken(tokenString)
	return userInfo, err
}

// UpdatePasswordRequest validates the request
func (s *Server) UpdatePasswordRequest(c *gin.Context, oldPassword, newPassword string) {

	userInfo, err := s.getUserFromToken(c.Request)
	if err != nil {
		s.redirectError(c, err)
		return
	}

	if oldPassword == "" || newPassword == "" {
		c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest)
		return
	}

	updatedUserInfo, err := s.Manager.UpdatePassword(false, oldPassword, newPassword, userInfo.GetUID())
	if err != nil {
		s.redirectError(c, err)
		return
	}
	s.redirect(c, updatedUserInfo)
	return
}

// ResetPasswordRequest validates the request
func (s *Server) ResetPasswordRequest(c *gin.Context, newPassword, userName string) {

	userInfo, err := s.getUserFromToken(c.Request)
	if err != nil {
		s.redirectError(c, err)
		return
	}

	if userName == "" || newPassword == "" {
		c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest)
		return
	}

	var updatedUserInfo *models.PublicUserInfo
	if userInfo.GetRole() == models.RoleAdmin {

		updatedUserInfo, err = s.Manager.UpdatePassword(true, "", newPassword, userName)
		if err != nil {
			s.redirectError(c, err)
			return
		}
	}
	s.redirect(c, updatedUserInfo)
	return
}

// UpdateUserDetailsRequest validates the request
func (s *Server) UpdateUserDetailsRequest(c *gin.Context, user *models.UserCredentials) {

	userInfo, err := s.getUserFromToken(c.Request)
	if err != nil {
		s.redirectError(c, err)
		return
	}

	user.ID = userInfo.GetID()
	updatedUserInfo, err := s.Manager.UpdateUserDetails(user)
	if err != nil {
		s.redirectError(c, err)
		return
	}
	s.redirect(c, updatedUserInfo)
	return
}

// CreateRequest validates the request
func (s *Server) CreateRequest(c *gin.Context, user *models.UserCredentials) {

	userInfo, err := s.getUserFromToken(c.Request)
	if err != nil {
		s.redirectError(c, err)
		return
	}

	if user.GetUserName() == "" || user.GetPassword() == "" {
		s.redirectError(c, errors.ErrInvalidRequest)
		return
	}

	var createdUserInfo *models.PublicUserInfo
	if userInfo.GetRole() == models.RoleAdmin {
		createdUserInfo, err = s.Manager.CreateUser(user)
		if err != nil {
			s.redirectError(c, err)
			return
		}
		s.redirect(c, createdUserInfo)
		return
	}
	s.redirectError(c, errors.ErrInvalidUser)
	return
}

// GetUsersRequest gets all the users
func (s *Server) GetUsersRequest(c *gin.Context) {

	_, err := s.getUserFromToken(c.Request)
	if err != nil {
		s.redirectError(c, err)
		return
	}

	users, err := s.Manager.GetAllUsers()
	if err != nil {
		s.redirectError(c, err)
		return
	}

	s.redirect(c, users)
	return
}

//GetUserRequest gets a particular user
func (s *Server) GetUserRequest(c *gin.Context, username string) {

	_, err := s.getUserFromToken(c.Request)
	if err != nil {
		s.redirectError(c, err)
		return
	}

	storedUser, err := s.Manager.GetUser(username)
	if err != nil {
		s.redirectError(c, err)
		return
	}
	s.redirect(c, storedUser.GetPublicInfo())
}
