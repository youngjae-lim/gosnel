package gosnel

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/CloudyKit/jet/v6"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/gomodule/redigo/redis"
	"github.com/joho/godotenv"
	"github.com/youngjae-lim/gosnel/cache"
	"github.com/youngjae-lim/gosnel/render"
	"github.com/youngjae-lim/gosnel/session"
)

const version = "1.0.0"

var myRedisCache *cache.RedisCache

type Gosnel struct {
	AppName       string
	Debug         bool
	Version       string
	ErrorLog      *log.Logger
	InfoLog       *log.Logger
	RootPath      string
	Routes        *chi.Mux
	Render        *render.Render
	Session       *scs.SessionManager
	DB            Database
	JetViews      *jet.Set
	config        config
	EncryptionKey string
	Cache         cache.Cache
}

type config struct {
	port        string
	renderer    string
	cookie      cookieConfig
	sessionType string
	database    databaseConfig
	redis       redisConfig
}

func (g *Gosnel) New(rootPath string) error {
	pathConfig := initPaths{
		rootPath:    rootPath,
		folderNames: []string{"handlers", "migrations", "views", "data", "public", "tmp", "logs", "middleware"},
	}

	err := g.Init(pathConfig)
	if err != nil {
		return err
	}

	err = g.checkDotEnv(rootPath)
	if err != nil {
		return err
	}

	// read .env file
	err = godotenv.Load(rootPath + "/.env")
	if err != nil {
		return err
	}

	// create loggers
	infoLog, errorLog := g.startLoggers()

	// connect to database if database type is specified in the .env
	if os.Getenv("DATABASE_TYPE") != "" {
		db, err := g.OpenDB(os.Getenv("DATABASE_TYPE"), g.BuildDSN())
		if err != nil {
			errorLog.Println(err)
			os.Exit(1)
		}
		g.DB = Database{
			DbType: os.Getenv("DATABASE_TYPE"),
			Pool:   db,
		}
	}

	// create a client redis cache if redis is specified in the .env
	if os.Getenv("CACHE") == "redis" || os.Getenv("SESSION_TYPE") == "redis" {
		myRedisCache = g.createClientRedisCache()
		g.Cache = myRedisCache
	}

	g.InfoLog = infoLog
	g.ErrorLog = errorLog

	g.Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
	g.Version = version
	g.RootPath = rootPath
	g.Routes = g.routes().(*chi.Mux)

	g.config = config{
		port:     os.Getenv("PORT"),
		renderer: os.Getenv("RENDERER"),
		cookie: cookieConfig{
			name:     os.Getenv("COOKIE_NAME"),
			lifetime: os.Getenv("COOKIE_LIFETIME"),
			persist:  os.Getenv("COOKIE_PERSIST"),
			secure:   os.Getenv("COOKIE_SECURE"),
			domain:   os.Getenv("COOKIE_DOMAIN"),
		},
		sessionType: os.Getenv("SESSION_TYPE"),
		database: databaseConfig{
			dsn:      g.BuildDSN(),
			database: os.Getenv("DATABASE_TYPE"),
		},
		redis: redisConfig{
			host:     os.Getenv("REDIS_HOST"),
			password: os.Getenv("REDIS_PASSWORD"),
			prefix:   os.Getenv("REDIS_PREFIX"),
		},
	}

	// create session
	sess := session.Session{
		CookieName:     g.config.cookie.name,
		CookieLifetime: g.config.cookie.lifetime,
		CookiePersist:  g.config.cookie.persist,
		CookieSecure:   g.config.cookie.secure,
		CookieDomain:   g.config.cookie.domain,
		SessionType:    g.config.sessionType,
	}

	switch g.config.sessionType {
	case "redis":
		sess.RedisPool = myRedisCache.Conn
	case "mysql", "postgres", "postgresql", "mariadb":
		sess.DBPool = g.DB.Pool
	}

	g.Session = sess.InitSession()
	g.EncryptionKey = os.Getenv("KEY")

	var views = jet.NewSet(
		jet.NewOSFileSystemLoader(fmt.Sprintf("%s/views", rootPath)),
		jet.InDevelopmentMode(),
	)

	g.JetViews = views

	g.createRenderer()

	return nil
}

// Init creates necessary directories for our Gosnel application
func (g *Gosnel) Init(p initPaths) error {
	root := p.rootPath
	for _, path := range p.folderNames {
		// create a folder if it doesn't exist
		err := g.CreateDirIfNotExist(root + "/" + path)
		if err != nil {
			return nil
		}
	}
	return nil
}

// ListenAndServe starts the web server
func (g *Gosnel) ListenAndServe() {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", os.Getenv("PORT")),
		ErrorLog:     g.ErrorLog,
		Handler:      g.Routes,
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 600 * time.Second,
	}

	defer g.DB.Pool.Close()

	g.InfoLog.Printf("Listening on port %s", os.Getenv("PORT"))
	err := srv.ListenAndServe()
	g.ErrorLog.Fatal(err)
}

func (g *Gosnel) checkDotEnv(path string) error {
	err := g.CreateFileIfNotExist(fmt.Sprintf("%s/.env", path))
	if err != nil {
		return err
	}
	return nil
}

func (g *Gosnel) startLoggers() (*log.Logger, *log.Logger) {
	var infoLog *log.Logger
	var errorLog *log.Logger

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	return infoLog, errorLog
}

func (g *Gosnel) createRenderer() {
	myRenderer := render.Render{
		Renderer: g.config.renderer,
		RootPath: g.RootPath,
		Port:     g.config.port,
		JetViews: g.JetViews,
		Session:  g.Session,
	}
	g.Render = &myRenderer
}

func (g *Gosnel) createClientRedisCache() *cache.RedisCache {
	cacheClient := cache.RedisCache{
		Conn:   g.createRedisPool(),
		Prefix: g.config.redis.prefix,
	}
	return &cacheClient
}

func (g *Gosnel) createRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     50,
		MaxActive:   10000,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial(
				"tcp",
				g.config.redis.host,
				redis.DialPassword(g.config.redis.password))
		},
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			_, err := conn.Do("PING")
			return err
		},
	}
}

// BuildDSN builds the datasource name for our database, and returns it as a string
func (g *Gosnel) BuildDSN() string {
	var dsn string

	switch os.Getenv("DATABASE_TYPE") {
	case "postgres", "postgresql":
		dsn = fmt.Sprintf("host=%s port=%s user=%s database=%s sslmode=%s timezone=UTC connect_timeout=5",
			os.Getenv("DATABASE_HOST"),
			os.Getenv("DATABASE_PORT"),
			os.Getenv("DATABASE_USER"),
			os.Getenv("DATABASE_NAME"),
			os.Getenv("DATABASE_SSL_MODE"),
		)

		// postgres://postgres:password@0.0.0.0:5432/gosnel?sslmode=disable
		if os.Getenv("DATABASE_PASS") != "" {
			dsn = fmt.Sprintf("%s password=%s", dsn, os.Getenv("DATABASE_PASS"))
		}

	default:
	}

	return dsn
}
