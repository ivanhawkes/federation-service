package characterappearances

import (
	"errors"
	"federation-services/accounts"
	"federation-services/resource"
	"net/http"
	"strconv"
	"time"

	"github.com/emicklei/go-restful"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

// The various states for a federation resource.
const (
	StatusActivationPending = iota
	StatusActive
	StatusDeactivated
	StatusDeletionPending
	StatusDeleted
)

const (
	Kind     = "characterappearances"
	RootPath = "/federation/" + Kind
)

func PreferredLink(k *datastore.Key) string {
	return RootPath + "/" + k.Encode()
}

type Api struct {
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

type ResourceMeta struct {
	Api
	resource.Meta
	Resource
}

type ResourceKey struct {
	Api
	Key datastore.Key `datastore:"-" json:"key" xml:"key"`
	Resource
}

type ResourceRequest struct {
	Api
	Resource
}

type ResourceResponse struct {
	ResourceMeta
}

type ListResource struct {
	Entry []Resource `json:"entry" xml:"entry"`
}

type ListResourceKey struct {
	Entry []ResourceKey `json:"entry" xml:"entry"`
}

type ListResourceMeta struct {
	Entry []ResourceMeta `json:"entry" xml:"entry"`
}

// Register the routes we require for this resource type.
func Register() {
	ws := new(restful.WebService)

	ws.
		Path(RootPath).
		Consumes(restful.MIME_JSON, restful.MIME_XML).
		Produces(restful.MIME_JSON, restful.MIME_XML).
		Doc("CharacterAppearance management.")

	ws.Route(ws.POST("").To(post).
		Doc("Create a new resource").
		Operation("postCharacterAppearance").
		Param(ws.BodyParameter("CharacterAppearance.Resource", "representation of a resource").DataType("CharacterAppearance.Resource")).
		Reads(ResourceRequest{}).
		Writes(ResourceResponse{}))

	ws.Route(ws.PUT("/{resource-id}").To(put).
		Doc("Update an existing resource").
		Operation("putCharacterAppearance").
		Param(ws.PathParameter("resource-id", "key for an existing resource").DataType("string")).
		Param(ws.BodyParameter("CharacterAppearance.Resource", "representation of a resource").DataType("CharacterAppearance.Resource")).
		Param(ws.HeaderParameter("If-Unmodified-Since", "Conditional modifier").DataType("RFC3339Nano Date")).
		Reads(ResourceRequest{}))

	ws.Route(ws.GET("/{resource-id}").To(get).
		Doc("Read a resource").
		Operation("getCharacterAppearance").
		Param(ws.PathParameter("resource-id", "key for an existing resource").DataType("string")).
		Param(ws.HeaderParameter("If-Modified-Since", "Optional conditional modifier").DataType("RFC3339Nano Date")).
		Writes(ResourceResponse{}))

	ws.Route(ws.HEAD("/{resource-id}").To(head).
		Doc("Returns the headers for a resource").
		Operation("headCharacterAppearance").
		Param(ws.PathParameter("resource-id", "key for an existing resource").DataType("string")).
		Param(ws.HeaderParameter("If-Modified-Since", "Optional conditional modifier").DataType("RFC3339Nano Date")))

	ws.Route(ws.GET("/list").To(listAll).
		Doc("Get a list of resources").
		Operation("listCharacterAppearance").
		Param(ws.HeaderParameter("If-Modified-Since", "Optional conditional modifier").DataType("RFC3339Nano Date")).
		Writes(ResourceResponse{}))

	restful.Add(ws)
}

// Attempts to create a valid key for a resource.
func getKey(r *restful.Request, w *restful.Response) (*datastore.Key, error) {
	// Decode the request parameter to determine the key for the entity.
	k, err := datastore.DecodeKey(r.PathParameter("resource-id"))
	if err != nil {
		resource.WriteError(w, resource.NewError(http.StatusBadRequest, "/html/error/invalidkey", "The key is not valid."))
		return nil, err
	}

	// Check for shenanigans with the key.
	if k.Kind() != Kind {
		resource.WriteError(w, resource.NewError(http.StatusBadRequest, "/html/error/invalidkeymalicious", "The key is not valid for this type of resource."))
		return nil, errors.New("The key is not valid for this type of resource.")
	}

	return k, nil
}

// Create a new resource.
func post(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	// Auth check.
	if err := accounts.AccessLevelGE(r, w, accounts.AccessLevelAdmin); err != nil {
		return
	}

	// Marshall the entity from the request into a struct.
	res := new(ResourceMeta)
	err := r.ReadEntity(res)
	if err != nil {
		resource.WriteError(w, resource.NewError(http.StatusNotAcceptable, "/html/error/statusnotacceptable", err.Error()))
		return
	}

	// Set the meta data for the resource.
	res.LastModified = time.Now()
	res.Status = StatusActive
	res.Revision = 1

	// TODO: The OwnerKey must be set at this point!

	// Store the resource.
	k, err := datastore.Put(c, datastore.NewIncompleteKey(c, Kind, nil), res)
	if err != nil {
		resource.WriteError(w, resource.NewError(http.StatusInternalServerError, "/html/error/statusinternalservererror", err.Error()))
		return
	}

	// The resource Key.
	res.Key = *k

	// Let them know the location of the newly created resource.
	w.AddHeader("Location", PreferredLink(k))

	// Provide a link for ease of API usage.
	res.Link.Rel = "self"
	res.Link.Href = PreferredLink(k)

	// Set the headers.
	w.WriteHeader(http.StatusCreated)
	w.AddHeader(restful.HEADER_LastModified, res.LastModified.Format(time.RFC3339Nano))
	w.AddHeader("ETag", strconv.Itoa(res.Revision))

	// Output the response body.
	w.WriteEntity(res)
}

// Update the resource.
func put(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	// Auth check.
	if err := accounts.AccessLevelGE(r, w, accounts.AccessLevelAdmin); err != nil {
		return
	}

	// Grab the key and validate it.
	if k, err := getKey(r, w); err != nil {
		return
	} else {

		// Marshall the entity from the request into a struct.
		res := new(ResourceResponse)
		err = r.ReadEntity(res)
		if err != nil {
			resource.WriteError(w, resource.NewError(http.StatusNotAcceptable, "/html/error/statusnotacceptable", err.Error()))
			return
		}

		// Retrieve the old entity from the datastore.
		old := new(ResourceResponse)
		if err := datastore.Get(c, k, old); err != nil {
			if err.Error() == "datastore: no such entity" {
				resource.WriteError(w, resource.NewError(http.StatusNotFound, "/html/error/statusnotfound", err.Error()))
				return
			} else {
				resource.WriteError(w, resource.NewError(http.StatusInternalServerError, "/html/error/statusinternalservererror", err.Error()))
				return
			}
			return
		}

		// Use conditional put - check LastModified before doing anything.
		if ifUnmodifiedSince := r.HeaderParameter("If-Unmodified-Since"); ifUnmodifiedSince == "" {
			resource.WriteError(w, resource.NewError(http.StatusForbidden, "/html/error/statusforbidden", "Unconditional updates are not supported. Please provide 'If-Unmodified-Since' headers."))
			return
		} else {
			if t, err := time.Parse(time.RFC3339Nano, ifUnmodifiedSince); err != nil {
				resource.WriteError(w, resource.NewError(http.StatusNotAcceptable, "/html/error/statusnotacceptable", err.Error()))
				return
			} else {
				if t.Before(old.LastModified) {
					resource.WriteError(w, resource.NewError(http.StatusPreconditionFailed, "/html/error/statuspreconditionfailed", "The resource has been modified recently. Refresh your copy and try again if updating is still desireable."))
					return
				}
			}
		}

		// Keep track of the last modification date.
		res.LastModified = time.Now()
		res.Revision = old.Revision + 1

		// Attempt to overwrite the old entity.
		_, err = datastore.Put(c, k, res)
		if err != nil {
			resource.WriteError(w, resource.NewError(http.StatusInternalServerError, "/html/error/statusinternalservererror", err.Error()))
			return
		}

		// Set the headers.
		w.AddHeader(restful.HEADER_LastModified, res.LastModified.Format(time.RFC3339Nano))
		w.AddHeader("ETag", strconv.Itoa(res.Revision))

		// Let them know it succeeded.
		w.WriteHeader(http.StatusNoContent)
		w.WriteEntity(nil)
	}
}

// Read a resource
func get(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	// Grab the key and validate it.
	if k, err := getKey(r, w); err != nil {
		return
	} else {
		// Handle cases where they want to check if the resource has been modified since...
		if ok, err := resource.IfModifiedSince(r, w, c, Kind, k); err != nil || ok != true {
			return
		}

		// Retrieve the entity from the datastore.
		res := new(ResourceResponse)
		if err := datastore.Get(c, k, res); err != nil {
			if err.Error() == "datastore: no such entity" {
				resource.WriteError(w, resource.NewError(http.StatusNotFound, "/html/error/statusnotfound", err.Error()))
				return
			} else {
				resource.WriteError(w, resource.NewError(http.StatusInternalServerError, "/html/error/statusinternalservererror", err.Error()))
				return
			}
			return
		}

		// Set their Key.
		res.Key = *k

		// Provide a link for ease of API usage.
		res.Link.Rel = "self"
		res.Link.Href = PreferredLink(k)

		// Set the headers.
		w.AddHeader(restful.HEADER_LastModified, res.LastModified.Format(time.RFC3339Nano))
		w.AddHeader("ETag", strconv.Itoa(res.Revision))

		// Cache Control: By allowing a short cache time here we can reduce database calls and cost.
		w.AddHeader("Cache-Control", "max-age=14400,must-revalidate")

		// Output the response body.
		w.WriteEntity(res)
	}
}

// Returns the headers for a resource
func head(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	// Grab the key and validate it.
	if k, err := getKey(r, w); err != nil {
		return
	} else {
		// Handle cases where they want to check if the resource has been modified since...
		if ok, err := resource.IfModifiedSince(r, w, c, Kind, k); err != nil || ok != true {
			return
		}

		// Retrieve the entity from the datastore.
		res := new(ResourceResponse)
		if err := datastore.Get(c, k, res); err != nil {
			if err.Error() == "datastore: no such entity" {
				resource.WriteError(w, resource.NewError(http.StatusNotFound, "/html/error/statusnotfound", err.Error()))
				return
			} else {
				resource.WriteError(w, resource.NewError(http.StatusInternalServerError, "/html/error/statusinternalservererror", err.Error()))
				return
			}
			return
		}

		// Set their Key.
		res.Key = *k

		// Provide a link for ease of API usage.
		res.Link.Rel = "self"
		res.Link.Href = PreferredLink(k)

		// Set the headers.
		w.AddHeader(restful.HEADER_LastModified, res.LastModified.Format(time.RFC3339Nano))
		w.AddHeader("ETag", strconv.Itoa(res.Revision))

		// Cache Control: By allowing a short cache time here we can reduce database calls and cost.
		w.AddHeader("Cache-Control", "max-age=14400,must-revalidate")

		// No response body required for this verb.
		w.WriteHeader(http.StatusNoContent)
		w.WriteEntity(nil)
	}
}

// Read a resource
func listAll(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	var result ListResourceKey
	var q *datastore.Query

	// Check if they want to limit the query using a modified since date.
	if ifModifiedSince := r.HeaderParameter("If-Modified-Since"); ifModifiedSince == "" {
		q = datastore.NewQuery(Kind).
			Project("Category", "Subcategory", "Description").
			Filter("Status =", StatusActive)
	} else {
		if t, err := time.Parse(time.RFC3339Nano, ifModifiedSince); err != nil {
			w.AddHeader("Content-Type", "text/plain")
			w.WriteErrorString(http.StatusNotAcceptable, err.Error())
			return
		} else {
			q = datastore.NewQuery(Kind).
				Project("Category", "Subcategory", "Description").
				Filter("Status =", StatusActive).
				Filter("LastModified >", t)
		}
	}

	if keys, err := q.GetAll(c, &result.Entry); err != nil {
		w.AddHeader("Content-Type", "text/plain")
		w.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	} else {
		for i, k := range keys {
			// TODO: why does XML not emit the key correctly (JSON does)?
			result.Entry[i].Key = *k
		}
	}

	// Cache Control: By allowing a short cache time here we can reduce database calls and cost.
	//	w.AddHeader("Cache-Control", "max-age=900,must-revalidate")
	//	w.AddHeader(restful.HEADER_LastModified, time.Now().Format(time.RFC3339Nano))

	w.WriteEntity(result)
}

// ***************************************
// SAVED FOR CUTTING OUT THE USEFUL PIECES
// ***************************************

// // Register the routes we require for this resource type.
// //
// func Register() {
// 	ws := new(restful.WebService)

// 	ws.
// 		Path(RootPath).
// 		Consumes(restful.MIME_JSON, restful.MIME_XML).
// 		Produces(restful.MIME_JSON, restful.MIME_XML).
// 		Doc("CharacterAppearance management.")

// 	ws.Route(ws.POST("").To(post).
// 		Doc("Create a new resource").
// 		Operation("postCharacterAppearance").
// 		Param(ws.BodyParameter("characterappearance.Resource", "representation of a resource").DataType("characterappearance.Resource")).
// 		Reads(ResourceRequest{}).
// 		Writes(ResourceResponse{}))

// 	ws.Route(ws.PUT("/{resource-id}").To(put).
// 		Doc("Update an existing resource").
// 		Operation("putCharacterAppearance").
// 		Param(ws.PathParameter("resource-id", "key for an existing resource").DataType("string")).
// 		Param(ws.BodyParameter("characterappearance.Resource", "representation of a resource").DataType("characterappearance.Resource")).
// 		Param(ws.HeaderParameter("If-Unmodified-Since", "Conditional modifier").DataType("RFC3339Nano Date")).
// 		Reads(ResourceRequest{}))

// 	ws.Route(ws.GET("/{resource-id}").To(get).
// 		Doc("Read a resource").
// 		Operation("getCharacterAppearance").
// 		Param(ws.PathParameter("resource-id", "key for an existing resource").DataType("string")).
// 		Param(ws.HeaderParameter("If-Modified-Since", "Optional conditional modifier").DataType("RFC3339Nano Date")).
// 		Writes(ResourceResponse{}))

// 	ws.Route(ws.HEAD("/{resource-id}").To(head).
// 		Doc("Returns the headers for a resource").
// 		Operation("headCharacterAppearance").
// 		Param(ws.PathParameter("resource-id", "key for an existing resource").DataType("string")).
// 		Param(ws.HeaderParameter("If-Modified-Since", "Optional conditional modifier").DataType("RFC3339Nano Date")))

// 	ws.Route(ws.GET("/list").To(listAll).
// 		Doc("Get a list of resources").
// 		Operation("listCharacterAppearance").
// 		Param(ws.HeaderParameter("If-Modified-Since", "Optional conditional modifier").DataType("RFC3339Nano Date")).
// 		Writes(ResourceResponse{}))

// 	restful.Add(ws)
// }

// // Attempts to create a valid key for a resource.
// //
// func (res *Resource) getKey(r *restful.Request, w *restful.Response) (*datastore.Key, error) {

// 	// Decode the request parameter to determine the key for the entity.
// 	k, err := datastore.DecodeKey(r.PathParameter("resource-id"))
// 	if err != nil {
// 		w.AddHeader("Content-Type", "text/plain")
// 		w.WriteErrorString(http.StatusBadRequest, "The key is not valid.\n")
// 		return nil, err
// 	}

// 	// Check for shenanigans with the key.
// 	if k.Kind() != res.Kind() {
// 		w.AddHeader("Content-Type", "text/plain")
// 		w.WriteErrorString(http.StatusBadRequest, "The key is not valid for this type of resource.\n")
// 		return nil, err
// 	}

// 	return k, nil
// }

// // Delete the resource.
// //
// func (res *Resource) delete(r *restful.Request, w *restful.Response) {
// 	c := appengine.NewContext(r.Request)

// 	// Grab the key and validate it.
// 	if k, err := res.getKey(r, w); err != nil {
// 		return
// 	} else {

// 		// Retrieve the old entity from the datastore.
// 		old := new(Resource)
// 		if err := datastore.Get(c, k, old); err != nil {
// 			if err.Error() == "datastore: no such entity" {
// 				w.AddHeader("Content-Type", "text/plain")
// 				w.WriteErrorString(http.StatusNotFound, err.Error())
// 				return
// 			} else {
// 				w.AddHeader("Content-Type", "text/plain")
// 				w.WriteErrorString(http.StatusInternalServerError, err.Error())
// 				return
// 			}
// 			return
// 		}

// 		// Use conditional delete - check LastModified before doing anything.
// 		if ifUnmodifiedSince := r.HeaderParameter("If-Unmodified-Since"); ifUnmodifiedSince == "" {
// 			w.AddHeader("Content-Type", "text/plain")
// 			w.WriteErrorString(http.StatusForbidden, "Unconditional deletes are not supported. Please provide 'If-Unmodified-Since' headers.")
// 			return
// 		} else {
// 			if t, err := time.Parse(time.RFC3339Nano, ifUnmodifiedSince); err != nil {
// 				w.AddHeader("Content-Type", "text/plain")
// 				w.WriteErrorString(http.StatusNotAcceptable, err.Error())
// 				return
// 			} else {
// 				if t.Before(old.LastModified) {
// 					w.AddHeader("Content-Type", "text/plain")
// 					w.WriteErrorString(http.StatusPreconditionFailed, "The resource has been modified recently. Refresh your copy and try again if deletion is still desireable.")
// 					return
// 				}
// 			}
// 		}

// 		// Delete the entity.
// 		if err := datastore.Delete(c, k); err != nil {
// 			w.AddHeader("Content-Type", "text/plain")
// 			w.WriteErrorString(http.StatusInternalServerError, err.Error())
// 			return
// 		}

// 		// Success notification.
// 		w.WriteHeader(http.StatusNoContent)
// 		w.WriteEntity(nil)
// 	}
// }

// // Create a new resource.
// //
// func (res *Resource) post(r *restful.Request, w *restful.Response) {
// 	c := appengine.NewContext(r.Request)

// 	// Marshall the entity from the request into a struct.
// 	err := r.ReadEntity(res)
// 	if err != nil {
// 		w.AddHeader("Content-Type", "text/plain")
// 		w.WriteErrorString(http.StatusNotAcceptable, err.Error())
// 		return
// 	}

// 	// Set some fields that need special handling.
// 	res.LastModified = time.Now()
// 	res.Status = StatusActive
// 	res.Revision = 1

// 	// Store the resource.
// 	k, err := datastore.Put(c, datastore.NewIncompleteKey(c, res.Kind(), nil), res)
// 	if err != nil {
// 		w.AddHeader("Content-Type", "text/plain")
// 		w.WriteErrorString(http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	// The resource Key.
// 	res.Key = k.Encode()

// 	// Let them know the location of the newly created resource.
// 	// TODO: Use a safe Url path append function.
// 	w.AddHeader("Location", res.PreferredLink(k))

// 	// Provide a link for ease of API usage.
// 	res.Link.Rel = "self"
// 	res.Link.Href = res.PreferredLink(k)

// 	// Set the headers.
// 	w.WriteHeader(http.StatusCreated)
// 	w.AddHeader(restful.HEADER_LastModified, res.LastModified.Format(time.RFC3339Nano))
// 	w.AddHeader("ETag", strconv.Itoa(res.Revision))

// 	// Output the response body.
// 	w.WriteEntity(res)
// }

// // Update the resource.
// //
// func (res *Resource) put(r *restful.Request, w *restful.Response) {
// 	c := appengine.NewContext(r.Request)

// 	// Grab the key and validate it.
// 	if k, err := res.getKey(r, w); err != nil {
// 		return
// 	} else {

// 		// Marshall the entity from the request into a struct.
// 		err = r.ReadEntity(res)
// 		if err != nil {
// 			w.AddHeader("Content-Type", "text/plain")
// 			w.WriteErrorString(http.StatusNotAcceptable, err.Error())
// 			return
// 		}

// 		// Retrieve the old entity from the datastore.
// 		old := new(Resource)
// 		if err := datastore.Get(c, k, old); err != nil {
// 			if err.Error() == "datastore: no such entity" {
// 				w.AddHeader("Content-Type", "text/plain")
// 				w.WriteErrorString(http.StatusNotFound, err.Error())
// 				return
// 			} else {
// 				w.AddHeader("Content-Type", "text/plain")
// 				w.WriteErrorString(http.StatusInternalServerError, err.Error())
// 				return
// 			}
// 			return
// 		}

// 		// Use conditional put - check LastModified before doing anything.
// 		if ifUnmodifiedSince := r.HeaderParameter("If-Unmodified-Since"); ifUnmodifiedSince == "" {
// 			w.AddHeader("Content-Type", "text/plain")
// 			w.WriteErrorString(http.StatusForbidden, "Unconditional updates are not supported. Please provide 'If-Unmodified-Since' headers.")
// 			return
// 		} else {
// 			if t, err := time.Parse(time.RFC3339Nano, ifUnmodifiedSince); err != nil {
// 				w.AddHeader("Content-Type", "text/plain")
// 				w.WriteErrorString(http.StatusNotAcceptable, err.Error())
// 				return
// 			} else {
// 				if t.Before(old.LastModified) {
// 					w.AddHeader("Content-Type", "text/plain")
// 					w.WriteErrorString(http.StatusPreconditionFailed, "The resource has been modified recently. Refresh your copy and try again if updating is still desireable.")
// 					return
// 				}
// 			}
// 		}

// 		// Keep track of the last modification date.
// 		res.LastModified = time.Now()
// 		res.Revision = old.Revision + 1

// 		// Attempt to overwrite the old entity.
// 		_, err = datastore.Put(c, k, res)
// 		if err != nil {
// 			w.AddHeader("Content-Type", "text/plain")
// 			w.WriteErrorString(http.StatusInternalServerError, err.Error())
// 			return
// 		}

// 		// Set the headers.
// 		w.AddHeader(restful.HEADER_LastModified, res.LastModified.Format(time.RFC3339Nano))
// 		w.AddHeader("ETag", strconv.Itoa(res.Revision))

// 		// Let them know it succeeded.
// 		w.WriteHeader(http.StatusNoContent)
// 		w.WriteEntity(nil)
// 	}
// }

// // Read a resource
// //
// func (res *Resource) get(r *restful.Request, w *restful.Response) {
// 	c := appengine.NewContext(r.Request)

// 	// Grab the key and validate it.
// 	if k, err := res.getKey(r, w); err != nil {
// 		return
// 	} else {

// 		// Check if they want to limit the query using a modified since date.
// 		if ifModifiedSince := r.HeaderParameter("If-Modified-Since"); ifModifiedSince != "" {
// 			if t, err := time.Parse(time.RFC3339Nano, ifModifiedSince); err != nil {
// 				w.AddHeader("Content-Type", "text/plain")
// 				w.WriteErrorString(http.StatusNotAcceptable, err.Error())
// 				return
// 			} else {
// 				// Check the versioning information and see if we can tell them the
// 				// resource is unmodified.
// 				q := datastore.NewQuery(res.Kind()).
// 					Filter("__key__ =", k).
// 					Filter("LastModified >", t).
// 					KeysOnly()

// 				if keys, err := q.GetAll(c, nil); err != nil {
// 					w.AddHeader("Content-Type", "text/plain")
// 					w.WriteErrorString(http.StatusInternalServerError, err.Error())
// 					return
// 				} else {

// 					if len(keys) == 0 {
// 						// Not modified.
// 						w.WriteHeader(http.StatusNotModified)
// 						w.WriteEntity(nil)
// 						return
// 					}
// 				}
// 			}
// 		}

// 		// Retrieve the entity from the datastore.
// 		if err := datastore.Get(c, k, res); err != nil {
// 			if err.Error() == "datastore: no such entity" {
// 				w.AddHeader("Content-Type", "text/plain")
// 				w.WriteErrorString(http.StatusNotFound, err.Error())
// 				return
// 			} else {
// 				w.AddHeader("Content-Type", "text/plain")
// 				w.WriteErrorString(http.StatusInternalServerError, err.Error())
// 				return
// 			}
// 			return
// 		}

// 		// Set their Key.
// 		res.Key = k.Encode()

// 		// Provide a link for ease of API usage.
// 		res.Link.Rel = "self"
// 		res.Link.Href = res.PreferredLink(k)

// 		// Set the headers.
// 		w.AddHeader(restful.HEADER_LastModified, res.LastModified.Format(time.RFC3339Nano))
// 		w.AddHeader("ETag", strconv.Itoa(res.Revision))

// 		// Cache Control: By allowing a short cache time here we can reduce database calls and cost.
// 		w.AddHeader("Cache-Control", "max-age=14400,must-revalidate")

// 		// Output the response body.
// 		w.WriteEntity(res)
// 	}
// }

// // Returns the headers for a resource
// //
// func (res *Resource) head(r *restful.Request, w *restful.Response) {
// 	c := appengine.NewContext(r.Request)

// 	// Grab the key and validate it.
// 	if k, err := res.getKey(r, w); err != nil {
// 		return
// 	} else {

// 		// Check if they want to limit the query using a modified since date.
// 		if ifModifiedSince := r.HeaderParameter("If-Modified-Since"); ifModifiedSince != "" {
// 			if t, err := time.Parse(time.RFC3339Nano, ifModifiedSince); err != nil {
// 				w.AddHeader("Content-Type", "text/plain")
// 				w.WriteErrorString(http.StatusNotAcceptable, err.Error())
// 				return
// 			} else {
// 				// Check the versioning information and see if we can tell them the
// 				// resource is unmodified.
// 				q := datastore.NewQuery(res.Kind()).
// 					Filter("__key__ =", k).
// 					Filter("LastModified >", t).
// 					KeysOnly()

// 				if keys, err := q.GetAll(c, nil); err != nil {
// 					w.AddHeader("Content-Type", "text/plain")
// 					w.WriteErrorString(http.StatusInternalServerError, err.Error())
// 					return
// 				} else {

// 					if len(keys) == 0 {
// 						// Not modified.
// 						w.WriteHeader(http.StatusNotModified)
// 						w.WriteEntity(nil)
// 						return
// 					}
// 				}
// 			}
// 		}

// 		// Retrieve the entity from the datastore.
// 		if err := datastore.Get(c, k, res); err != nil {
// 			if err.Error() == "datastore: no such entity" {
// 				w.AddHeader("Content-Type", "text/plain")
// 				w.WriteErrorString(http.StatusNotFound, err.Error())
// 				return
// 			} else {
// 				w.AddHeader("Content-Type", "text/plain")
// 				w.WriteErrorString(http.StatusInternalServerError, err.Error())
// 				return
// 			}
// 			return
// 		}

// 		// Set their Key.
// 		res.Key = k.Encode()

// 		// Provide a link for ease of API usage.
// 		res.Link.Rel = "self"
// 		res.Link.Href = res.PreferredLink(k)

// 		// Set the headers.
// 		w.AddHeader(restful.HEADER_LastModified, res.LastModified.Format(time.RFC3339Nano))
// 		w.AddHeader("ETag", strconv.Itoa(res.Revision))

// 		// Cache Control: By allowing a short cache time here we can reduce database calls and cost.
// 		w.AddHeader("Cache-Control", "max-age=14400,must-revalidate")

// 		// No response body required for this verb.
// 		w.WriteHeader(http.StatusNoContent)
// 		w.WriteEntity(nil)
// 	}
// }

// // Summary list of all resources
// //
// func (res *Resource) listSummary(r *restful.Request, w *restful.Response) {
// 	c := appengine.NewContext(r.Request)

// 	var q *datastore.Query
// 	var result ListSummary

// 	// Check if they want to limit the query using a modified since date.
// 	if ifModifiedSince := r.HeaderParameter("If-Modified-Since"); ifModifiedSince == "" {
// 		q = datastore.NewQuery(res.Kind()).
// 			Project("LastModified", "Revision", "Status", "Name")
// 	} else {
// 		if t, err := time.Parse(time.RFC3339Nano, ifModifiedSince); err != nil {
// 			w.AddHeader("Content-Type", "text/plain")
// 			w.WriteErrorString(http.StatusNotAcceptable, err.Error())
// 			return
// 		} else {
// 			q = datastore.NewQuery(res.Kind()).
// 				Project("LastModified", "Revision", "Status", "Name").
// 				Filter("LastModified >", t)
// 		}
// 	}

// 	if keys, err := q.GetAll(c, &result.Entry); err != nil {
// 		w.AddHeader("Content-Type", "text/plain")
// 		w.WriteErrorString(http.StatusInternalServerError, err.Error())
// 		return
// 	} else {
// 		for i, k := range keys {
// 			result.Entry[i].Key = k.Encode()
// 			result.Entry[i].Link.Rel = "self"
// 			result.Entry[i].Link.Href = result.Entry[i].PreferredLink(k)
// 		}
// 	}

// 	// Cache Control: By allowing a short cache time here we can reduce database calls and cost.
// 	w.AddHeader("Cache-Control", "max-age=900,must-revalidate")
// 	w.AddHeader(restful.HEADER_LastModified, time.Now().Format(time.RFC3339Nano))

// 	w.WriteEntity(result)

// }

// // Comprehensive list of all resources
// //
// func (res *Resource) listAll(r *restful.Request, w *restful.Response) {
// 	c := appengine.NewContext(r.Request)

// 	var result ListComprehensive
// 	var q *datastore.Query

// 	// Check if they want to limit the query using a modified since date.
// 	if ifModifiedSince := r.HeaderParameter("If-Modified-Since"); ifModifiedSince == "" {
// 		q = datastore.NewQuery(res.Kind())
// 	} else {
// 		if t, err := time.Parse(time.RFC3339Nano, ifModifiedSince); err != nil {
// 			w.AddHeader("Content-Type", "text/plain")
// 			w.WriteErrorString(http.StatusNotAcceptable, err.Error())
// 			return
// 		} else {
// 			q = datastore.NewQuery(res.Kind()).
// 				Filter("LastModified >", t)
// 		}
// 	}

// 	if keys, err := q.GetAll(c, &result.Entry); err != nil {
// 		w.AddHeader("Content-Type", "text/plain")
// 		w.WriteErrorString(http.StatusInternalServerError, err.Error())
// 		return
// 	} else {
// 		for i, k := range keys {
// 			result.Entry[i].Key = k.Encode()
// 			result.Entry[i].Link.Rel = "self"
// 			result.Entry[i].Link.Href = result.Entry[i].PreferredLink(k)
// 		}
// 	}

// 	// Cache Control: By allowing a short cache time here we can reduce database calls and cost.
// 	w.AddHeader("Cache-Control", "max-age=900,must-revalidate")
// 	w.AddHeader(restful.HEADER_LastModified, time.Now().Format(time.RFC3339Nano))

// 	w.WriteEntity(result)
// }
