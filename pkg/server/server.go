package server

import (
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	log "github.com/golang/glog"
	"github.com/imdario/mergo"

	"github.com/mayadata-io/kubera-auth/manager/emailmanager"
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
		GoogleConfig:   oauth.NewGoogleConfig(),
	}
	srv.MustUserStorage(store.NewUserStore(userStoreCfg, store.NewDefaultUserConfig()))

	return srv
}

// Server Provide authorization server
type Server struct {
	Config         *Config
	GithubConfig   oauth.SocialAuthConfig
	GoogleConfig   oauth.SocialAuthConfig
	accessGenerate *generates.JWTAccessGenerate
	userStore      *store.UserStore
}

// MustUserStorage mandatory mapping the user store interface
func (s *Server) MustUserStorage(stor *store.UserStore, err error) {
	if err != nil {
		panic(err)
	}
	s.userStore = stor
	_, err = usermanager.CreateUser(stor, models.DefaultUser, false)
	if err != nil {
		log.Infoln("Unable to create default user with error:", err)
	}
}

func (s *Server) errorResponse(c *gin.Context, err error) {
	data, code, _ := s.getErrorData(err)
	c.JSON(code, data)
}

func (s *Server) successResponse(c *gin.Context, data interface{}) {
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
		s.errorResponse(c, err)
		return
	}
	s.successResponse(c, s.getTokenData(tokenInfo))
}

// SocialLoginRequest logs in the user with github or gmail
func (s *Server) SocialLoginRequest(c *gin.Context, user *models.UserCredentials, urlString string) {
	values := url.Values{}
	tokenInfo, err := loginmanager.SocialLoginUser(s.userStore, s.accessGenerate, user)
	if err != nil {
		log.Errorln("Error logging in ", err)
		s.errorResponse(c, err)
		return
	}

	values.Set("access_token", tokenInfo.GetAccess())
	c.Redirect(http.StatusFound, urlString+values.Encode())
}

// LogoutRequest the authorization request handling
func (s *Server) LogoutRequest(c *gin.Context) {
	jwtUser, exists := c.Get(types.JWTUserCredentialsKey)
	if !exists {
		s.errorResponse(c, errors.ErrInvalidAccessToken)
		return
	}
	jwtUserCredentials := jwtUser.(*models.UserCredentials)

	err := loginmanager.LogoutUser(s.userStore, jwtUserCredentials.ID)
	if err != nil {
		s.errorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "LoggedOut successfully",
	})
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
// nolint: cyclop
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
func (s *Server) GetUserFromToken(token string) (*models.UserCredentials, error) {
	return jwtmanager.ParseToken(s.userStore, s.accessGenerate, token)
}

// UpdatePasswordRequest validates the request
func (s *Server) UpdatePasswordRequest(c *gin.Context, newPassword string) {
	jwtUser, exists := c.Get(types.JWTUserCredentialsKey)
	if !exists {
		s.errorResponse(c, errors.ErrInvalidAccessToken)
		return
	}
	jwtUserCredentials := jwtUser.(*models.UserCredentials)

	if newPassword == "" {
		c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest)
		return
	}

	updatedUserInfo, err := usermanager.UpdatePassword(s.userStore, newPassword, jwtUserCredentials.UserName)
	if err != nil {
		s.errorResponse(c, err)
		return
	}
	s.successResponse(c, updatedUserInfo)
}

// ResetPasswordRequest validates the request
func (s *Server) ResetPasswordRequest(c *gin.Context, newPassword, userName string) {
	jwtUser, exists := c.Get(types.JWTUserCredentialsKey)
	if !exists {
		s.errorResponse(c, errors.ErrInvalidAccessToken)
		return
	}
	jwtUserCredentials := jwtUser.(*models.UserCredentials)

	if userName == "" || newPassword == "" {
		c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest)
		return
	}

	var updatedUserInfo *models.PublicUserInfo
	var err error
	if jwtUserCredentials.Role == models.RoleAdmin {
		updatedUserInfo, err = usermanager.UpdatePassword(s.userStore, newPassword, userName)
		if err != nil {
			s.errorResponse(c, err)
			return
		}
	} else {
		s.errorResponse(c, errors.ErrInvalidUser)
		return
	}
	s.successResponse(c, updatedUserInfo)
}

