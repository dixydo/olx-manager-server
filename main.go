package main

import (
	"github.com/dixydo/olxmanager-server/controllers"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/middleware/basicauth"
	"github.com/iris-contrib/middleware/cors"

)

type User struct {
	Username string
	Password string
}

func main() {
	app := newApp()
	app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
}

func newApp() *iris.Application {
	app := iris.New()

	opts := basicauth.Options{
		Allow: basicauth.AllowUsersFile("users.yml", basicauth.BCRYPT),
		Realm: basicauth.DefaultRealm,
		// [...more options]
	}
	
	auth := basicauth.New(opts)

	app.Use(auth)
	app.UseError(auth)

	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	app.Use(crs)

	adverts := mvc.New(app.Party("/adverts"))
	adverts.Handle(new(controllers.AdvertController))

	return app
}
