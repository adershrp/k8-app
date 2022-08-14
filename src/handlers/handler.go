package handlers

// Handler
type Handler interface {
	ObjectCreated(obj interface{})
	ObjectDeleted(obj interface{})
	ObjectUpdated(oldObj, newObj interface{})
}

// HandlerFunc
type HandlerFunc struct {
	handle func(obj interface{})
}
