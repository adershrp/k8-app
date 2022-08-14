package config

import (
	"time"
)

type Config struct {
	NamespaceToWatch string
	WatchPods        bool
	WatchJobs        bool
	WatchServices    bool
	WatchSecrets     bool
	ResyncPeriod     time.Duration
	LabelsToWatch    map[string]string
	NamesToWatch     []string
}
