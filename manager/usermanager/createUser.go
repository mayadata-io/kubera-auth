package usermanager

import (
	"github.com/globalsign/mgo/bson"
	"golang.org/x/crypto/bcrypt"

	"github.com/mayadata-io/kubera-auth/pkg/errors"
	"github.com/mayadata-io/kubera-auth/pkg/models"
	"github.com/mayadata-io/kubera-auth/pkg/store"
	"github.com/mayadata-io/kubera-auth/pkg/types"
	"github.com/mayadata-io/kubera-auth/pkg/utils/uuid"
)

// CreateUser builds a user entry from the provided details about the user
// such as username, password etc. for insertion. These values are embedded
// inside the usercredentials struct.
// `isSignup` is a bool value used to detect whether this user creation is being
// done via a local auth signup form or through an admin and will accordingly set
// the values for the user to be created.
func CreateUser(userStore *store.UserStore, user *models.UserCredentials, isSignup bool) (*models.PublicUserInfo, error) {
	exists, err := IsUserExists(userStore, user)
	if err != nil {
		return nil, err
	} else if exists {
		return nil, errors.ErrUserExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), types.PasswordEncryptionCost)
	if err != nil {
		return nil, err
	}

	var newUser *models.UserCredentials
	if isSignup {
		newUser = &models.UserCredentials{
			UID:             uuid.Must(uuid.NewRandom()).String(),
			UserName:        user.UserName,
			Password:        string(hashedPassword),
			Name:            user.Name,
			UnverifiedEmail: user.UserName,
			Kind:            models.LocalAuth,
			Role:            models.RoleUser,
			State:           models.StateCreated,
			OnBoardingState: models.BoardingStateSignup,
		}
	} else {
		newUser = &models.UserCredentials{
			UID:             uuid.Must(uuid.NewRandom()).String(),
			UserName:        user.UserName,
			Password:        string(hashedPassword),
			Name:            user.Name,
			UnverifiedEmail: user.UnverifiedEmail,
			Kind:            models.LocalAuth,
			State:           models.StateCreated,
			OnBoardingState: models.BoardingStateUnverifiedAndComplete,
		}

		if user.Role != "" {
			newUser.Role = user.Role
		} else {
			newUser.Role = models.RoleUser
		}
	}

	err = userStore.Set(newUser)
	return newUser.GetPublicInfo(), err
}

// CreateSocialUser creates a user if the user opts logging in with some oauth
func CreateSocialUser(userStore *store.UserStore, user *models.UserCredentials) error {
	userWithSameEmail, err := GetUser(userStore, bson.M{"email": user.Email})
	if err == nil && userWithSameEmail != nil {
		return errors.ErrUserExists
	} else if err != errors.ErrInvalidUser {
		return err
	}

	for {
		user.UserName = generateUserName(user.Name)
		_, err = GetUserByUserName(userStore, user.UserName)
		if err == errors.ErrInvalidUser {
			break
		}
	}
	user.UID = uuid.Must(uuid.NewRandom()).String()

	return userStore.Set(user)
}
