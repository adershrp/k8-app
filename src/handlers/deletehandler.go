package handlers

import (
	"zts-upgrade-handler/config"

	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
)

// LogHandler
type DeleteHandler struct {
	kubeclient    *kubernetes.Clientset
	label         string
	config        *config.Config
	createHandler HandlerFunc
	updateHandler HandlerFunc
	deleteHandler HandlerFunc
}

// Init
func (s *DeleteHandler) Init(conf *config.Config, kubeclient *kubernetes.Clientset) error {
	s.kubeclient = kubeclient
	s.config = conf
	s.prepareCreateHandler()
	s.prepareUpdateHandler()
	s.prepareDeleteHandler()
	return nil
}

// prepareCreateHandler
func (s *DeleteHandler) prepareCreateHandler() {
	s.createHandler = HandlerFunc{
		handle: func(obj interface{}) {
			switch v := obj.(type) {
			case *batchv1.Job:
				s.handleJobDelete(v)
			}
		},
	}
}

// prepareUpdateHandler
func (s *DeleteHandler) prepareUpdateHandler() {
	s.updateHandler = HandlerFunc{
		handle: func(obj interface{}) {
			switch v := obj.(type) {
			case *batchv1.Job:
				s.handleJobDelete(v)
			}
		},
	}
}

// prepareDeleteHandler
func (s *DeleteHandler) prepareDeleteHandler() {
	s.deleteHandler = HandlerFunc{
		handle: func(obj interface{}) {
			switch v := obj.(type) {
			case *batchv1.Job:
				s.handleJobDelete(v)
			}
		},
	}
}
func (s *DeleteHandler) handleJobDelete(job *batchv1.Job) {
	klog.Infof("deleting job %s namespace %s", job.Name, job.Namespace)
	return
}

// ObjectCreated
func (s *DeleteHandler) ObjectCreated(obj interface{}) {
	// if s.shouldProcessEvent(obj) {
	s.handleEvent(obj, s.createHandler)
	// }
}

// ObjectDeleted
func (s *DeleteHandler) ObjectDeleted(obj interface{}) {
	// if s.shouldProcessEvent(obj) {
	s.handleEvent(obj, s.deleteHandler)
	// }
}

// ObjectUpdated
func (s *DeleteHandler) ObjectUpdated(oldObj, newObj interface{}) {
	// if s.shouldProcessEvent(newObj) {
	s.handleEvent(newObj, s.updateHandler)
	// }
}

// handleEvent
func (s *DeleteHandler) handleEvent(obj interface{}, handler HandlerFunc) {
	handler.handle(obj)
}
