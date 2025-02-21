package entity

import "time"

type FileMeta struct {
	Path    string    `json:"path"`
	ModTime time.Time `json:"mod_time"`
}
