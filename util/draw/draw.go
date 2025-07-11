package draw

import (
	"example.com/m/util/utilerror"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

func DrawPoint(points plotter.XYs, xMax, yMax float64) *utilerror.UtilError {

	scatter, err := plotter.NewScatter(points)
	if err != nil {
		return utilerror.NewError(err.Error())
	}
	scatter.Shape = draw.CircleGlyph{}

	plt := plot.New()
	plt.Y.Min, plt.X.Min, plt.Y.Max, plt.X.Max = 0, 0, yMax, xMax

	plt.Add(scatter)

	if err := plt.Save(10*vg.Inch, 25*vg.Inch, "04-draw-dot.png"); err != nil {
		return utilerror.NewError(err.Error())
	}
	return nil
}
