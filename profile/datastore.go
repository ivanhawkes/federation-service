package profile

import (
	"appengine"
	"appengine/datastore"
	"appengine/user"
	"fmt"
	"time"
)

// guestbookKey returns the key used for all guestbook entries.
func profileKey(c appengine.Context) *datastore.Key {
	// The string "default_guestbook" here could be varied to have multiple guestbooks.
	return datastore.NewKey(c, "ShatteredScreens", "stringID", 0, nil)
}


func (p *Profile) Create() int {
	fmt.Println("create")
	return 1
}

func (p *Profile) Read() int {
	fmt.Println("read")
	return 1
}

func (p *Profile) Update() int {
	fmt.Println("update")
	return 1
}

func (p *Profile) Delete() int {
	fmt.Println("delete")
	return 1
}
