package incidents

type DecisionType string

const (
	DecisionNoop    DecisionType = "NOOP"
	DecisionOpen    DecisionType = "OPEN_INCIDENT"
	DecisionUpdate  DecisionType = "UPDATE_INCIDENT"
	DecisionResolve DecisionType = "RESOLVE_INCIDENT"
)

type Decision struct {
	Type DecisionType
}
