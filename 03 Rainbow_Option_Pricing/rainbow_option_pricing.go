// rainbow_option_pricing
package main

import (
	"fmt"
	"math"

	"../rainbow_option_pricing/est"

	"github.com/markcheno/go-quote"
	"gonum.org/v1/gonum/stat"
	"github.com/montanaflynn/stats"

	r1 "math/rand"

	r2 "golang.org/x/exp/rand"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat/distmv"

	"image/color"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"

)

type ROP struct {
	rho, r, Par, r_mon, T  float64
	S0, K, auto_out, sigma [2]float64
	// rho 兩資產相關係數, S0 兩資產期初價格, K 兩資產執行價, auto_out 兩資產自動出場價, r 無風險利率, Par 面額, r_mon 月配息率
}

type ModelParams struct {
	dist                  *distmv.Normal
	steps, d_day, pre_day int
	// 天數:steps, 評價日:v_day, 付息日:i_day, 評價日之前的天數:pre_day
}

func MCS(contr ROP, model ModelParams, path_num int) []float64 {

	delta := contr.T / float64(model.steps)
	p_day := float64(model.pre_day) / float64(model.steps)

	// 每月自動出場次數
	var mon int
	auto_out_num := [12]int{}
	mon_interest := [12]float64{}

	// 每條path現值
	PV := make([]float64, path_num)
	
	filename := "Simulation S1"
	// plot
	p, _ := plot.New()

	p.Title.Text = filename
	p.X.Label.Text = "t"
	p.Y.Label.Text = "St"
	

	for path := 0; path < path_num; path++ {

		var n, _n, t, S0_1, S0_2, N_ratio float64
		
		// dW := make([][]float64 , model.steps+model.pre_day)
		// for i := 0; i < model.steps+model.pre_day; i++  {
		// 	dW[i] = make([]float64, 2)
		// 	dW[i] = model.dist.Rand(nil)
		// }
		
		// dW := model.dist.Rand(nil)
		
		// plot
		Path := make(plotter.XYs, model.steps+model.pre_day+1)
		Path[0].X = 0
		Path[0].Y = contr.S0[0]

		
		if model.pre_day != 0 {
			S0_1 = contr.S0[0]
			S0_2 = contr.S0[1]
			for i := 0; i < model.pre_day; i++ {
				S0_1, S0_2 = est.Simulation(S0_1, S0_2, delta, contr.r, contr.sigma, model.dist)
				
				j := i+1
				Path[j].X = 0
				Path[j].Y = S0_1
				
				// S0_1 = math.Exp(math.Log(S0_1) + (contr.r-(math.Pow(contr.sigma[0], 2)/2.))*delta + contr.sigma[0]*math.Sqrt(delta)*dW[i][0])
				// S0_2 = math.Exp(math.Log(S0_2) + (contr.r-(math.Pow(contr.sigma[1], 2)/2.))*delta + contr.sigma[1]*math.Sqrt(delta)*dW[i][1])
			}
		}

		// fmt.Printf("\nS1=%.4f", S0_1)
		// fmt.Printf("\nS2=%.4f", S0_2)
		
		
		cash_flow := [][]float64{[]float64{0, 0}}

		S1 := []float64{S0_1}
		S2 := []float64{S0_2}
		// fmt.Println("S1=", S0_1, "\nS2=", contr.S0[1])

		// dW := make([][]float64, model.steps)
		// for i := 0; i < model.steps; i++ {
		// 	dW[i] = make([]float64, 2)
		// 	dW[i] = model.dist.Rand(nil)
		// }

		for i := 0; i < model.steps; i++ {
			j := i + model.pre_day+1
			St_1, St_2 := est.Simulation(S1[i], S2[i], delta, contr.r, contr.sigma, model.dist)
			
			// PLOT
			Path[j].X = float64(i+1)*delta
			Path[j].Y = St_1

			// St_1 := math.Exp(math.Log(S1[i]) + (contr.r-(math.Pow(contr.sigma[0], 2)/2.))*delta + contr.sigma[0]*math.Sqrt(delta)*dW[i+model.pre_day][0])
			// St_2 := math.Exp(math.Log(S2[i]) + (contr.r-(math.Pow(contr.sigma[1], 2)/2.))*delta + contr.sigma[1]*math.Sqrt(delta)*dW[i+model.pre_day][1])

			// St_1 := S1[i] * math.Exp((contr.r-(math.Pow(contr.sigma[0], 2)/2.))*delta+contr.sigma[0]*math.Sqrt(delta)*dW[0])
			// St_2 := S2[i] * math.Exp((contr.r-(math.Pow(contr.sigma[1], 2)/2.))*delta+contr.sigma[1]*math.Sqrt(delta)*dW[1])

			// St_1 := S1[i] + contr.r*delta + contr.sigma[0]*math.Sqrt(delta)*dW[0]
			// St_2 := S2[i] + contr.r*delta + contr.sigma[1]*math.Sqrt(delta)*dW[1]

			// fmt.Println("S1=", St_1, "\nS2=", St_2)
			S1 = append(S1, St_1)
			S2 = append(S2, St_2)

			// count n
			if i%model.d_day == 0 {
				n = 0.
				_n = 0.
				// month
				mon = i / model.d_day
				// fmt.Println("n0=", n, "\n_n0=", _n)
			} else if St_1 >= contr.K[0] && St_2 >= contr.K[1] {
				n += 1.
			} else {
				_n += 1.
			}

			// fmt.Println("n=", n, "\n_n=", _n)

			// cal interest
			// FIRST month
			if int((i+1)%model.d_day) == 0 {

				t = delta * float64(i+1)
				N_ratio = n / (n + _n)

				if mon == 0 {
					if St_1 >= contr.auto_out[0] && St_2 >= contr.auto_out[1] {
						auto_out_num[mon] += 1
						// fmt.Println("n=", n)
						mon_interest[mon] += contr.Par * contr.r_mon * (1. / float64(path_num))
						cash_flow = append(cash_flow, []float64{t, contr.Par * (1 + contr.r_mon)})
						break

					} else {

						mon_interest[mon] += contr.Par * contr.r_mon * (1. / float64(path_num))
						cash_flow = append(cash_flow, []float64{t, contr.Par * contr.r_mon})
					}

				} else if 0 < mon && mon < 11 {
					if St_1 >= contr.auto_out[0] && St_2 >= contr.auto_out[1] {
						// fmt.Println("n=", n)
						auto_out_num[mon] += 1
						mon_interest[mon] += contr.Par * contr.r_mon * N_ratio * (1. / float64(path_num))
						cash_flow = append(cash_flow, []float64{t, contr.Par * (1 + contr.r_mon*N_ratio)})
						break

					} else {
						mon_interest[mon] += contr.Par * contr.r_mon * N_ratio * (1. / float64(path_num))
						cash_flow = append(cash_flow, []float64{t, contr.Par * contr.r_mon * N_ratio})
					}

				} else {
					auto_out_num[mon] += 1
					mon_interest[mon] += contr.Par * contr.r_mon * N_ratio * (1. / float64(path_num))

					R1 := math.Log(S1[len(S1)-1] / contr.S0[0])
					R2 := math.Log(S2[len(S2)-1] / contr.S0[1])

					if R1 < R2 {
						
						if St_1 >= contr.K[0] {
							cash_flow = append(cash_flow, []float64{t, contr.Par * (1 + contr.r_mon*N_ratio)})
						} else {
							cash_flow = append(cash_flow, []float64{t, contr.Par * (St_1 / contr.K[0] + contr.r_mon*N_ratio)})
						}

					} else {
						
						if St_2 >= contr.K[1] {
							cash_flow = append(cash_flow, []float64{t, contr.Par * (1 + contr.r_mon*N_ratio)})
						} else {
							cash_flow = append(cash_flow, []float64{t, contr.Par * (St_2/contr.K[1] + contr.r_mon*N_ratio)})
						}
						
					}

				}
			}//if


		}//path loop
		
		// plot
		if path < 10 {
			l, _ := plotter.NewLine(Path)
			l.Color = color.RGBA{R: uint8(r1.Intn(255)), G: uint8(r1.Intn(255)), B: uint8(r1.Intn(255))}
			p.Add(l)
		} else if path == 10 {
			p.Save(6*vg.Inch, 3*vg.Inch, filename+".png")
		}


		// fmt.Println("cash_flow:", cash_flow)
		// pv path
		PV[path] = math.Exp(-contr.r*p_day) * est.PV(cash_flow, contr.r)

		// s path
		// fmt.Printf("\nS1=%.4f", S1)
		// fmt.Println("\nlen S1=", len(S1))
		// fmt.Printf("S2=%.4f", S2)
		// fmt.Println("\nlen S2=", len(S2))
		
	}// steps loop

	s := 0
	for i := range auto_out_num {
		s += auto_out_num[i]
	}

	fmt.Println("\nauto_out_num:", auto_out_num, " sum=", s)
	fmt.Println("mon_interest:", mon_interest)

	return PV
}//func

