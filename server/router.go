package server

type route struct {
	method  string
	handler interface{}
}

type router struct {
	routes []*route
}

func NewRouter() *router {
	return &router{}
}

func (r *router) List(pattern string, handler Handler) {
	r.AddRoute("list", pattern, handler)
}

func (r *router) Retrieve(pattern string, handler Handler) {
	r.AddRoute("retrieve", pattern, handler)
}

func (r *router) Update(pattern string, handler Handler) {
	r.AddRoute("update", pattern, handler)
}

func (r *router) Create(pattern string, handler Handler) {
	r.AddRoute("create", pattern, handler)
}

// Add a new route to the Router
func (r *router) AddRoute(method string, pattern string, handler interface{}) {

}

// find a matching route and invoke its handler
func (r *router) handleEvent(message *Message) {

}
