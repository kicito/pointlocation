package pointlocation

import (
	"fmt"
	"log"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func plotAddTr(p *plot.Plot, tr trapezoid) (err error) {
	lp, rp, t, b, bl, br, err := tr.plotData()
	if err != nil {
		return
	}
	p.Add(lp, rp, t, b, bl, br)
	return
}

func PlotTr(filename string, tr trapezoid) {

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
	p.Add(plotter.NewGrid())

	lp, rp, t, b, bl, br, err := tr.plotData()
	if err != nil {
		log.Fatal("unable to get plot data ", err)
	}
	p.Add(lp, rp, t, b, bl, br)
	if filename == "" {
		filename = "points"
	}
	p.Title.Text = filename

	// Save the plot to a PNG file.
	if err := p.Save(10*vg.Inch, 10*vg.Inch, fmt.Sprintf("%v.png", filename)); err != nil {
		panic(err)
	}
}

func PlotTrNode(filename string, trNode node) {

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
	p.Add(plotter.NewGrid())

	lp, rp, t, b, bl, br, err := trNode.(*trapezoidNode).tr.plotData()
	if err != nil {
		log.Fatal("unable to get plot data ", err)
	}
	p.Add(lp, rp, t, b, bl, br)
	if filename == "" {
		filename = "points"
	}
	p.Title.Text = filename
	// Save the plot to a PNG file.
	if err := p.Save(10*vg.Inch, 10*vg.Inch, fmt.Sprintf("%v.png", filename)); err != nil {
		panic(err)
	}
}

func PlotTrNodeWithPoint(filename string, trNode node, pp Point) {

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
	p.Add(plotter.NewGrid())

	lp, rp, t, b, bl, br, err := trNode.(*trapezoidNode).tr.plotData()
	if err != nil {
		log.Fatal("unable to get plot data ", err)
	}
	p.Add(lp, rp, t, b, bl, br)

	l, err := pp.scatter()
	if err != nil {
		return
	}
	p.Add(l)

	if filename == "" {
		filename = "points"
	}
	p.Title.Text = filename
	// Save the plot to a PNG file.
	if err := p.Save(10*vg.Inch, 10*vg.Inch, fmt.Sprintf("%v.png", filename)); err != nil {
		panic(err)
	}
}
