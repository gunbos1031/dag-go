package dag_go

type Dag struct {
	Id string

	nodes   map[string]*Node
	hasEdge bool
}

/*
	노드는 기본적으로 두개의 슬라이스를 가지는데
	next, prev 이다.

	next, prev 같은 경우 AddEdge 시 next, prev 를 추가해준다.
*/

type Node struct {
	Id string

	children  []*Node // children
	dependsOn []*Node // parents

	next []*Node
	prev []*Node

	parentDag *Dag // 자신이 소속되어 있는 Dag
}

func NewDag() *Dag {
	this := new(Dag)
	this.nodes = make(map[string]*Node)
	this.Id = "root"
	return this
}

/*
	AddVertex 를 통해서 Node 를 생성하고 동시에 Dag struct 의 nodes 에 집어 넣어준다.
	현재 시점에서 Node 는 아직 아무것도 하지 않는 녀석이다. 추후 job 에 대한 부분을 넣어 주어야 한다.
*/
func (dag *Dag) AddVertex(id string) *Node {

	node := &Node{Id: id}
	node.parentDag = dag

	dag.nodes[id] = node
	return node
}

/*
	AddEdge 는 기본적으로 노드가 없을 경우 노드를 생성해준다.
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

	dag.hasEdge = true

	return nil
}
