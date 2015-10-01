package error

// Struct tags are used to map struct fields to fields in the database
type ErrorException struct {
	Id      string `gorethink:"id,omitempty"`
	Name    string `gorethink:"name"`
	Message string `gorethink:"message"`
}
