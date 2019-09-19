package extractors

import (
	"context"

	"go.aporeto.io/trireme-lib/common"
	"go.aporeto.io/trireme-lib/policy"

	corev1 "github.com/kubernetes/core/v1"
	"github.com/kubernetes/apimachinery/pkg/runtime"

	"github.com/kubernetes-sigs/controller-runtime/pkg/client"
)

// EventMetadataExtractor is a function used to extract a *policy.PURuntime from a given
// EventInfo. The EventInfo is generic and is provided over the RPC interface
type EventMetadataExtractor func(*common.EventInfo) (*policy.PURuntime, error)

// PodMetadataExtractor is a function used to extract a *policy.PURuntime from a given
// Kubernetes pod. It can furthermore extract more information using the client.
// The 5th argument (bool) indicates if a network namespace should get extracted
type PodMetadataExtractor func(context.Context, client.Client, *runtime.Scheme, *corev1.Pod, bool) (*policy.PURuntime, error)

// PodSandboxExtractor is a function used to extract the SandboxID from a given pod.
type PodSandboxExtractor func(context.Context, *corev1.Pod) (string, error)

// PodNetclsProgrammer is a function used to program the net_cls cgroup of a pod for Trireme.
// This has to be used when Trireme is used in conjunction with pods that are in HostNetwork=true mode.
type PodNetclsProgrammer func(context.Context, *corev1.Pod, policy.RuntimeReader) error

// ResetNetclsKubepods is a function which must implement to reset all netcls cgroup programming of Trireme.
// It is called during Resync events in monitors and guarantees a fresh slate for the monitors for Kubernetes.
type ResetNetclsKubepods func(context.Context) error
