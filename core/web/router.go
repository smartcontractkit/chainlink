package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"net/http/pprof"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/Depado/ginprom"
	helmet "github.com/danielkov/gin-helmet"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/expvar"
	limits "github.com/gin-contrib/size"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
	"github.com/ulule/limiter"
	mgin "github.com/ulule/limiter/drivers/middleware/gin"
	"github.com/ulule/limiter/drivers/store/memory"
	"github.com/unrolled/secure"
)

var prometheus *ginprom.Prometheus

func init() {
	gin.DebugPrintRouteFunc = printRoutes

	// ensure metrics are regsitered once per instance to avoid registering
	// metrics multiple times (panic)
	prometheus = ginprom.New(ginprom.Namespace("service"))
}

func printRoutes(httpMethod, absolutePath, handlerName string, nuHandlers int) {
	logger.Debugf("%-6s %-25s --> %s (%d handlers)", httpMethod, absolutePath, handlerName, nuHandlers)
}

const (
	// SessionName is the session name
	SessionName = "clsession"
	// SessionIDKey is the session ID key in the session map
	SessionIDKey = "clsession_id"
	// SessionUserKey is the User key in the session map
	SessionUserKey = "user"
	// SessionExternalInitiatorKey is the External Initiator key in the session map
	SessionExternalInitiatorKey = "external_initiator"
)

func explorerStatus(app chainlink.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		es := presenters.NewExplorerStatus(app.GetStatsPusher())
		b, err := json.Marshal(es)
		if err != nil {
			panic(err)
		}

		c.SetCookie("explorer", (string)(b), 0, "", "", http.SameSiteStrictMode, false, false)
		c.Next()
	}
}

// Router listens and responds to requests to the node for valid paths.
func Router(app chainlink.Application) *gin.Engine {
	engine := gin.New()
	store := app.GetStore()
	config := store.Config
	secret, err := config.SessionSecret()
	if err != nil {
		logger.Panic(err)
	}
	sessionStore := sessions.NewCookieStore(secret)
	sessionStore.Options(config.SessionOptions())
	cors := uiCorsHandler(config)

	prometheus.Use(engine)
	engine.Use(
		limits.RequestSizeLimiter(config.DefaultHTTPLimit()),
		loggerFunc(),
		gin.Recovery(),
		cors,
		secureMiddleware(config),
		prometheus.Instrument(),
	)
	engine.Use(helmet.Default())

	api := engine.Group(
		"/",
		rateLimiter(
			config.AuthenticatedRateLimitPeriod().Duration(),
			config.AuthenticatedRateLimit(),
		),
		sessions.Sessions(SessionName, sessionStore),
		explorerStatus(app),
	)

	metricRoutes(app, api)
	healthRoutes(app, api)
	sessionRoutes(app, api)
	v2Routes(app, api)

	guiAssetRoutes(app.NewBox(), engine, config)

	return engine
}

func rateLimiter(period time.Duration, limit int64) gin.HandlerFunc {
	store := memory.NewStore()
	rate := limiter.Rate{
		Period: period,
		Limit:  limit,
	}
	return mgin.NewMiddleware(limiter.New(store, rate))
}

// secureOptions configure security options for the secure middleware, mostly
// for TLS redirection
func secureOptions(config orm.ConfigReader) secure.Options {
	return secure.Options{
		FrameDeny:     true,
		IsDevelopment: config.Dev(),
		SSLRedirect:   config.TLSRedirect(),
		SSLHost:       config.TLSHost(),
	}
}

