package profile

import ( 
	"time"
)

type Profile struct {
	id        int // just to remind me it should have it somewhere
	createDate time.Time
	firstName string
	lastName  string
	nickName  string
	Account string
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
