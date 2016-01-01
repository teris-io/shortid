// Copyright 2015 Ventu.io. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file

package shortid_test

import (
	"github.com/op/go-logging"
	"os"
)

var log = logging.MustGetLogger("shortid_test")

func init() {
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	leveled := logging.AddModuleLevel(backend)
	leveled.SetLevel(logging.DEBUG, "")
	logging.SetBackend(leveled)

	format := logging.MustStringFormatter("%{color}%{time:15:04:05.000} %{longfunc}: %{level:.6s}%{color:reset} %{message}")
	logging.SetFormatter(format)
}
