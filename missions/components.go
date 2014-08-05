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
	X      float32 `json:"x" xml:"X"`
	Y      float32 `json:"y" xml:"Y"`
	Z      float32 `json:"Z" xml:"Z"`
	Radius float32 `json:"radius" xml:"Radius"`
}

type LocationArea struct {
	LocationZone
	Area string `json:"area" xml:"Area"`
}

type LocationCharacter struct {
	LocationZone
	CharacterId int64 `json:"area" xml:"Area"`
}
