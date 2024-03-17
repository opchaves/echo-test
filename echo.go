package main

import (
	"context"
	"echo-test/config"
	"echo-test/model"
	"echo-test/pkg/password"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// jwtCustomClaims are custom claims extending default ones.
// See https://github.com/golang-jwt/jwt for more examples
type jwtCustomClaims struct {
	UserId string `json:"id"`
	Admin  bool   `json:"admin"`
	jwt.RegisteredClaims
}

type Server struct {
	Db *pgxpool.Pool
	Q  *model.Queries
}

func EchoHandler() *echo.Echo {
	db, err := pgxpool.New(context.Background(), config.DatabaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	log.Println("Connected to database")

	s := &Server{
		Db: db,
		Q:  model.New(db),
	}

	e := echo.New()
	e.HideBanner = true
	e.Logger.SetLevel(config.LogLevel)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Login route
	e.POST("/login", s.login)

	// Unauthenticated route
	e.GET("/", accessible)

	// Restricted group
	r := e.Group("/restricted")

	// Configure middleware with the custom claims type
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtCustomClaims)
		},
		SigningKey: []byte(config.JwtSecret),
	}
	r.Use(echojwt.WithConfig(config))
	r.GET("/hello", restricted)

	return e
}

func (s *Server) login(c echo.Context) error {
	email := c.FormValue("email")
	pass := c.FormValue("password")

	user, err := s.Q.GetUserByEmail(context.Background(), email)

	if err != nil {
		c.Logger().Error(err)
		return echo.ErrUnauthorized
	}

	match, err := password.Compare(user.Password, pass)

	if !match || err != nil {
		c.Logger().Warn("Password mismatch")
		return echo.ErrUnauthorized
	}

	// Set custom claims
	claims := &jwtCustomClaims{
		UserId: user.ID.String(),
		Admin:  true,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.JwtExpiry)),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(config.JwtSecret))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

func accessible(c echo.Context) error {
	return c.String(http.StatusOK, "Accessible")
}

func restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	id := claims.UserId
	return c.String(http.StatusOK, "Welcome "+id+"!")
}
