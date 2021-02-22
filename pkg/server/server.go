package server

import (
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	log "github.com/golang/glog"

	"github.com/mayadata-io/kubera-auth/manager/jwtmanager"
	"github.com/mayadata-io/kubera-auth/manager/loginmanager"
	"github.com/mayadata-io/kubera-auth/manager/usermanager"
	"github.com/mayadata-io/kubera-auth/pkg/errors"
	"github.com/mayadata-io/kubera-auth/pkg/generates"
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

	userStoreCfg := store.NewConfig(types.DefaultDBServerURL, types.DefaultAuthDB)
	srv := &Server{
		Config:         cfg,
		accessGenerate: generates.NewJWTAccessGenerate(jwt.SigningMethodHS512),
		GithubConfig:   oauth.NewGithubConfig(),
	}
	srv.MustUserStorage(store.NewUserStore(userStoreCfg, store.NewDefaultUserConfig()))

	return srv
}

// Server Provide authorization server
type Server struct {
	Config         *Config
	GithubConfig   oauth.SocialAuthConfig
	accessGenerate *generates.JWTAccessGenerate
	userStore      *store.UserStore
}

// MustUserStorage mandatory mapping the user store interface
func (s *Server) MustUserStorage(stor *store.UserStore, err error) {
	if err != nil {
		panic(err)
	}
	s.userStore = stor
	_, err = usermanager.CreateUser(stor, models.DefaultUser)
	if err != nil {
		log.Infoln("Unable to create default user with error:", err)
	}
}

func (s *Server) redirectError(c *gin.Context, err error) {
	data, code, _ := s.getErrorData(err)
	c.JSON(code, data)
}

func (s *Server) redirect(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

// LocalLoginRequest the local authentication request handling
func (s *Server) LocalLoginRequest(c *gin.Context, username, password string) {

	if username == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Username or password cannot be empty",
		})
		return
	}

	tokenInfo, err := loginmanager.LocalLoginUser(s.userStore, s.accessGenerate, username, password)
	if err != nil {
		s.redirectError(c, err)
		return
	}

	s.redirect(c, s.getTokenData(tokenInfo))
	return
}

//SocialLoginRequest logs in the user with github or gmail
func (s *Server) SocialLoginRequest(c *gin.Context, user *models.UserCredentials, urlString string) {

	values := url.Values{}
	tokenInfo, err := loginmanager.SocialLoginUser(s.userStore, s.accessGenerate, user)
	if err != nil {
		log.Errorln("Error logging in ", err)
		s.redirectError(c, err)
		return
	}

	values.Set("access_token", tokenInfo.GetAccess())
	c.Redirect(http.StatusFound, urlString+values.Encode())
	return
}

// LogoutRequest the authorization request handling
func (s *Server) LogoutRequest(c *gin.Context) {

	jwtUser, exists := c.Get(types.UserInfoKey)
	if !exists {
		s.redirectError(c, errors.ErrInvalidAccessToken)
		return
	}
	jwtUserInfo := jwtUser.(*models.PublicUserInfo)

	err := loginmanager.LogoutUser(s.userStore, jwtUserInfo.GetID())
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

// GetUserFromToken gets the user from token
func (s *Server) GetUserFromToken(token string) (*models.PublicUserInfo, error) {
	return jwtmanager.ParseToken(s.userStore, s.accessGenerate, token)
}

// UpdatePasswordRequest validates the request
func (s *Server) UpdatePasswordRequest(c *gin.Context, oldPassword, newPassword string) {

	jwtUser, exists := c.Get(types.UserInfoKey)
	if !exists {
		s.redirectError(c, errors.ErrInvalidAccessToken)
		return
	}
	jwtUserInfo := jwtUser.(*models.PublicUserInfo)

	if oldPassword == "" || newPassword == "" {
		c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest)
		return
	}

	updatedUserInfo, err := usermanager.UpdatePassword(s.userStore, false, oldPassword, newPassword, jwtUserInfo.GetUID())
	if err != nil {
		s.redirectError(c, err)
		return
	}
	s.redirect(c, updatedUserInfo)
	return
}

