package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
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

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"

	"github.com/gin-gonic/gin"
	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/pkg/errors"
	"github.com/ulule/limiter"
	mgin "github.com/ulule/limiter/drivers/middleware/gin"
	"github.com/ulule/limiter/drivers/store/memory"
	"github.com/unrolled/secure"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/web/auth"
	"github.com/smartcontractkit/chainlink/core/web/loader"
	"github.com/smartcontractkit/chainlink/core/web/resolver"
	"github.com/smartcontractkit/chainlink/core/web/schema"
)

// NewRouter returns *gin.Engine router that listens and responds to requests to the node for valid paths.
func NewRouter(app chainlink.Application, prometheus *ginprom.Prometheus) (*gin.Engine, error) {
	engine := gin.New()
	config := app.GetConfig()
	secret, err := app.SecretGenerator().Generate(config.RootDir())
	if err != nil {
		return nil, err
	}
	sessionStore := cookie.NewStore(secret)
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

	debugRoutes(app, api)
	healthRoutes(app, api)
	sessionRoutes(app, api)
	v2Routes(app, api)

	guiAssetRoutes(engine, config.Dev(), app.GetLogger())

	api.POST("/query",
		auth.AuthenticateGQL(app.SessionORM(), app.GetLogger().Named("GQLHandler")),
		loader.Middleware(app),
		graphqlHandler(app),
	)

	return engine, nil
}

