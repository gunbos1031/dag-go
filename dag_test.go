package dag_go

import (
	"fmt"
	"strings"
	"testing"
)

//TODO 12/6 보강 해야함.
func TestDag_AddEdge(t *testing.T) {
	// 그래프를 살펴보지는 않는다.
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

// TODO 12/15 이거 먼저 테스트 하고 가자.
func TestDag_GetLefMostNode(t *testing.T) {

}

func TestDag_SetNodes(t *testing.T) {

}

// TODO 12/15 Circle 관련 테스트 진행하면서 수정해야 함.
// no cycle 떠야 하는데 cycle 뜸. 흠흠.
func TestDag_CheckCircle(t *testing.T) {
	dag := NewDag()

	// dag.startNode.Id 는 시작노드 이다.
	// dag.startNode.Id 는 이미 노드가 있다. 이때 정상적인지도 살펴봐야 한다.
	err := dag.AddEdge(dag.startNode.Id, "test2")
	if err != nil {
		t.Fatalf("Failed to AddEdge: %v", err)
	}

	err = dag.AddEdge(dag.startNode.Id, "test3")
	if err != nil {
		t.Fatalf("Failed to AddEdge: %v", err)
	}

	str, err := dag.checkCycle()
	if err != nil {
		t.Fatalf("checkCycle: %v", err)
	}
	// 이 테스트 코드에서는 cycle 이면 에러임.
	if err == nil && strings.Contains(str, "Detect Cycle") {
		t.Fatal("Detect Cycle 이면 에러")
	}

	fmt.Println(str)
}
