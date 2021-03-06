package controller

import (
	"github.com/sambacha/besu-kubernetes/besu-operator/pkg/controller/prometheus"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, prometheus.Add)
}