// secureMiddleware adds a TLS handler and redirector, to button up security
// for this node
func secureMiddleware(config orm.ConfigReader) gin.HandlerFunc {
	secureMiddleware := secure.New(secureOptions(config))
	secureFunc := func() gin.HandlerFunc {
		return func(c *gin.Context) {
			err := secureMiddleware.Process(c.Writer, c.Request)

			// If there was an error, do not continue.
			if err != nil {
				c.Abort()
				return
			}

			// Avoid header rewrite if response is a redirection.
			if status := c.Writer.Status(); status > 300 && status < 399 {
				c.Abort()
			}
		}
	}()

	return secureFunc
}
func metricRoutes(app chainlink.Application, r *gin.RouterGroup) {
	group := r.Group("/debug", RequireAuth(app.GetStore(), AuthenticateBySession))
	group.GET("/vars", expvar.Handler())

	if app.GetStore().Config.Dev() {
		// No authentication because `go tool pprof` doesn't support it
		pprofGroup := r.Group("/debug/pprof")
		pprofGroup.GET("/", pprofHandler(pprof.Index))
		pprofGroup.GET("/cmdline", pprofHandler(pprof.Cmdline))
		pprofGroup.GET("/profile", pprofHandler(pprof.Profile))
		pprofGroup.POST("/symbol", pprofHandler(pprof.Symbol))
		pprofGroup.GET("/symbol", pprofHandler(pprof.Symbol))
		pprofGroup.GET("/trace", pprofHandler(pprof.Trace))
		pprofGroup.GET("/allocs", pprofHandler(pprof.Handler("allocs").ServeHTTP))
		pprofGroup.GET("/block", pprofHandler(pprof.Handler("block").ServeHTTP))
		pprofGroup.GET("/goroutine", pprofHandler(pprof.Handler("goroutine").ServeHTTP))
		pprofGroup.GET("/heap", pprofHandler(pprof.Handler("heap").ServeHTTP))
		pprofGroup.GET("/mutex", pprofHandler(pprof.Handler("mutex").ServeHTTP))
		pprofGroup.GET("/threadcreate", pprofHandler(pprof.Handler("threadcreate").ServeHTTP))
	}
}

func pprofHandler(h http.HandlerFunc) gin.HandlerFunc {
	handler := http.HandlerFunc(h)
	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}

func sessionRoutes(app chainlink.Application, r *gin.RouterGroup) {
	config := app.GetStore().Config
	unauth := r.Group("/", rateLimiter(
		config.UnAuthenticatedRateLimitPeriod().Duration(),
		config.UnAuthenticatedRateLimit(),
	))
	sc := SessionsController{app}
	unauth.POST("/sessions", sc.Create)
	auth := r.Group("/", RequireAuth(app.GetStore(), AuthenticateBySession))
	auth.DELETE("/sessions", sc.Destroy)
}

func healthRoutes(app chainlink.Application, r *gin.RouterGroup) {
	hc := HealthController{app}
	r.GET("/readyz", hc.Readyz)
	r.GET("/health", hc.Health)
}

