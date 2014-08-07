package missions

import (
	"appengine/datastore"
)

type LocationZone struct {
	ZoneId datastore.Key `json:"zone_id" xml:"zone-id"`
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
	CharacterId datastore.Key `json:"character_id" xml:"character-id"`
}
