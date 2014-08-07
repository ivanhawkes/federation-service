package main

import (
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
	"zones"
)

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
}

func registerMyTest() {
	ws := new(restful.WebService)
	ws.Path("/test")
	ws.Route(ws.GET("/my").To(myTestGet).
		Doc("A test page for logged in users").
		Operation("A test page for logged in users"))
	restful.Add(ws)
}

// Present a page for federated logins.
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
