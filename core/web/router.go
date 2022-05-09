package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"math"
	"net/http"
	"net/http/pprof"
	"net/url"
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
	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/pkg/errors"
	"github.com/ulule/limiter"
	mgin "github.com/ulule/limiter/drivers/middleware/gin"
	"github.com/ulule/limiter/drivers/store/memory"
	"github.com/unrolled/secure"

	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/web/auth"
	"github.com/smartcontractkit/chainlink/core/web/loader"
	"github.com/smartcontractkit/chainlink/core/web/resolver"
	"github.com/smartcontractkit/chainlink/core/web/schema"
)

// Router listens and responds to requests to the node for valid paths.
func Router(app chainlink.Application, prometheus *ginprom.Prometheus) *gin.Engine {
	engine := gin.New()
	config := app.GetConfig()
	secret, err := config.SessionSecret()
	if err != nil {
		app.GetLogger().Panic(err)
	}
	sessionStore := sessions.NewCookieStore(secret)
	sessionStore.Options(config.SessionOptions())
	cors := uiCorsHandler(config)
	if prometheus != nil {
		prometheus.Use(engine)
	}

	engine.Use(
		limits.RequestSizeLimiter(config.DefaultHTTPLimit()),
		loggerFunc(app.GetLogger()),
		gin.Recovery(),
		cors,
		secureMiddleware(config),
	)
	if prometheus != nil {
		engine.Use(prometheus.Instrument())
	}
	engine.Use(helmet.Default())

	api := engine.Group(
		"/",
		rateLimiter(
			config.AuthenticatedRateLimitPeriod().Duration(),
			config.AuthenticatedRateLimit(),
		),
		sessions.Sessions(auth.SessionName, sessionStore),
	)

	unauthenticatedDevOnlyMetricRoutes(app, api)
	healthRoutes(app, api)
	sessionRoutes(app, api)
	v2Routes(app, api)

	guiAssetRoutes(engine, config, app.GetLogger())

	api.POST("/query",
		auth.AuthenticateGQL(app.SessionORM()),
		loader.Middleware(app),
		graphqlHandler(app),
	)

	return engine
}

