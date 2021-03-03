package usermanager

import (
	"strconv"
	"strings"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/mayadata-io/kubera-auth/pkg/models"
	"github.com/mayadata-io/kubera-auth/pkg/store"
	"github.com/mayadata-io/kubera-auth/pkg/utils/random"
)

// IsUserExists get the user information
func IsUserExists(userStore *store.UserStore, user *models.UserCredentials) (bool, error) {
	exists := true
	_, err := userStore.GetUser(bson.M{"username": user.GetUserName()})
	if err != nil && err == mgo.ErrNotFound {
		exists = false
	} else if err != nil {
		return false, err
	}

	if !exists && user.Email != nil {
		_, err := userStore.GetUser(bson.M{"email": *user.Email})
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

func generateUserName(name string) *string {
	var username string
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
		username = fname + appendString
	}
	username = lname + appendString

	return &username
}
