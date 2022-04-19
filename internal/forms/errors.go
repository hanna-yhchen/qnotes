package forms

type errors map[string][]string

// Add appends an error message for a given form field.
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// Get returns the first error message related to a given field.
func (e errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}
