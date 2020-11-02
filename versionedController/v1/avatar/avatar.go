package avatar

import (
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/golang/glog"

	controller "github.com/mayadata-io/kubera-auth/versionedController/v1"
)

const (
	extnpng  = "png"
	extnjpeg = "jpeg"
	extnjpg  = "jpg"
)

// Controller contains common stuff for avatar controller
type Controller struct {
	controller.GenericController
	path string
	dir  string
}

// New returns new instance of controller
func New() *Controller {
	avatar := "avatar"
	if os.Getenv("AVATAR_DIR") != "" {
		avatar = os.Getenv("AVATAR_DIR")
	} else {
		log.Warningf("Environment variable AVATAR_DIR is not set assigning default value `%s`", avatar)
	}
	return &Controller{
		path: controller.AvatarRoute,
		dir:  avatar,
	}
}

// Get lists all the avatar images insisde a given dir
func (c *Controller) Get(context *gin.Context) {
	items, err := ioutil.ReadDir(c.dir)
	if err != nil {
		context.String(http.StatusInternalServerError, "Error reading avatar dir error : %s", err)
		log.Errorf("Error reading avatar dir error : %s", err)
		return
	}
	files := make([]string, 0)
	for _, item := range items {
		if item.IsDir() {
			continue
		}
		if len(strings.Split(item.Name(), ".")) != 2 {
			continue
		}
		ext := strings.Split(item.Name(), ".")[1]
		if ext != "jpeg" && ext != "jpg" && ext != "png" {
			continue
		}
		files = append(files, item.Name())
	}
	context.JSON(http.StatusOK, files)
}

// GetByID returns avatar image insisde a given dir for a given avatar
func (c *Controller) GetByID(context *gin.Context) {
	name := context.Param("id")
	if len(strings.Split(name, ".")) != 2 {
		context.String(http.StatusInternalServerError, "Unsupported avatar name : %s", name)
		log.Errorf("Unsupported avatar name : %s", name)
		return
	}
	ext := strings.Split(name, ".")[1]
	if ext != extnjpeg && ext != extnjpg && ext != extnpng {
		context.String(http.StatusInternalServerError, "Unsupported avatar extension : %s", ext)
		log.Errorf("Unsupported avatar extension : %s", ext)
		return
	}
	fp, err := os.Open(c.dir + "/" + name)
	if err != nil {
		context.String(http.StatusNotFound, "Avatar file not found error : %s", err)
		log.Errorf("Avatar file not found error : %s", err)
		return
	}
	defer fp.Close()
	img, format, err := image.Decode(fp)
	if err != nil {
		context.String(http.StatusInternalServerError, "Unable to decode avatar error : %s", err)
		log.Errorf("Unable to decode avatar error : %s", err)
		return
	}
	switch format {
	case extnjpg, extnjpeg:
		var rgba *image.RGBA
		if nrgba, ok := img.(*image.NRGBA); ok {
			if nrgba.Opaque() {
				rgba = &image.RGBA{
					Pix:    nrgba.Pix,
					Stride: nrgba.Stride,
					Rect:   nrgba.Rect,
				}
			}
		}
		if rgba != nil {
			err = jpeg.Encode(context.Writer, rgba, &jpeg.Options{Quality: 100})
		} else {
			err = jpeg.Encode(context.Writer, img, &jpeg.Options{Quality: 100})
		}
	case extnpng:
		err = png.Encode(context.Writer, img)
	}
	if err != nil {
		context.String(http.StatusInternalServerError, "Unable to write avatar error : %s", err)
		log.Errorf("Unable to write avatar in respose writter err : %s", err)
		return
	}
}

// Register will rsgister this controller to the specified router
func (c *Controller) Register(router *gin.RouterGroup) {
	controller.RegisterController(router, c, c.path)
}