// Defining the Graphql handler
func graphqlHandler(app chainlink.Application) gin.HandlerFunc {
	rootSchema := schema.MustGetRootSchema()

	// Disable introspection and set a max query depth in production.
	schemaOpts := []graphql.SchemaOpt{}
	if !app.GetConfig().Dev() {
		schemaOpts = append(schemaOpts,
			graphql.MaxDepth(10),
		)
	}

	schema := graphql.MustParseSchema(rootSchema,
		&resolver.Resolver{
			App: app,
		},
		schemaOpts...,
	)

	h := relay.Handler{Schema: schema}

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func rateLimiter(period time.Duration, limit int64) gin.HandlerFunc {
	store := memory.NewStore()
	rate := limiter.Rate{
		Period: period,
		Limit:  limit,
	}
	return mgin.NewMiddleware(limiter.New(store, rate))
}

type WebSecurityConfig interface {
	AllowOrigins() string
	Dev() bool
	TLSRedirect() bool
	TLSHost() string
}

// secureOptions configure security options for the secure middleware, mostly
// for TLS redirection
func secureOptions(cfg WebSecurityConfig) secure.Options {
	return secure.Options{
		FrameDeny:     true,
		IsDevelopment: cfg.Dev(),
		SSLRedirect:   cfg.TLSRedirect(),
		SSLHost:       cfg.TLSHost(),
	}
}

// secureMiddleware adds a TLS handler and redirector, to button up security
// for this node
func secureMiddleware(cfg WebSecurityConfig) gin.HandlerFunc {
	secureMiddleware := secure.New(secureOptions(cfg))
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
func unauthenticatedDevOnlyMetricRoutes(app chainlink.Application, r *gin.RouterGroup) {
	group := r.Group("/debug", auth.Authenticate(app.SessionORM(), auth.AuthenticateBySession))
	group.GET("/vars", expvar.Handler())

	if app.GetConfig().Dev() {
		// No authentication because `go tool pprof` doesn't support it
		metricRoutes(r)
	}
}

func metricRoutes(r *gin.RouterGroup) {
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

func pprofHandler(h http.HandlerFunc) gin.HandlerFunc {
	handler := http.HandlerFunc(h)
	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}

func sessionRoutes(app chainlink.Application, r *gin.RouterGroup) {
	config := app.GetConfig()
	unauth := r.Group("/", rateLimiter(
		config.UnAuthenticatedRateLimitPeriod().Duration(),
		config.UnAuthenticatedRateLimit(),
	))
	sc := NewSessionsController(app)
	unauth.POST("/sessions", sc.Create)
	auth := r.Group("/", auth.Authenticate(app.SessionORM(), auth.AuthenticateBySession))
	auth.DELETE("/sessions", sc.Destroy)
}

func healthRoutes(app chainlink.Application, r *gin.RouterGroup) {
	hc := HealthController{app}
	r.GET("/readyz", hc.Readyz)
	r.GET("/health", hc.Health)
}

func v2Routes(app chainlink.Application, r *gin.RouterGroup) {
	unauthedv2 := r.Group("/v2")

	prc := PipelineRunsController{app}
	psec := PipelineJobSpecErrorsController{app}
	unauthedv2.PATCH("/resume/:runID", prc.Resume)

	authv2 := r.Group("/v2", auth.Authenticate(app.SessionORM(),
		auth.AuthenticateByToken,
		auth.AuthenticateBySession,
	))
	{
		uc := UserController{app}
		authv2.PATCH("/user/password", uc.UpdatePassword)
		authv2.POST("/user/token", uc.NewAPIToken)
		authv2.POST("/user/token/delete", uc.DeleteAPIToken)

		wa := NewWebAuthnController(app)
		authv2.GET("/enroll_webauthn", wa.BeginRegistration)
		authv2.POST("/enroll_webauthn", wa.FinishRegistration)

		eia := ExternalInitiatorsController{app}
		authv2.GET("/external_initiators", paginatedRequest(eia.Index))
		authv2.POST("/external_initiators", eia.Create)
		authv2.DELETE("/external_initiators/:Name", eia.Destroy)

		bt := BridgeTypesController{app}
		authv2.GET("/bridge_types", paginatedRequest(bt.Index))
		authv2.POST("/bridge_types", bt.Create)
		authv2.GET("/bridge_types/:BridgeName", bt.Show)
		authv2.PATCH("/bridge_types/:BridgeName", bt.Update)
		authv2.DELETE("/bridge_types/:BridgeName", bt.Destroy)

		ets := EVMTransfersController{app}
		authv2.POST("/transfers", ets.Create)
		authv2.POST("/transfers/evm", ets.Create)
		tts := TerraTransfersController{app}
		authv2.POST("/transfers/terra", tts.Create)
		sts := SolanaTransfersController{app}
		authv2.POST("/transfers/solana", sts.Create)

		cc := ConfigController{app}
		authv2.GET("/config", cc.Show)
		authv2.PATCH("/config", cc.Patch)

		tas := TxAttemptsController{app}
		authv2.GET("/tx_attempts", paginatedRequest(tas.Index))
		authv2.GET("/tx_attempts/evm", paginatedRequest(tas.Index))

		txs := TransactionsController{app}
		authv2.GET("/transactions/evm", paginatedRequest(txs.Index))
		authv2.GET("/transactions/evm/:TxHash", txs.Show)
		authv2.GET("/transactions", paginatedRequest(txs.Index))
		authv2.GET("/transactions/:TxHash", txs.Show)

		rc := ReplayController{app}
		authv2.POST("/replay_from_block/:number", rc.ReplayFromBlock)

		csakc := CSAKeysController{app}
		authv2.GET("/keys/csa", csakc.Index)
		authv2.POST("/keys/csa", csakc.Create)
		authv2.POST("/keys/csa/import", csakc.Import)
		authv2.POST("/keys/csa/export/:ID", csakc.Export)

		ekc := ETHKeysController{app}
		authv2.GET("/keys/eth", ekc.Index)
		authv2.POST("/keys/eth", ekc.Create)
		authv2.PUT("/keys/eth/:keyID", ekc.Update)
		authv2.DELETE("/keys/eth/:keyID", ekc.Delete)
		authv2.POST("/keys/eth/import", ekc.Import)
		authv2.POST("/keys/eth/export/:address", ekc.Export)

		ocrkc := OCRKeysController{app}
		authv2.GET("/keys/ocr", ocrkc.Index)
		authv2.POST("/keys/ocr", ocrkc.Create)
		authv2.DELETE("/keys/ocr/:keyID", ocrkc.Delete)
		authv2.POST("/keys/ocr/import", ocrkc.Import)
		authv2.POST("/keys/ocr/export/:ID", ocrkc.Export)

		ocr2kc := OCR2KeysController{app}
		authv2.GET("/keys/ocr2", ocr2kc.Index)
		authv2.POST("/keys/ocr2/:chainType", ocr2kc.Create)
		authv2.DELETE("/keys/ocr2/:keyID", ocr2kc.Delete)
		authv2.POST("/keys/ocr2/import", ocr2kc.Import)
		authv2.POST("/keys/ocr2/export/:ID", ocr2kc.Export)

		p2pkc := P2PKeysController{app}
		authv2.GET("/keys/p2p", p2pkc.Index)
		authv2.POST("/keys/p2p", p2pkc.Create)
		authv2.DELETE("/keys/p2p/:keyID", p2pkc.Delete)
		authv2.POST("/keys/p2p/import", p2pkc.Import)
		authv2.POST("/keys/p2p/export/:ID", p2pkc.Export)

		solkc := SolanaKeysController{app}
		authv2.GET("/keys/solana", solkc.Index)
		authv2.POST("/keys/solana", solkc.Create)
		authv2.DELETE("/keys/solana/:keyID", solkc.Delete)
		authv2.POST("/keys/solana/import", solkc.Import)
		authv2.POST("/keys/solana/export/:ID", solkc.Export)

		terkc := TerraKeysController{app}
		authv2.GET("/keys/terra", terkc.Index)
		authv2.POST("/keys/terra", terkc.Create)
		authv2.DELETE("/keys/terra/:keyID", terkc.Delete)
		authv2.POST("/keys/terra/import", terkc.Import)
		authv2.POST("/keys/terra/export/:ID", terkc.Export)

		vrfkc := VRFKeysController{app}
		authv2.GET("/keys/vrf", vrfkc.Index)
		authv2.POST("/keys/vrf", vrfkc.Create)
		authv2.DELETE("/keys/vrf/:keyID", vrfkc.Delete)
		authv2.POST("/keys/vrf/import", vrfkc.Import)
		authv2.POST("/keys/vrf/export/:keyID", vrfkc.Export)

		jc := JobsController{app}
		authv2.GET("/jobs", paginatedRequest(jc.Index))
		authv2.GET("/jobs/:ID", jc.Show)
		authv2.POST("/jobs", jc.Create)
		authv2.DELETE("/jobs/:ID", jc.Delete)

		// PipelineRunsController
		authv2.GET("/pipeline/runs", paginatedRequest(prc.Index))
		authv2.GET("/jobs/:ID/runs", paginatedRequest(prc.Index))
		authv2.GET("/jobs/:ID/runs/:runID", prc.Show)

		// FeaturesController
		fc := FeaturesController{app}
		authv2.GET("/features", fc.Index)

		// PipelineJobSpecErrorsController
		authv2.DELETE("/pipeline/job_spec_errors/:ID", psec.Destroy)

		lgc := LogController{app}
		authv2.GET("/log", lgc.Get)
		authv2.PATCH("/log", lgc.Patch)

		chains := authv2.Group("chains")
		for _, chain := range []struct {
			path string
			cc   ChainsController
		}{
			{"evm", NewEVMChainsController(app)},
			{"solana", NewSolanaChainsController(app)},
			{"terra", NewTerraChainsController(app)},
		} {
			chains.GET(chain.path, paginatedRequest(chain.cc.Index))
			chains.POST(chain.path, chain.cc.Create)
			chains.GET(chain.path+"/:ID", chain.cc.Show)
			chains.PATCH(chain.path+"/:ID", chain.cc.Update)
			chains.DELETE(chain.path+"/:ID", chain.cc.Delete)
		}

		nodes := authv2.Group("nodes")
		for _, chain := range []struct {
			path string
			nc   NodesController
		}{
			{"evm", NewEVMNodesController(app)},
			{"solana", NewSolanaNodesController(app)},
			{"terra", NewTerraNodesController(app)},
		} {
			if chain.path == "evm" {
				// TODO still EVM only https://app.shortcut.com/chainlinklabs/story/26276/multi-chain-type-ui-node-chain-configuration
				nodes.GET("", paginatedRequest(chain.nc.Index))
				nodes.POST("", chain.nc.Create)
				nodes.DELETE("/:ID", chain.nc.Delete)
			}
			nodes.GET(chain.path, paginatedRequest(chain.nc.Index))
			chains.GET(chain.path+"/:ID/nodes", paginatedRequest(chain.nc.Index))
			nodes.POST(chain.path, chain.nc.Create)
			nodes.DELETE(chain.path+"/:ID", chain.nc.Delete)
		}

		efc := EVMForwardersController{app}
		authv2.GET("/nodes/evm/forwarders", paginatedRequest(efc.Index))
		authv2.POST("/nodes/evm/forwarders", efc.Create)
		authv2.DELETE("/nodes/evm/forwarders/:fwdID", efc.Delete)

		build_info := BuildInfoController{app}
		authv2.GET("/build_info", build_info.Show)

		// Debug routes accessible via authentication
		metricRoutes(authv2)
	}

	ping := PingController{app}
	userOrEI := r.Group("/v2", auth.Authenticate(app.SessionORM(),
		auth.AuthenticateExternalInitiator,
		auth.AuthenticateByToken,
		auth.AuthenticateBySession,
	))
	userOrEI.GET("/ping", ping.Show)
	userOrEI.POST("/jobs/:ID/runs", prc.Create)
}

// This is higher because it serves main.js and any static images. There are
// 5 assets which must be served, so this allows for 20 requests/min
var staticAssetsRateLimit = int64(100)
var staticAssetsRateLimitPeriod = 1 * time.Minute
var indexRateLimit = int64(20)
var indexRateLimitPeriod = 1 * time.Minute

// guiAssetRoutes serves the operator UI static files and index.html. Rate
// limiting is disabled when in dev mode.
func guiAssetRoutes(engine *gin.Engine, config config.GeneralConfig, lggr logger.Logger) {
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
		ServeGzippedAssets("/assets", assetFs, lggr),
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
		file, err := assetFs.Open("index.html")
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				c.AbortWithStatus(http.StatusNotFound)
			} else {
				lggr.Errorf("failed to open static file '%s': %+v", path, err)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
			return
		}
		defer lggr.ErrorIfClosing(file, "file")

		http.ServeContent(c.Writer, c.Request, path, time.Time{}, file)
	})

	engine.NoRoute(noRouteHandlers...)
}

// Inspired by https://github.com/gin-gonic/gin/issues/961
func loggerFunc(lggr logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		buf, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			lggr.Error("Web request log error: ", err.Error())
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

		lggr.Infow(fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.Path),
			"method", c.Request.Method,
			"status", c.Writer.Status(),
			"path", c.Request.URL.Path,
			"query", redact(c.Request.URL.Query()),
			"body", readBody(rdr, lggr),
			"clientIP", c.ClientIP(),
			"errors", c.Errors.String(),
			"servedAt", end.Format("2006-01-02 15:04:05"),
			"latency", fmt.Sprintf("%v", end.Sub(start)),
		)
	}
}

// Add CORS headers so UI can make api requests
func uiCorsHandler(config WebSecurityConfig) gin.HandlerFunc {
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

func readBody(reader io.Reader, lggr logger.Logger) string {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(reader)
	if err != nil {
		lggr.Warn("unable to read from body for sanitization: ", err)
		return "*FAILED TO READ BODY*"
	}

	if buf.Len() == 0 {
		return ""
	}

	s, err := readSanitizedJSON(buf)
	if err != nil {
		lggr.Warn("unable to sanitize json for logging: ", err)
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
