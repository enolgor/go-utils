package examples

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/enolgor/go-utils/parse"
	"github.com/enolgor/go-utils/sec"
	"github.com/enolgor/go-utils/server"
	"github.com/go-utils/jwtauth"
	"github.com/google/uuid"
)

var hashedPasswords map[string]string = map[string]string{}

func init() {
	cost := sec.OptimalCost(250 * time.Millisecond)
	hashedPasswords["admin"], _ = sec.HashPassword("test", cost)
}

func Server(port int) {
	key := parse.Must(parse.HexBytes)("bc27bec0c4291b4e43a2ec657d8afc9b668e158c6acd4004ffb1faa16c5b88bf")
	signer := jwtauth.NewSigner(key, 5*time.Minute)

	login := loginHandler(signer)
	hello := helloHandler(signer)
	router := server.NewRouter().
		SubRoute("/users", server.NewRouter().
			Get("/login", form).
			Post("/login", login)).
		Get("/(.*)", hello).
		PreFilters(requestLog, getUserFilter(signer)).
		PostFilters(responseLog)
	http.ListenAndServe(fmt.Sprintf(":%d", port), router)
}

func form(w http.ResponseWriter, req *http.Request) {
	server.Response(w).Status(http.StatusOK).WithBody(`
	<!DOCTYPE html>
	<html lang="en">
		<head>
			<meta charset="UTF-8">
		</head>
		<body>
			<h1>Login</h1>
			<form action="/users/login" method="post">
				<label for="user">User:</label>
				<input type="text" id="user" name="user"><br><br>
				<label for="pass">Password:</label>
				<input type="password" id="pass" name="pass"><br><br>
				<input type="submit" value="Authenticate">
			</form>
		</body>
	</html>
	`).AsHtml()
}

func loginHandler(signer *jwtauth.Signer) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		if err := req.ParseForm(); err != nil {
			server.Response(w).Status(http.StatusBadRequest).WithBody(err).AsTextPlain()
			return
		}
		var errs error
		if !req.Form.Has("user") {
			errs = errors.Join(errs, errors.New("user missing in form"))
		}
		if !req.Form.Has("pass") {
			errs = errors.Join(errs, errors.New("pass missing in form"))
		}
		if errs != nil {
			server.Response(w).Status(http.StatusBadRequest).WithBody(errs).AsTextPlain()
			return
		}
		var hash string
		var ok bool
		if hash, ok = hashedPasswords[req.Form.Get("user")]; !ok {
			server.Response(w).Status(http.StatusBadRequest).WithBody(errors.New("user not found")).AsTextPlain()
			return
		}
		if sec.ComparePassword(hash, req.Form.Get("pass")) != nil {
			server.Response(w).Status(http.StatusBadRequest).WithBody(errors.New("wrong password")).AsTextPlain()
			return
		}
		token, exp, err := signer.ForgeToken(req.Form.Get("user"))
		if err != nil {
			server.Response(w).Status(http.StatusInternalServerError).WithBody(err).AsTextPlain()
			return
		}
		server.Response(w).WithCookie(signer.CreateCookie(token, exp)).Redirect("/hello")
	}
}

type contextKey int

const (
	userContextKey contextKey = iota
	requestIDKey
)

func getUserFilter(signer *jwtauth.Signer) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		token, err := signer.GetFromRequest(req)
		if err != nil {
			return
		}
		subject, err := token.Claims.GetSubject()
		if err != nil {
			return
		}
		server.AddContextValue(req, userContextKey, subject)
	}
}

func requestLog(w http.ResponseWriter, req *http.Request) {
	id := uuid.NewString()
	server.AddContextValue(req, requestIDKey, id)
	fmt.Printf("[%s] %s %s\n", id, req.Method, req.URL.Path)
}

func responseLog(w http.ResponseWriter, req *http.Request) {
	var id string
	server.GetContextValue(req, requestIDKey, &id)
	rw := w.(server.ResponseWriter)
	fmt.Printf("[%s] %d - %d\n", id, rw.Status(), rw.Size())
}

func helloHandler(signer *jwtauth.Signer) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		var subject string
		if ok := server.GetContextValue(req, userContextKey, &subject); !ok {
			server.Response(w).Status(http.StatusUnauthorized).WithBody(errors.New("unauthorized")).AsTextPlain()
			return
		}
		pathParams := server.PathParams(req)
		if pathParams[0] == "panic" {
			panic("example panic!")
		}
		data := struct {
			User string `json:"user"`
		}{User: subject}
		server.Response(w).WithBody(data).AsJson()
	}

}
