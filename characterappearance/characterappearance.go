package characterappearance

import (
	"appengine/datastore"
	"github.com/emicklei/go-restful"
	//	"log"
	"net/http"
	"resource"
)

// Status values for records of this resource type.
const (
	StatusActive = iota
	StatusDeactivated
	StatusPendingActivation
	StatusDeletionPending
	StatusDeleted
)

type Shallow struct {
	resource.BaseResource
}

type Mouth struct {
	Position          float32 `json:"position" xml:"position"`
	Width             float32 `json:"width" xml:"width"`
	LipThicknessUpper float32 `json:"lip_thickness_upper" xml:"lip-thickness-upper"`
	LipThicknessLower float32 `json:"lip_thickness_lower" xml:"lip_thickness_lower"`
	Gape              float32 `json:"gape" xml:"gape"`
	Pucker            float32 `json:"pucker" xml:"pucker"`
	MouthCorners      float32 `json:"mouth_corners" xml:"mouth-corners"`
}

type Nose struct {
	Extension    float32 `json:"extension" xml:"extension"`
	NozeSize     float32 `json:"noze_size" xml:"noze-size"`
	Bridge       float32 `json:"bridge" xml:"bridge"`
	NostrilWidth float32 `json:"nostril_width" xml:"nostril-width"`
	TipWidth     float32 `json:"tip_width" xml:"tip-width"`
	NoseTip      float32 `json:"nose_tip" xml:"nose-tip"`
	NostrilFlare float32 `json:"nostril_flare" xml:"nostril-flare"`
	Bend         float32 `json:"bend" xml:"bend"`
}

type Ears struct {
	EarRotation  float32 `json:"ear_rotation" xml:"ear-rotation"`
	EarExtension float32 `json:"ear_extension" xml:"ear-extension"`
	EarTrim      float32 `json:"ear_trim" xml:"ear-trim"`
	EarSize      float32 `json:"ear_size" xml:"ear-size"`
}

type Eyes struct {
	Colour     float32 `json:"colour" xml:"colour"`
	Width      float32 `json:"width" xml:"width"`
	Height     float32 `json:"height" xml:"height"`
	Shape      float32 `json:"shape" xml:"shape"`
	Separation float32 `json:"separation" xml:"separation"`
	Angle      float32 `json:"angle" xml:"angle"`
	InnerBrow  float32 `json:"inner_brow" xml:"inner-brow"`
	OuterBrow  float32 `json:"outer_brow" xml:"outer-brow"`
}

type Eyebrows struct {
	BrowThickness float32 `json:"brow_thickness" xml:"brow-thickness"`
	BrowPlacement float32 `json:"brow_placement" xml:"brow-placement"`
}

type Skin struct {
	TextureHead string `json:"texture_head" xml:"texture-head"`
	TextureBody string `json:"texture_body" xml:"texture-body"`
}

type Skar struct {
	SkarHead string `json:"skar_head" xml:"skar-head"`
	SkarBody string `json:"skar_body" xml:"skar-body"`
}

type Tattoo struct {
	TattooHead string `json:"tattoo_head" xml:"tattoo-head"`
	TattooBody string `json:"tattoo_body" xml:"tattoo-body"`
}

type Model struct {
	HeadModel  int32   `json:"head_model" xml:"head-model"`
	TorsoModel int32   `json:"torso_model" xml:"torso-model"`
	Height     float32 `json:"height" xml:"height"`
	Weight     float32 `json:"weight" xml:"weight"`
	Phsyique   float32 `json:"phsyique" xml:"phsyique"`
}

type Jaw struct {
	Width  float32 `json:"width" xml:"width"`
	Length float32 `json:"length" xml:"length"`
	Jut    float32 `json:"jut" xml:"jut"`
}

type Cheeks struct {
	Width  float32 `json:"width" xml:"width"`
	Height float32 `json:"height" xml:"height"`
}

type Hair struct {
	HairStyle       int32 `json:"hair_style" xml:"hair-style"`
	PrimaryColour   int32 `json:"primary_colour" xml:"primary-colour"`
	SecondaryColour int32 `json:"secondary_colour" xml:"secondary-colour"`
}

type Resource struct {
	Shallow
	CharacterKey datastore.Key `json:"character_key" xml:"character-key"`
	Model        Model         `json:"model" xml:"model"`
	Mouth        Mouth         `json:"mouth" xml:"mouth"`
	Nose         Nose          `json:"nose" xml:"nose"`
	Ears         Ears          `json:"ears" xml:"ears"`
	Eyes         Eyes          `json:"eyes" xml:"eyes"`
	Skar         Skar          `json:"skar" xml:"skar"`
	Skin         Skin          `json:"skin" xml:"skin"`
	Tattoo       Tattoo        `json:"tattoo" xml:"tattoo"`
	Jaw          Jaw           `json:"jaw" xml:"jaw"`
	Eyebrows     Eyebrows      `json:"eyebrows" xml:"eyebrows"`
	Cheeks       Cheeks        `json:"cheeks" xml:"cheeks"`
	Hair         Hair          `json:"hair" xml:"hair"`
	Voice        int32         `json:"voice" xml:"voice"`
}

type ListSummary struct {
	Entry []Shallow `json:"entry" xml:"entry"`
}

type ListComprehensive struct {
	Entry []Resource `json:"entry" xml:"entry"`
}

func (s Shallow) Kind() string {
	return "characterappearance"
}

func (s Shallow) AdminRootPath() string {
	return "/api/client/admin/" + s.Kind()
}

func (s Shallow) ShardRootPath() string {
	return "/api/shard/" + s.Kind()
}

func (s Shallow) ClientRootPath() string {
	return "/api/client/" + s.Kind()
}

func init() {
}

// Register the routes we require for this resource type.
//
func (res Resource) Register() {
	res.RegisterAdmin()
	res.RegisterShard()
	res.RegisterClient()
}

// Attempts to create a valid key for a resource.
//
func (res *Resource) getKey(r *restful.Request, w *restful.Response) (*datastore.Key, error) {

	// Decode the request parameter to determine the key for the entity.
	k, err := datastore.DecodeKey(r.PathParameter("resource-id"))
	if err != nil {
		w.AddHeader("Content-Type", "text/plain")
		w.WriteErrorString(http.StatusBadRequest, "The key is not valid.\n")
		return nil, err
	}

	// Check for shenanigans with the key.
	if k.Kind() != res.Kind() {
		w.AddHeader("Content-Type", "text/plain")
		w.WriteErrorString(http.StatusBadRequest, "The key is not valid for this type of resource.\n")
		return nil, err
	}

	return k, nil
}
