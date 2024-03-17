package app_test

import (
	"echo-test/app"
	"echo-test/pkg/testutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

func TestMain(m *testing.M) {
	testutil.ResetDB()

	testutil.CreateUser("test@example.com", "password12")

	// Run tests
	code := m.Run()

	os.Exit(code)
}

func TestEchoClient(t *testing.T) {
	srv := app.EchoHandler()

	server := httptest.NewServer(srv.Mux)
	defer server.Close()

	e := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  server.URL,
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})

	testEcho(e)
}

func TestEchoHandler(t *testing.T) {
	srv := app.EchoHandler()

	e := httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewBinder(srv.Mux),
			Jar:       httpexpect.NewCookieJar(),
		},
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})

	testEcho(e)
}

// Echo JWT token authentication tests.
//
// This test is executed for the EchoHandler() in two modes:
//   - via http client
//   - via http.Handler
func testEcho(e *httpexpect.Expect) {
	type Login struct {
		Email    string `form:"email"`
		Password string `form:"password"`
	}

	e.POST("/login").WithForm(Login{"ford", "<bad password>"}).
		Expect().
		Status(http.StatusUnauthorized)

	r := e.POST("/login").WithForm(Login{"test@example.com", "password12"}).
		Expect().
		Status(http.StatusOK).JSON().Object()

	r.Keys().ContainsOnly("token")

	token := r.Value("token").String().Raw()

	e.GET("/restricted/hello").
		Expect().
		Status(http.StatusUnauthorized)

	e.GET("/restricted/hello").WithHeader("Authorization", "Bearer <bad token>").
		Expect().
		Status(http.StatusUnauthorized)

	e.GET("/restricted/hello").WithHeader("Authorization", "Bearer "+token).
		Expect().
		Status(http.StatusOK).Body().Contains("Welcome")

	auth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+token)
	})

	auth.GET("/restricted/hello").
		Expect().
		Status(http.StatusOK).Body().Contains("Welcome")
}
