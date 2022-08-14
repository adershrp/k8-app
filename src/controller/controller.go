// Copyright © 2018 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

package controller

import (
	"zts-upgrade-handler/config"
	"zts-upgrade-handler/handlers"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
)

func StartLogController(informerFactory informers.SharedInformerFactory, eventHandler handlers.Handler, config *config.Config) error {
	if config.WatchPods {
		watchPods(informerFactory, eventHandler, config)
	}
	if config.WatchJobs {
		watchJobs(informerFactory, eventHandler, config)
	}
	return nil
}

// watchPods
func watchPods(informerFactory informers.SharedInformerFactory, eventHandler handlers.Handler, config *config.Config) cache.Store {
	informer := informerFactory.Core().V1().Pods().Informer()
	informer.AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    eventHandler.ObjectCreated,
			UpdateFunc: eventHandler.ObjectUpdated,
			DeleteFunc: eventHandler.ObjectDeleted,
		},
	)
	go informer.Run(wait.NeverStop)
	klog.Infof("Waiting for pods to be synced")
	cache.WaitForCacheSync(wait.NeverStop, informer.HasSynced)
	klog.Infof("synced pods")
	return nil
}

// watchJobs
func watchJobs(informerFactory informers.SharedInformerFactory, eventHandler handlers.Handler, config *config.Config) cache.Store {
	informer := informerFactory.Batch().V1().Jobs().Informer()
	informer.AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    eventHandler.ObjectCreated,
			UpdateFunc: eventHandler.ObjectUpdated,
			DeleteFunc: eventHandler.ObjectDeleted,
		},
	)
	go informer.Run(wait.NeverStop)
	klog.Infof("Waiting for jobs to be synced")
	cache.WaitForCacheSync(wait.NeverStop, informer.HasSynced)
	klog.Infof("synced jobs")
	return nil
}
