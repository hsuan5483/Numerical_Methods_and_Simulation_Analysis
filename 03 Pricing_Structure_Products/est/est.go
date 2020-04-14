// est
package est

import (
	"math"

	// "image/color"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"

	"gonum.org/v1/gonum/stat/distmv"
)

// flow:[t][c]
func PV(flow [][]float64, r float64) float64 {
	pv := 0.
	for i := range flow {
		f := flow[i]
		pv += f[1] * math.Exp(-r*f[0])
	}
	return pv
}

func Hist(data []float64, title string, normalize bool) {
	// 轉格式
	v := make(plotter.Values, len(data))
	for i := range v {
		v[i] = data[i]
	}

	// Make a plot and set its title.
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = title

	// Create a histogram of our values drawn from the data.
	// n := int(1 + 3.3*math.Log(float64(len(data))))
	h, err := plotter.NewHist(v, 30)
	if err != nil {
		panic(err)
	}

	if normalize {
		h.Normalize(1)
	}

	p.Add(h)

	p.Save(4*vg.Inch, 4*vg.Inch, title+".png")

}

func Scatter(x, y []float64, x_name, y_name, title string) {
	// 轉格式
	data := make(plotter.XYs, len(x))
	for i := range data {
		data[i].X = x[i]
		data[i].Y = y[i]
	}

	// Make a plot and set its title.
	p, _ := plot.New()

	p.Title.Text = title
	p.X.Label.Text = x_name
	p.Y.Label.Text = y_name

	scatter, _ := plotter.NewScatter(data)
	p.Add(scatter)

	p.Save(4*vg.Inch, 4*vg.Inch, title+".png")
}

func Line (x, lineData []float64, x_name, y_name, title string) {
	// 轉格式
	data := make(plotter.XYs, len(lineData))
	for i,j := range lineData {
		data[i].X = x[i]
		data[i].Y = j
	}	
	
	p, _ := plot.New()
	
	p.Title.Text = title
	p.X.Label.Text = x_name
	p.Y.Label.Text = y_name
	
	l, _ := plotter.NewLine(data)
	p.Add(l)
	
	p.Save(4*vg.Inch, 4*vg.Inch, title+".png")
}

func Simulation(x, y, delta, r float64, sigma [2]float64, dist *distmv.Normal) (float64, float64) {
	dW := dist.Rand(nil)
	x_1 := math.Exp(math.Log(x) + (r-(math.Pow(sigma[0], 2)/2.))*delta + sigma[0]*math.Sqrt(delta)*dW[0])
	y_1 := math.Exp(math.Log(y) + (r-(math.Pow(sigma[0], 2)/2.))*delta + sigma[0]*math.Sqrt(delta)*dW[1])
	
	return x_1, y_1
}

