package shortid

import (
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("longpoll")

const (
	no int32 = iota
	yes
)

const Version = 1.0
