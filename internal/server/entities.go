package server

import "github.com/dtbead/wc-maps-archive/internal/entities"

type Video struct {
	Id       int64
	Duration int32
	Hashes   entities.Hashes
}
