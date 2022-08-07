package contoller

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	coreinformers "k8s.io/client-go/informers/core/v1"
	klog "k8s.io/klog/v2"

	//"k8s.io/client-go/pkg/api/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
)

// PodLoggingController logs the name and namespace of pods that are added,
// deleted, or updated
type PodLoggingController struct {
	informerFactory informers.SharedInformerFactory
	podInformer     coreinformers.PodInformer
}

// Run starts shared informers and waits for the shared informer cache to
// synchronize.
func (c *PodLoggingController) Run(stopCh chan struct{}) error {
	// Starts all the shared informers that have been created by the factory so
	// far.
	c.informerFactory.Start(stopCh)
	// wait for the initial synchronization of the local cache.
	if !cache.WaitForCacheSync(stopCh, c.podInformer.Informer().HasSynced) {
		return fmt.Errorf("Failed to sync")
	}
	return nil
}

func (c *PodLoggingController) podAdd(obj interface{}) {
	pod := obj.(*v1.Pod)
	var podScheduledTime, podReadyTime metav1.Time
	for _, conditions := range pod.Status.Conditions {
		if conditions.Type == v1.PodScheduled {
			podScheduledTime = conditions.LastTransitionTime
		}
		if conditions.Type == v1.PodReady {
			podReadyTime = conditions.LastTransitionTime
		}
	}
	if !podScheduledTime.IsZero() && !podReadyTime.IsZero() && pod.Status.Phase == v1.PodRunning {
		klog.Infof(
			"POD UPDATED. %s/%s %s %s",
			pod.Namespace, pod.Name, pod.Status.Phase, podReadyTime.Sub(podScheduledTime.Time),
		)
	}
}

func (c *PodLoggingController) podUpdate(old, new interface{}) {
	oldPod := old.(*v1.Pod)
	newPod := new.(*v1.Pod)
	var podScheduledTime, podReadyTime metav1.Time
	for _, conditions := range newPod.Status.Conditions {
		if conditions.Type == v1.PodScheduled {
			podScheduledTime = conditions.LastTransitionTime
		}
		if conditions.Type == v1.PodReady {
			podReadyTime = conditions.LastTransitionTime
		}
	}
	if !podScheduledTime.IsZero() && !podReadyTime.IsZero() && newPod.Status.Phase == v1.PodRunning {
		klog.Infof(
			"POD UPDATED. %s/%s %s %s",
			oldPod.Namespace, oldPod.Name, newPod.Status.Phase, podReadyTime.Sub(podScheduledTime.Time),
		)
	}
}

func (c *PodLoggingController) podDelete(obj interface{}) {
	pod := obj.(*v1.Pod)
	klog.Infof("POD DELETED: %s/%s", pod.Namespace, pod.Name)
}

// NewPodLoggingController creates a PodLoggingController
func NewPodLoggingController(informerFactory informers.SharedInformerFactory) *PodLoggingController {
	podInformer := informerFactory.Core().V1().Pods()
	c := &PodLoggingController{
		informerFactory: informerFactory,
		podInformer:     podInformer,
	}
	podInformer.Informer().AddEventHandler(
		// Your custom resource event handlers.
		cache.ResourceEventHandlerFuncs{
			// Called on creation
			AddFunc: c.podAdd,
			// Called on resource update and every resyncPeriod on existing resources.
			UpdateFunc: c.podUpdate,
			// Called on resource deletion.
			DeleteFunc: c.podDelete,
		},
	)
	return c
}
