package store

// Template .
type Template struct {
	domain   string
	Section  string
	Context  string
	KeyValue string
}

// NewTemplate .
func NewTemplate(section string, keyValue string, context string) *Template {
	return &Template{
		domain:   "system",
		Section:  section,
		Context:  context,
		KeyValue: keyValue,
	}
}

// key .
func (t *Template) key() string {
	if t.Context != "" {
		return Key(t.domain, t.KeyValue, t.Context)
	}
	return Key(t.domain, t.KeyValue)
}

// Delete .
func (t *Template) Delete() error {
	return Delete(t.key())
}
