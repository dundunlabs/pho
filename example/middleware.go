package main

import "github.com/dundunlabs/tra"

type APIResult struct {
	Data  any   `json:"data"`
	Error error `json:"error"`
}

func apiMiddleware(next tra.Handler) tra.Handler {
	return func(ctx *tra.Context) any {
		switch v := next(ctx).(type) {
		case error:
			return APIResult{Error: v}
		default:
			return APIResult{Data: v}
		}
	}
}
