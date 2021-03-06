package store

import (
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	log "github.com/golang/glog"

	"github.com/mayadata-io/kubera-auth/pkg/models"
	"github.com/mayadata-io/kubera-auth/pkg/types"
)

// UserConfig user configuration parameters
type UserConfig struct {
	// store users data collection name(The default is user)
	UsersCName string
}

// NewDefaultUserConfig create a default user configuration
func NewDefaultUserConfig() *UserConfig {
	return &UserConfig{
		UsersCName: types.DefaultLocalAuthCollection,
	}
}

// NewUserStore create a user store instance based on mongodb
func NewUserStore(cfg *Config, ucfgs ...*UserConfig) (*UserStore, error) {
	session, err := mgo.Dial(cfg.URL)
	if err != nil {
		return nil, err
	}

	if types.DBUser != "" && types.DBPassword != "" {
		cred := mgo.Credential{
			Username: types.DBUser,
			Password: types.DBPassword,
		}
		err = session.Login(&cred)
		if err != nil {
			log.Errorln("Error connecting database error", err)
			return nil, err
		}
	}

	return NewUserStoreWithSession(session, cfg.DB, ucfgs...)
}

// NewUserStoreWithSession create a User store instance based on mongodb
func NewUserStoreWithSession(session *mgo.Session, dbName string, ucfgs ...*UserConfig) (*UserStore, error) {
	us := &UserStore{
		dbName:  dbName,
		session: session,
		ucfg:    NewDefaultUserConfig(),
	}
	if len(ucfgs) > 0 {
		us.ucfg = ucfgs[0]
	}

	return us, nil
}

// UserStore MongoDB storage for OAuth 2.0
type UserStore struct {
	ucfg    *UserConfig
	dbName  string
	session *mgo.Session
}

// Close close the mongo session
func (us *UserStore) Close() {
	us.session.Close()
}

// nolint: unused
func (us *UserStore) c(name string) *mgo.Collection {
	return us.session.DB(us.dbName).C(name)
}

func (us *UserStore) cHandler(name string, handler func(c *mgo.Collection)) {
	session := us.session.Clone()
	defer session.Close()
	handler(session.DB(us.dbName).C(name))
}

// Set set user information
func (us *UserStore) Set(user *models.UserCredentials) (err error) {
	us.cHandler(us.ucfg.UsersCName, func(c *mgo.Collection) {
		// user.UID = uuid.Must(uuid.NewRandom()).String()
		t := time.Now()
		user.CreatedAt = &t
		if cerr := c.Insert(user); cerr != nil {
			err = cerr
			return
		}
	})
	return
}

// GetAllUsers according to the ID for the user information
func (us *UserStore) GetAllUsers() (users []*models.UserCredentials, err error) {
	us.cHandler(us.ucfg.UsersCName, func(c *mgo.Collection) {
		if cerr := c.Find(bson.M{}).All(&users); cerr != nil {
			err = cerr
			return
		}
	})

	return
}

//UpdateUser updates the user
func (us *UserStore) UpdateUser(user *models.UserCredentials) (err error) {
	us.cHandler(us.ucfg.UsersCName, func(c *mgo.Collection) {
		t := time.Now()
		user.UpdatedAt = &t
		if cerr := c.UpdateId(user.ID, user); cerr != nil {
			err = cerr
			return
		}
	})
	return
}

// RemoveByUserName use the user id to delete the user information
func (us *UserStore) RemoveByUserName(username string) (err error) {
	us.cHandler(us.ucfg.UsersCName, func(c *mgo.Collection) {
		if cerr := c.Remove(bson.M{"username": username}); cerr != nil {
			err = cerr
			return
		}
	})

	return
}

// GetUser according to the whatever passed
func (us *UserStore) GetUser(query interface{}) (user *models.UserCredentials, err error) {
	us.cHandler(us.ucfg.UsersCName, func(c *mgo.Collection) {
		user = new(models.UserCredentials)
		if cerr := c.Find(query).One(user); cerr != nil {
			err = cerr
			return
		}
	})

	return
}

// GetUserByID according to the whatever passed
func (us *UserStore) GetUserByID(id bson.ObjectId) (user *models.UserCredentials, err error) {
	return us.GetUser(bson.M{"_id": id})
}
