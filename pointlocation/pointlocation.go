package pointlocation

import "fmt"

func TrapezoidMap(ss []segment) (tr trapezoid, d *dag) {
	// Determine a bounding box R that contains all segments of S, and initialize
	// the trapezoidal map structure T and search structure D for it.
	tr = boundingBox(ss)
	d = newDAG(&tr)
	// fmt.Println("root main func", &d.root)

	for index := range ss {
		intersectedTrapeziods := d.findIntersectedTrapeziodFromSegment(ss[index])
		fmt.Println(&d.root, &intersectedTrapeziods[0])

		for tIndex := range intersectedTrapeziods {
			oldTr := intersectedTrapeziods[tIndex]
			newTrs := oldTr.addSegment(ss[index])
			newNode := createLeaves(newTrs, ss[index])
			currNode := intersectedTrapeziods[tIndex].dagRef
			fmt.Println("currNode", currNode)
			*currNode = newNode
		}

		fmt.Println(d)
	}
	return
}
