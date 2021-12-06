package dag_go

import (
	"testing"
)

//TODO 12/6 보강 해야함.
func TestDag_AddEdge(t *testing.T) {

	dag := NewDag()
	err := dag.AddEdge("test1", "test2")

	if err != nil {
		t.Fatalf("Failed to AddEdge #1: %v", err)
	}

	err = dag.AddEdge("test2", "test3")

	if err != nil {
		t.Fatalf("Failed to AddEdge #2: %v", err)
	}

	err = dag.AddEdge("test3", "test4")

	if err != nil {
		t.Fatalf("Failed to AddEdge #3: %v", err)
	}
}
