package store

// Route .
type Route struct {
	Domain string `json:"domain"`
}

// NewRoute .
func NewRoute(domain string) *Route {
	return &Route{Domain: domain}
}

// key .
func (r *Route) key() string {
	return Key(r.Domain, "route")
}




// Delete .
func (r *Route) Delete() error {
	return Delete(r.key())
}
