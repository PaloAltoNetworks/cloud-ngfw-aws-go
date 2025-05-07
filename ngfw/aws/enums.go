package aws

import "fmt"

type CommitState int

const (
	CommitStateFailed CommitState = iota
	CommitStateSuccess
	CommitStateCommitting
	CommitStatePrecommitDone
	CommitStateUnknown
	CommitStatePrevalidateDone
	CommitStateValidating
	CommitStateUncommitted
)

func (t CommitState) String() string {
	return CommitStateToString[t]
}

var CommitStateToString = map[CommitState]string{
	CommitStateFailed:          "Failed",
	CommitStateSuccess:         "Success",
	CommitStateCommitting:      "Committing",
	CommitStatePrecommitDone:   "PrecommitDone",
	CommitStatePrevalidateDone: "PrevalidateDone",
	CommitStateUnknown:         "Unknown",
	CommitStateValidating:      "Validating",
	CommitStateUncommitted:     "Uncommitted",
}

var StringToCommitState = map[string]CommitState{
	"Failed":          CommitStateFailed,
	"Success":         CommitStateSuccess,
	"Committing":      CommitStateCommitting,
	"Unknown":         CommitStateUnknown,
	"PrevalidateDone": CommitStatePrevalidateDone,
	"Validating":      CommitStateValidating,
	"Uncommitted":     CommitStateUncommitted,
}

func CommitStateFromString(state string) (CommitState, error) {
	if commitState, ok := StringToCommitState[state]; ok {
		return commitState, nil
	}
	return CommitStateUnknown, fmt.Errorf("invalid commit state: %s", state)
}

type FirewallStatus int

const (
	FirewallStatusCreating FirewallStatus = iota
	FirewallStatusUpdating
	FirewallStatusDeleting
	FirewallStatusCreateComplete
	FirewallStatusUpdateComplete
	FirewallStatusCreateFail
	FirewallStatusUpdateFail
	FirewallStatusDeleteFail
	FirewallStatusDeleteComplete
	FirewallStatusUnknown
)

func (t FirewallStatus) String() string {
	return FirewallStatusToString[t]
}

var FirewallStatusToString = map[FirewallStatus]string{
	FirewallStatusCreating:       "CREATING",
	FirewallStatusUpdating:       "UPDATING",
	FirewallStatusDeleting:       "DELETING",
	FirewallStatusCreateComplete: "CREATE_COMPLETE",
	FirewallStatusUpdateComplete: "UPDATE_COMPLETE",
	FirewallStatusCreateFail:     "CREATE_FAIL",
	FirewallStatusUpdateFail:     "UPDATE_FAIL",
	FirewallStatusDeleteFail:     "DELETE_FAIL",
	FirewallStatusDeleteComplete: "DELETE_COMPLETE",
	FirewallStatusUnknown:        "UNKNOWN",
}

func FirewallStatusFromString(status string) (FirewallStatus, error) {
	if fwStatus, ok := StringToFirewallStatus[status]; ok {
		return fwStatus, nil
	}
	return FirewallStatusUnknown, fmt.Errorf("invalid firewall status: %s", status)
}

var StringToFirewallStatus = map[string]FirewallStatus{
	"CREATING":        FirewallStatusCreating,
	"UPDATING":        FirewallStatusUpdating,
	"DELETING":        FirewallStatusDeleting,
	"CREATE_COMPLETE": FirewallStatusCreateComplete,
	"UPDATE_COMPLETE": FirewallStatusUpdateComplete,
	"CREATE_FAIL":     FirewallStatusCreateFail,
	"UPDATE_FAIL":     FirewallStatusUpdateFail,
	"DELETE_FAIL":     FirewallStatusDeleteFail,
	"DELETE_COMPLETE": FirewallStatusDeleteComplete,
	"Unknown":         FirewallStatusUnknown,
}
