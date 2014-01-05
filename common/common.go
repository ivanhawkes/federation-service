package common

import (
	"time"
)

type Link struct {
	Rel  string `datastore:"-" json:"rel" xml:"rel"`
	Href string `datastore:"-" json:"href" xml:"href"`
}

type Versioning struct {
	LastModified time.Time `json:"last_modified" xml:"last-modified"`
	Revision      int       `json:"revision" xml:"revision"`
}

type BaseResource struct {
	Versioning
	Key          string    `datastore:"-" json:"key" xml:"key"`
	Status       int       `json:"status" xml:"status"`
	Link         Link      `datastore:"-" json:"link" xml:"link"`
}

