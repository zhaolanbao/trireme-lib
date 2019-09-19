package podmonitor

import (
	corev1 "github.com/kubernetes/core/v1"
	"github.com/kubernetes/apimachinery/pkg/types"
	"github.com/kubernetes-sigs/controller-runtime/pkg/client"
	"github.com/kubernetes-sigs/controller-runtime/pkg/handler"
	"github.com/kubernetes-sigs/controller-runtime/pkg/reconcile"
)

// WatchPodMapper determines if we want to reconcile on a pod event. There are two limitiations:
// - the pod must be schedule on a matching nodeName
// - if the pod requests host networking, only reconcile if we want to enable host pods
type WatchPodMapper struct {
	client         client.Client
	nodeName       string
	enableHostPods bool
}

// Map implements the handler.Mapper interface to emit reconciles for corev1.Pods. It effectively
// filters the pods by looking for a matching nodeName and filters them out if host networking is requested,
// but we don't want to enable those.
func (w *WatchPodMapper) Map(obj handler.MapObject) []reconcile.Request {
	pod, ok := obj.Object.(*corev1.Pod)
	if !ok {
		return nil
	}

	if pod.Spec.NodeName != w.nodeName {
		return nil
	}

	if pod.Spec.HostNetwork && !w.enableHostPods {
		return nil
	}

	return []reconcile.Request{
		{
			NamespacedName: types.NamespacedName{
				Name:      pod.Name,
				Namespace: pod.Namespace,
			},
		},
	}
}
