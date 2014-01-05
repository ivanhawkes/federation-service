package common

import (
	"time"
)

type Link struct {
	Rel  string `datastore:"-" json:"rel" xml:"rel"`
	Href string `datastore:"-" json:"href" xml:"href"`
}

type BaseResource struct {
	Key          string    `datastore:"-" json:"key" xml:"key"`
	LastModified time.Time `json:"last_modified" xml:"last-modified"`
	Version      int       `json:"version" xml:"version"`
	Status       int       `json:"status" xml:"status"`
	Link         Link      `datastore:"-" json:"link" xml:"link"`
}

