package MLE

import (
	"crypto/md5"
	"encoding/binary"
	"math/big"
)

func fiatShamir(msg []byte) int64 {
	// high collision function, may consider better mapping
	hashed := md5.Sum(msg)
	var b = hashed[:8]
	challenge := int64(binary.LittleEndian.Uint64(b))
	//	challenge = 1
	return challenge
}

// index: the non-zero locations of the polynomial
// val: the value of the polynomial at index
// x: the evaluation point
func MultilinearExtension(index, val, x []int64) int64 {
	ret := int64(0)
	for i := 0; i < len(index); i++ {
		prd := int64(1)
		for j := 0; j < len(x); j++ {
			ei := (index[i] & (1 << j)) >> j
			prd *= x[j]*ei + (1-x[j])*(1-ei)
		}
		ret += val[i] * prd
	}
	return ret
}

// x: the multivatiates in binary form
// E.g., x0 + x0*x1 + x1*x2 = [001, 011, 110] = [1, 3, 6]
func SumPoly(coeff, x []int64, deglow, deghigh int) int64 {
	totalx := int64(1)
	ret := int64(0)
	for i := 0; i < deghigh-deglow; i++ {
		totalx <<= 1
		totalx |= 1
	}
	for i := 0; i < deglow; i++ {
		totalx <<= 1
	}

	for i := 0; i < len(x); i++ {
		ones := totalx - x[i]
		times := int64(1)
		for j := 0; j < deghigh; j++ {
			if (ones&(1<<j))>>j == 1 {
				times <<= 1
			}
		}
		times *= coeff[i]
		ret += times
	}
	return ret
}

// E.g., evaluate except x3 from (x0, x1, x2, x3, x4), set ex = 4 (#4 element).
func SumPolyExceptX(coeff, x []int64, ex, deg int) (int64, int64) {
	x0 := make([]int64, len(x))
	x1 := make([]int64, len(x))
	coeff0 := make([]int64, len(x))
	coeff1 := make([]int64, len(x))
	for i := 0; i < len(x); i++ {
		if x[i]&(1<<ex) == 0 {
			coeff0[i] = coeff[i]
			x0[i] = x[i]
		} else {
			coeff1[i] = coeff[i]
			x1[i] = x[i]
		}
	}
	c0 := SumPoly(coeff0, x0, ex+1, deg)
	c1 := SumPoly(coeff1, x1, ex, deg)

	return c0, c1
}

func EvalPoly(coeff, x []int64, xloc int, r int64) ([]int64, []int64) {
	retc := make([]int64, len(x))
	retx := make([]int64, len(x))
	for i := 0; i < len(x); i++ {
		if x[i]&(1<<xloc) != 0 {
			retc[i] = coeff[i] * r
			retx[i] = x[i] & ^(1 << xloc)
		} else {
			retc[i] = coeff[i]
			retx[i] = x[i]
		}
	}
	return retc, retx
}

func sumcheckOneRound(coeff, x []int64, xloc, deg int) (int64, int64, []int64, []int64) {
	ci0, ci1 := SumPolyExceptX(coeff, x, xloc, deg)
	big0 := new(big.Int)
	big0.SetInt64(ci0)
	big1 := new(big.Int)
	big1.SetInt64(ci1)
	pu := append(big0.Bytes(), big1.Bytes()...)
	r := fiatShamir(pu)
	newc, newx := EvalPoly(coeff, x, xloc, r)
	return ci0, ci1, newc, newx
}

func SumcheckProver(coeff, x []int64, deg int) ([]int64, []int64) {
	c0 := make([]int64, deg)
	c1 := make([]int64, deg)

	cpoly := coeff
	xpoly := x
	for i := 0; i < deg; i++ {
		c0[i], c1[i], cpoly, xpoly = sumcheckOneRound(cpoly, xpoly, i, deg)
	}

	return c0, c1
}

func SumcheckVerifier(c0, c1 []int64, sum int64) bool {
	for i := 0; i < len(c0); i++ {
		if 2*c0[i]+c1[i] != sum {
			//			fmt.Printf("\t2 * %v + %v neq %v \n", c0[i], c1[i], sum)
			return false
		}
		//		fmt.Printf("\t2 * %v + %v = %v \n", c0[i], c1[i], sum)
		big0 := new(big.Int)
		big0.SetInt64(c0[i])
		big1 := new(big.Int)
		big1.SetInt64(c1[i])
		pu := append(big0.Bytes(), big1.Bytes()...)
		r := fiatShamir(pu)
		sum = c0[i] + r*c1[i]
	}
	return true
}
