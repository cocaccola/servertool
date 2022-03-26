package reconciler

import "slack/servertool/internal/resources"

type Reconciler interface {
	Reconcile() error
}

type ResourceReconciler struct {
	OrderedResources []*resources.Resource
	ResourceMap      *resources.ResourceMap
}
