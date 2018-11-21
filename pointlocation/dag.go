package pointlocation

import (
	"fmt"
	"log"
	"strings"

	"github.com/pkg/errors"
)

// node either by trapeziod, xnode or ynode
type node interface {
	check(Point) (node, error)
	leftChild() node
	rightChild() node
	assignLeft(node)
	assignRight(node)
	replaceWith(node)
}

type trapezoidNode struct {
	tr      *trapezoid
	parents []node
}

func (n trapezoidNode) String() string {
	return fmt.Sprintf("[%v]node tr = %v, parents = %v", "tr", n.tr, n.parents)
}

// implementation of node interface
func (trapezoidNode) check(i Point) (node, error) {
	return nil, nil
}

func (trapezoidNode) leftChild() node {
	return nil
}

func (trapezoidNode) rightChild() node {
	return nil
}

func (trapezoidNode) assignLeft(n node) {
	log.Fatal("this node shouldn't called this method")
}

func (trapezoidNode) assignRight(n node) {
	log.Fatal("this node shouldn't called this method")
}

func (tr *trapezoidNode) replaceWith(n node) {
	for index := range tr.parents {
		parent := tr.parents[index]
		if parent.leftChild() == tr {
			parent.assignLeft(n)
		} else {
			parent.assignRight(n)
		}
	}
}
func (n *trapezoidNode) removeParent(nn node) {
	for pIndex := range n.parents {
		if n.parents[pIndex] == nn {
			n.parents = append(n.parents[:pIndex], n.parents[pIndex+1:]...)
		}
	}
}

type xNode struct {
	xCoordinate float64
	lChild      node
	rChild      node
	parent      node
}

func (n xNode) String() string {
	return fmt.Sprintf("[%c]node xcoor = %v, parent = %v", 'x', n.xCoordinate, n.parent)
}

func (n xNode) check(p Point) (node, error) {
	if p.x < n.xCoordinate {
		return n.lChild, nil
	}
	return n.rChild, nil
}

func (n xNode) leftChild() node {
	return n.lChild
}
func (n xNode) rightChild() node {
	return n.rChild
}

func (n *xNode) assignLeft(nn node) {
	//fix current lchild
	if trNode, ok := n.lChild.(*trapezoidNode); ok {
		trNode.removeParent(n)
	}
	n.lChild = nn

	if trapezoidNode, ok := nn.(*trapezoidNode); ok {
		trapezoidNode.parents = append(trapezoidNode.parents, n)
	} else if xNode, ok := nn.(*xNode); ok {
		xNode.parent = n
	} else if yNode, ok := nn.(*yNode); ok {
		yNode.parent = n
	}
}

func (n *xNode) assignRight(nn node) {
	if trNode, ok := n.rChild.(*trapezoidNode); ok {
		trNode.removeParent(n)
	}
	n.rChild = nn
	if trapezoidNode, ok := nn.(*trapezoidNode); ok {
		trapezoidNode.parents = append(trapezoidNode.parents, n)
	} else if xNode, ok := nn.(*xNode); ok {
		xNode.parent = n
	} else if yNode, ok := nn.(*yNode); ok {
		yNode.parent = n
	}
}

func (n *xNode) replaceWith(nn node) {
	if n.parent.leftChild() == n {
		n.parent.assignLeft(nn)
	} else {
		if n.parent.rightChild() != n {
			log.Fatal("can't find right child to replace with")
		}
		n.parent.assignRight(nn)
	}
}

type yNode struct {
	s      Segment
	lChild node
	rChild node
	parent node
}

func (n yNode) String() string {
	return fmt.Sprintf("[%c]node segment = %v, parent = %v", 'y', n.s, n.parent)
}

func (n yNode) check(p Point) (node, error) {
	if n.s.startPoint.sameCoordinate(p) || n.s.endPoint.sameCoordinate(p) {
		if *n.s.slope > *p.s.slope {
			return n.rChild, nil
		}
		return n.lChild, nil
	}
	pos, err := p.positionBySegment(n.s)
	if err != nil {
		err = errors.Wrapf(err, "node:%v", n)
		return nil, err
	}
	if pos == lower {
		return n.rChild, nil
	}
	return n.lChild, nil
}

func (n yNode) leftChild() node {
	return n.lChild
}
func (n yNode) rightChild() node {
	return n.rChild
}

func (n *yNode) assignLeft(nn node) {
	if trNode, ok := n.lChild.(*trapezoidNode); ok {
		trNode.removeParent(n)
	}
	n.lChild = nn
	if trapezoidNode, ok := nn.(*trapezoidNode); ok {
		trapezoidNode.parents = append(trapezoidNode.parents, n)
	} else if xNode, ok := nn.(*xNode); ok {
		xNode.parent = n
	} else if yNode, ok := nn.(*yNode); ok {
		yNode.parent = n
	}
}

func (n *yNode) assignRight(nn node) {
	if trNode, ok := n.rChild.(*trapezoidNode); ok {
		trNode.removeParent(n)
	}
	n.rChild = nn
	if trapezoidNode, ok := nn.(*trapezoidNode); ok {
		trapezoidNode.parents = append(trapezoidNode.parents, n)
	} else if xNode, ok := nn.(*xNode); ok {
		xNode.parent = n
	} else if yNode, ok := nn.(*yNode); ok {
		yNode.parent = n
	}
}

func (n *yNode) replaceWith(nn node) {
	if n.parent.leftChild() == n {
		n.parent.assignLeft(nn)
	} else {
		if n.parent.rightChild() != n {
			log.Fatal("can't find right child to replace with")
		}
		n.parent.assignRight(n)
	}
}

type dag struct {
	root node
}

func newDAG(bb *trapezoid) dag {
	n := trapezoidNode{tr: bb}
	d := dag{root: &n}
	bb.dagRef = &n
	return d
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
		// fmt.Println(reflect.TypeOf(curr))
		// fmt.Println(curr.leftChild())
		// fmt.Println("_____________________________")
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
	return fmt.Sprintf("%v", strb.String())
}

// FindPoint function finds trapeziod that Point is inside
func (d *dag) FindPoint(p Point) (tr node, err error) {
	// if there is only boundery trapeziod,return boundery trapeziodal
	if d.root.leftChild() == nil && d.root.rightChild() == nil {
		// fmt.Println("root", &d.root)
		tr = d.root
		// fmt.Println("root casted", tr)
		return
	}

	// else, try to drill down the dag and find such trapeziodal
	curr := d.root
	for {
		if curr, err = curr.check(p); err != nil {
			err = errors.Wrapf(err, "Finding point:%v\n", p)
			return
		}
		if _, ok := curr.(*trapezoidNode); ok {
			break
		}
	}
	tr = curr
	return
}
