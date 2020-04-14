// scatter
package scatter

import (
	"fmt"
	"math/rand"

	"image/color"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func UFRanNums(N, trials int, filename string) {
	var r1, r2 float64
	var n = 0
	data := make(plotter.XYs, N)
	for {
		r2 = rand.Float64()
		if (r2 > 0.700001) && (r2 < 0.700009) {
			data[n].X = r2
			data[n].Y = r1
			n += 1
		}

		r1 = r2

		if n == N {
			PlotScatter(data, filename)
			fmt.Println(filename + ".png is done!")
			fmt.Println("Length of data：", len(data), "\n")
			break
		}
	}
}

func FRanNums(N, trials int, Seed int64, filename string) {
	var r1, r2 float64
	var n = 0
	data := make(plotter.XYs, N)
	s := rand.New(rand.NewSource(Seed))
	for {
		r2 = s.Float64()
		if (r2 > 0.700001) && (r2 < 0.700009) {
			data[n].X = r2
			data[n].Y = r1
			n += 1
		}

		r1 = r2

		if n == N {
			PlotScatter(data, filename)
			fmt.Println(filename + ".png is done!")
			fmt.Println("Length of data：", len(data), "\n")
			break
		}
	}
}

func PlotScatter(data plotter.XYs, filename string) {

	p, _ := plot.New()

	p.Title.Text = "Lattice -- " + filename
	p.X.Label.Text = "Xn"
	p.Y.Label.Text = "Xn-1"

	scatter, _ := plotter.NewScatter(data)
	scatter.Color = color.RGBA
	p.Add(scatter)

	p.Save(4*vg.Inch, 4*vg.Inch, "plot/"+filename+".png")

}
