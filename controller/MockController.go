package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zc2638/gotool/curlx"
	"github.com/zctod/go-tool/common/utils"
	"io"
	"io/ioutil"
	"math/rand"
	"mock/config"
	"mock/lib/network"
	"net/http"
	"strconv"
	"time"
)

/**
 * Created by zc on 2019-11-20.
 */
type MockController struct{ BaseController }

type Call struct {
	Address  string `json:"address"`
	Method   string `json:"method"`
	Sleep    int64  `json:"sleep"`     // ms
	BodySize int    `json:"body_size"` // KB
	Extend   Extend `json:"extend"`
}

type Extend struct {
	Headers []string `json:"headers"`
}

func (t *MockController) Any(c *gin.Context) {
	// 增加异常模拟方法
	if errorPercent, err := strconv.Atoi(c.Query("error_percent")); err == nil {
		rand.Seed(time.Now().Unix())
		i := rand.Intn(100)
		fmt.Printf("%d < %d\n", i, errorPercent)
		if i <= errorPercent {
			var errorCode int
			if c.Query("error_code") != "" {
				errorCode, err = strconv.Atoi(c.Query("error_code"))
				if err != nil {
					t.ErrData(c, err)
					return
				}
			} else {
				switch i % 4 {
				case 0:
					errorCode = http.StatusUnauthorized
				case 1:
					errorCode = http.StatusNotFound
				case 2:
					errorCode = http.StatusInternalServerError
				case 3:
					errorCode = http.StatusPermanentRedirect
				default:
					errorCode = http.StatusBadRequest
				}
			}
			c.JSON(errorCode, gin.H{
				"timestamp":    time.Now().UnixNano() / 1e3,
				"Request-Host": c.Request.Host,
				"URL":          c.Request.URL.String(),
				"RequestURI":   c.Request.RequestURI,
				"RemoteAddr":   c.Request.RemoteAddr,
				"Method":       c.Request.Method,
				"Header":       c.Request.Header,
				"Server-Host":  network.Hostname(),
				"Server-Ip":    network.IP(),
				"Body":         nil,
				"Next":         nil,
			})
			return
		}
	}

	b, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		t.ErrData(c, err)
		return
	}
	var result interface{}
	var result_code int

	if len(b) > 0 {
		var calls []Call
		if err := json.Unmarshal(b, &calls); err != nil {
			t.ErrData(c, err)
			return
		}

		l := len(calls)
		if l > 0 {
			call := calls[0]
			// sleep
			if call.Sleep > 0 {
				time.Sleep(time.Millisecond * time.Duration(call.Sleep))
			}
			if call.BodySize > 0 {
				var buffer bytes.Buffer
				currentSize := call.BodySize * 1024
				for i := 0; i < currentSize; i++ {
					buffer.WriteString("a")
				}
				c.String(http.StatusOK, buffer.String())
				return
			}
			// call
			r := curlx.NewRequest()
			r.Url = call.Address
			r.Method = call.Method
			if l > 1 {
				newCalls := calls[1:]
				nb, err := json.Marshal(newCalls)
				if err != nil {
					t.ErrData(c, err)
					return
				}
				r.Body = nb
			}
			r.Header = make(http.Header)
			for _, h := range config.OpenTracingHeaders {
				if hv := c.GetHeader(h); hv != "" {
					r.Header.Set(h, hv)
				}
			}
			for _, h := range call.Extend.Headers {
				if hv := c.GetHeader(h); hv != "" {
					r.Header.Set(h, hv)
				}
			}
			res, err := r.Do()
			if res != nil {
				result_code = res.StatusCode
			}
			if err != nil {
				result = gin.H{
					"status":  "error",
					"message": err.Error(),
				}
			} else if err := res.ParseJSON(&result); err != nil {
				result = gin.H{
					"status":  "parse json error",
					"message": string(res.Result),
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"timestamp":     time.Now().UnixNano() / 1e3,
		"Request-Host":  c.Request.Host,
		"URL":           c.Request.URL.String(),
		"RequestURI":    c.Request.RequestURI,
		"RemoteAddr":    c.Request.RemoteAddr,
		"Method":        c.Request.Method,
		"Header":        c.Request.Header,
		"Server-Host":   network.Hostname(),
		"Server-Ip":     network.IP(),
		"Body":          string(b),
		"Next-Code":     result_code,
		"Next-Response": result,
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
