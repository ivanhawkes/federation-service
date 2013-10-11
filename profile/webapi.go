package profile

import (
	"appengine"
	"appengine/datastore"
	"appengine/user"
	"fmt"
	"time"
)

func handlePut(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	p := Employee{
		firstName:     "Ivan",
		nickName:     "Socks",
		lastName:     "Hawkes",
		Date: time.Now(),
		Account:  user.Current(c).String(),
	}

	key, err := datastore.Put(c, datastore.NewIncompleteKey(c, "employee", nil), &p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	}

	fmt.Fprintf(w, "Stored and retrieved the Employee named %q", p.firstName)
}
