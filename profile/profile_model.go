package profile

import (
	//"appengine/datastore"
	"time"
)

type profile struct {
	LastModified  time.Time     `json:"last_modified"`
	//ApplicationId datastore.Key `json:"application_id"`
	//AccountId     datastore.Key `json:"account_id"`
	FirstName     string        `json:"first_name"`
	NickName      string        `json:"nick_name"`
	LastName      string        `json:"last_name"`
	// We might need an int to append to nickname to make it unique like in GW2 e.g. socks.451
}

func NewProfile(FirstName, NickName, LastName string) *profile {
	return &profile{time.Now(),FirstName, NickName, LastName}
}
