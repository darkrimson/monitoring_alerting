package incidents

type Evaluator struct {
	FailureThreshold int
}

func NewEvaluator(threshold int) *Evaluator {
	return &Evaluator{
		FailureThreshold: threshold,
	}
}

type EvaluateInput struct {
	HasOpenIncident bool
	FailureCount    int // подряд идущие DOWN до этого check

	CheckStatus string // UP / DOWN
}

const (
	up   = "UP"
	down = "DOWN"
)

func (e *Evaluator) Evaluate(input EvaluateInput) Decision {
	switch input.CheckStatus {

	case up:
		if input.HasOpenIncident {
			return Decision{Type: DecisionResolve}
		}
		return Decision{Type: DecisionNoop}

	case down:
		if input.HasOpenIncident {
			return Decision{Type: DecisionUpdate}
		}

		if input.FailureCount+1 >= e.FailureThreshold {
			return Decision{Type: DecisionOpen}
		}

		return Decision{Type: DecisionNoop}

	default:
		return Decision{Type: DecisionNoop}
	}
}
