package usermanager

import (
	"github.com/globalsign/mgo/bson"
	"golang.org/x/crypto/bcrypt"

	"github.com/mayadata-io/kubera-auth/pkg/errors"
	"github.com/mayadata-io/kubera-auth/pkg/models"
	"github.com/mayadata-io/kubera-auth/pkg/store"
	"github.com/mayadata-io/kubera-auth/pkg/types"
)

// UpdateUserDetails updates the user information
func UpdateUserDetails(userStore *store.UserStore, user *models.UserCredentials) (*models.PublicUserInfo, error) {
	// There will be following possible transitions of OnboardingState
	// BoardingStateSignup -> BoardingStateEmailVerified -> BoardingStateVerifiedAndComplete
	// BoardingStateSignup -> BoardingStateUnverifiedAndComplete -> BoardingStateVerifiedAndComplete
	switch user.OnBoardingState {
	case models.BoardingStateSignup:
		{
			if user.Email != nil {
				user.OnBoardingState = models.BoardingStateEmailVerified
			} else if user.Company != nil {
				user.OnBoardingState = models.BoardingStateUnverifiedAndComplete
			}
		}
	case models.BoardingStateEmailVerified:
		{
			if user.Company != nil {
				user.OnBoardingState = models.BoardingStateVerifiedAndComplete
			}
		}
	case models.BoardingStateUnverifiedAndComplete:
		{
			if user.Email != nil {
				user.OnBoardingState = models.BoardingStateVerifiedAndComplete
			}
		}
	}

	err := userStore.UpdateUser(user)
	return user.GetPublicInfo(), err
}

// UpdatePassword sets the user password
func UpdatePassword(userStore *store.UserStore, reset bool, oldPassword, newPassword, userID string) (*models.PublicUserInfo, error) {
	var storedUser *models.UserCredentials
	var err error
	if reset {
		storedUser, err = userStore.GetUser(bson.M{"username": userID, "kind": models.LocalAuth})
		if err != nil {
			return nil, err
		}
	} else {
		storedUser, err = userStore.GetUser(bson.M{"uid": userID, "kind": models.LocalAuth})
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

	password := string(hashedPassword)
	storedUser.Password = &password

	err = userStore.UpdateUser(storedUser)
	return storedUser.GetPublicInfo(), err
}
