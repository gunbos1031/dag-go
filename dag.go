package dag_go

import (
	"fmt"
)

/*
	아직 작성중이라서 너무 심각하게 보지 말아주세요. by seoyhaein
*/

type tuple2 struct {
	id   int
	node *Node
}

type Dag struct {
	Id string

	/*
		이름 순서대로 순회하는 것을 생각해보자.
		현재 사용자가 집어넣어주는 임의의 string 으로 맵을 구성했지만
		id 에 특정 순서를 기록하여서 이것을 맵에 넣어주는 방향으로 작성한다.

		그럼 순회 하는데 보다 효과적일 것이다. 다만, 이런 이름을 순서대로 만들어서 넣어주는 문제가 발생한다.
	*/
	nodes map[string]*Node

	// 순회용도
	travel []*tuple2
	/*
		일단 시작노드와 끝노드는 반드시 하나여야 한다.(향후 추후 업데이트 할때 복수의 시작을 만들 수는 있음)
	*/
	startNode *Node // 시작 노드
	endNode   *Node // 끝 노드

	validated bool // dag 체크
	hasEdge   bool
}

/*
	노드는 기본적으로 두개의 슬라이스를 가지는데
	next, prev 이다.

	next, prev 같은 경우 AddEdge 시 next, prev 를 추가해준다.
*/

type Node struct {
	Id      string
	iter_id string // 순회하기 위한 ID

	children  []*Node // children
	dependsOn []*Node // parents

	indegree  int // indegree
	outdegree int // outdegree

	next []*Node
	prev []*Node

	parentDag *Dag // 자신이 소속되어 있는 Dag
}

/*
	최초로 만들어준다. 한번만 만들어 준다.
*/

func NewDag() *Dag {
	this := new(Dag)

	// 이녀석을 어떻게 할지 일단 고민중
	this.nodes = make(map[string]*Node)

	this.Id = "root"
	this.validated = false
	this.hasEdge = false

	return this
}

/*
	AddVertex 를 통해서 Node 를 생성하고 동시에 Dag struct 의 nodes 에 집어 넣어준다.
	현재 시점에서 Node 는 아직 아무것도 하지 않는 녀석이다. 추후 job 에 대한 부분을 넣어 주어야 한다.
*/
func (dag *Dag) AddVertex(id string) *Node {

	node := &Node{Id: id}
	//node.parentDag = dag
	/*
		일단은 그냥 넣는 방식입.
		child 를 알아서 상관은 없지만, travel 할때 음..
	*/
	dag.nodes[id] = node
	return node
}

/*
	AddEdge 는 기본적으로 노드가 없을 경우 노드를 생성해준다.
	에러처리는 해줘야 한다.

	가령 from 과 to 같은 경우,
	원을 이루는 경우, 등 기타 등등

*/

func (dag *Dag) AddEdge(from, to string) error {
	fromNode := dag.nodes[from]
	if fromNode == nil {
		fromNode = dag.AddVertex(from)
	}
	toNode := dag.nodes[to]
	if toNode == nil {
		toNode = dag.AddVertex(to)
	}

	/*
		from, to 는 일단 신경쓰지 말아주세요.
	*/
	fromNode.next = append(fromNode.next, toNode)
	fromNode.next = append(fromNode.next, toNode.next...)
	for _, b := range fromNode.prev {
		b.next = append(b.next, toNode)
		b.next = append(b.next, toNode.next...)
	}

	toNode.prev = append(toNode.prev, fromNode)
	toNode.prev = append(toNode.prev, fromNode.prev...)
	for _, b := range toNode.next {
		b.prev = append(b.prev, fromNode)
		b.prev = append(b.prev, fromNode.prev...)
	}

	fromNode.children = append(fromNode.children, toNode)
	toNode.dependsOn = append(toNode.dependsOn, fromNode)

	fromNode.outdegree++
	toNode.indegree++
	dag.hasEdge = true

	return nil
}

/*
	두가지 일을 해준다.
	첫째는 시작노드가 복수일 경우 에러처리한다. (복수의 시작노드는 제외하는 것으로 정함. 일단 추후 업데이트 할때 생각해보자.)
	끝노드들을 하나의 가상 노드로 모아준다.

	이 메서드는 반드시 dag 를 완성하기 위해서는 호출해주어야 한다.
	이 메서드는 사용자가 최종적으로 dag 를 완성한 후에 내부적으로 사용한다. 사용자가 사용하지 않는다.

	내부적으로 dag 가 맞는지는 판단하지 않고 있다.
*/

func (dag *Dag) validate() error {

	var (
		startNodes []*Node
		endNodes   []*Node
	)

	if dag.validated {
		return nil
	}

	if len(dag.nodes) == 0 {
		return fmt.Errorf("no vertex")
	}

	for _, b := range dag.nodes {
		if b.indegree == 0 {
			startNodes = append(startNodes, b)
		}
		if b.outdegree == 0 {
			endNodes = append(endNodes, b)
		}
	}

	if len(startNodes) > 1 {
		return fmt.Errorf("%v, dag: %s", "multiple start", dag.Id)
	}
	// start node 구별하기 위해
	startNodes[0].iter_id = "start_node"
	dag.startNode = startNodes[0]

	// 모든 끝 노드들을 하나의  dag id 를 노드 id 로 갖 노드로 만들고 그 노드로 모은다.
	endNodeId := fmt.Sprintf("end_%s", dag.Id) // 마지막 노드는 dag id 를 가짐.
	endNode := dag.AddVertex(endNodeId)

	for _, b := range endNodes {
		dag.AddEdge(b.Id, endNodeId)
	}

	dag.endNode = endNode
	// end node 구별하기 위해서, 안넣어줘도 될것 같지만 일단은 지금은 넣어줌.
	endNode.iter_id = "end_node"
	dag.validated = true

	return nil
}

/*
	Validate 호출 후 사용해야 한다.

*/

// TODO 이름 일단 바꾸자 지금은 그냥 이렇게 그리고 에러 및 예외 처리는 일단 나중에
func (dag *Dag) setTuple2() ([]*tuple2, error) {

	var t []*tuple2
	t = make([]*tuple2, len(dag.nodes)) // len, cap 같게 설정

	if !dag.validated {
		return nil, fmt.Errorf("%s", "validation needed")
	}
	// nodes 의 수가 하나면 싱글 노드이다.
	if len(dag.nodes) < 2 {
		t[0] = &tuple2{0, dag.startNode}
		return t, nil
	}
	// 루프를 안돌려도 될 거 같은데... 루프로 돌리지 말고 시작노드에서 자식 노드의 indegree 를 확인하고, 자식노드의 수에 따라서 순회하는 방식으로 돌아야 할 것 같다.
	for _, b := range dag.nodes {
		if dag.startNode.Id == b.Id {
			// slice 의 처음 값은 반드시 start node
			t[0] = &tuple2{0, b}
		}

		// t[1:] 부터 채워 넣는 과정이 필요하다.
		// tuple2 의 int 가 순서이다. 이것을 채워 넣어야 한다.
		// start node 자식 노드들을 확인해야 한다. 해당 노드의 indegree, outdegree 를 확인해야 한다.
		// 싱글 노드를 이미 체크 했기 때문에 여기서는 반드시 자식 노드가 하나는 있어야 한다. 만약 자식 노드가 없다면 이건 에러
		// 불필요한 체크일지라도 일단 체크를 하고 추후에 성능 개선할 때 지워나가도록 하겠다.

	}

	return nil, nil
}

// 일단 테스트용
func (dag *Dag) travelDag() {

}
