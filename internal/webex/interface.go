package webex

// Reader performs read operations
type Reader interface {
	Get(endpoint string, params map[string]string) (map[string]interface{}, error)
}

// Writer performs write operations
type Writer interface {
	Post(endpoint string, data interface{}) (map[string]interface{}, error)
	Put(endpoint string, data interface{}) (map[string]interface{}, error)
}

// Deleter performs delete operations
type Deleter interface {
	Delete(endpoint string) error
}

// HTTPClient combines all HTTP operations
// This maintains backward compatibility while providing segregated interfaces
type HTTPClient interface {
	Reader
	Writer
	Deleter
}
