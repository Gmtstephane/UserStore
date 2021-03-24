package userStoreAPi

import (
	"UserStore/pkg/config"
	database "UserStore/pkg/database"
	redis "UserStore/pkg/store"
	"fmt"
	"net/http"

	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/labstack/echo-contrib/session"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rbcervilla/redisstore/v8"
	"gorm.io/gorm"
)

const (
	BASEPATH = "api/v1"
)

type Userstore struct {
	db     *gorm.DB
	redis  *redisstore.RedisStore
	router *echo.Echo
}

func (u *Userstore) RegisterRoutes() {
	fmt.Println("Registering routes")
	//r := mux.NewRouter()
	e := echo.New()

	e.Use(middleware.Logger())

	e.Use(session.MiddlewareWithConfig(session.Config{
		Store: u.redis,
	}))

	e.Use(middleware.Recover())
	e.POST(BASEPATH+"/login", u.loginHandler)
	e.POST(BASEPATH+"/signup", u.signUpHandler)
	e.GET("/swagger/", echoSwagger.WrapHandler)
	e.GET(BASEPATH+"/secret", u.secret)
	u.router = e
}
func (u *Userstore) Serve() {
	//log.Fatal(http.ListenAndServe(":8080", u.router))
	u.router.Logger.Fatal(u.router.Start(":8080"))
}

func NewStore(appconfig config.Config) Userstore {
	u := Userstore{}
	u.db = database.SetupOrDie(&appconfig.Database)
	u.redis = redis.InitOrDie(&appconfig.Redis)
	u.RegisterRoutes()
	return u
}

func (u *Userstore) secret(c echo.Context) error {
	session, _ := u.redis.Get(c.Request(), "Auth-cookie")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(c.Response().Writer, "Forbidden", http.StatusForbidden)
		return nil
	}
	fmt.Fprintln(c.Response().Writer, "The cake is a lie!")
	fmt.Fprintln(c.Response().Writer, session.Values["email"])
	fmt.Fprintln(c.Response().Writer, session.Values["id"])
	return nil
}
