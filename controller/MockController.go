package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/zctod/go-tool/common/utils"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

/**
 * Created by zc on 2019-11-20.
 */
type MockController struct{ BaseController }

func (t *MockController) Any(c *gin.Context) {

	b, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		t.ErrData(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"timestamp":  time.Now().UnixNano() / 1e3,
		"Host":       c.Request.Host,
		"URL":        c.Request.URL.String(),
		"RequestURI": c.Request.RequestURI,
		"RemoteAddr": c.Request.RemoteAddr,
		"Method":     c.Request.Method,
		"Header":     c.Request.Header,
		"Body":       string(b),
	})
}

func (t *MockController) Upload(c *gin.Context) {

	// 上传图片
	file, info, err := c.Request.FormFile("file")
	if err != nil {
		t.ErrData(c, err)
		return
	}

	filename := info.Filename
	imagePath := "uploads/" + filename
	out, err := utils.CreateFile(imagePath)
	if err != nil {
		t.ErrData(c, err)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		t.ErrData(c, err)
		return
	}

	imageUrl := c.Request.Host + "/" + imagePath
	c.String(http.StatusCreated, imageUrl)
}