/*
test
t := make([][]float64,1)
fmt.Println(t)
t[0] = []float64{0,5}
fmt.Println(t)
t = append(t, []float64{1,8})
fmt.Println(t)
*/

func main() {
	// get data and estimate params
	asset1, _ := quote.NewQuoteFromYahoo("COP", "2014-01-01", "2015-04-14", quote.Daily, true)
	asset2, _ := quote.NewQuoteFromYahoo("MSFT", "2014-01-01", "2015-04-14", quote.Daily, true)

	asset1_close := asset1.Close
	asset2_close := asset2.Close

	// 計算log return
	asset1_logr := make([]float64, len(asset1_close)-1)
	for i := range asset1_logr {
		asset1_logr[i] = math.Log(asset1_close[i+1] / asset1_close[i])
	}

	asset2_logr := make([]float64, len(asset2_close)-1)
	for i := range asset1_logr {
		asset2_logr[i] = math.Log(asset2_close[i+1] / asset2_close[i])
	}

	// 散佈圖
	est.Scatter(asset1_logr, asset2_logr, "R(COP)", "R(MSFT)", "COP vs MSFT")
	
	// 標準差
	asset1_std := stat.StdDev(asset1_logr, nil) * math.Sqrt(252)
	asset2_std := stat.StdDev(asset2_logr, nil) * math.Sqrt(252)
	sigma := [2]float64{asset1_std, asset2_std}
	fmt.Println("sigma=", sigma)

	// 相關係數
	rho := stat.Correlation(asset1_logr, asset2_logr, nil)

	// create dist
	mu := []float64{0., 0.}
	// rho := assets_corr
	data := []float64{1, rho, rho, 1}
	cov := mat.NewSymDense(2, data)

	seed := 123457
	src := r2.New(r2.NewSource(uint64(seed)))

	normal, _ := distmv.NewNormal(mu, cov, src)
	// rn := normal.Rand(nil)
	// fmt.Println(rn)

	// model params
	model := ModelParams{dist: normal, steps: 360, d_day: 30, pre_day: 15}

	// treasury rate:https://finance.yahoo.com/bonds/
	// treasury_rate, _ := quote.NewQuoteFromYahoo("^IRX", "2015-03-26", "2015-04-14", quote.Daily, true)
	// r := treasury_rate.Close
	// last:r[len(r)-1]

	// fmt.Println("r=", r[len(r)-1]*.01)
	// fmt.Printf("r type = %T", r)

	fmt.Println("std1=", asset1_std, "\nstd2=", asset2_std, "\ncorr=", rho)

	// input params
	p1 := asset1_close[len(asset1_close)-1]
	p2 := asset2_close[len(asset2_close)-1]

	fmt.Println("p1=", p1, "\np2=", p2)

	S0 := [2]float64{p1, p2}
	K := [2]float64{}
	auto_out := [2]float64{}

	for i := range S0 {
		// 執行價為期初價格的85%
		K[i] = S0[i] * 0.85
		// 自動提前出場價為期初價格的100%
		auto_out[i] = S0[i] * 1.
	}

	params := ROP{rho: rho, r: 0.0022, Par: 10000., r_mon: 0.0106, T: 1., S0: S0, K: K, auto_out: auto_out, sigma: sigma}

	path_PV := MCS(params, model, 1000)
	fmt.Println("path_PV=", path_PV, "\nmean PV=", stat.Mean(path_PV, nil))

	est.Hist(path_PV, "PV Hist", false)
	
	// VaR
	alpha := 0.01
	VaR,_ := stats.Percentile(path_PV,1-alpha)
	fmt.Println("VaR=",VaR*(-1))

	// // save cash flow
	// f := make([][]float64, 2)
	// for i := range f {
	// 	f[i] = make([]float64, 2)
	// 	f[i][0] = float64(i) / 5.
	// 	f[i][1] = 100. * r2.Float64()
	// }

	// pv := est.PV(f, 0.02)
	// fmt.Println("cash flow:", f)
	// fmt.Println("PV=", pv)
}
