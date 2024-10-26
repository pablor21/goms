package database

import "github.com/pablor21/goms/pkg/errors"

var ErrUnsupportedDriver = errors.NewAppError("Unsupported database driver", 500)
