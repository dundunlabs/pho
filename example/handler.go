package main

import "github.com/dundunlabs/tra"

func homeHandler(ctx *tra.Context) any {
	return "Hi! Welcome to trà"
}

func userListHandler(ctx *tra.Context) any {
	return []string{"user 1", "user 2"}
}
