package rest

import (
	"errors"
)

var (
	ErrNotExist = errors.New("Entity does not exist")
)

type RESTType struct {
	Plural   string
	Singular string
}

type RESTConfig struct {
	BaseUrl string
	Type    RESTType
}
