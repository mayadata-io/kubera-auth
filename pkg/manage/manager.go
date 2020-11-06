package manage

import (
	"strconv"
	"strings"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	log "github.com/golang/glog"
	"golang.org/x/crypto/bcrypt"

	"github.com/mayadata-io/kubera-auth/pkg/errors"
	"github.com/mayadata-io/kubera-auth/pkg/generates"
	"github.com/mayadata-io/kubera-auth/pkg/models"
	"github.com/mayadata-io/kubera-auth/pkg/store"
	"github.com/mayadata-io/kubera-auth/pkg/types"
	"github.com/mayadata-io/kubera-auth/pkg/utils/random"
	"github.com/mayadata-io/kubera-auth/pkg/utils/uuid"
)

// NewManager create to authorization management instance
func NewManager() *Manager {
	return &Manager{}
}

// Manager provide authorization management
type Manager struct {
	accessGenerate *generates.JWTAccessGenerate
	userStore      *store.UserStore
}

// MapAccessGenerate mapping the access token generate interface
func (m *Manager) MapAccessGenerate(gen *generates.JWTAccessGenerate) {
	m.accessGenerate = gen
}

// MustUserStorage mandatory mapping the user store interface
func (m *Manager) MustUserStorage(stor *store.UserStore, err error) {
	if err != nil {
		panic(err)
	}
	m.userStore = stor
	_, err = m.CreateUser(models.DefaultUser)
	if err != nil {
		log.Infoln("Unable to create default user with error:", err)
	}
}

// GetUser get the user information
func (m *Manager) GetUser(userName string) (user *models.UserCredentials, err error) {
	query := bson.M{"username": userName}
	user, err = m.userStore.GetUser(query)
	if err != nil && err == mgo.ErrNotFound {
		err = errors.ErrInvalidUser
	}
	return
}

// GetAllUsers get the user information
func (m *Manager) GetAllUsers() ([]*models.PublicUserInfo, error) {
	users, err := m.userStore.GetAllUsers()
	if err != nil {
		return nil, err
	}

	var allUsers []*models.PublicUserInfo
	for _, user := range users {
		allUsers = append(allUsers, user.GetPublicInfo())
	}
	return allUsers, nil
}

// CheckUserExists get the user information
func (m *Manager) CheckUserExists(user *models.UserCredentials) (bool, error) {

	exists := true
	_, err := m.userStore.GetUser(bson.M{"username": user.GetUserName()})
	if err != nil && err == mgo.ErrNotFound {
		exists = false
	} else if err != nil {
		return false, err
	}

	if user.GetEmail() != nil && exists == false {
		_, err := m.userStore.GetUser(bson.M{"email": *user.Email})
		if err != nil && err == mgo.ErrNotFound {
			exists = false
		} else if err != nil {
			return false, err
		} else {
			exists = true
		}
	}
	return exists, nil
}

// VerifyUserPassword verifies user password
func (m *Manager) VerifyUserPassword(username, password string) (*models.PublicUserInfo, error) {
	user, err := m.userStore.GetUser(bson.M{"username": username, "kind": models.LocalAuth})
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.GetPassword()), []byte(password))
	if err != nil {
		return user.GetPublicInfo(), errors.ErrInvalidPassword
	}
	return user.GetPublicInfo(), nil
}

// LocalLoginUser verifies user password
func (m *Manager) LocalLoginUser(username string) error {
	storedUser, err := m.userStore.GetUser(bson.M{"username": username, "kind": models.LocalAuth})
	if err != nil {
		return err
	}
	storedUser.LoggedIn = true
	return m.userStore.UpdateUser(storedUser)
}

// SocialLoginUser get the user information
func (m *Manager) SocialLoginUser(user *models.UserCredentials) (*models.UserCredentials, error) {

	query := bson.M{"social_auth_id": user.SocialAuthID}
	storedUser, err := m.userStore.GetUser(query)
	if err != nil && err == mgo.ErrNotFound {
		createErr := m.CreateSocialUser(user)
		if createErr != nil {
			return nil, createErr
		}
		return user, nil
	} else if err != nil {
		return nil, err
	}

	storedUser.Email = user.Email
	storedUser.Name = user.Name
	storedUser.LoggedIn = true

	return storedUser, m.userStore.UpdateUser(storedUser)
}

