package main

import (
	"github.com/fasthttp/router"
	"github.com/juliovcruz/user-register/cmd/api/handlers"
	_ "github.com/juliovcruz/user-register/docs"
	"github.com/juliovcruz/user-register/internal/mailvalidation"
	"github.com/juliovcruz/user-register/internal/mailvalidation/sender"
	"github.com/juliovcruz/user-register/internal/platform/database"
	"github.com/juliovcruz/user-register/internal/security/hash"
	"github.com/juliovcruz/user-register/internal/security/token"
	"github.com/juliovcruz/user-register/internal/settings"
	"github.com/juliovcruz/user-register/internal/users"
	"github.com/juliovcruz/user-register/internal/users/zipcode"
	"github.com/juliovcruz/user-register/internal/users/zipcode/viacep"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

// @title User Register API
// @version 1.0
// @description API para registro de usuários
// @host localhost:8080
// @BasePath /
func main() {
	sett, err := settings.LoadSettings(settings.Local)
	if err != nil {
		panic(err)
	}

	tokenService := token.NewService(sett)
	zipCodeService := zipcode.NewService(viacep.NewClient(sett.ZipCodeSettings))
	hashService := hash.NewService(sett)

	db, err := database.NewDatabase(sett.Database.FilePath, sett.Database.Driver)
	if err != nil {
		panic(err)
	}

	userRepository, err := users.NewSQLiteRepository(db)
	if err != nil {
		panic(err)
	}

	mailValidationService := mailvalidation.NewService(mailvalidation.NewRepository(db), sender.NewClient(), sett.MailValidationExpirationTime)

	userService := users.NewService(userRepository, tokenService, zipCodeService, hashService, mailValidationService)
	userHandler := handlers.NewUserHandler(userService, tokenService)
	r := router.New()

	r.POST("/users", userHandler.CreateUser)
	r.GET("/users", userHandler.JWTMiddleware(userHandler.ListUsers))
	r.PUT("/users/password", userHandler.UpdatePassword)
	r.POST("/users/forgot_password", userHandler.ForgotPassword)

	r.POST("/login", userHandler.Login)

	r.GET("/{filepath:*}", fasthttpadaptor.NewFastHTTPHandler(httpSwagger.WrapHandler))

	println("Server on port 8080")
	if err := fasthttp.ListenAndServe(":8080", handlers.CorsMiddleware(r.Handler)); err != nil {
		panic(err)
	}
}
