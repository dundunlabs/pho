package tra

type route struct {
	path    string
	method  string
	handler Handler
}
