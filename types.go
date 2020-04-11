package community_bl

import "github.com/satori/go.uuid"

type ApplicationState string

func (s ApplicationState) Valid() bool {

	switch s {
	case ApplicationStateApproved:
		return true
	case ApplicationStateRejected:
		return true
	case ApplicationStatePending:
		return true
	default:
		return false
	}

}

var ApplicationStateRejected = ApplicationState("Rejected")
var ApplicationStateApproved = ApplicationState("Approved")
var ApplicationStatePending = ApplicationState("Pending")

type MemberIdentifier = uuid.UUID
type ApplicationID = uuid.UUID