// Defining the Graphql handler
func graphqlHandler(app chainlink.Application) gin.HandlerFunc {
	rootSchema := schema.MustGetRootSchema()

	// Disable introspection and set a max query depth in production.
	var schemaOpts []graphql.SchemaOpt
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

type SecurityConfig interface {
	AllowOrigins() string
	Dev() bool
	TLSRedirect() bool
	TLSHost() string
}

// secureOptions configure security options for the secure middleware, mostly
// for TLS redirection
func secureOptions(cfg SecurityConfig) secure.Options {
	return secure.Options{
		FrameDeny:     true,
		IsDevelopment: cfg.Dev(),
		SSLRedirect:   cfg.TLSRedirect(),
		SSLHost:       cfg.TLSHost(),
	}
}

// secureMiddleware adds a TLS handler and redirector, to button up security
// for this node
func secureMiddleware(cfg SecurityConfig) gin.HandlerFunc {
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

func debugRoutes(app chainlink.Application, r *gin.RouterGroup) {
	group := r.Group("/debug", auth.Authenticate(app.SessionORM(), auth.AuthenticateBySession))
	group.GET("/vars", expvar.Handler())
}

func metricRoutes(r *gin.RouterGroup, includeHeap bool) {
	pprofGroup := r.Group("/debug/pprof")
	pprofGroup.GET("/", ginHandlerFromHTTP(pprof.Index))
	pprofGroup.GET("/cmdline", ginHandlerFromHTTP(pprof.Cmdline))
	pprofGroup.GET("/profile", ginHandlerFromHTTP(pprof.Profile))
	pprofGroup.POST("/symbol", ginHandlerFromHTTP(pprof.Symbol))
	pprofGroup.GET("/symbol", ginHandlerFromHTTP(pprof.Symbol))
	pprofGroup.GET("/trace", ginHandlerFromHTTP(pprof.Trace))
	pprofGroup.GET("/allocs", ginHandlerFromHTTP(pprof.Handler("allocs").ServeHTTP))
	pprofGroup.GET("/block", ginHandlerFromHTTP(pprof.Handler("block").ServeHTTP))
	pprofGroup.GET("/goroutine", ginHandlerFromHTTP(pprof.Handler("goroutine").ServeHTTP))
	if includeHeap {
		pprofGroup.GET("/heap", ginHandlerFromHTTP(pprof.Handler("heap").ServeHTTP))
	}
	pprofGroup.GET("/mutex", ginHandlerFromHTTP(pprof.Handler("mutex").ServeHTTP))
	pprofGroup.GET("/threadcreate", ginHandlerFromHTTP(pprof.Handler("threadcreate").ServeHTTP))
}

func ginHandlerFromHTTP(h http.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
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
		authv2.GET("/users", auth.RequiresAdminRole(uc.Index))
		authv2.POST("/users", auth.RequiresAdminRole(uc.Create))
		authv2.PATCH("/users", auth.RequiresAdminRole(uc.UpdateRole))
		authv2.DELETE("/users/:email", auth.RequiresAdminRole(uc.Delete))
		authv2.PATCH("/user/password", uc.UpdatePassword)
		authv2.POST("/user/token", uc.NewAPIToken)
		authv2.POST("/user/token/delete", uc.DeleteAPIToken)

		wa := NewWebAuthnController(app)
		authv2.GET("/enroll_webauthn", wa.BeginRegistration)
		authv2.POST("/enroll_webauthn", wa.FinishRegistration)

		eia := ExternalInitiatorsController{app}
		authv2.GET("/external_initiators", paginatedRequest(eia.Index))
		authv2.POST("/external_initiators", auth.RequiresEditRole(eia.Create))
		authv2.DELETE("/external_initiators/:Name", auth.RequiresEditRole(eia.Destroy))

		bt := BridgeTypesController{app}
		authv2.GET("/bridge_types", paginatedRequest(bt.Index))
		authv2.POST("/bridge_types", auth.RequiresEditRole(bt.Create))
		authv2.GET("/bridge_types/:BridgeName", bt.Show)
		authv2.PATCH("/bridge_types/:BridgeName", auth.RequiresEditRole(bt.Update))
		authv2.DELETE("/bridge_types/:BridgeName", auth.RequiresEditRole(bt.Destroy))

		ets := EVMTransfersController{app}
		authv2.POST("/transfers", auth.RequiresAdminRole(ets.Create))
		authv2.POST("/transfers/evm", auth.RequiresAdminRole(ets.Create))
		sts := SolanaTransfersController{app}
		authv2.POST("/transfers/solana", auth.RequiresAdminRole(sts.Create))

		cc := ConfigController{app}
		authv2.GET("/config", cc.Show)
		authv2.PATCH("/config", auth.RequiresAdminRole(cc.Patch))
		authv2.GET("/config/dump-v1-as-v2", cc.Dump)
		authv2.GET("/config/v2", cc.Show2)

		tas := TxAttemptsController{app}
		authv2.GET("/tx_attempts", paginatedRequest(tas.Index))
		authv2.GET("/tx_attempts/evm", paginatedRequest(tas.Index))

		txs := TransactionsController{app}
		authv2.GET("/transactions/evm", paginatedRequest(txs.Index))
		authv2.GET("/transactions/evm/:TxHash", txs.Show)
		authv2.GET("/transactions", paginatedRequest(txs.Index))
		authv2.GET("/transactions/:TxHash", txs.Show)

		rc := ReplayController{app}
		authv2.POST("/replay_from_block/:number", auth.RequiresRunRole(rc.ReplayFromBlock))

		csakc := CSAKeysController{app}
		authv2.GET("/keys/csa", csakc.Index)
		authv2.POST("/keys/csa", auth.RequiresEditRole(csakc.Create))
		authv2.POST("/keys/csa/import", auth.RequiresAdminRole(csakc.Import))
		authv2.POST("/keys/csa/export/:ID", auth.RequiresAdminRole(csakc.Export))

		ekc := NewETHKeysController(app)
		authv2.GET("/keys/eth", ekc.Index)
		authv2.POST("/keys/eth", auth.RequiresEditRole(ekc.Create))
		authv2.PUT("/keys/eth/:keyID", auth.RequiresAdminRole(ekc.Update))
		authv2.DELETE("/keys/eth/:keyID", auth.RequiresAdminRole(ekc.Delete))
		authv2.POST("/keys/eth/import", auth.RequiresAdminRole(ekc.Import))
		authv2.POST("/keys/eth/export/:address", auth.RequiresAdminRole(ekc.Export))
		// duplicated from above, with `evm` instead of `eth`
		// legacy ones remain for backwards compatibility
		authv2.GET("/keys/evm", ekc.Index)
		authv2.POST("/keys/evm", auth.RequiresEditRole(ekc.Create))
		authv2.PUT("/keys/evm/:keyID", auth.RequiresAdminRole(ekc.Update))
		authv2.DELETE("/keys/evm/:keyID", auth.RequiresAdminRole(ekc.Delete))
		authv2.POST("/keys/evm/import", auth.RequiresAdminRole(ekc.Import))
		authv2.POST("/keys/evm/export/:address", auth.RequiresAdminRole(ekc.Export))
		authv2.POST("/keys/evm/chain", auth.RequiresAdminRole(ekc.Chain))

		ocrkc := OCRKeysController{app}
		authv2.GET("/keys/ocr", ocrkc.Index)
		authv2.POST("/keys/ocr", auth.RequiresEditRole(ocrkc.Create))
		authv2.DELETE("/keys/ocr/:keyID", auth.RequiresAdminRole(ocrkc.Delete))
		authv2.POST("/keys/ocr/import", auth.RequiresAdminRole(ocrkc.Import))
		authv2.POST("/keys/ocr/export/:ID", auth.RequiresAdminRole(ocrkc.Export))

		ocr2kc := OCR2KeysController{app}
		authv2.GET("/keys/ocr2", ocr2kc.Index)
		authv2.POST("/keys/ocr2/:chainType", auth.RequiresEditRole(ocr2kc.Create))
		authv2.DELETE("/keys/ocr2/:keyID", auth.RequiresAdminRole(ocr2kc.Delete))
		authv2.POST("/keys/ocr2/import", auth.RequiresAdminRole(ocr2kc.Import))
		authv2.POST("/keys/ocr2/export/:ID", auth.RequiresAdminRole(ocr2kc.Export))

		p2pkc := P2PKeysController{app}
		authv2.GET("/keys/p2p", p2pkc.Index)
		authv2.POST("/keys/p2p", auth.RequiresEditRole(p2pkc.Create))
		authv2.DELETE("/keys/p2p/:keyID", auth.RequiresAdminRole(p2pkc.Delete))
		authv2.POST("/keys/p2p/import", auth.RequiresAdminRole(p2pkc.Import))
		authv2.POST("/keys/p2p/export/:ID", auth.RequiresAdminRole(p2pkc.Export))

		for _, keys := range []struct {
			path string
			kc   KeysController
		}{
			{"solana", NewSolanaKeysController(app)},
			{"starknet", NewStarkNetKeysController(app)},
			{"dkgsign", NewDKGSignKeysController(app)},
			{"dkgencrypt", NewDKGEncryptKeysController(app)},
		} {
			authv2.GET("/keys/"+keys.path, keys.kc.Index)
			authv2.POST("/keys/"+keys.path, auth.RequiresEditRole(keys.kc.Create))
			authv2.DELETE("/keys/"+keys.path+"/:keyID", auth.RequiresAdminRole(keys.kc.Delete))
			authv2.POST("/keys/"+keys.path+"/import", auth.RequiresAdminRole(keys.kc.Import))
			authv2.POST("/keys/"+keys.path+"/export/:ID", auth.RequiresAdminRole(keys.kc.Export))
		}

		vrfkc := VRFKeysController{app}
		authv2.GET("/keys/vrf", vrfkc.Index)
		authv2.POST("/keys/vrf", auth.RequiresEditRole(vrfkc.Create))
		authv2.DELETE("/keys/vrf/:keyID", auth.RequiresAdminRole(vrfkc.Delete))
		authv2.POST("/keys/vrf/import", auth.RequiresAdminRole(vrfkc.Import))
		authv2.POST("/keys/vrf/export/:keyID", auth.RequiresAdminRole(vrfkc.Export))

		jc := JobsController{app}
		authv2.GET("/jobs", paginatedRequest(jc.Index))
		authv2.GET("/jobs/:ID", jc.Show)
		authv2.POST("/jobs", auth.RequiresEditRole(jc.Create))
		authv2.PUT("/jobs/:ID", auth.RequiresEditRole(jc.Update))
		authv2.DELETE("/jobs/:ID", auth.RequiresEditRole(jc.Delete))

		// PipelineRunsController
		authv2.GET("/pipeline/runs", paginatedRequest(prc.Index))
		authv2.GET("/jobs/:ID/runs", paginatedRequest(prc.Index))
		authv2.GET("/jobs/:ID/runs/:runID", prc.Show)

		// FeaturesController
		fc := FeaturesController{app}
		authv2.GET("/features", fc.Index)

		// PipelineJobSpecErrorsController
		authv2.DELETE("/pipeline/job_spec_errors/:ID", auth.RequiresEditRole(psec.Destroy))

		lgc := LogController{app}
		authv2.GET("/log", lgc.Get)
		authv2.PATCH("/log", auth.RequiresAdminRole(lgc.Patch))

		chains := authv2.Group("chains")
		for _, chain := range []struct {
			path string
			cc   ChainsController
		}{
			{"evm", NewEVMChainsController(app)},
			{"solana", NewSolanaChainsController(app)},
			{"starknet", NewStarkNetChainsController(app)},
		} {
			chains.GET(chain.path, paginatedRequest(chain.cc.Index))
			chains.POST(chain.path, auth.RequiresEditRole(chain.cc.Create))
			chains.GET(chain.path+"/:ID", chain.cc.Show)
			chains.PATCH(chain.path+"/:ID", auth.RequiresEditRole(chain.cc.Update))
			chains.DELETE(chain.path+"/:ID", auth.RequiresEditRole(chain.cc.Delete))
		}

		nodes := authv2.Group("nodes")
		for _, chain := range []struct {
			path string
			nc   NodesController
		}{
			{"evm", NewEVMNodesController(app)},
			{"solana", NewSolanaNodesController(app)},
			{"starknet", NewStarkNetNodesController(app)},
		} {
			if chain.path == "evm" {
				// TODO still EVM only https://app.shortcut.com/chainlinklabs/story/26276/multi-chain-type-ui-node-chain-configuration
				nodes.GET("", paginatedRequest(chain.nc.Index))
				nodes.POST("", auth.RequiresEditRole(chain.nc.Create))
				nodes.DELETE("/:ID", auth.RequiresEditRole(chain.nc.Delete))
			}
			nodes.GET(chain.path, paginatedRequest(chain.nc.Index))
			chains.GET(chain.path+"/:ID/nodes", paginatedRequest(chain.nc.Index))
			nodes.POST(chain.path, auth.RequiresEditRole(chain.nc.Create))
			nodes.DELETE(chain.path+"/:ID", auth.RequiresEditRole(chain.nc.Delete))
		}

		efc := EVMForwardersController{app}
		authv2.GET("/nodes/evm/forwarders", paginatedRequest(efc.Index))
		authv2.POST("/nodes/evm/forwarders/track", auth.RequiresEditRole(efc.Track))
		authv2.DELETE("/nodes/evm/forwarders/:fwdID", auth.RequiresEditRole(efc.Delete))

		buildInfo := BuildInfoController{app}
		authv2.GET("/build_info", buildInfo.Show)

		// Debug routes accessible via authentication
		metricRoutes(authv2, app.GetConfig().Dev())
	}

	ping := PingController{app}
	userOrEI := r.Group("/v2", auth.Authenticate(app.SessionORM(),
		auth.AuthenticateExternalInitiator,
		auth.AuthenticateByToken,
		auth.AuthenticateBySession,
	))
	userOrEI.GET("/ping", ping.Show)
	userOrEI.POST("/jobs/:ID/runs", auth.RequiresRunRole(prc.Create))
}

// This is higher because it serves main.js and any static images. There are
// 5 assets which must be served, so this allows for 20 requests/min
var staticAssetsRateLimit = int64(100)
var staticAssetsRateLimitPeriod = 1 * time.Minute
var indexRateLimit = int64(20)
var indexRateLimitPeriod = 1 * time.Minute

// guiAssetRoutes serves the operator UI static files and index.html. Rate
// limiting is disabled when in dev mode.
func guiAssetRoutes(engine *gin.Engine, devMode bool, lggr logger.SugaredLogger) {
	// Serve static files
	var assetsRouterHandlers []gin.HandlerFunc
	if !devMode {
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
	var noRouteHandlers []gin.HandlerFunc
	if !devMode {
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
		defer lggr.ErrorIfFn(file.Close, "Error closing file")

		http.ServeContent(c.Writer, c.Request, path, time.Time{}, file)
	})

	engine.NoRoute(noRouteHandlers...)
}

// Inspired by https://github.com/gin-gonic/gin/issues/961
func loggerFunc(lggr logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		buf, err := io.ReadAll(c.Request.Body)
		if err != nil {
			lggr.Error("Web request log error: ", err.Error())
			// Implicitly relies on limits.RequestSizeLimiter
			// overriding of c.Request.Body to abort gin's Context
			// inside io.ReadAll.
			// Functions as we would like, but horrible from an architecture
			// and design pattern perspective.
			if !c.IsAborted() {
				c.AbortWithStatus(http.StatusBadRequest)
			}
			return
		}
		rdr := bytes.NewBuffer(buf)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(buf))

		start := time.Now()
		c.Next()
		end := time.Now()

		lggr.Debugw(fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.Path),
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
func uiCorsHandler(config SecurityConfig) gin.HandlerFunc {
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