//CreateSocialUser ...
func (m *Manager) CreateSocialUser(user *models.UserCredentials) error {

	query := bson.M{"email": user.Email, "kind": models.LocalAuth}
	storedUser, err := m.userStore.GetUser(query)
	if err != nil && err == mgo.ErrNotFound {
		user.UserName = generateUserName(user.Name)
		user.UID = uuid.Must(uuid.NewRandom()).String()
	} else if err != nil {
		return err
	} else {
		user.UserName = storedUser.UserName
		user.UID = storedUser.UID
	}

	return m.userStore.Set(user)
}

func generateUserName(name string) string {
	names := strings.Split(name, " ")
	fname := names[0]
	var lname string
	if len(names) > 1 {
		lname = names[1]
	} else {
		lname = names[0]
	}

	appendString := random.GetRandomNumbers(5)
	choose, err := strconv.Atoi(random.GetRandomNumbers(1))

	if err != nil {
		choose = 0
	}
	if choose < 5 {
		return fname + appendString
	}

	return lname + appendString
}

// LogoutUser verifies user password
func (m *Manager) LogoutUser(id bson.ObjectId) error {
	storedUser, err := m.userStore.GetUserByID(id)
	if err != nil {
		return err
	}
	storedUser.LoggedIn = false
	return m.userStore.UpdateUser(storedUser)
}

// CreateUser get the user information
func (m *Manager) CreateUser(user *models.UserCredentials) (*models.PublicUserInfo, error) {

	exists, err := m.CheckUserExists(user)
	if err != nil {
		return nil, err
	} else if exists == true {
		return nil, errors.ErrUserExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), types.PasswordEncryptionCost)
	if err != nil {
		return nil, err
	}

	user.Password = string(hashedPassword)
	user.UID = uuid.Must(uuid.NewRandom()).String()
	err = m.userStore.Set(user)
	return user.GetPublicInfo(), err
}

// GenerateAuthToken generate the authorization token(code)
func (m *Manager) GenerateAuthToken(tgr *TokenGenerateRequest) (*models.Token, error) {

	ti := models.NewToken()

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

	tv, err := m.accessGenerate.Token(td)
	if err != nil {
		return nil, err
	}
	ti.SetAccess(tv)
	return ti, nil
}

// ValidateToken validates the token
func (m *Manager) ValidateToken(tokenString string) (valid bool, err error) {
	valid, err = m.accessGenerate.Validate(tokenString)
	return
}

// ParseToken validates the token
func (m *Manager) ParseToken(tokenString string) (userInfo *models.PublicUserInfo, err error) {
	userInfo, err = m.accessGenerate.Parse(tokenString)
	return
}

// UpdateUserDetails get the user information
func (m *Manager) UpdateUserDetails(user *models.UserCredentials) (*models.PublicUserInfo, error) {

	storedUser, err := m.userStore.GetUserByID(user.GetID())
	if err != nil {
		return nil, errors.ErrInvalidUser
	}

	if user.GetEmail() != nil {
		storedUser.Email = user.GetEmail()
	}
	if user.GetName() != "" {
		storedUser.Name = user.GetName()
	}
	if user.GetUserName() != "" {
		storedUser.UserName = user.GetUserName()
	}

	exists, err := m.CheckUserExists(storedUser)
	if err != nil {
		return nil, err
	} else if exists == true {
		return nil, errors.ErrUserExists
	}

	err = m.userStore.UpdateUser(storedUser)
	return storedUser.GetPublicInfo(), err
}

// UpdatePassword get the user information
func (m *Manager) UpdatePassword(reset bool, oldPassword, newPassword, userID string) (*models.PublicUserInfo, error) {

	var storedUser *models.UserCredentials
	var err error
	if reset == true {
		storedUser, err = m.userStore.GetUser(bson.M{"username": userID, "kind": models.LocalAuth})
		if err != nil {
			return nil, err
		}
	} else {
		storedUser, err = m.userStore.GetUser(bson.M{"uid": userID, "kind": models.LocalAuth})
		if err != nil {
			return nil, err
		}

		err = bcrypt.CompareHashAndPassword([]byte(storedUser.GetPassword()), []byte(oldPassword))
		if err != nil {
			return storedUser.GetPublicInfo(), errors.ErrInvalidPassword
		}
		// Updating the state when the user changes the password himself. Will be useful for the first time
		storedUser.State = models.StateActive
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), types.PasswordEncryptionCost)
	if err != nil {
		return nil, err
	}
	storedUser.Password = string(hashedPassword)

	err = m.userStore.UpdateUser(storedUser)
	return storedUser.GetPublicInfo(), err
}
