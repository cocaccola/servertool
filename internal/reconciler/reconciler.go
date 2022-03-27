package reconciler

import "slack/servertool/internal/resources"

type ResourceReconciler struct {
	OrderedResources resources.Resources
	ResourceMap      resources.ResourceMap
}

func NewResourceReconciler(orderedResources resources.Resources, resourceMap resources.ResourceMap) *ResourceReconciler {
	return &ResourceReconciler{
		OrderedResources: orderedResources,
		ResourceMap:      resourceMap,
	}
}

func (rr *ResourceReconciler) ReconcileAll() error {
	for _, r := range rr.OrderedResources {
		err := r.Reconcile(rr.ResourceMap)
		if err != nil {
			return err
		}
	}
	return nil
}
