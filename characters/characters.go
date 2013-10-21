package characters

import (
	"appengine"
	"appengine/datastore"
	"appengine/user"
	"github.com/emicklei/go-restful"
	"log"
	"net/http"
	"time"
)

const (
	rootPath = "/characters"
)

// The various states for a character resource.
const (
	StatusActive = iota
	StatusDeactivated
	StatusDeletionPending
	StatusDeleted
	StatusBanned
	StatusPermanentBan
)

type CharacterShallow struct {
	Id        string `datastore:"-" json:"id" xml:"id"`
	FirstName string `json:"first_name" xml:"first-name"`
	NickName  string `json:"nick_name" xml:"nick-name"`
	LastName  string `json:"last_name" xml:"last-name"`
	Link      string `datastore:"-" json:"link" xml:"link"`
	FactionId int64  `datastore:"FactionId" json:"faction_id" xml:"faction-id"`
}

type Character struct {
	CharacterShallow
	UserId       string    `datastore:"UserId" json:"-" xml:"-"`
	LastModified time.Time `json:"-" xml:"-"`
	Status       int       `json:"status" xml:"status"`
}

type CharacterApi struct {
	Path string
}

func init() {
	log.Printf("Characters: Register")
}

// Register the routes we require for this resource type.
//
func (api CharacterApi) Register() {
	ws := new(restful.WebService)

	ws.
		Path(rootPath).
		Consumes(restful.MIME_JSON, restful.MIME_XML).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	ws.Route(ws.POST("").To(api.create).
		// Swagger documentation.
		Doc("create a new character").
		Param(ws.BodyParameter("Character", "representation of a character").DataType("characters.Character")).
		Reads(Character{}))

	ws.Route(ws.GET("/{character-id}").To(api.read).
		// Swagger documentation.
		Doc("read a character").
		Param(ws.PathParameter("character-id", "identifier for a character").DataType("string")).
		Writes(Character{}))

	ws.Route(ws.PUT("/{character-id}").To(api.update).
		// Swagger documentation.
		Doc("update an existing character").
		Param(ws.PathParameter("character-id", "identifier for a character").DataType("string")).
		Param(ws.BodyParameter("Character", "representation of a character").DataType("characters.Character")).
		Reads(Character{}))

	ws.Route(ws.DELETE("/{character-id}").To(api.delete).
		// Swagger documentation.
		Doc("delete a character").
		Param(ws.PathParameter("character-id", "identifier for a character").DataType("string")))

	restful.Add(ws)
}

// Create a new resource.
//
func (api *CharacterApi) create(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	// Marshall the entity from the request into a struct.
	character := new(Character)
	err := r.ReadEntity(&character)
	if err != nil {
		w.WriteError(http.StatusNotAcceptable, err)
		return
	}

	// Set some fields that need special handling.
	character.LastModified = time.Now()
	character.Status = StatusActive

	// The resource belongs to this character.
	character.UserId = user.Current(c).ID

	// Store the character.
	k, err := datastore.Put(c, datastore.NewIncompleteKey(c, "characters", nil /*ancestor*/), character)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// The resource Id.
	character.Id = k.Encode()

	// Let them know the location of the newly created resource.
	// TODO: Use a safe Url path append function.
	w.AddHeader("Location", rootPath+"/"+k.Encode())

	// Provide a link for ease of API usage.
	// TODO: This should be a fully qualified path.
	character.Link = rootPath + "/" + k.Encode()

	// Return the resultant entity.
	w.WriteHeader(http.StatusCreated)
	w.WriteEntity(character)
}

// Read the resource.
//
func (api CharacterApi) read(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	// Decode the request parameter to determine the key for the entity.
	k, err := datastore.DecodeKey(r.PathParameter("character-id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Retrieve the entity from the datastore.
	character := Character{}
	if err := datastore.Get(c, k, &character); err != nil {
		if err.Error() == "datastore: no such entity" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check we own the resource before allowing them to view it.
	// Optionally, return a 404 instead to help prevent guessing ids.
	// TODO: Allow admins access.
	//if character.UserId != user.Current(c).ID {
	//	http.Error(w, "You do not have access to this resource", http.StatusForbidden)
	//	return
	//}

	// Set their Id.
	character.Id = k.Encode()

	// Provide a link for ease of API usage.
	// TODO: This should be a fully qualified path.
	character.Link = rootPath + "/" + k.Encode()

	w.WriteEntity(character)
}

// Update the resource.
//
func (api *CharacterApi) update(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	// Decode the request parameter to determine the key for the entity.
	k, err := datastore.DecodeKey(r.PathParameter("character-id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Marshall the entity from the request into a struct.
	character := new(Character)
	err = r.ReadEntity(&character)
	if err != nil {
		w.WriteError(http.StatusNotAcceptable, err)
		return
	}

	// Retrieve the old entity from the datastore.
	old := Character{}
	if err := datastore.Get(c, k, &old); err != nil {
		if err.Error() == "datastore: no such entity" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check we own the resource before allowing them to update it.
	// Optionally, return a 404 instead to help prevent guessing ids.
	// TODO: Allow admins access.
	if old.UserId != user.Current(c).ID {
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	// Since the whole entity is re-written, we need to assign any invariant fields again
	// e.g. the owner of the entity.
	character.UserId = user.Current(c).ID

	// Keep track of the last modification date.
	character.LastModified = time.Now()

	// Attempt to overwrite the old entity.
	_, err = datastore.Put(c, k, character)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Let them know it succeeded.
	w.WriteHeader(http.StatusNoContent)
}

// Delete the resource.
//
func (api *CharacterApi) delete(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	// Decode the request parameter to determine the key for the entity.
	k, err := datastore.DecodeKey(r.PathParameter("character-id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Retrieve the old entity from the datastore.
	old := Character{}
	if err := datastore.Get(c, k, &old); err != nil {
		if err.Error() == "datastore: no such entity" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check we own the resource before allowing them to delete it.
	// Optionally, return a 404 instead to help prevent guessing ids.
	// TODO: Allow admins access.
	if old.UserId != user.Current(c).ID {
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	// Delete the entity.
	if err := datastore.Delete(c, k); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Success notification.
	w.WriteHeader(http.StatusNoContent)
}
