package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store"
)

const (
	// SessionName is the session name
	SessionName = "clsession"
	// SessionIDKey is the session ID key in the session map
	SessionIDKey = "clsession_id"
)

// Router listens and responds to requests to the node for valid paths.
func Router(app *services.ChainlinkApplication) *gin.Engine {
	engine := gin.New()
	config := app.Store.Config
	secret, err := config.SessionSecret()
	if err != nil {
		logger.Panic(err)
	}
	sessionStore := sessions.NewCookieStore(secret)
	sessionStore.Options(config.SessionOptions())
	cors := uiCorsHandler(config)

	engine.Use(
		loggerFunc(),
		gin.Recovery(),
		cors,
		sessions.Sessions(SessionName, sessionStore),
	)

	sessionRoutes(app, engine)
	v1Routes(app, engine)
	v2Routes(app, engine)
	guiAssetRoutes(engine)

	return engine
}

func authRequired(store *store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		sessionID, ok := session.Get(SessionIDKey).(string)
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
		} else if _, err := store.AuthorizedUserWithSession(sessionID); err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		} else {
			c.Next()
		}
	}
}

func sessionRoutes(app *services.ChainlinkApplication, engine *gin.Engine) {
	sc := SessionsController{app}
	engine.POST("/sessions", sc.Create)
	auth := engine.Group("/", authRequired(app.Store))
	auth.DELETE("/sessions", sc.Destroy)
}

func v1Routes(app *services.ChainlinkApplication, engine *gin.Engine) {
	v1 := engine.Group("/v1")
	v1.Use(authRequired(app.Store))

	ac := AssignmentsController{app}
	v1.POST("/assignments", ac.Create)
	v1.GET("/assignments/:ID", ac.Show)

	sc := SnapshotsController{app}
	v1.POST("/assignments/:AID/snapshots", sc.CreateSnapshot)
	v1.GET("/snapshots/:ID", sc.ShowSnapshot)
}

func v2Routes(app *services.ChainlinkApplication, engine *gin.Engine) {
	v2 := engine.Group("/v2")
	jr := JobRunsController{app}

	ab := AccountBalanceController{app}
	v2.GET("/account_balance", ab.Show)

	sa := ServiceAgreementsController{app}
	v2.POST("/service_agreements", sa.Create)
	v2.PATCH("/runs/:RunID", jr.Update)

	authv2 := engine.Group("/v2", authRequired(app.Store))
	{
		j := JobSpecsController{app}
		authv2.GET("/specs", j.Index)
		authv2.POST("/specs", j.Create)
		authv2.GET("/specs/:SpecID", j.Show)

		authv2.GET("/specs/:SpecID/runs", jr.Index)
		authv2.POST("/specs/:SpecID/runs", jr.Create)
		authv2.GET("/runs/:RunID", jr.Show)

		authv2.GET("/service_agreements/:SAID", sa.Show)

		bt := BridgeTypesController{app}
		authv2.GET("/bridge_types", bt.Index)
		authv2.POST("/bridge_types", bt.Create)
		authv2.GET("/bridge_types/:BridgeName", bt.Show)
		authv2.DELETE("/bridge_types/:BridgeName", bt.Destroy)

		backup := BackupController{app}
		authv2.GET("/backup", backup.Show)

		cc := ConfigController{app}
		authv2.GET("/config", cc.Show)
	}
}

func guiAssetRoutes(engine *gin.Engine) {
	box := NewBox()
	boxList := box.List()

	engine.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		matchedBoxPath := MatchExactBoxPath(boxList, path)

		if matchedBoxPath == "" {
			if filepath.Ext(path) == "" {
				matchedBoxPath = MatchWildcardBoxPath(
					boxList,
					path,
					"index.html",
				)
			} else if filepath.Ext(path) == ".json" {
				matchedBoxPath = MatchWildcardBoxPath(
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
			"query", redact(c.Request.URL.Query()),
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
		AllowMethods:     []string{"GET", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           math.MaxInt32,
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

	s, err := readSanitizedJSON(buf)
	if err != nil {
		return buf.String()
	}
	return s
}

var blacklist = map[string]struct{}{"password": struct{}{}}

func readSanitizedJSON(buf *bytes.Buffer) (string, error) {
	var dst map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &dst)
	if err != nil {
		return "", err
	}

	cleaned := map[string]interface{}{}
	for k, v := range dst {
		if _, ok := blacklist[strings.ToLower(k)]; ok {
			cleaned[k] = "*REDACTED*"
			continue
		}
		cleaned[k] = v
	}

	b, err := json.Marshal(cleaned)
	if err != nil {
		return "", err
	}
	return string(b), err
}

func redact(values url.Values) string {
	cleaned := url.Values{}
	for k, v := range values {
		if _, ok := blacklist[strings.ToLower(k)]; ok {
			cleaned[k] = []string{"REDACTED"}
			continue
		}
		cleaned[k] = v
	}
	return cleaned.Encode()
}
