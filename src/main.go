package main

import (
	"flag"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/klog/v2"
	"k8s.io/utils/env"

	"zts-upgrade-handler/config"
	"zts-upgrade-handler/controller"
	"zts-upgrade-handler/handlers"
)

func main() {
	c, err := loadConfig()
	if err != nil {
		return
	}
	kubeClient, err := getKubeClient()
	if err != nil {
		klog.Fatalf("Error: %s", err)
	}
	nsOption := informers.WithNamespace(c.NamespaceToWatch)
	factory := informers.NewSharedInformerFactoryWithOptions(kubeClient, c.ResyncPeriod, nsOption)
	// creating log handler --
	logHandler := &handlers.LogHandler{}
	if handlerErr := logHandler.Init(c, kubeClient); handlerErr != nil {
		klog.Errorf("failed to initialize handler %v", handlerErr)
		return
	}
	// delHandler := &handlers.DeleteHandler{}

	go func() {
		err := controller.StartLogController(factory, logHandler, c)
		if err != nil {
			klog.Errorf("failed to start log controller %v", err)
			return
		}
	}()

	// go controller.StartLogController(factory, delHandler, c)
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	<-sigterm
}

// getKubeClient
func getKubeClient() (*kubernetes.Clientset, error) {
	klog.Infof("building client")
	var kubeConfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeConfig = flag.String("kubeConfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeConfig file")
	} else {
		kubeConfig = flag.String("kubeConfig", "", "absolute path to the kubeConfig file")
	}
	flag.Parse()

	var config *rest.Config
	var err error
	// use the current context in kubeConfig
	config, err = clientcmd.BuildConfigFromFlags("", *kubeConfig)
	if config == nil {
		klog.Error(err)
		// creates the in-cluster config
		config, err = rest.InClusterConfig()
		if err != nil {
			klog.Errorf("Error creating client with inclusterConfig, %s", err)
			panic(err.Error())
		}
		klog.Info("kubeclient created using InClusterConfig")
	}
	if client, err := kubernetes.NewForConfig(config); err != nil {
		klog.Errorf("Error creating client with inclusterConfig, %s", err)
		return nil, err
	} else {
		return client, nil
	}
}

// loadConfig
func loadConfig() (*config.Config, error) {
	conf := &config.Config{
		NamespaceToWatch: env.GetString("CURRENT_NAMESPACE", v1.NamespaceDefault),
		ResyncPeriod:     0 * time.Minute,
		WatchPods:        true,
		WatchJobs:        true,
		WatchServices:    true,
		WatchSecrets:     true,
		LabelsToWatch:    map[string]string{"app": "nginx"},
		NamesToWatch:     []string{"nginx"},
	}
	return conf, nil
}
