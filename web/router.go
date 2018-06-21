package web

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store"
)

// Router listens and responds to requests to the node for valid paths.
func Router(app *services.ChainlinkApplication) *gin.Engine {
	engine := gin.New()
	config := app.Store.Config
	cors := uiCorsHandler(config)
	engine.Use(
		loggerFunc(),
		gin.Recovery(),
		cors,
	)

	v1Routes(app, engine)
	v2Routes(app, engine)
	guiAssetRoutes(engine)

	return engine
}

func basicAuth(config store.Config) gin.HandlerFunc {
	return gin.BasicAuth(gin.Accounts{config.BasicAuthUsername: config.BasicAuthPassword})
}

func v1Routes(app *services.ChainlinkApplication, engine *gin.Engine) {
	v1 := engine.Group("/v1")
	v1.Use(basicAuth(app.Store.Config))
	ac := AssignmentsController{app}
	v1.POST("/assignments", ac.Create)
	v1.GET("/assignments/:ID", ac.Show)

	sc := SnapshotsController{app}
	v1.POST("/assignments/:AID/snapshots", sc.CreateSnapshot)
	v1.GET("/snapshots/:ID", sc.ShowSnapshot)
}

func v2Routes(app *services.ChainlinkApplication, engine *gin.Engine) {
	v2 := engine.Group("/v2")
	v2.Use(basicAuth(app.Store.Config))
	ab := AccountBalanceController{app}
	v2.GET("/account_balance", ab.Show)

	j := JobSpecsController{app}
	v2.GET("/specs", j.Index)
	v2.POST("/specs", j.Create)
	v2.GET("/specs/:SpecID", j.Show)

	jr := JobRunsController{app}
	v2.GET("/specs/:SpecID/runs", jr.Index)
	v2.POST("/specs/:SpecID/runs", jr.Create)
	v2.PATCH("/runs/:RunID", jr.Update)
	v2.GET("/runs/:RunID", jr.Show)

	tt := BridgeTypesController{app}
	v2.GET("/bridge_types", tt.Index)
	v2.POST("/bridge_types", tt.Create)
	v2.GET("/bridge_types/:BridgeName", tt.Show)
	v2.DELETE("/bridge_types/:BridgeName", tt.Destroy)

	backup := BackupController{app}
	v2.GET("/backup", backup.Show)

	cc := ConfigController{app}
	v2.GET("/config", cc.Show)
}

func guiAssetRoutes(engine *gin.Engine) {
	box := NewBox()
	boxList := box.List()

	engine.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		matchedBoxPath := matchPath(boxList, path)

		if matchedBoxPath == "" {
			if filepath.Ext(path) == "" {
				matchedBoxPath = wildcardMatchPath(
					boxList,
					path,
					"index.html",
				)
			} else if filepath.Ext(path) == ".json" {
				matchedBoxPath = wildcardMatchPath(
					boxList,
					filepath.Dir(path),
					filepath.Base(path),
				)
			}
		}

		if matchedBoxPath != "" {
			file, err := box.Open(matchedBoxPath)
			if err != nil {
				if err == os.ErrNotExist {
					c.AbortWithStatus(http.StatusNotFound)
				} else {
					err := fmt.Errorf("failed to open static file '%s': %+v", path, err)
					logger.Error(err.Error())
					c.AbortWithError(http.StatusInternalServerError, err)
				}
				return
			}

			http.ServeContent(c.Writer, c.Request, path, time.Time{}, file)
		}
	})
}

func wildcardMatchPath(boxList []string, path string, file string) (matchedPath string) {
	idPattern := regexp.MustCompile(`(_[a-zA-Z0-9]+_)`)
	pathSeparator := (string)(os.PathSeparator)
	escapedPathSeparator := pathSeparator
	if pathSeparator == `\` {
		escapedPathSeparator = `\\`
	}
	pathAndFile := filepath.Clean(strings.Join(
		[]string{path, file},
		pathSeparator,
	))
	normalizedPathAndFile := strings.Replace(
		strings.TrimPrefix(pathAndFile, `/`),
		`/`,
		pathSeparator,
		-1,
	)

	for i := 0; i < len(boxList) && matchedPath == ""; i++ {
		boxPathWithIDPattern := idPattern.ReplaceAllString(boxList[i], `[a-zA-Z0-9]+`)
		pathPattern := fmt.Sprintf(
			`^%s$`,
			strings.Replace(boxPathWithIDPattern, `\`, escapedPathSeparator, -1),
		)
		match, _ := regexp.MatchString(pathPattern, normalizedPathAndFile)

		if match {
			matchedPath = boxList[i]
		}
	}

	return matchedPath
}

func matchPath(boxList []string, path string) (matchedPath string) {
	pathWithoutPrefix := strings.TrimPrefix(path, "/")

	for i := 0; i < len(boxList) && matchedPath == ""; i++ {
		boxPath := boxList[i]
		if boxPath == pathWithoutPrefix {
			matchedPath = boxPath
		}
	}

	return matchedPath
}

// Inspired by https://github.com/gin-gonic/gin/issues/961
func loggerFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		buf, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			logger.Error("Web request log error: ", err.Error())
			c.Next()
			return
		}
		rdr := bytes.NewBuffer(buf)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(buf))

		start := time.Now()
		c.Next()
		end := time.Now()

		logger.Infow("Web request",
			"method", c.Request.Method,
			"status", c.Writer.Status(),
			"path", c.Request.URL.Path,
			"query", c.Request.URL.RawQuery,
			"body", readBody(rdr),
			"clientIP", c.ClientIP(),
			"errors", c.Errors.String(),
			"servedAt", end.Format("2006/01/02 - 15:04:05"),
			"latency", fmt.Sprintf("%v", end.Sub(start)),
		)
	}
}

// Add CORS headers so UI can make api requests
func uiCorsHandler(config store.Config) gin.HandlerFunc {
	c := cors.Config{
		AllowMethods:     []string{"GET"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	if config.AllowOrigins == "*" {
		c.AllowAllOrigins = true
	} else {
		allowOrigins := strings.Split(config.AllowOrigins, ",")
		if len(allowOrigins) > 0 {
			c.AllowOrigins = allowOrigins
		}
	}
	return cors.New(c)
}

func readBody(reader io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)

	s := buf.String()
	return s
}
