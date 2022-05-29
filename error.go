package main

import "errors"

var (
	ErrEmptyName    = errors.New("name cannot be empty")
	ErrEmptyAddress = errors.New("address cannot be empty")

	ErrInvalidAddress = errors.New("address must be url style")

	ErrDuplicateName = errors.New("name cannot be duplicated")

	ErrGameStarted = errors.New("game is already started")
)
