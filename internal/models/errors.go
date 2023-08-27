package models

import "errors"

var (
	ErrShorteningNotFound = errors.New("shortening not found")
	ErrShorteningExists   = errors.New("shortening already exist")
)