// ResetPasswordRequest validates the request
func (s *Server) ResetPasswordRequest(c *gin.Context, newPassword, userName string) {

	jwtUser, exists := c.Get(types.UserInfoKey)
	if !exists {
		s.redirectError(c, errors.ErrInvalidAccessToken)
		return
	}
	jwtUserInfo := jwtUser.(*models.PublicUserInfo)

	if userName == "" || newPassword == "" {
		c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest)
		return
	}

	var updatedUserInfo *models.PublicUserInfo
	var err error
	if jwtUserInfo.GetRole() == models.RoleAdmin {

		updatedUserInfo, err = usermanager.UpdatePassword(s.userStore, true, "", newPassword, userName)
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

	jwtUser, exists := c.Get(types.UserInfoKey)
	if !exists {
		s.redirectError(c, errors.ErrInvalidAccessToken)
		return
	}
	jwtUserInfo := jwtUser.(*models.PublicUserInfo)

	user.ID = jwtUserInfo.GetID()
	updatedUserInfo, err := usermanager.UpdateUserDetails(s.userStore, user)
	if err != nil {
		s.redirectError(c, err)
		return
	}
	s.redirect(c, updatedUserInfo)
	return
}

// CreateRequest validates the request
func (s *Server) CreateRequest(c *gin.Context, user *models.UserCredentials) {

	jwtUser, exists := c.Get(types.UserInfoKey)
	if !exists {
		s.redirectError(c, errors.ErrInvalidAccessToken)
		return
	}
	jwtUserInfo := jwtUser.(*models.PublicUserInfo)

	if user.GetUserName() == "" || user.GetPassword() == "" {
		s.redirectError(c, errors.ErrInvalidRequest)
		return
	}

	var createdUserInfo *models.PublicUserInfo
	var err error
	if jwtUserInfo.GetRole() == models.RoleAdmin {
		createdUserInfo, err = usermanager.CreateUser(s.userStore, user)
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
	users, err := usermanager.GetAllUsers(s.userStore)
	if err != nil {
		s.redirectError(c, err)
		return
	}

	s.redirect(c, users)
	return
}

//GetUserByUID gets a particular user
func (s *Server) GetUserByUID(c *gin.Context, userID string) {
	storedUser, err := usermanager.GetUserByUID(s.userStore, userID)
	if err != nil {
		s.redirectError(c, err)
		return
	}
	s.redirect(c, storedUser.GetPublicInfo())
}

//GetUserByUserName gets a particular user
func (s *Server) GetUserByUserName(c *gin.Context, userID string) {
	storedUser, err := usermanager.GetUserByUserName(s.userStore, userID)
	if err != nil {
		s.redirectError(c, err)
		return
	}
	s.redirect(c, storedUser.GetPublicInfo())
}

// SendVerificationLink sends the verification link in the desired email
func (s *Server) SendVerificationLink(c *gin.Context, email string) {
	jwtUser, exists := c.Get(types.UserInfoKey)
	if !exists {
		s.redirectError(c, errors.ErrInvalidAccessToken)
		return
	}
	jwtUserInfo := jwtUser.(*models.PublicUserInfo)

	jwtUserInfo.Email = &email
	updatedUserInfo, err := usermanager.UpdateUserDetails(s.userStore, jwtUserInfo.GetUserCredentials())
	if err != nil {
		s.redirectError(c, err)
		return
	}

	tgr := &jwtmanager.TokenGenerateRequest{
		UserInfo:       updatedUserInfo,
		AccessTokenExp: time.Minute * types.VerificationLinkExpirationTimeUnit,
	}

	tokenInfo, err := jwtmanager.GenerateAuthToken(s.accessGenerate, tgr, models.TokenVerify)
	if err != nil {
		s.redirectError(c, err)
		return
	}

	link := types.PortalURL + "/api/auth/v1/email?access=" + tokenInfo.Access

	buf, err := generates.GetEmailBody(jwtUserInfo.GetName(), link)
	if err != nil {
		log.Error("Error occurred while getting email body for user: " + jwtUserInfo.GetUID() + "error: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = generates.SendEmail(email, "Email Verification", buf.String())
	if err != nil {
		log.Error("Error occurred while sending email for user: " + jwtUserInfo.GetUID() + "error: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Verification email sent",
	})
}

// VerifyEmail marks a user email as verified
func (s *Server) VerifyEmail(c *gin.Context) {

	jwtUser, exists := c.Get(types.UserInfoKey)
	if !exists {
		s.redirectError(c, errors.ErrInvalidAccessToken)
		return
	}
	jwtUserInfo := jwtUser.(*models.PublicUserInfo)

	usermanager.UpdateUserDetails(s.userStore, jwtUserInfo.GetUserCredentials())

	tgr := &jwtmanager.TokenGenerateRequest{
		UserInfo:       jwtUserInfo,
		AccessTokenExp: time.Minute * 10,
	}

	ti, err := jwtmanager.GenerateAuthToken(s.accessGenerate, tgr, models.TokenVerify)
	if err != nil {
		s.redirectError(c, err)
		return
	}

	s.redirect(c, s.getTokenData(ti))
}