func v2Routes(app chainlink.Application, r *gin.RouterGroup) {
	unauthedv2 := r.Group("/v2")

	jr := JobRunsController{app}
	unauthedv2.PATCH("/runs/:RunID", jr.Update)

	sa := ServiceAgreementsController{app}
	unauthedv2.POST("/service_agreements", sa.Create)

	j := JobSpecsController{app}
	jsec := JobSpecErrorsController{app}

	authv2 := r.Group("/v2", RequireAuth(app.GetStore(), AuthenticateByToken, AuthenticateBySession))
	{
		uc := UserController{app}
		authv2.PATCH("/user/password", uc.UpdatePassword)
		authv2.POST("/user/token", uc.NewAPIToken)
		authv2.POST("/user/token/delete", uc.DeleteAPIToken)

		eia := ExternalInitiatorsController{app}
		authv2.POST("/external_initiators", eia.Create)
		authv2.DELETE("/external_initiators/:Name", eia.Destroy)

		authv2.POST("/specs", j.Create)
		authv2.GET("/specs", paginatedRequest(j.Index))
		authv2.GET("/specs/:SpecID", j.Show)
		authv2.DELETE("/specs/:SpecID", j.Destroy)

		authv2.GET("/runs", paginatedRequest(jr.Index))
		authv2.GET("/runs/:RunID", jr.Show)
		authv2.PUT("/runs/:RunID/cancellation", jr.Cancel)

		authv2.DELETE("/job_spec_errors/:jobSpecErrorID", jsec.Destroy)

		authv2.GET("/service_agreements/:SAID", sa.Show)

		bt := BridgeTypesController{app}
		authv2.GET("/bridge_types", paginatedRequest(bt.Index))
		authv2.POST("/bridge_types", bt.Create)
		authv2.GET("/bridge_types/:BridgeName", bt.Show)
		authv2.PATCH("/bridge_types/:BridgeName", bt.Update)
		authv2.DELETE("/bridge_types/:BridgeName", bt.Destroy)

		ts := TransfersController{app}
		authv2.POST("/transfers", ts.Create)

		cc := ConfigController{app}
		authv2.GET("/config", cc.Show)
		authv2.PATCH("/config", cc.Patch)

		tas := TxAttemptsController{app}
		authv2.GET("/tx_attempts", paginatedRequest(tas.Index))

		txs := TransactionsController{app}
		authv2.GET("/transactions", paginatedRequest(txs.Index))
		authv2.GET("/transactions/:TxHash", txs.Show)

		bdc := BulkDeletesController{app}
		authv2.DELETE("/bulk_delete_runs", bdc.Delete)

		ekc := ETHKeysController{app}
		authv2.GET("/keys/eth", ekc.Index)
		authv2.POST("/keys/eth", ekc.Create)
		authv2.DELETE("/keys/eth/:keyID", ekc.Delete)
		authv2.POST("/keys/eth/import", ekc.Import)
		authv2.POST("/keys/eth/export/:address", ekc.Export)

		ocrkc := OCRKeysController{app}
		authv2.GET("/keys/ocr", ocrkc.Index)
		authv2.POST("/keys/ocr", ocrkc.Create)
		authv2.DELETE("/keys/ocr/:keyID", ocrkc.Delete)
		authv2.POST("/keys/ocr/import", ocrkc.Import)
		authv2.POST("/keys/ocr/export/:ID", ocrkc.Export)

		p2pkc := P2PKeysController{app}
		authv2.GET("/keys/p2p", p2pkc.Index)
		authv2.POST("/keys/p2p", p2pkc.Create)
		authv2.DELETE("/keys/p2p/:keyID", p2pkc.Delete)
		authv2.POST("/keys/p2p/import", p2pkc.Import)
		authv2.POST("/keys/p2p/export/:ID", p2pkc.Export)

		jc := JobsController{app}
		authv2.GET("/jobs", jc.Index)
		authv2.GET("/jobs/:ID", jc.Show)
		authv2.POST("/jobs", jc.Create)
		authv2.DELETE("/jobs/:ID", jc.Delete)

		mc := MigrateController{app}
		authv2.POST("/migrate/:ID", mc.Migrate)

		prc := PipelineRunsController{app}
		authv2.GET("/jobs/:ID/runs", paginatedRequest(prc.Index))
		authv2.GET("/jobs/:ID/runs/:runID", prc.Show)
		authv2.POST("/jobs/:ID/runs", prc.Create)

		lgc := LogController{app}
		authv2.GET("/log", lgc.Get)
		authv2.PATCH("/log", lgc.Patch)
	}

	ping := PingController{app}
	userOrEI := r.Group("/v2", RequireAuth(app.GetStore(),
		AuthenticateExternalInitiator,
		AuthenticateByToken,
		AuthenticateBySession,
	))
	userOrEI.POST("/specs/:SpecID/runs", jr.Create)
	userOrEI.GET("/ping", ping.Show)
}

// This is higher because it serves main.js and any static images. There are
// 5 assets which must be served, so this allows for 20 requests/min
var staticAssetsRateLimit = int64(100)
var staticAssetsRateLimitPeriod = 1 * time.Minute
var indexRateLimit = int64(20)
var indexRateLimitPeriod = 1 * time.Minute