// UpdateUserDetailsRequest validates the request
func (s *Server) UpdateUserDetailsRequest(c *gin.Context, user *models.UserCredentials) {
	jwtUser, exists := c.Get(types.JWTUserCredentialsKey)
	if !exists {
		s.errorResponse(c, errors.ErrInvalidAccessToken)
		return
	}
	jwtUserCredentials := jwtUser.(*models.UserCredentials)

	// It will override the `jwtUserCredentials` with values filled in `user` and preserve the other values of `storedUser`
	err := mergo.Merge(jwtUserCredentials, user, mergo.WithOverride)
	if err != nil {
		s.errorResponse(c, err)
	}

	updatedUserInfo, err := usermanager.UpdateUserDetails(s.userStore, jwtUserCredentials)
	if err != nil {
		s.errorResponse(c, err)
		return
	}
	s.successResponse(c, updatedUserInfo)
}

// CreateRequest validates the request
func (s *Server) CreateRequest(c *gin.Context, user *models.UserCredentials) {
	jwtUser, exists := c.Get(types.JWTUserCredentialsKey)
	if !exists {
		s.errorResponse(c, errors.ErrInvalidAccessToken)
		return
	}
	jwtUserCredentials := jwtUser.(*models.UserCredentials)

	if user.UserName == "" || user.Password == "" {
		s.errorResponse(c, errors.ErrInvalidRequest)
		return
	}

	var createdUserInfo *models.PublicUserInfo
	var err error
	if jwtUserCredentials.Role == models.RoleAdmin {
		createdUserInfo, err = usermanager.CreateUser(s.userStore, user, false)
		if err != nil {
			s.errorResponse(c, err)
			return
		}
		s.successResponse(c, createdUserInfo)
		return
	}
	s.errorResponse(c, errors.ErrInvalidUser)
}

// SelfSignupUser lets a user to signup into kubera by filling a signup form
func (s *Server) SelfSignupUser(c *gin.Context, user *models.UserCredentials) {
	if user.Password == "" || user.UserName == "" {
		s.errorResponse(c, errors.ErrInvalidRequest)
		return
	}

	createdUserInfo, err := usermanager.CreateUser(s.userStore, user, true)
	if err != nil {
		s.errorResponse(c, err)
		return
	}

	err = emailmanager.SendEmail(s.accessGenerate, createdUserInfo, emailmanager.VerificationEmail)
	if err != nil {
		s.errorResponse(c, err)
		return
	}

	tokenInfo, err := loginmanager.LocalLoginUser(s.userStore, s.accessGenerate, user.UserName, user.Password)
	if err != nil {
		s.errorResponse(c, err)
		return
	}
	s.successResponse(c, s.getTokenData(tokenInfo))
}

// GetUsersRequest gets all the users
func (s *Server) GetUsersRequest(c *gin.Context) {
	users, err := usermanager.GetAllUsers(s.userStore)
	if err != nil {
		s.errorResponse(c, err)
		return
	}
	s.successResponse(c, users)
}

//GetUserByUID gets a particular user
func (s *Server) GetUserByUID(c *gin.Context, userID string) {
	storedUser, err := usermanager.GetUserByUID(s.userStore, userID)
	if err != nil {
		s.errorResponse(c, err)
		return
	}
	s.successResponse(c, storedUser.GetPublicInfo())
}

//GetUserByUserName gets a particular user
func (s *Server) GetUserByUserName(c *gin.Context, userID string) {
	storedUser, err := usermanager.GetUserByUserName(s.userStore, userID)
	if err != nil {
		s.errorResponse(c, err)
		return
	}
	s.successResponse(c, storedUser.GetPublicInfo())
}

// SendVerificationLink sends the verification link in the desired email
func (s *Server) SendVerificationLink(c *gin.Context, resend bool, unverifiedEmail string) {
	jwtUser, exists := c.Get(types.JWTUserCredentialsKey)
	if !exists {
		s.errorResponse(c, errors.ErrInvalidAccessToken)
		return
	}
	jwtUserCredentials := jwtUser.(*models.UserCredentials)

	userWithSameEmail, err := usermanager.GetUser(s.userStore, bson.M{"email": unverifiedEmail})
	if err == nil && userWithSameEmail != nil {
		s.errorResponse(c, errors.ErrUserExists)
		return
	} else if err != errors.ErrInvalidUser {
		s.errorResponse(c, err)
		return
	}

	var userInfo *models.PublicUserInfo
	if resend {
		userInfo = jwtUserCredentials.GetPublicInfo()
	} else {
		jwtUserCredentials.UnverifiedEmail = unverifiedEmail
		userInfo, err = usermanager.UpdateUserDetails(s.userStore, jwtUserCredentials)
		if err != nil {
			s.errorResponse(c, err)
			return
		}
	}

	err = emailmanager.SendEmail(s.accessGenerate, userInfo, emailmanager.VerificationEmail)
	if err != nil {
		s.errorResponse(c, err)
		return
	}

	s.successResponse(c, userInfo)
}

