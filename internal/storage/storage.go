package storage

// Storage - encapsulate work with storage e.g file or storage
type Query struct {
	// main data
	data []byte
	// field for filter data
	filter map[string]any
}
type StorageProvider interface {
	// Get - get data from storage
	Get(Query) ([]byte, error)

	// Add - add data to storage
	Add(Query) (string, error)

	// Update - update data in storage
	Update(Query) error

	// Delete - delete data from storage
	Delete(Query) error
}