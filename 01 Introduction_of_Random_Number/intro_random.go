// Go's `math/rand` package provides
// [pseudorandom number](http://en.wikipedia.org/wiki/Pseudorandom_number_generator)
// generation.

package main

import (
	"fmt"
	"math/rand"
	"time"

	// "gonum.org/v1/gonum/stat"
	RS "../rand/scatter"
)

func main() {

	// keyword[golang math rand] -> 第二個 https://golang.org/src/math/rand/rand.go

	// 產生整數型態的亂數：`rand.Intn(n)`、`rand.Int31n(n)`...
	// rand.Intn(n)：returns a random `int` between 0 ~ n-1.
	fmt.Println("產生整數型態亂數")
	a1 := rand.Intn(100)
	a2 := rand.Int31()

	fmt.Println(a1, " , ", a2, "\n")

	// 產生浮點數型態的亂數：`rand.Float64()`、`rand.Float32()`...
	// rand.Float64()：returns a `float64` `f`,`0.0 <= f < 1.0`.
	fmt.Println("產生浮點數型態亂數")
	b1 := rand.Float32()
	b2 := rand.Float64()
	fmt.Println(b1, " , ", b2, "\n")

	// 刮號裡可以加數字? https://play.golang.org/p/sqNzxdrN2nQ

	// 如何產生不同範圍的亂數?
	// example：`5.0 <= f' < 10.0`.
	fmt.Println("如何產生不同範圍的亂數?")
	b3 := rand.Float64()*5 + 5
	fmt.Println(b3, "\n")

	// 產生標準常態的亂數：`rand.NormFloat64()`
	fmt.Println("產生標準常態的亂數")
	c1 := rand.NormFloat64()
	fmt.Println(c1)

	// example：產生10個N(10 , 4)的亂數
	// N(10 , 4) -> rand.NormFloat64() * Std + Mean = rand.NormFloat64() * 2 + 10
	fmt.Println("\nexample：產生10個X~N(10 , 4)的亂數")
	type S []float64
	N := 10

	c2 := make(S, N)
	for i := range c2 {
		c2[i] = rand.NormFloat64()*2 + 10
	}

	fmt.Println(c2, "\n")

	// 產生指數分布Exp(1)的隨機亂數：`rand.ExpFloat64()`
	fmt.Println("\n產生指數分布Exp(1)的隨機亂數")
	d1 := rand.ExpFloat64()
	fmt.Println(d1)

	// 產生X~Exp(2)的亂數
	fmt.Println("\nexample：產生X~Exp(2)的亂數")
	d2 := rand.ExpFloat64() / 2
	fmt.Println(d2, "\n")

	// rand.Perm(n)：把0~n-1打亂 (n只能是整數)
	e := rand.Perm(10)
	fmt.Println(e)

	// rand.shuffle()：把Array或Slice打亂
	sli := []float64{1.3, 2.8, 4.2, 2.5, 3.4, 7.1}
	fmt.Println(sli)
	rand.Shuffle(len(sli), func(i, j int) { sli[i], sli[j] = sli[j], sli[i] })
	fmt.Println(sli)

	// 利用時間產生亂數
	fmt.Println("\n利用時間產生亂數")
	// 1：NewSource and New
	fmt.Println("1：NewSource and New")

	// step1：利用NewSource(YOUR_SEED)產生一個亂數指標序列
	//    註：YOUR_SEED的格式必須為int64(時間或整數)
	// https://play.golang.org/p/8vn_exm4u_U

	fmt.Println("step1：")
	t := time.Now().UnixNano()
	s1 := rand.NewSource(t)
	fmt.Println(t)

	// step2：用New函數從NewSource產生新的亂數
	fmt.Println("step2：")
	r1 := rand.New(s1)
	fmt.Println(r1)

	// step3：給定亂數的型態(ex：int、float...)
	fmt.Println("step3：")
	fmt.Println("Int：", r1.Int())
	fmt.Println("Float：", r1.Float64())

	// example：從source 's1'中產生10個介於0~99的整數亂數
	fmt.Println("\nexample：從source 's1'中產生10個介於0~99的整數亂數")
	rands := make([]int, 10)
	for i := range rands {
		rands[i] = r1.Intn(100)
	}
	fmt.Println(rands)

	// 相同的亂數種子產生相同的亂數序列
	fmt.Println("\n相同的亂數種子產生相同的亂數序列")
	r2 := rand.New(rand.NewSource(42))

	rands2 := make([]float64, 3)
	for i := range rands2 {
		rands2[i] = r2.Float64()
	}

	fmt.Println("rands2：", rands2)

	r3 := rand.New(rand.NewSource(42))

	rands3 := make([]float64, 3)
	for i := range rands3 {
		rands3[i] = r3.Float64()
	}

	fmt.Println("rands3：", rands3)

	// 2：給定全域的亂數種子
	fmt.Println("\n2：給定全域的亂數種子")
	rand.Seed(123457)
	fmt.Println(rand.Intn(100), ",", rand.Intn(100)) // 77 , 52

	// Plot Scatter
	// 要產生多少亂數? https://play.golang.org/p/VKPkmY6QG-T
	fmt.Println("\nPlot Scatter")
	RS.UFRanNums(100, 100000000, "UnfixedSeed")
	RS.FRanNums(100, 100000000, time.Now().UnixNano(), "FixedSeed(time)")
	RS.FRanNums(100, 100000000, 123457, "FixedSeed")

}
