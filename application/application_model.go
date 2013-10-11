package application

import (
	"time"
)

// This will need to tie to a user somehow, or a complete company struct or something. We
// can slack off to start and just make sure the key is in there.
type profile struct {
	LastModified          time.Time `json:"last_modified"`
	CompanyName           string    `json:"company_name"`
	ApplicationName       string    `json:"application_name"`
	ApplicationPrivateKey string    `json:"application_private_key"`
	ApplicationPublicKey  string    `json:"application_public_key"`
}
