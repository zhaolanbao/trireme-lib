package podmonitor

import (
	"context"
	"errors"
	"fmt"

	corev1 "github.com/kubernetes/core/v1"

	"github.com/kubernetes-sigs/controller-runtime/pkg/client"
	"github.com/kubernetes-sigs/controller-runtime/pkg/event"
)

// ResyncWithAllPods is called from the implemented resync, it will list all pods
// and fire them down the event source (the generic event channel)
func ResyncWithAllPods(ctx context.Context, c client.Client, evCh chan<- event.GenericEvent) error {
	if c == nil {
		return errors.New("pod: no client available")
	}

	if evCh == nil {
		return errors.New("pod: no event source available")
	}

	list := &corev1.PodList{}
	if err := c.List(ctx, &client.ListOptions{}, list); err != nil {
		return fmt.Errorf("pod: %s", err.Error())
	}

	for _, pod := range list.Items {
		p := pod.DeepCopy()
		evCh <- event.GenericEvent{
			Meta:   p.GetObjectMeta(),
			Object: p,
		}
	}

	return nil
}
