package serializer

import (
	"encoding/gob"

	"github.com/pablor21/goms/app/dtos"
	"github.com/pablor21/goms/pkg/auth"
)

func init() {
	gob.Register(auth.DflPrincipal{})
	gob.Register(dtos.UserDTO{})
}
