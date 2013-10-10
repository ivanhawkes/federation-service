package profile

import "fmt"

func Hello() {
	fmt.Println ("profile library")
}

type Profile struct {
	id int // just to remind me it should have it somewhere
	firstName string
	lastName string
	nickName string
}

func (p *Profile) FirstName() string {
	return p.firstName}

func (p *Profile) SetFirstName(name string) string {
	p.firstName = name
	return name}

func (p *Profile) LastName() string {
	return p.lastName}

func (p *Profile) NickName() string {
	return p.nickName}

func (p *Profile) Create() int {
	fmt.Println ("create")
	return 1}

func (p *Profile) Read() int {
	fmt.Println ("read")
	return 1}

func (p *Profile) Update() int {
	fmt.Println ("update")
	return 1}

func (p *Profile) Delete() int {
	fmt.Println ("save")
	return 1}

