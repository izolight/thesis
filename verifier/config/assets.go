// +build dev

package config

import "net/http"

var Assets http.FileSystem = http.Dir("../config/files")