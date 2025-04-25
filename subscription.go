package quelldb

// Subscribe adds a new subscriber to the database.
// The subscriber is a function that takes a ChangeEvent as an argument.
// The subscriber will be notified of changes to the database, such as PUT or DELETE operations.
func (db *DB) publish(event ChangeEvent) {
	db.subLock.RLock()
	defer db.subLock.RUnlock()
	for _, handler := range db.subscribers {
		go handler(event)
	}
}
