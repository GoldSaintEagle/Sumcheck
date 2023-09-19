package MLE

import (
	"fmt"
	"testing"
)

func TestMultilinearExtension(t *testing.T) {
	index := []int64{0, 1, 63}
	val := []int64{1, 5, 7}
	x := []int64{7, 10, 2, 1, 7, 8}
	ret := MultilinearExtension(index, val, x)
	fmt.Printf("%d\n", ret)
}

func TestSumPoly(t *testing.T) {
	coeff := []int64{1, 1, 1, 1, 1, 1, 1}
	x := []int64{1, 2, 3, 4, 5, 6, 7}
	ret := SumPoly(coeff, x, 0, 3)
	fmt.Printf("SumPoly: %d\n", ret)

	c0, c1 := SumPolyExceptX(coeff, x, 0, 3)
	fmt.Printf("SumPolyExceptX: %d %d\n", c0, c1)
}

func TestSumcheck(t *testing.T) {
	coeff := []int64{1, 1, 1, 1, 1, 1}
	x := []int64{1, 2, 3, 4, 5, 6}
	c0, c1 := SumcheckProver(coeff, x, 10)
	fmt.Printf("SumcheckProver: \n\t%v\n\t%v\n", c0, c1)

	fmt.Printf("SumcheckVerifier: \n")
	sum := SumPoly(coeff, x, 0, 10)
	fmt.Printf("\t%v\n", SumcheckVerifier(c0, c1, sum))
}

func BenchmarkSumcheckProver(b *testing.B) {
	size := 20
	coeff := make([]int64, (1<<size)-1)
	x := make([]int64, (1<<size)-1)
	for i := 0; i < (1<<size)-1; i++ {
		coeff[i] = int64(1)
		x[i] = int64(i)
	}
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = SumcheckProver(coeff, x, size)
		}
	})
}

func BenchmarkSumcheckVerifier(b *testing.B) {
	size := 12
	coeff := make([]int64, (1<<size)-1)
	x := make([]int64, (1<<size)-1)
	for i := 0; i < (1<<size)-1; i++ {
		coeff[i] = int64(1)
		x[i] = int64(i)
	}
	c0, c1 := SumcheckProver(coeff, x, size)
	sum := SumPoly(coeff, x, 0, size)

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = SumcheckVerifier(c0, c1, sum)
		}
	})
}
