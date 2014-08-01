package missions

import (
	"appengine"
	"appengine/datastore"
)

type LocationZone struct {
	Zone string `json:"zone" xml:"Zone"`
}

type LocationXYZ struct {
	LocationZone
	X      real `json:"x" xml:"X"`
	Y      real `json:"y" xml:"Y"`
	Z      real `json:"Z" xml:"Z"`
	Radius real `json:"radius" xml:"Radius"`
}

type LocationArea struct {
	LocationZone
	Area string `json:"area" xml:"Area"`
}

type LocationCharacter struct {
	LocationZone
	CharacterId datastore.Key `json:"area" xml:"Area"`
}