// guiAssetRoutes serves the operator UI static files and index.html. Rate
// limiting is disabled when in dev mode.
func guiAssetRoutes(box packr.Box, engine *gin.Engine, config *orm.Config) {
	// Serve static files
	assetsRouterHandlers := []gin.HandlerFunc{}
	if !config.Dev() {
		assetsRouterHandlers = append(assetsRouterHandlers, rateLimiter(
			staticAssetsRateLimitPeriod,
			staticAssetsRateLimit,
		))
	}

	assetsRouterHandlers = append(
		assetsRouterHandlers,
		ServeGzippedAssets("/assets", &BoxFileSystem{Box: box}),
	)

	// Get Operator UI Assets
	//
	// We have to use a route here because a RouterGroup only runs middlewares
	// when a route matches exactly. See https://github.com/gin-gonic/gin/issues/531
	engine.GET("/assets/:file", assetsRouterHandlers...)

	// Serve the index HTML file unless it is an api path
	noRouteHandlers := []gin.HandlerFunc{}
	if !config.Dev() {
		noRouteHandlers = append(noRouteHandlers, rateLimiter(
			indexRateLimitPeriod,
			indexRateLimit,
		))
	}
	noRouteHandlers = append(noRouteHandlers, func(c *gin.Context) {
		path := c.Request.URL.Path

		// Return a 404 if the path is an unmatched API path
		if match, _ := regexp.MatchString(`^/v[0-9]+/.*`, path); match {
			c.AbortWithStatus(http.StatusNotFound)

			return
		}

		// Return a 404 for unknown extensions
		if filepath.Ext(path) != "" {
			c.AbortWithStatus(http.StatusNotFound)

			return
		}

		// Render the React index page for any other unknown requests
		file, err := box.Open("index.html")
		if err != nil {
			if err == os.ErrNotExist {
				c.AbortWithStatus(http.StatusNotFound)
			} else {
				logger.Errorf("failed to open static file '%s': %+v", path, err)
				c.AbortWithStatus(http.StatusInternalServerError)
			}

		}
		defer logger.ErrorIfCalling(file.Close, "failed when close file")

		http.ServeContent(c.Writer, c.Request, path, time.Time{}, file)
	})

	engine.NoRoute(noRouteHandlers...)
}

// Inspired by https://github.com/gin-gonic/gin/issues/961
func loggerFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		buf, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			logger.Error("Web request log error: ", err.Error())
			// Implicitly relies on limits.RequestSizeLimiter
			// overriding of c.Request.Body to abort gin's Context
			// inside ioutil.ReadAll.
			// Functions as we would like, but horrible from an architecture
			// and design pattern perspective.
			if !c.IsAborted() {
				c.AbortWithStatus(http.StatusBadRequest)
			}
			return
		}
		rdr := bytes.NewBuffer(buf)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(buf))

		start := time.Now()
		c.Next()
		end := time.Now()

		logger.Infow(fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.Path),
			"method", c.Request.Method,
			"status", c.Writer.Status(),
			"path", c.Request.URL.Path,
			"query", redact(c.Request.URL.Query()),
			"body", readBody(rdr),
			"clientIP", c.ClientIP(),
			"errors", c.Errors.String(),
			"servedAt", end.Format("2006-01-02 15:04:05"),
			"latency", fmt.Sprintf("%v", end.Sub(start)),
		)
	}
}

// Add CORS headers so UI can make api requests
func uiCorsHandler(config orm.ConfigReader) gin.HandlerFunc {
	c := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           math.MaxInt32,
	}
	if config.AllowOrigins() == "*" {
		c.AllowAllOrigins = true
	} else if allowOrigins := strings.Split(config.AllowOrigins(), ","); len(allowOrigins) > 0 {
		c.AllowOrigins = allowOrigins
	}
	return cors.New(c)
}

func readBody(reader io.Reader) string {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(reader)
	if err != nil {
		logger.Warn("unable to read from body for sanitization: ", err)
		return "*FAILED TO READ BODY*"
	}

	if buf.Len() == 0 {
		return ""
	}

	s, err := readSanitizedJSON(buf)
	if err != nil {
		logger.Warn("unable to sanitize json for logging: ", err)
		return "*FAILED TO READ BODY*"
	}
	return s
}

func readSanitizedJSON(buf *bytes.Buffer) (string, error) {
	var dst map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &dst)
	if err != nil {
		return "", err
	}

	cleaned := map[string]interface{}{}
	for k, v := range dst {
		if isBlacklisted(k) {
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
		if isBlacklisted(k) {
			cleaned[k] = []string{"REDACTED"}
			continue
		}
		cleaned[k] = v
	}
	return cleaned.Encode()
}

// NOTE: keys must be in lowercase for case insensitive match
var blacklist = map[string]struct{}{
	"password":             {},
	"newpassword":          {},
	"oldpassword":          {},
	"current_password":     {},
	"new_account_password": {},
}

func isBlacklisted(k string) bool {
	lk := strings.ToLower(k)
	if _, ok := blacklist[lk]; ok || strings.Contains(lk, "password") {
		return true
	}
	return false
}
