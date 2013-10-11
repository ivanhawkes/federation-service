package main

import (
	//"bitbucket.org/shatteredscreens/federationservices/character"
	//"bitbucket.org/shatteredscreens/federationservices/federation"
	//"bitbucket.org/shatteredscreens/federationservices/loot"
	//"bitbucket.org/shatteredscreens/federationservices/profile"
	//"bitbucket.org/shatteredscreens/federationservices/realm"
	//"bitbucket.org/shatteredscreens/federationservices/zone"
	//"html/template"
	//"profile"
	"appengine"
	"fmt"
	"net/http"
)

func init() {
	http.HandleFunc("/", handleMainPage)
}

func WebapiProfile(w http.ResponseWriter, r *http.Request) {
	//c := appengine.NewContext(r)

	//p := profile.NewProfile("Ivan", "Socks", "Hawkes", user.Current(c).String())

	//_, err := datastore.Put(c, datastore.NewIncompleteKey(c, "Profile", nil), &p)
	//if err != nil {
	//http.Error(w, err.Error(), http.StatusInternalServerError)
	//return
	//}

	return
}

func handleMainPage(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	fmt.Fprintf(w, "appengine.AppID(c) = %q\n", appengine.AppID(c))
	fmt.Fprintf(w, "appengine.VersionID(c) = %q\n", appengine.VersionID(c))
}
