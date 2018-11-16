package pointlocation

import (
	"fmt"
	"log"
	"strings"
)

// node either by trapeziod, xnode or ynode
type node interface {
	check(Point) node
	leftChild() node
	rightChild() node
	assignLeft(node)
	assignRight(node)
	replaceWith(node) node
}

type trapezoidNode struct {
	tr      *trapezoid
	parents []node
}

func (n trapezoidNode) String() string {
	return fmt.Sprintf("[%v]node tr = %v, parents = %v", "tr", n.tr, n.parents)
}

// implementation of node interface
func (trapezoidNode) check(i Point) node {
	return nil
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

func (tr *trapezoidNode) replaceWith(n node) node {
	for index := range tr.parents {
		parent := tr.parents[index]
		if parent.leftChild() == tr {
			parent.assignLeft(n)
		} else {
			if parent.rightChild() != tr {
				fmt.Println("---------------------------")
				fmt.Println(parent.leftChild())
				fmt.Println("---------------------------")
				fmt.Println(parent.rightChild())
				fmt.Println("---------------------------")
				fmt.Println("---------------------------")
				log.Fatal("can't find right child to replace with")
			}
			parent.assignRight(n)
		}
	}
	return n
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

func (n xNode) check(p Point) node {
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

func (n *xNode) assignLeft(nn node) {
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
	n.rChild = nn
	if trapezoidNode, ok := nn.(*trapezoidNode); ok {
		trapezoidNode.parents = append(trapezoidNode.parents, n)
	} else if xNode, ok := nn.(*xNode); ok {
		xNode.parent = n
	} else if yNode, ok := nn.(*yNode); ok {
		yNode.parent = n
	}
}

func (n *xNode) replaceWith(nn node) node {
	if n.parent.leftChild() == n {
		n.parent.assignLeft(nn)
	} else {
		if n.parent.rightChild() != n {
			log.Fatal("can't find right child to replace with")
		}
		n.parent.assignRight(nn)
	}
	return nn
}

type yNode struct {
	s      Segment
	lChild node
	rChild node
	parent node
}

func (n yNode) String() string {
	return fmt.Sprintf("[%c]node xcoor = %v, parent = %v", 'y', n.s, n.parent)
}

func (n yNode) check(p Point) node {
	pos, err := p.positionBySegment(n.s)
	if err != nil {
		log.Fatalln(err)
		return nil
	}
	if pos == lower {
		return n.rChild
	}
	return n.lChild
}

func (n yNode) leftChild() node {
	return n.lChild
}
func (n yNode) rightChild() node {
	return n.rChild
}

func (n *yNode) assignLeft(nn node) {
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
	n.rChild = nn
	if trapezoidNode, ok := nn.(*trapezoidNode); ok {
		trapezoidNode.parents = append(trapezoidNode.parents, n)
	} else if xNode, ok := nn.(*xNode); ok {
		xNode.parent = n
	} else if yNode, ok := nn.(*yNode); ok {
		yNode.parent = n
	}
}

func (n *yNode) replaceWith(nn node) node {
	if n.parent.leftChild() == n {
		n.parent.assignLeft(nn)
	} else {
		if n.parent.rightChild() != n {
			log.Fatal("can't find right child to replace with")
		}
		n.parent.assignRight(n)
	}
	return nn
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
func (d *dag) FindPoint(p Point) (tr node) {
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
		curr = curr.check(p)
		if _, ok := curr.(*trapezoidNode); ok {
			break
		}
	}
	tr = curr
	return
}

// func (d *dag) addSegment(s Segment) {
// 	startNode := d.FindPoint(s.startPoint)
// 	endNode := d.FindPoint(s.endPoint)
// 	trStart := startNode.(*trapezoidNode).tr
// 	trEnd := endNode.(*trapezoidNode).tr

// 	// case Segment connecting 2 other segments
// 	if trStart.leftp.x == s.startPoint.x &&
// 		trStart.leftp.y == s.startPoint.y &&
// 		trStart.rightp.x == s.endPoint.x &&
// 		trStart.rightp.y == s.endPoint.y {
// 		d.addConnectingCase(startNode, &s)
// 	} else if trStart == trEnd {
// 		d.addSimpleCase(startNode, &s)
// 	} else {
// 		d.addHardCase(startNode, endNode, &s)
// 	}
// }

// func (d *dag) addConnectingCase(n node, s *Segment) {
// 	tr := n.(*trapezoidNode).tr
// 	trNewTop := &trapezoid{
// 		leftp:       &s.startPoint,
// 		rightp:      tr.rightp,
// 		top:         tr.top,
// 		bottom:      s,
// 		upperLeftN:  tr.upperLeftN,
// 		upperRightN: tr.upperRightN,
// 	}

// 	trNewBot := &trapezoid{
// 		leftp:       &s.startPoint,
// 		rightp:      tr.rightp,
// 		top:         s,
// 		bottom:      tr.bottom,
// 		lowerLeftN:  tr.lowerLeftN,
// 		lowerRightN: tr.lowerRightN,
// 	}
// 	utNode := trapezoidNode{tr: trNewTop, parents: make([]node, 0)}
// 	btNode := trapezoidNode{tr: trNewBot, parents: make([]node, 0)}
// 	trNewTop.dagRef = &utNode
// 	trNewBot.dagRef = &btNode
// 	y := yNode{s: *s}
// 	y.assignLeft(&utNode)
// 	y.assignRight(&btNode)

// 	n.replaceWith(&y)

// }

// func (d *dag) addSimpleCase(n node, s *Segment) {
// 	tr := n.(*trapezoidNode).tr
// 	lt, ut, rt, bt := tr.addSegmentInside(*s)

// 	// if it's shared endpoints, fix neighbors and add only x node, y node
// 	if (*tr).leftp.x == s.startPoint.x && (*tr).leftp.y == s.startPoint.y {
// 		// fixing neighbors
// 		if tr.upperLeftN != nil && tr.upperLeftN.upperRightN == tr {
// 			tr.upperLeftN.upperRightN = ut
// 		}
// 		if tr.lowerLeftN != nil && tr.lowerLeftN.lowerRightN == tr {
// 			tr.lowerLeftN.lowerRightN = bt
// 		}

// 		utNode := trapezoidNode{tr: ut, parents: make([]node, 0)}
// 		rtNode := trapezoidNode{tr: rt, parents: make([]node, 0)}
// 		btNode := trapezoidNode{tr: bt, parents: make([]node, 0)}
// 		ut.dagRef = &utNode
// 		rt.dagRef = &rtNode
// 		bt.dagRef = &btNode

// 		y := yNode{s: *s}
// 		y.assignLeft(&utNode)
// 		y.assignRight(&btNode)
// 		x2 := xNode{xCoordinate: s.endPoint.x}
// 		x2.assignLeft(&y)
// 		x2.assignRight(&rtNode)
// 		n.replaceWith(&x2)
// 		return
// 	}

// 	if tr.upperLeftN != nil && tr.upperLeftN.upperRightN == tr {
// 		tr.upperLeftN.upperRightN = lt
// 		lt.upperLeftN = tr.upperLeftN
// 	}
// 	if tr.lowerLeftN != nil && tr.lowerLeftN.lowerRightN == tr {
// 		tr.lowerLeftN.lowerRightN = lt
// 		lt.lowerLeftN = tr.lowerLeftN
// 	}

// 	ltNode := trapezoidNode{tr: lt, parents: make([]node, 0)}
// 	utNode := trapezoidNode{tr: ut, parents: make([]node, 0)}
// 	rtNode := trapezoidNode{tr: rt, parents: make([]node, 0)}
// 	btNode := trapezoidNode{tr: bt, parents: make([]node, 0)}
// 	lt.dagRef = &ltNode
// 	ut.dagRef = &utNode
// 	rt.dagRef = &rtNode
// 	bt.dagRef = &btNode

// 	x1 := xNode{xCoordinate: s.startPoint.x}
// 	y := yNode{s: *s}
// 	x2 := xNode{xCoordinate: s.endPoint.x}
// 	x1.assignLeft(&ltNode)
// 	x1.assignRight(&x2)
// 	y.assignLeft(&utNode)
// 	y.assignRight(&btNode)
// 	x2.assignLeft(&y)
// 	x2.assignRight(&rtNode)

// 	if n == d.root {
// 		d.root = &x1
// 	} else {
// 		n.replaceWith(&x1)
// 	}
// }

// func (d *dag) addHardCase(startNode, endNode node, s *Segment) {
// 	trStart := startNode.(*trapezoidNode).tr
// 	trEnd := endNode.(*trapezoidNode).tr
// 	var trNewTop, trNewBot *trapezoid
// 	trNewTop, trNewBot = d.addHardCaseStart(trStart, s)

// 	for trCurr := trStart.nextIntersect(*s); trCurr != nil; trCurr = trCurr.nextIntersect(*s) {
// 		// check if both trapeziodal has same top or bottom, merge trapezoid
// 		// merging trapeziod: set right p of previous to new one
// 		if trNewTop.top == trCurr.top || trNewTop.bottom == trCurr.bottom {
// 			if s.endPoint.x < trCurr.rightp.x {
// 				trNewTop.rightp = &s.endPoint
// 			} else {
// 				trNewTop.rightp = trCurr.rightp
// 			}
// 		} else {
// 			trFixedTop := &trapezoid{
// 				leftp:  trCurr.leftp,
// 				top:    trCurr.top,
// 				bottom: s,
// 			}
// 			if s.endPoint.x < trCurr.rightp.x {
// 				trFixedTop.rightp = &s.endPoint
// 			} else {
// 				trFixedTop.rightp = trCurr.rightp
// 			}
// 			trFixedTop.lowerLeftN = trNewTop
// 			trNewTop.lowerRightN = trFixedTop
// 			trNewTop = trFixedTop
// 		}

// 		if trNewBot.top == trCurr.top || trNewBot.bottom == trCurr.bottom {
// 			if s.endPoint.x < trCurr.rightp.x {
// 				trNewBot.rightp = &s.endPoint
// 			} else {
// 				trNewBot.rightp = trCurr.rightp
// 			}
// 		} else {
// 			trFixedBot := &trapezoid{
// 				leftp:  trCurr.leftp,
// 				top:    s,
// 				bottom: trCurr.bottom,
// 			}
// 			if s.endPoint.x < trCurr.rightp.x {
// 				trFixedBot.rightp = &s.endPoint
// 			} else {
// 				trFixedBot.rightp = trCurr.rightp
// 			}
// 			trFixedBot.upperLeftN = trNewBot
// 			trNewBot.upperRightN = trFixedBot
// 			trNewBot = trFixedBot
// 		}
// 		trNewTop.upperRightN = trCurr.upperRightN
// 		trNewBot.lowerRightN = trCurr.lowerRightN

// 		if trCurr == trEnd {
// 			break
// 		}
// 		fmt.Println("trNewTop", trNewTop)
// 		fmt.Println("trNewBot", trNewBot)

// 		nodeTrNewTop := trapezoidNode{tr: trNewTop}
// 		nodeTrNewBot := trapezoidNode{tr: trNewBot}
// 		trNewTop.dagRef = &nodeTrNewTop
// 		trNewBot.dagRef = &nodeTrNewBot
// 		// create new node to replace trapeziodal node in dag
// 		newY := yNode{s: *s}
// 		newY.assignLeft(&nodeTrNewTop)
// 		newY.assignRight(&nodeTrNewBot)

// 		fmt.Println("trCurr.dagRef before replace", trCurr.dagRef)
// 		fmt.Println("replacing with", newY)

// 		trCurr.dagRef = trCurr.dagRef.replaceWith(&newY)
// 		fmt.Println("trCurr.dagRef after replace", trCurr.dagRef)
// 	}

// 	d.addHardCaseEnd(trEnd, trNewTop, trNewBot, s)

// }

// func (d *dag) addHardCaseEnd(trEnd, trNewTop, trNewBot *trapezoid, s *Segment) {

// 	nodeNewTopEnd := trapezoidNode{tr: trNewTop}
// 	nodeNewBotEnd := trapezoidNode{tr: trNewBot}
// 	trNewTop.dagRef = &nodeNewTopEnd
// 	trNewBot.dagRef = &nodeNewBotEnd
// 	yEnd := yNode{s: *s}
// 	yEnd.assignLeft(&nodeNewTopEnd)
// 	yEnd.assignRight(&nodeNewBotEnd)

// 	// if leftp is same as endpoint Segment meaning that it has shared endpoints then fix neighbors
// 	if (*trEnd).leftp.x == s.endPoint.x && (*trEnd).leftp.y == s.endPoint.y {

// 		// line Segment is intersected in lowerpart, else, line Segment is intersect upperpart
// 		// fix neighbors
// 		if *trEnd.top == *trNewTop.top {
// 			trNewTop.lowerRightN = trEnd
// 			trEnd.lowerLeftN = trNewBot
// 		} else if *trEnd.bottom != *trNewBot.bottom {
// 			// if *trEnd.bottom != *trNewBot.bottom {
// 			// 	log.Fatal("something went wrong, this case shouldn't happen!")
// 			// }
// 			trNewBot.lowerRightN = trEnd
// 			trEnd.lowerLeftN = trNewBot
// 		}

// 		return
// 	}
// 	trNewEnd := trapezoid{
// 		leftp:  &s.endPoint,
// 		rightp: trEnd.rightp,
// 		top:    trEnd.top,
// 		bottom: trEnd.bottom,
// 	}

// 	trNewEnd.upperLeftN = trNewTop
// 	trNewEnd.lowerLeftN = trNewBot

// 	trNewEnd.upperRightN = trEnd.upperRightN
// 	trNewEnd.lowerRightN = trEnd.lowerRightN
// 	// fixing exist neighbors

// 	fmt.Println("trNewTop", trNewTop)
// 	fmt.Println("trNewBot", trNewBot)
// 	trNewTop.upperRightN = &trNewEnd
// 	trNewBot.lowerRightN = &trNewEnd
// 	nodeNewEnd := trapezoidNode{tr: &trNewEnd}
// 	trNewEnd.dagRef = &nodeNewEnd
// 	xEnd := xNode{xCoordinate: s.endPoint.x}
// 	xEnd.assignLeft(&yEnd)
// 	xEnd.assignRight(&nodeNewEnd)
// 	fmt.Println("endNode before replace", trEnd.dagRef)
// 	fmt.Println("replacing with", yEnd)
// 	trEnd.dagRef = trEnd.dagRef.replaceWith(&xEnd)
// 	fmt.Println("endNode after replace", trEnd.dagRef)

// }

// func (d *dag) addHardCaseStart(trStart *trapezoid, s *Segment) (trNewTop, trNewBot *trapezoid) {

// 	// in case of trapeziod shared same endpoint, replace node with ynode
// 	if (*trStart).leftp.x == s.startPoint.x && (*trStart).leftp.y == s.startPoint.y {
// 		trNewTop, trNewBot = d.sameEndpoint(trStart, s)

// 		nodeTrNewTop := trapezoidNode{tr: trNewTop}
// 		nodeTrNewBot := trapezoidNode{tr: trNewBot}
// 		trNewTop.dagRef = &nodeTrNewTop
// 		trNewBot.dagRef = &nodeTrNewBot
// 		// create new node to replace trapeziodal node in dag
// 		newY := yNode{s: *s}
// 		newY.assignLeft(&nodeTrNewTop)
// 		newY.assignRight(&nodeTrNewBot)

// 		fmt.Println("trStart.dagRef before replace", trStart.dagRef)
// 		fmt.Println("replacing with", newY)

// 		trStart.dagRef = trStart.dagRef.replaceWith(&newY)
// 		fmt.Println("trStart.dagRef after replace", trStart.dagRef)
// 		return
// 	}
// 	trNewStart := trapezoid{
// 		leftp:  trStart.leftp,
// 		rightp: &s.startPoint,
// 		top:    trStart.top,
// 		bottom: trStart.bottom,
// 	}

// 	trNewTop = &trapezoid{
// 		leftp:  &s.startPoint,
// 		rightp: trStart.rightp,
// 		top:    trStart.top,
// 		bottom: s,
// 	}

// 	trNewBot = &trapezoid{
// 		leftp:  &s.startPoint,
// 		rightp: trStart.rightp,
// 		top:    s,
// 		bottom: trStart.bottom,
// 	}

// 	trNewStart.upperRightN = trNewTop
// 	trNewStart.lowerRightN = trNewBot
// 	trNewStart.upperLeftN = trStart.upperLeftN
// 	trNewStart.lowerLeftN = trStart.lowerLeftN

// 	trNewTop.upperLeftN = &trNewStart
// 	trNewBot.lowerLeftN = &trNewStart

// 	// fixing exist neighbors
// 	trStart.replaceLeftNeighborsWith(&trNewStart)

// 	nodeNewStart := trapezoidNode{tr: &trNewStart}
// 	nodeNewTop := trapezoidNode{tr: trNewTop}
// 	nodeNewBot := trapezoidNode{tr: trNewBot}
// 	trNewStart.dagRef = &nodeNewStart
// 	trNewTop.dagRef = &nodeNewTop
// 	trNewBot.dagRef = &nodeNewBot
// 	x1 := xNode{xCoordinate: s.startPoint.x}
// 	y := yNode{s: *s}
// 	x1.assignLeft(&nodeNewStart)
// 	x1.assignRight(&y)
// 	y.assignLeft(&nodeNewTop)
// 	y.assignRight(&nodeNewBot)

// 	fmt.Println("startnode before replace", trStart.dagRef)
// 	fmt.Println("replacing with", x1)

// 	trStart.dagRef = trStart.dagRef.replaceWith(&x1)

// 	fmt.Println("startnode after replace", trStart.dagRef)
// 	return
// }

// func (d *dag) sameEndpoint(trCurr *trapezoid, s *Segment) (trNewTop, trNewBot *trapezoid) {

// 	if (*trCurr).leftp.x == s.startPoint.x && (*trCurr).leftp.y == s.startPoint.y {
// 		trNewTop = &trapezoid{
// 			leftp:  &s.startPoint,
// 			rightp: trCurr.rightp,
// 			top:    trCurr.top,
// 			bottom: s,
// 		}

// 		trNewBot = &trapezoid{
// 			leftp:  &s.startPoint,
// 			rightp: trCurr.rightp,
// 			top:    s,
// 			bottom: trCurr.bottom,
// 		}

// 		if trCurr.upperLeftN != nil && trCurr.upperLeftN.upperRightN == trCurr {
// 			trCurr.upperLeftN.upperRightN = trNewTop
// 			trNewTop.upperLeftN = trCurr.upperLeftN
// 		}

// 		if trCurr.lowerLeftN != nil && trCurr.lowerLeftN.lowerRightN == trCurr {
// 			trCurr.lowerLeftN.lowerRightN = trNewBot
// 			trNewBot.lowerLeftN = trCurr.lowerLeftN
// 		}

// 		trNewTop.upperRightN = trCurr.upperRightN
// 		trNewBot.lowerRightN = trCurr.lowerRightN

// 	} else if *trCurr.rightp == s.endPoint {

// 	} else {
// 		fmt.Println("---------------------")
// 		fmt.Println("something went wrong at sameendpoint")
// 		fmt.Println(*trCurr)
// 		fmt.Println(s)
// 		log.Fatal("--------------------")
// 	}

// 	return
// }
