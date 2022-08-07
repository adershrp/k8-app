package contoller

import (
	v1 "k8s.io/api/batch/v1"
	"k8s.io/client-go/informers"
	jobinformers "k8s.io/client-go/informers/batch/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
)

// JobController - can be used to monitor and delete certain jobs
type JobDeleteController struct {
	informerFactory informers.SharedInformerFactory
	jobInformer     jobinformers.JobInformer
}

// addFunc
func (j *JobDeleteController) addFunc(obj interface{}) {
	job := obj.(*v1.Job)

	klog.Info("Created %s job", job.GetName())
}

// updateFunc
func (j *JobDeleteController) updateFunc(oldObj, newObj interface{}) {
	job := newObj.(*v1.Job)
	klog.Info("Updated %s job", job.GetName())
}

// deleteFunc
func (j *JobDeleteController) deleteFunc(obj interface{}) {
	job := obj.(*v1.Job)
	klog.Info("Deleted %s job", job.GetName())
}

// NewJobDeleteController
func NewJobDeleteController(informerFactory informers.SharedInformerFactory) *JobDeleteController {
	jobinformer := informerFactory.Batch().V1().Jobs()
	j := &JobDeleteController{
		informerFactory: informerFactory,
		jobInformer:     jobinformer,
	}
	jobinformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: j.addFunc,
		// UpdateFunc: j.updateFunc,
		DeleteFunc: j.deleteFunc,
	})
	return j
}
