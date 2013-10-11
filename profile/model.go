package profile

import ( 
	"fmt"
	"time"
)

type Profile struct {
	id        int // just to remind me it should have it somewhere
	Date time.Time
	firstName string
	lastName  string
	nickName  string
}

func (p *Profile) FirstName() string {
	return p.firstName
}

func (p *Profile) SetFirstName(name string) string {
	p.firstName = name
	return name
}

func (p *Profile) LastName() string {
	return p.lastName
}

func (p *Profile) NickName() string {
	return p.nickName
}
