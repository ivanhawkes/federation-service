package missions

// import (
// 	"appengine"
// 	"appengine/datastore"
// )

type LocationZone struct {
	ZoneId int64 `json:"zone-id" xml:"ZoneId"`
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
	CharacterId int64 `json:"area" xml:"Area"`
}
