package main

import (
	"accounts"
	"appengine"
	"appengine/user"
	"buffs"
	"characterappearance"
	"characterlocation"
	"characters"
	"debuffs"
	"factions"
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
	"items"
	"itemtemplates"
	"loottables"
	"missions"
	"net/http"
	"profiles"
	"realms"
	"skills"
	"storagecontainers"
	"storageitems"
	"users"
	"verbs"
	"zones"
)

func gaeUrl() string {
	if appengine.IsDevAppServer() {
		//		return "http://localhost:8080"
		return "http://localhost:8080"
	} else {
		// Include your URL on App Engine here.
		// I found no way to get AppID without appengine.Context and that is always
		// based on a http.Request.
		return "http://api.shatteredscreens.com"
	}
}

func init() {
	// Enable CORS on the default container.
	cors := restful.CrossOriginResourceSharing{
		ExposeHeaders:  []string{"X-My-Header"},
		AllowedHeaders: []string{"Content-Type"},
		CookiesAllowed: false,
		Container:      restful.DefaultContainer}
	restful.Filter(cors.Filter)

	// Add container filter to respond to OPTIONS
	//wsContainer.Filter(wsContainer.OPTIONSFilter)

	// Register all the routes we need and document the interface.
	accounts.Register()
	characters.Register()
	factions.Register()
	profiles.Register()
	realms.Register()
	users.Register()
	zones.Register()

	characterappearance.Register()
	characterlocation.Register()
	loottables.Register()
	storagecontainers.Register()
	storageitems.Register()

	// Game design elements.
	buffs.Register()
	debuffs.Register()
	items.Register()
	itemtemplates.Register()
	missions.Register()
	skills.Register()

	// Test page.
	registerMyTest()

	// Register the verbs.
	verbs.RegisterVerbRegister()

	// Optionally, you can install the Swagger Service which provides a nice Web UI on your REST API
	// You need to download the Swagger HTML5 assets and change the FilePath location in the config below.
	// Open <your_app_id>.appspot.com/apidocs and enter
	// http://<your_app_id>.appspot.com/apidocs.json in the api input field.
	config := swagger.Config{
		// You control what services are visible
		WebServices:    restful.RegisteredWebServices(),
		WebServicesUrl: gaeUrl(),
		ApiPath:        "/apidocs.json",

		// Optionally, specifiy where the UI is located
		SwaggerPath: "/apidoc/",

		// GAE support static content which is configured in your app.yaml.
		// This example expect the swagger-ui in static/swagger so you should place it there :)
		SwaggerFilePath: "static/swagger/"}
	swagger.InstallSwaggerService(config)

}

func registerLoginRequired() {
	ws := new(restful.WebService)
	ws.Path("/_ah")
	ws.Route(ws.GET("/login_required").To(loginRequired).
		Doc("Federated user login page").
		Operation("Federated user login page"))
	restful.Add(ws)
}

// Present a page for federated logins.
//
func loginRequired(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	providers := map[string]string{
		"Google": "www.google.com/accounts/o8/id",
		"Steam":  "steamcommunity.com/openid", // no email with this one
		"Yahoo":  "yahoo.com",
		"AOL":    "aol.com",
	}

	// Provide support for redirection back to the page they came from.
	cont := r.QueryParameter("continue")
	if len(cont) == 0 {
		cont = "/"
	}

	// Header.
	fmt.Fprintf(w, "<html><head></head><body>")

	for name, url := range providers {
		login_url, err := user.LoginURLFederated(c, cont, url)
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(w, "<a href='%s'>%s</a><br>", login_url, name)
	}

	// Footer.
	fmt.Fprintf(w, "</body></html>")
}

func registerMyTest() {
	ws := new(restful.WebService)
	ws.Path("/test")
	ws.Route(ws.GET("/my").To(myTestGet).
		Doc("A test page for logged in users").
		Operation("A test page for logged in users"))
	restful.Add(ws)
}

// Test page.
//
func myTestGet(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)
	// Header.
	fmt.Fprintf(w, "<html><head></head><body>")

	// fmt.Fprintf(w, "ID = %s<br>Email = %s<br>AuthDomain = %s<br>Admin = %t<br>FederatedIdentity = %s<br>FederatedProvider = %s<br>",
	// 	user.Current(c).ID,
	// 	user.Current(c).Email,
	// 	user.Current(c).AuthDomain,
	// 	user.Current(c).Admin,
	// 	user.Current(c).FederatedIdentity,
	// 	user.Current(c).FederatedProvider)

	u, err := user.CurrentOAuth(c, "")
	if err != nil {
		http.Error(w, "OAuth Authorization header required", http.StatusUnauthorized)
		return
	}
	// if !u.Admin {
	//     http.Error(w, "Admin login only", http.StatusUnauthorized)
	//     return
	// }
	fmt.Fprintf(w, `Welcome, user %s!`, u)

	// Footer.
	fmt.Fprintf(w, "</body></html>")
}
