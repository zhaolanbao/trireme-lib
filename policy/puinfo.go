package policy

import "go.aporeto.io/trireme-lib/common"

// PUInfo  captures all policy information related to a connection as well as runtime.
// It makes passing data around simpler.
type PUInfo struct {
	// ContextID is the ID of the container that the policy applies to
	ContextID string
	// Policy is an instantiation of the container policy
	Policy *PUPolicy
	// RunTime captures all data that are captured from the container
	Runtime *PURuntime
}

// NewPUInfo instantiates a new ContainerPolicy
func NewPUInfo(contextID, namespace string, puType common.PUType) *PUInfo {
	policy := NewPUPolicy(contextID, namespace, AllowAll, nil, nil, nil, nil, nil, nil, nil, nil, nil, 0, 0, nil, nil, []string{})
	runtime := NewPURuntime("", 0, "", nil, nil, puType, nil)
	return PUInfoFromPolicyAndRuntime(contextID, policy, runtime)
}

// PUInfoFromPolicyAndRuntime generates a ContainerInfo Struct from an existing RuntimeInfo and PolicyInfo
func PUInfoFromPolicyAndRuntime(contextID string, policyInfo *PUPolicy, runtimeInfo *PURuntime) *PUInfo {
	return &PUInfo{
		ContextID: contextID,
		Policy:    policyInfo,
		Runtime:   runtimeInfo,
	}
}
