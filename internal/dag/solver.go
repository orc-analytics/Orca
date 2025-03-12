package dag

type Node struct {
	from_algorithm uint32
	to_algorithm   uint32
	from_processor uint32
	to_processor   uint32
}

type DAG struct {
	Nodes []Node
}

// SolveDag - perform the dag solving activity
// given a set of algorithms, their triggering windows, processors
// and dependencies, find the execution paths broken down by processor
func SolveDag(dag DAG) (*DAG, error) {
	return nil, nil
}
