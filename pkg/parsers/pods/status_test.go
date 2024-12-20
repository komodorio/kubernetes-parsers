package parsers_pods

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PodWithStatus creates a Pod with provided attributes for testing.
func PodWithStatus(phase corev1.PodPhase, reason string, initStatuses, containerStatuses []corev1.ContainerStatus,
	deletionTimestamp *metav1.Time, ready corev1.ConditionStatus, customSpec *corev1.PodSpec) corev1.Pod {
	spec := corev1.PodSpec{InitContainers: []corev1.Container{{}}}
	if customSpec != nil {
		spec = *customSpec
	}

	return corev1.Pod{
		Spec: spec,
		Status: corev1.PodStatus{
			Phase:                 phase,
			Reason:                reason,
			InitContainerStatuses: initStatuses,
			ContainerStatuses:     containerStatuses,
			Conditions: []corev1.PodCondition{
				{Type: corev1.PodReady, Status: ready},
			},
		},
		ObjectMeta: metav1.ObjectMeta{
			DeletionTimestamp: deletionTimestamp,
		},
	}
}

func TestFindStatusForPod(t *testing.T) {
	now := metav1.Now()
	tests := []struct {
		name     string
		pod      corev1.Pod
		expected string
	}{
		{
			name:     "Phase Reason",
			pod:      PodWithStatus(corev1.PodRunning, "SomeReason", nil, nil, nil, corev1.ConditionTrue, nil),
			expected: "SomeReason",
		},
		{
			name: "InitContainer Terminated",
			pod: PodWithStatus(corev1.PodPending, "", []corev1.ContainerStatus{
				{State: corev1.ContainerState{Terminated: &corev1.ContainerStateTerminated{ExitCode: 1, Reason: "Error"}}},
			}, nil, nil, corev1.ConditionFalse, nil),
			expected: "Init:Error",
		},
		{
			name: "InitContainer Waiting",
			pod: PodWithStatus(corev1.PodPending, "", []corev1.ContainerStatus{
				{State: corev1.ContainerState{Waiting: &corev1.ContainerStateWaiting{Reason: "PodInitializing"}}},
			}, nil, nil, corev1.ConditionFalse, nil),
			expected: "Init:0/1",
		},
		{
			name: "Container Terminated",
			pod: PodWithStatus(corev1.PodRunning, "", nil, []corev1.ContainerStatus{
				{State: corev1.ContainerState{Terminated: &corev1.ContainerStateTerminated{Reason: "Error", ExitCode: 1}}},
			}, nil, corev1.ConditionFalse, nil),
			expected: "Error",
		},
		{
			name:     "Pod Deletion Timestamp",
			pod:      PodWithStatus(corev1.PodRunning, "", nil, nil, &now, corev1.ConditionFalse, nil),
			expected: "Terminating",
		},
		{
			name:     "Pod Node Unreachable",
			pod:      PodWithStatus(corev1.PodRunning, "NodeLost", nil, nil, &now, corev1.ConditionUnknown, nil),
			expected: "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := FindStatusForPod(tt.pod)
			assert.Equal(t, tt.expected, status)
		})
	}
}

func TestFindStatusForPod_InitContainerRestartableBeforeVersion_1_29(t *testing.T) {
	initStatuses := []corev1.ContainerStatus{
		{State: corev1.ContainerState{Running: &corev1.ContainerStateRunning{StartedAt: metav1.Now()}}},
	}
	pod := PodWithStatus(corev1.PodPending, "", initStatuses, nil, nil, corev1.ConditionTrue, nil)
	status := FindStatusForPod(pod)
	assert.Equal(t, "Init:0/1", status)
}

// TestFindStatusForPod_InitContainerRestartableFromVersion_1_29 tests the FindStatusForPod function with a Pod that has a restartable init container.
// https://kubernetes.io/docs/concepts/workloads/pods/sidecar-containers/#pod-sidecar-containers
func TestFindStatusForPod_InitContainerRestartableFromVersion_1_29(t *testing.T) {
	restartPolicy := corev1.ContainerRestartPolicyAlways
	spec := &corev1.PodSpec{InitContainers: []corev1.Container{{RestartPolicy: &restartPolicy}}}

	initStatuses := []corev1.ContainerStatus{
		{State: corev1.ContainerState{Running: &corev1.ContainerStateRunning{StartedAt: metav1.Now()}}},
	}
	pod := PodWithStatus(corev1.PodRunning, "", initStatuses, nil, nil, corev1.ConditionTrue, spec)
	status := FindStatusForPod(pod)
	assert.Equal(t, "Running", status)
}

func TestFindStatusForPod_ContainerRunning(t *testing.T) {
	containerStatuses := []corev1.ContainerStatus{
		{State: corev1.ContainerState{Running: &corev1.ContainerStateRunning{}}, Ready: true},
	}
	pod := PodWithStatus(corev1.PodRunning, "Completed", nil, containerStatuses, nil, corev1.ConditionTrue, nil)
	status := FindStatusForPod(pod)
	assert.Equal(t, "Running", status)
}
