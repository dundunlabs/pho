package tra

type Middleware func(next Handler) Handler
