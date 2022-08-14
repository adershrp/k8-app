package utils

import "strings"

// FilterByLabels
func FilterByLabels(filter, k8Labels map[string]string) bool {
	if len(filter) <= 0 {
		return true
	}

	var labelFound bool = true
	for sKey, sVal := range filter {
		if w, ok := k8Labels[sKey]; ok && sVal != w {
			labelFound = false
		}
	}
	return labelFound
}

// FilterByNames
func FilterByNames(filter []string, k8Name string) bool {
	for _, v := range filter {
		if strings.HasPrefix(k8Name, v) {
			return true
		}
	}
	return len(filter) <= 0
}
