package profile

import (
	"time"
)

type profile struct {
	created   time.Time
	firstName string
	lastName  string
	nickName  string
	Account   string
}

func NewProfile(FirstName, NickName, LastName, Account string) *profile {
	return &profile{time.Now(), FirstName, NickName, LastName, Account}
}

func (p *profile) FirstName() string {
	return p.firstName
}

func (p *profile) LastName() string {
	return p.lastName
}

func (p *profile) NickName() string {
	return p.nickName
}
