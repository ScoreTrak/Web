package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	forcessl "github.com/gobuffalo/mw-forcessl"
	paramlogger "github.com/gobuffalo/mw-paramlogger"
	"github.com/unrolled/secure"

	"github.com/L1ghtman2k/ScoreTrakWeb/models"
	"github.com/gobuffalo/buffalo-pop/v2/pop/popmw"
	contenttype "github.com/gobuffalo/mw-contenttype"
	"github.com/gobuffalo/x/sessions"
	"github.com/rs/cors"

	csrf "github.com/gobuffalo/mw-csrf"

)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
//
// Routing, middleware, groups, etc... are declared TOP -> DOWN.
// This means if you add a middleware to `app` *after* declaring a
// group, that group will NOT have that new middleware. The same
// is true of resource declarations as well.
//
// It also means that routes are checked in the order they are declared.
// `ServeFiles` is a CATCH-ALL route, so it should always be
// placed last in the route declarations, as it will prevent routes
// declared after it to never be called.
func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{
            Env:          ENV,
            SessionStore: sessions.Null{},
            PreWares: []buffalo.PreWare{
                cors.Default().Handler,
            },
            SessionName: "_scoretrak_session",
        })

		// Automatically redirect to SSL
		app.Use(forceSSL())

		// Log request parameters (filters apply).
		app.Use(paramlogger.ParameterLogger)

		app.Use(contenttype.Set("application/json"))

		// Protect against CSRF attacks. https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)
		// Remove to disable this.
		app.Use(csrf.New)

		// Wraps each request in a transaction.
		//  c.Value("tx").(*pop.Connection)
		// Remove to disable this.
		app.Use(popmw.Transaction(models.DB))

		app.GET("/", HomeHandler)

		//AuthMiddlewares
		app.Use(SetCurrentUser)
		app.Use(Authorize)
		app.Use(AuthorizeBlackTeam)


		//Routes for Auth
		auth := app.Group("/auth")
		auth.GET("/", AuthLanding)
		auth.GET("/new", AuthNew)
		auth.POST("/", AuthCreate)
		auth.DELETE("/", AuthDestroy)

		auth.Middleware.Skip(Authorize, AuthLanding, AuthNew, AuthCreate)
		auth.Middleware.Skip(AuthorizeBlackTeam, AuthLanding, AuthNew, AuthCreate, AuthDestroy)

		//Routes for User registration (Disabled. Let admin create new accounts)
		users := app.Group("/users")
		users.GET("/new", UsersNew)
		users.GET("/", UsersList)
		users.GET("/{user_id}", UsersShow)
		users.PUT("/{user_id}", UsersUpdate)
		users.DELETE("/{user_id}", UsersDestroy)
		users.GET("/{user_id}/edit", UsersEdit)
		users.POST("/", UsersCreate)
		res := TeamsResource{}
		tr := app.Resource("/teams", res)
		tr.Middleware.Skip(AuthorizeBlackTeam, res.Show)
	}
	return app
}


// forceSSL will return a middleware that will redirect an incoming request
// if it is not HTTPS. "http://example.com" => "https://example.com".
// This middleware does **not** enable SSL. for your application. To do that
// we recommend using a proxy: https://gobuffalo.io/en/docs/proxy
// for more information: https://github.com/unrolled/secure/
func forceSSL() buffalo.MiddlewareFunc {
	return forcessl.Middleware(secure.Options{
		SSLRedirect:     ENV == "production",
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
	})
}
