package main

import (
	"bitbucket.org/shatteredscreens/federalservices/character"
	"bitbucket.org/shatteredscreens/federalservices/federation"
	"bitbucket.org/shatteredscreens/federalservices/loot"
	"bitbucket.org/shatteredscreens/federalservices/profile"
	"bitbucket.org/shatteredscreens/federalservices/realm"
	"bitbucket.org/shatteredscreens/federalservices/zone"
	"fmt"
	"html/template"
	"net/http"
)

func init() {
	http.HandleFunc("/", handleMainPage)
	http.HandleFunc("/profile", handlePut)
}

func main() {
	fmt.Println("main")

	p := new(profile.Profile)
	p.SetFirstName("Ivan")
	fmt.Println(p.FirstName())
	p.Create()
}
