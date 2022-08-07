package main

import (
	"flag"
	"path/filepath"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	klog "k8s.io/klog/v2"
	"k8s.io/utils/env"

	//"k8s.io/client-go/pkg/api/v1"
	controller "zts-upgrade-handler/controller"

	v1 "k8s.io/api/core/v1"
)

var namespace string

func init() {
	namespace = env.GetString("CURRENT_NAMESPACE", v1.NamespaceDefault)
	klog.Info("CURRENT_NAMESPACE ", namespace)
}

func main() {

	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	var config *rest.Config
	var err error
	// use the current context in kubeconfig
	config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if config == nil {
		klog.Error(err)
		// creates the in-cluster config
		config, err = rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Fatal(err)
	}
	// namespace
	informerOption := informers.WithNamespace(namespace)
	// labels := informers.WithTweakListOptions(func(lo *metav1.ListOptions) {
	// 	lo.LabelSelector = "app=nats-box"
	// })
	// create shared factory
	factory := informers.NewSharedInformerFactoryWithOptions(clientset, 0, informerOption)
	_ = controller.NewPodLoggingController(factory)
	// stop := make(chan struct{})
	// defer close(stop)
	// if err = podController.Run(stop); err != nil {
	// 	klog.Fatal(err)
	// }
	_ = controller.NewJobDeleteController(factory)
	// if err := jobController.Run(stop); err != nil {
	// 	klog.Fatal(err)
	// }
	// select {}
	factory.Start(wait.NeverStop)
	factory.WaitForCacheSync(wait.NeverStop)
	select {}
}
