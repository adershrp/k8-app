package handlers

import (
	"strings"

	"zts-upgrade-handler/config"
	"zts-upgrade-handler/utils"

	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
)

// LogHandler
type LogHandler struct {
	kubeclient    *kubernetes.Clientset
	label         string
	config        *config.Config
	createHandler HandlerFunc
	updateHandler HandlerFunc
	deleteHandler HandlerFunc
}

// Init
func (s *LogHandler) Init(conf *config.Config, kubeclient *kubernetes.Clientset) error {
	s.kubeclient = kubeclient
	s.config = conf
	s.prepareCreateHandler()
	s.prepareUpdateHandler()
	s.prepareDeleteHandler()
	return nil
}

// prepareCreateHandler
func (s *LogHandler) prepareCreateHandler() {
	s.createHandler = HandlerFunc{
		handle: func(obj interface{}) {
			switch v := obj.(type) {
			case *v1.Pod:
				s.handlePodCreate(v)
			case *batchv1.Job:
				s.handleJobCreate(v)
			}
		},
	}
}

// prepareUpdateHandler
func (s *LogHandler) prepareUpdateHandler() {
	s.updateHandler = HandlerFunc{
		handle: func(obj interface{}) {
			switch v := obj.(type) {
			case *v1.Pod:
				s.handlePodUpdate(v)
			case *batchv1.Job:
				s.handleJobUpdate(v)
			}
		},
	}
}

// prepareDeleteHandler
func (s *LogHandler) prepareDeleteHandler() {
	s.deleteHandler = HandlerFunc{
		handle: func(obj interface{}) {
			switch v := obj.(type) {
			case *v1.Pod:
				s.handlePodDelete(v)
			case *batchv1.Job:
				s.handleJobDelete(v)
			}
		},
	}
}

// shouldProcessEvent
func (s *LogHandler) shouldProcessEvent(obj interface{}) bool {
	switch v := obj.(type) {
	case *v1.Pod:
		return utils.FilterByNames(s.config.NamesToWatch, v.Name) && utils.FilterByLabels(s.config.LabelsToWatch, v.Labels)
	case *batchv1.Job:
		return utils.FilterByNames(s.config.NamesToWatch, v.Name) && utils.FilterByLabels(s.config.LabelsToWatch, v.Labels)
	case *v1.Service:
		if strings.HasSuffix(v.Name, "-syndicate") {
			return false
		}
		return utils.FilterByLabels(s.config.LabelsToWatch, v.Labels)
	}
	return false
}

// handleEvent
func (s *LogHandler) handleEvent(obj interface{}, handler HandlerFunc) {
	handler.handle(obj)
}

// ObjectCreated
func (s *LogHandler) ObjectCreated(obj interface{}) {
	if s.shouldProcessEvent(obj) {
		s.handleEvent(obj, s.createHandler)
	}
}

// ObjectDeleted
func (s *LogHandler) ObjectDeleted(obj interface{}) {
	if s.shouldProcessEvent(obj) {
		s.handleEvent(obj, s.deleteHandler)
	}
}

// ObjectUpdated
func (s *LogHandler) ObjectUpdated(oldObj, newObj interface{}) {
	if s.shouldProcessEvent(newObj) {
		s.handleEvent(newObj, s.updateHandler)
	}
}

// handlePodUpdate
func (s *LogHandler) handlePodUpdate(pod *v1.Pod) {
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
		klog.Infof("updating pod %s namespace %s", pod.Name, pod.Namespace)
		klog.Infof(
			"POD UPDATED. %s/%s %s %s",
			pod.Namespace, pod.Name, pod.Status.Phase, podReadyTime.Sub(podScheduledTime.Time),
		)
	}
	return
}

// handlePodDelete
func (s *LogHandler) handlePodDelete(pod *v1.Pod) {
	// klog.Infof("deleting pod %s namespace %s", pod.Name, pod.Namespace)
	return
}

// handlePodCreate
func (s *LogHandler) handlePodCreate(pod *v1.Pod) {
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
		klog.Infof("creating pod %s namespace %s", pod.Name, pod.Namespace)
		klog.Infof(
			"POD UPDATED. %s/%s %s %s",
			pod.Namespace, pod.Name, pod.Status.Phase, podReadyTime.Sub(podScheduledTime.Time),
		)
	}
	return
}

// handleJobUpdate
func (s *LogHandler) handleJobUpdate(job *batchv1.Job) {
	var jobScheduledTime, jobCompleteTime metav1.Time
	jobScheduledTime = *job.Status.StartTime
	for _, conditions := range job.Status.Conditions {
		if conditions.Type == batchv1.JobComplete {
			jobCompleteTime = conditions.LastTransitionTime
		}
	}
	if !jobScheduledTime.IsZero() && !jobCompleteTime.IsZero() && job.Status.Succeeded == 1 {
		klog.Infof("updating pod %s namespace %s", job.Name, job.Namespace)
		klog.Infof(
			"POD UPDATED. %s/%s %s %s",
			job.Namespace, job.Name, job.Status.Succeeded, jobCompleteTime.Sub(jobScheduledTime.Time),
		)
	}
	return
}

// handleJobDelete
func (s *LogHandler) handleJobDelete(job *batchv1.Job) {
	klog.Infof("deleting job %s namespace %s", job.Name, job.Namespace)
	return
}

// handleJobCreate
func (s *LogHandler) handleJobCreate(job *batchv1.Job) {
	klog.Infof("creating job %s namespace %s", job.Name, job.Namespace)
	return
}
