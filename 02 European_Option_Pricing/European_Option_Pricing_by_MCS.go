// European_Option_Pricing_by_MCS.go
package main

import (
	"fmt"
	"image/color"

	"math"
	r1 "math/rand"

	r2 "golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"

	"gonum.org/v1/gonum/stat"
)

type MCSparams struct {
	steps              int32
	S0, K, r, T, sigma float64
	dist               distuv.Normal
}

func MCS_1(bm MCSparams, num_rep, num_path int32, filename string) {

	reps := make([]float64, num_rep)
	for rep := 0; rep < len(reps); rep++ {
		St := make([]float64, num_path)
		P := make([]float64, num_path)

		// plot
		p, _ := plot.New()

		p.Title.Text = filename
		p.X.Label.Text = "t"
		p.Y.Label.Text = "St"

		for i := range St {
			path := make(plotter.XYs, bm.steps)
			path[0].X = 0
			path[0].Y = bm.S0
			for j := range path {

				// dW := rand.NormFloat64() * math.Sqrt(bm.T/float64(bm.steps))
				dW := math.Sqrt(bm.T/float64(bm.steps)) * bm.dist.Rand()
				dXt := bm.r*path[j].Y*bm.T/float64(bm.steps) + bm.sigma*path[j].Y*dW

				if j < len(path)-1 {
					k := j + 1
					path[k].X = float64(k) * bm.T / float64(bm.steps)
					path[k].Y = path[j].Y + dXt
				} else {
					break
				}
			}

			if i < 10 {
				l, _ := plotter.NewLine(path)
				l.Color = color.RGBA{R: uint8(r1.Intn(255)), G: uint8(r1.Intn(255)), B: uint8(r1.Intn(255))}
				p.Add(l)
			} else if i == 10 {
				p.Save(6*vg.Inch, 3*vg.Inch, filename+".png")
			}

			St[i] = path[bm.steps-1].Y
			P[i] = math.Max(path[bm.steps-1].Y-bm.K, 0) * math.Exp(-bm.r*bm.T)

		}

		// https://godoc.org/github.com/gonum/stat#MeanStdDev
		mean, std := stat.MeanStdDev(P, nil)

		if num_rep == 1 {
			std = std / math.Sqrt(float64(num_path))
		}

		fmt.Println("Repaet :", rep+1)
		fmt.Println("Mean =", mean)
		fmt.Println("Sigma =", std)
		fmt.Println("Confident Interval:[", mean-1.96*std, ",", mean+1.96*std, "]")
	}
}

func main() {
	// rand.Seed(123457)
	seed := 123457
	dist := distuv.Normal{
		Mu:    0,
		Sigma: 1,
		Src:   r2.New(r2.NewSource(uint64(seed))),
	}

	BM := MCSparams{steps: 500, S0: 10520., K: 10500., r: 0.02, T: 20. / 252., sigma: 0.01, dist: dist}
	MCS_1(BM, 1, 10000, "MCS_1")
}
