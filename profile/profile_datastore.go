package profile

import (
	"appengine"
	"appengine/datastore"
	//"appengine/user"
	"fmt"
	//"time"
)

// guestbookKey returns the key used for all guestbook entries.
func profileKey(c appengine.Context) *datastore.Key {
	// The string "default_guestbook" here could be varied to have multiple guestbooks.
	return datastore.NewKey(c, "ShatteredScreens", "stringID", 0, nil)
}

func (p *profile) Get() int {
	fmt.Println("get")
	return 1
}

func (p *profile) Put() int {
	fmt.Println("put")
	return 1
}

func (p *profile) Post() int {
	fmt.Println("post")
	return 1
}

func (p *profile) Delete() int {
	fmt.Println("delete")
	return 1
}
