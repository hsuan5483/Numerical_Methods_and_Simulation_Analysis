// main
package main

import (
	"fmt"
	"math"

	"../MidtermExam/est"

	"github.com/markcheno/go-quote"
	"gonum.org/v1/gonum/stat"
	"github.com/montanaflynn/stats"

	r2 "golang.org/x/exp/rand"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat/distmv"

)

type ROP struct {
	rho, r, Par, r_mon  float64
	S0, K, lower_bond, sigma [2]float64
	// rho 兩資產相關係數, S0 兩資產期初價格, K 兩資產執行價, lower_bond 兩資產自動出場價, r 無風險利率, Par 面額, r_mon 月配息率
}

type ModelParams struct {
	dist                  *distmv.Normal
	steps, d_day, pre_day int
	auto_out bool
	// 天數:steps, 評價日:v_day, 付息日:i_day, 評價日之前的天數:pre_day
}

func MCS(contr ROP, model ModelParams, path_num int) ([]float64,float64,float64,float64) {

	delta := 1. / float64(model.steps)
	p_day := float64(model.pre_day-1) / float64(model.steps)

	// 每月自動出場次數
	var mon int
	auto_out_num := [11]int{}
	
	// 每條path現值
	PV := make([]float64, path_num)
	
	// auto out PV
	auto_out_PV := 0.
	n_auto_out_PV := 0.
	
	for path := 0; path < path_num; path++ {

		var n, _n, t, S0_1, S0_2, N_ratio float64
		
		if model.pre_day != 0 {
			S0_1 = contr.S0[0]
			S0_2 = contr.S0[1]
			for i := 0; i < model.pre_day; i++ {
				S0_1, S0_2 = est.Simulation(S0_1, S0_2, delta, contr.r, contr.sigma, model.dist)
				// dW := model.dist.Rand(nil)
				// S0_1 = S0_1 * math.Exp((contr.r-(math.Pow(contr.sigma[0], 2)/2.))*delta+contr.sigma[0]*math.Sqrt(delta)*dW[0])
				// S0_2 = S0_2 * math.Exp((contr.r-(math.Pow(contr.sigma[1], 2)/2.))*delta+contr.sigma[1]*math.Sqrt(delta)*dW[1])

			}
		}

		// 利息現金流
		cash_flow_i := [][]float64{[]float64{0, 0}}
		// 最終金額現金流
		cash_flow_p := [][]float64{[]float64{0, 0}}
		
		S1 := []float64{S0_1}
		S2 := []float64{S0_2}

		for i := 0; i < model.steps; i++ {
			St_1, St_2 := est.Simulation(S1[i], S2[i], delta, contr.r, contr.sigma, model.dist)
			// dW := model.dist.Rand(nil)
			// St_1 := S1[i] * math.Exp((contr.r-(math.Pow(contr.sigma[0], 2)/2.))*delta+contr.sigma[0]*math.Sqrt(delta)*dW[0])
			// St_2 := S2[i] * math.Exp((contr.r-(math.Pow(contr.sigma[1], 2)/2.))*delta+contr.sigma[1]*math.Sqrt(delta)*dW[1])
			
			S1 = append(S1, St_1)
			S2 = append(S2, St_2)
			
			// 期初判斷月份
			if i%model.d_day == 0 {
				n = 0.
				_n = 0.
				// month
				mon = i / model.d_day
			}
			
			// 計算計息天數
			if St_1 >= contr.lower_bond[0] && St_2 >= contr.lower_bond[1] {
				n += 1.
			} else {
				_n += 1.
			}

			// calulate cash flow
			// 評價日
			if int((i+1)%model.d_day) == 0 {

				t = delta * float64(i+1)
				N_ratio = n / (n + _n)
				// 第一個月
				if mon == 0 {
					
					// 利息
					cash_flow_i = append(cash_flow_i, []float64{t, contr.Par * contr.r_mon})
					// auto out
					if St_1 >= contr.S0[0] && St_2 >= contr.S0[1] { 
						auto_out_num[mon] += 1
						cash_flow_p = append(cash_flow_p, []float64{t, contr.Par})
						auto_out_PV += math.Exp(-contr.r*p_day) * (est.PV(cash_flow_i, contr.r) + est.PV(cash_flow_p, contr.r))
						break
					}
					
				// 第2~11個月
				} else if 0 < mon && mon < 11 {
					
					// 利息
					cash_flow_i = append(cash_flow_i, []float64{t, contr.Par * contr.r_mon * N_ratio})
					// auto out
					if St_1 >= contr.S0[0] && St_2 >= contr.S0[1] { 
						auto_out_num[mon] += 1
						cash_flow_p = append(cash_flow_p, []float64{t, contr.Par})
						auto_out_PV += math.Exp(-contr.r*p_day) * (est.PV(cash_flow_i, contr.r) + est.PV(cash_flow_p, contr.r))
						break
					}
					
				// 最終評價日
				} else {
					// 計算兩資產績效
					R1 := (S1[len(S1)-1] / contr.S0[0]) - 1
					R2 := (S2[len(S2)-1] / contr.S0[1]) - 1
					
					// 決定最終給付
					if St_1 >= contr.K[0] && St_2 >= contr.K[1]{
						cash_flow_i = append(cash_flow_i, []float64{t, contr.Par * contr.r_mon * N_ratio})
						cash_flow_p = append(cash_flow_p, []float64{t, contr.Par})
					} else if R1 < R2 {
						
						cash_flow_p = append(cash_flow_p, []float64{t, contr.Par * St_1 / contr.K[0]})
						
					} else if R1 > R2 {
						
						cash_flow_p = append(cash_flow_p, []float64{t, contr.Par * St_2 / contr.K[1]})						
					}
					
					n_auto_out_PV += math.Exp(-contr.r*p_day) * (est.PV(cash_flow_i, contr.r) + est.PV(cash_flow_p, contr.r))

				}
			}//if
		
		if path == 0 {
			times := []float64{}
			for t := range S1 {
				times = append(times, float64(t))
			}
			
			est.Line(times, S1,"Time","Price","Price Simulation - GDX")
			est.Line(times, S2,"Time","Price","Price Simulation - XME")

		}

		}// steps loop
		
		// pv path
		PV[path] = math.Exp(-contr.r*p_day) * (est.PV(cash_flow_i, contr.r) + est.PV(cash_flow_p, contr.r))
		
		
	}// path loop
	
	// 計算自動出場機率及畫折線圖
	x := []float64{}
	auto_out_p := []float64{}
	p_auto_out := 0.
	for i,j := range auto_out_num {
		x = append(x, float64(i+1))
		auto_out_p = append(auto_out_p, float64(j)/float64(path_num))
		p_auto_out += float64(j)/float64(path_num)
	}
	
	est.Line(x, auto_out_p,"Month","P","auto out probability")
	
	return PV,auto_out_PV/(p_auto_out*float64(path_num)),n_auto_out_PV/((1-p_auto_out)*float64(path_num)),p_auto_out
	
}//func

