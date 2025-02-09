package main

import "time"

const (
	baseURL        = "localhost:8081"
	usersPostfix   = "/users"
	userPostfix    = usersPostfix + "/{id}"
	defaultTimeout = time.Second * 5
)

const (
	// AdminRole is ...
	AdminRole Role = iota
	// UserRole is ...
	UserRole
)