// VerifyEmail marks a user email as verified
func (s *Server) VerifyEmail(c *gin.Context, redirectURL string) {
	jwtUser, exists := c.Get(types.JWTUserCredentialsKey)
	if !exists {
		s.errorResponse(c, errors.ErrInvalidAccessToken)
		// Redirecting user to UI if the user is not authorized
		c.Redirect(http.StatusPermanentRedirect, redirectURL)
		return
	}
	jwtUserCredentials := jwtUser.(*models.UserCredentials)

	if jwtUserCredentials.UnverifiedEmail != "" {
		userWithSameEmail, err := usermanager.GetUser(s.userStore, bson.M{"email": jwtUserCredentials.UnverifiedEmail})
		if err == nil && userWithSameEmail != nil {
			log.Errorln("Email already in use with another user for user uid: ", jwtUserCredentials.UID)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Email already in use with another user",
			})
			c.Redirect(http.StatusPermanentRedirect, redirectURL)
			return
		} else if err != errors.ErrInvalidUser {
			s.errorResponse(c, err)
			// Redirecting user to UI if email found in field `email` is already present for some user
			c.Redirect(http.StatusPermanentRedirect, redirectURL)
			return
		}

		jwtUserCredentials.Email = jwtUserCredentials.UnverifiedEmail
		jwtUserCredentials.UnverifiedEmail = ""
		if jwtUserCredentials.Kind == models.LocalAuth {
			jwtUserCredentials.UserName = jwtUserCredentials.Email
		}
	} else {
		log.Errorln("No email found to be verified for user uid: ", jwtUserCredentials.UID)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No email found to be verified",
		})
		// Redirecting user to UI if no email found in field `unverified_email`
		c.Redirect(http.StatusPermanentRedirect, redirectURL)
		return
	}

	_, err := usermanager.UpdateUserDetails(s.userStore, jwtUserCredentials)
	if err != nil {
		s.errorResponse(c, err)
		// Redirecting user to UI if updating the database fails
		c.Redirect(http.StatusPermanentRedirect, redirectURL)
		return
	}

	log.Infoln("Email: ", jwtUserCredentials.Email, " is verified successfully for user uid: ", jwtUserCredentials.UID)
	c.Redirect(http.StatusPermanentRedirect, redirectURL)
}

// RestoreEmail makes the `UnverifiedEmail` as blank.
// A user might want to change his email and later changes his mind and
// wants to continue with his previous email only. Since change of email in DB
// is withhold until the user verifies his new email. So this just removes
// his unverified email from DB and let the user continue with the verified one.
func (s *Server) RestoreEmail(c *gin.Context) {
	jwtUser, exists := c.Get(types.JWTUserCredentialsKey)
	if !exists {
		s.errorResponse(c, errors.ErrInvalidAccessToken)
		return
	}
	jwtUserCredentials := jwtUser.(*models.UserCredentials)

	jwtUserCredentials.UnverifiedEmail = ""
	userInfo, err := usermanager.UpdateUserDetails(s.userStore, jwtUserCredentials)
	if err != nil {
		s.errorResponse(c, err)
		return
	}

	s.successResponse(c, userInfo)
}

// ForgotPasswordRequest validates the request
func (s *Server) ForgotPasswordRequest(c *gin.Context, email string) {
	storedUser, err := usermanager.GetUserByUserName(s.userStore, email)
	if err != nil {
		s.errorResponse(c, err)
		return
	}

	err = emailmanager.SendEmail(s.accessGenerate, storedUser.GetPublicInfo(), emailmanager.ResetPasswordEmail)
	if err != nil {
		s.errorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset email sent",
	})
}
