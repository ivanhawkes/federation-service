package components

import (
	"google.golang.org/appengine/datastore"
)

type ComponentMap map[string]interface{}

type Component struct {
	// Prototype of this component.
	Prototype string `datastore:",noindex" json:"prototype" xml:"prototype"`

	Parts ComponentMap `datastore:"-,noindex" json:"parts" xml:"parts"`
}

func (com Component) Load(c <-chan datastore.Property) error {
	// Get the prototype first.
	// if err := datastore.LoadStruct(com, c); err != nil {
	// 	return err
	// }

	// for p := range c {
	// 	if p.Multiple {
	// 		value := reflect.ValueOf(com[p.Name])
	// 		if value.Kind() != reflect.Slice {
	// 			com[p.Name] = p.Value.(string)
	// 		} /*else {
	// 			com[p.Name] = append(com[p.Name].([]string), p.Value)
	// 		}*/
	// 	} else {
	// 		com[p.Name] = p.Value.(string)
	// 	}
	// }
	return nil
}

func (com Component) Save(c chan<- datastore.Property) error {
	//defer close(c)

	//datastore.SaveStruct(com, c)

	c <- datastore.Property{
		Name:  "TEST",
		Value: "A Test",
	}

	for k, v := range com.Parts {
		c <- datastore.Property{
			Name:  "Parts." + com.Prototype + "." + k,
			Value: v,
		}
	}

	return nil
}