func main() {
	// get data and estimate params
	// VanEck Vectors Gold Miners ETF (GDX)
	// SPDR S&P Metals and Mining ETF (XME)
	asset1, _ := quote.NewQuoteFromYahoo("GDX", "2017-05-15", "2018-05-15", quote.Daily, true)
	asset2, _ := quote.NewQuoteFromYahoo("XME", "2017-05-15", "2018-05-15", quote.Daily, true)

	asset1_close := asset1.Close
	asset2_close := asset2.Close

	// 計算log return
	x1 := []float64{}
	asset1_logr := make([]float64, len(asset1_close)-1)
	for i := range asset1_logr {
		x1 = append(x1, float64(i+1))
		asset1_logr[i] = math.Log(asset1_close[i+1] / asset1_close[i])
	}
	
	x2 := []float64{}
	asset2_logr := make([]float64, len(asset2_close)-1)
	for i := range asset2_logr {
		x2 = append(x2, float64(i+1))
		asset2_logr[i] = math.Log(asset2_close[i+1] / asset2_close[i])
	}

	// 散佈圖
	est.Scatter(asset1_logr, asset2_logr, "R(GDX)", "R(XME)", "GDX vs XME")
	
	// 畫log return
	est.Line(x1,asset1_logr,"t","log retuen","GDX Log Return")
	est.Line(x1,asset2_logr,"t","log retuen","XME Log Return")

	// 年化標準差
	asset1_std := stat.StdDev(asset1_logr, nil) * math.Sqrt(252)
	asset2_std := stat.StdDev(asset2_logr, nil) * math.Sqrt(252)
	sigma := [2]float64{asset1_std, asset2_std}

	// 相關係數
	rho := stat.Correlation(asset1_logr, asset2_logr, nil)

	// create dist
	mu := []float64{0., 0.}
	data := []float64{1, rho, rho, 1}
	cov := mat.NewSymDense(2, data)

	seed := 123457
	src := r2.New(r2.NewSource(uint64(seed)))

	normal, _ := distmv.NewNormal(mu, cov, src)

	// 設定模型參數
	model := ModelParams{dist: normal, steps: 240, d_day: 20, pre_day: 5, auto_out: true}

	fmt.Println("參數估計：")
	fmt.Println("std1=", asset1_std, "\nstd2=", asset2_std, "\ncorr=", rho)

	// input params
	P1, _ := quote.NewQuoteFromYahoo("GDX", "2018-05-01", "2018-05-15", quote.Daily, false)
	P2, _ := quote.NewQuoteFromYahoo("XME", "2018-05-01", "2018-05-15", quote.Daily, false)
	
	p1 := P1.Close[len(P1.Close)-1]
	p2 := P2.Close[len(P2.Close)-1]

	fmt.Println("p1=", p1, "\np2=", p2)
	fmt.Println("===============================")
	
	// Par value
	Par := 10000.
	
	// 交易日收盤價
	S0 := [2]float64{p1, p2}
	// 轉換價
	K := [2]float64{}
	lower_bond := [2]float64{}

	for i := range S0 {
		// 執行價為期初價格的87%
		K[i] = S0[i] * 0.87
		// 自動提前出場價為期初價格的87%
		lower_bond[i] = S0[i] * 0.87
	}
	
	// 設定合約參數值
	params := ROP{rho: rho, r: 0.0228, Par: Par, r_mon: .07/12., S0: S0, K: K, lower_bond: lower_bond, sigma: sigma}
	
	// 模擬
	path_PV, auto_out_PV, n_auto_out_PV, p_auto_out := MCS(params, model, 10000)
	
	fmt.Println("模擬結果:")
	fmt.Println("(a)")
	fmt.Printf("mean PV=%.4f", stat.Mean(path_PV, nil))
	std := stat.StdDev(path_PV, nil)
	stdErr := stat.StdErr(std, float64(len(path_PV)))
	fmt.Printf("\nstd PV=%.4f", stdErr)
	
	fmt.Printf("\n(b)\n提前出場機率=%.4f", p_auto_out)
	
	// 保本機率
	p_keep := 0.
	for _ , i:= range path_PV {
		if i > Par {
			p_keep += 1/float64(len(path_PV))
		}
	}
	fmt.Printf("\n(c)\np_keep=%.4f",p_keep)
	
	// VaR
	fmt.Println("\n(d)")
	VaR_1 , _ := stats.Percentile(path_PV,.99)
	fmt.Printf("VaR(0.01)=%.4f",VaR_1)
	
	fmt.Printf("\n(e)\n提前出場價值=%.4f",auto_out_PV)
	fmt.Printf("\n(f)\n未提前出場價值=%.4f",n_auto_out_PV)
	
	fmt.Printf("\n(g)\nE(PV)=%.4f\n", auto_out_PV * p_auto_out + n_auto_out_PV * (1. - p_auto_out))

	est.Hist(path_PV, "PV Hist", false)
	
}
