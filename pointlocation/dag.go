package pointlocation

import (
	"fmt"
	"log"
	"reflect"
	"strings"
)

// node either by trapeziod, xnode or ynode
type node interface {
	check(point) node
	leftChild() node
	rightChild() node
}

type xNode struct {
	xCoordinate float64
	indegree    int
	lChild      node
	rChild      node
}

func (n xNode) String() string {
	return fmt.Sprintf("[%c]node xcoor = %v", 'x', n.xCoordinate)
}

func (n xNode) check(p point) node {
	if p.x < n.xCoordinate {
		return n.lChild
	}
	return n.rChild
}

func (n xNode) leftChild() node {
	return n.lChild
}
func (n xNode) rightChild() node {
	return n.rChild
}

type yNode struct {
	s        segment
	indegree int
	lChild   node
	rChild   node
}

func (n yNode) String() string {
	return fmt.Sprintf("[%c]node xcoor = %v", 'y', n.s)
}

func (n yNode) check(p point) node {
	pos, err := p.positionBySegment(n.s)
	if err != nil {
		log.Fatalln(err)
		return nil
	}
	if pos == lower {
		return n.lChild
	}
	return n.rChild
}

func (n yNode) leftChild() node {
	return n.lChild
}
func (n yNode) rightChild() node {
	return n.rChild
}

type dag struct {
	root node
}

func newDAG(bb *trapezoid) *dag {
	n := node(bb)
	d := dag{root: n}
	bb.dagRef = &d.root
	return &d
}

// inorder printing
func (d dag) String() string {
	stack := make([]node, 0)
	curr := d.root

	var strb strings.Builder

	for curr != nil || len(stack) > 0 {
		for curr != nil {
			stack = append(stack, curr)
			curr = curr.leftChild()
		}
		lastIndex := len(stack) - 1
		curr = stack[lastIndex] // Top element
		fmt.Println(reflect.TypeOf(curr))
		fmt.Println(curr)
		fmt.Println("_____________________________")
		_, err := strb.WriteString(fmt.Sprintf("%v\n\n", curr))
		if err != nil {
			log.Fatalln("unable to append string")
			return ""
		}
		stack = stack[:lastIndex] // Pop
		curr = curr.rightChild()
	}
	// fmt.Println("______________333_______________")
	// fmt.Printf("dag, %v\n", strb.String())
	// fmt.Println("______________333_______________")
	return fmt.Sprintf("dag, %v", strb.String())
}

// findPoint function finds trapeziod that point is inside
func (d *dag) findPoint(p point) (tr *node) {
	// if there is only boundery trapeziod,return boundery trapeziodal
	if d.root.leftChild() == nil && d.root.rightChild() == nil {
		fmt.Println("root", &d.root)
		tr = &d.root
		fmt.Println("root casted", tr)
		return
	}

	// else, try to drill down the dag and find such trapeziodal
	curr := &d.root
	for {
		if curr := (*curr).check(p); curr == nil {
			break
		}
	}
	tr = curr
	return
}

// func (d *dag) addSegment(s segment) *dag {
// 	// find start segment
// 	d.findPoint(s.startPoint)
// 	return d
// }

func (d *dag) findIntersectedTrapeziodFromSegment(s segment) (trs []*trapezoid) {
	trs = make([]*trapezoid, 0)
	startNode := d.findPoint(s.startPoint)
	startTr := (*startNode).(*trapezoid)
	trs = append(trs, startTr)

	trNode := (*startNode).(*trapezoid)
	fmt.Println("tr0", startNode)
	fmt.Println("startTr", &startTr)
	q := s.endPoint
	currRightP := trNode.rightp
	for trNode != nil && q.x > currRightP.x {
		if currRightP.orientationFromSegment(s) == counterclockwise {
			trs = append(trs, trNode.lowerRightN)
			trNode = trNode.lowerRightN
		} else {
			trs = append(trs, trNode.upperRightN)
			trNode = trNode.upperRightN
		}
		currRightP = trNode.rightp
	}
	return
}

func createLeaves(trs []trapezoid, s segment) (n node) {
	if len(trs) == 4 {
		x1 := xNode{xCoordinate: s.startPoint.x}
		y := yNode{s: s}
		x2 := xNode{xCoordinate: s.endPoint.x}
		x1.lChild = &trs[0]
		x1.rChild = &x2
		y.lChild = &trs[1]
		y.rChild = &trs[2]
		x2.lChild = &y
		x2.rChild = &trs[3]
		n1 := node(&trs[0])
		n2 := node(&trs[1])
		n3 := node(&trs[2])
		n4 := node(&trs[3])

		trs[0].dagRef = &n1
		trs[1].dagRef = &n2
		trs[2].dagRef = &n3
		trs[3].dagRef = &n4
		return &x1
	}
	return
}
