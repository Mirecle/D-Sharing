// Copyright 2020 ConsenSys AG
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bls24315

import (
	"errors"

	"github.com/consensys/gnark-crypto/ecc/bls24-315/internal/fptower"
)

// GT target group of the pairing
type GT = fptower.E24

type lineEvaluation struct {
	r0 fptower.E4
	r1 fptower.E4
	r2 fptower.E4
}

// Pair calculates the reduced pairing for a set of points
func Pair(P []G1Affine, Q []G2Affine) (GT, error) {
	f, err := MillerLoop(P, Q)
	if err != nil {
		return GT{}, err
	}
	return FinalExponentiation(&f), nil
}

// PairingCheck calculates the reduced pairing for a set of points and returns True if the result is One
func PairingCheck(P []G1Affine, Q []G2Affine) (bool, error) {
	f, err := Pair(P, Q)
	if err != nil {
		return false, err
	}
	var one GT
	one.SetOne()
	return f.Equal(&one), nil
}

// FinalExponentiation computes the final expo x**(p**12-1)(p**4+1)(p**8 - p**4 +1)/r
func FinalExponentiation(z *GT, _z ...*GT) GT {

	var result GT
	result.Set(z)

	for _, e := range _z {
		result.Mul(&result, e)
	}

	// https://eprint.iacr.org/2012/232.pdf, section 7
	var t [9]GT

	// easy part
	t[0].Conjugate(&result)
	result.Inverse(&result)
	t[0].Mul(&t[0], &result)
	result.FrobeniusQuad(&t[0]).
		Mul(&result, &t[0])

	// hard part (up to permutation)
	// Daiki Hayashida and Kenichiro Hayasaka
	// and Tadanori Teruya
	// https://eprint.iacr.org/2020/875.pdf
	// 3*Phi_24(p)/r = (u-1)^2 * (u+p) * (u^2+p^2) * (u^4+p^4-1) + 3
	t[0].CyclotomicSquare(&result)
	t[1].Expt(&result)
	t[2].InverseUnitary(&result)
	t[1].Mul(&t[1], &t[2])
	t[2].Expt(&t[1])
	t[1].InverseUnitary(&t[1])
	t[1].Mul(&t[1], &t[2])
	t[2].Expt(&t[1])
	t[1].Frobenius(&t[1])
	t[1].Mul(&t[1], &t[2])
	result.Mul(&result, &t[0])
	t[0].Expt(&t[1])
	t[2].Expt(&t[0])
	t[0].FrobeniusSquare(&t[1])
	t[2].Mul(&t[0], &t[2])
	t[1].Expt(&t[2])
	t[1].Expt(&t[1])
	t[1].Expt(&t[1])
	t[1].Expt(&t[1])
	t[0].FrobeniusQuad(&t[2])
	t[0].Mul(&t[0], &t[1])
	t[2].InverseUnitary(&t[2])
	t[0].Mul(&t[0], &t[2])
	result.Mul(&result, &t[0])

	return result
}

// MillerLoop Miller loop
func MillerLoop(P []G1Affine, Q []G2Affine) (GT, error) {
	// check input size match
	n := len(P)
	if n == 0 || n != len(Q) {
		return GT{}, errors.New("invalid inputs sizes")
	}

	// filter infinity points
	p := make([]G1Affine, 0, n)
	q := make([]G2Affine, 0, n)

	for k := 0; k < n; k++ {
		if P[k].IsInfinity() || Q[k].IsInfinity() {
			continue
		}
		p = append(p, P[k])
		q = append(q, Q[k])
	}

	n = len(p)

	// projective points for Q
	qProj := make([]g2Proj, n)
	qNeg := make([]G2Affine, n)
	for k := 0; k < n; k++ {
		qProj[k].FromAffine(&q[k])
		qNeg[k].Neg(&q[k])
	}

	var result GT
	result.SetOne()

	var l lineEvaluation

	// i == 31
	for k := 0; k < n; k++ {
		qProj[k].DoubleStep(&l)
		// line evaluation
		l.r0.MulByElement(&l.r0, &p[k].Y)
		l.r2.MulByElement(&l.r2, &p[k].X)
		result.MulBy012(&l.r0, &l.r1, &l.r2)
	}

	for i := len(loopCounter) - 3; i >= 0; i-- {
		result.Square(&result)

		for k := 0; k < n; k++ {
			qProj[k].DoubleStep(&l)
			// line evaluation
			l.r0.MulByElement(&l.r0, &p[k].Y)
			l.r2.MulByElement(&l.r2, &p[k].X)
			result.MulBy012(&l.r0, &l.r1, &l.r2)

			if loopCounter[i] == 1 {
				qProj[k].AddMixedStep(&l, &q[k])
				// line evaluation
				l.r0.MulByElement(&l.r0, &p[k].Y)
				l.r2.MulByElement(&l.r2, &p[k].X)
				result.MulBy012(&l.r0, &l.r1, &l.r2)

			} else if loopCounter[i] == -1 {
				qProj[k].AddMixedStep(&l, &qNeg[k])
				// line evaluation
				l.r0.MulByElement(&l.r0, &p[k].Y)
				l.r2.MulByElement(&l.r2, &p[k].X)
				result.MulBy012(&l.r0, &l.r1, &l.r2)
			}
		}
	}

	result.Conjugate(&result)

	return result, nil
}

// DoubleStep doubles a point in Homogenous projective coordinates, and evaluates the line in Miller loop
// https://eprint.iacr.org/2013/722.pdf (Section 4.3)
func (p *g2Proj) DoubleStep(evaluations *lineEvaluation) {

	// get some Element from our pool
	var t1, A, B, C, D, E, EE, F, G, H, I, J, K fptower.E4
	A.Mul(&p.x, &p.y)
	A.Halve()
	B.Square(&p.y)
	C.Square(&p.z)
	D.Double(&C).
		Add(&D, &C)
	E.Mul(&D, &bTwistCurveCoeff)
	F.Double(&E).
		Add(&F, &E)
	G.Add(&B, &F)
	G.Halve()
	H.Add(&p.y, &p.z).
		Square(&H)
	t1.Add(&B, &C)
	H.Sub(&H, &t1)
	I.Sub(&E, &B)
	J.Square(&p.x)
	EE.Square(&E)
	K.Double(&EE).
		Add(&K, &EE)

	// X, Y, Z
	p.x.Sub(&B, &F).
		Mul(&p.x, &A)
	p.y.Square(&G).
		Sub(&p.y, &K)
	p.z.Mul(&B, &H)

	// Line evaluation
	evaluations.r0.Neg(&H)
	evaluations.r1.Set(&I)
	evaluations.r2.Double(&J).
		Add(&evaluations.r2, &J)
}

// AddMixedStep point addition in Mixed Homogenous projective and Affine coordinates
// https://eprint.iacr.org/2013/722.pdf (Section 4.3)
func (p *g2Proj) AddMixedStep(evaluations *lineEvaluation, a *G2Affine) {

	// get some Element from our pool
	var Y2Z1, X2Z1, O, L, C, D, E, F, G, H, t0, t1, t2, J fptower.E4
	Y2Z1.Mul(&a.Y, &p.z)
	O.Sub(&p.y, &Y2Z1)
	X2Z1.Mul(&a.X, &p.z)
	L.Sub(&p.x, &X2Z1)
	C.Square(&O)
	D.Square(&L)
	E.Mul(&L, &D)
	F.Mul(&p.z, &C)
	G.Mul(&p.x, &D)
	t0.Double(&G)
	H.Add(&E, &F).
		Sub(&H, &t0)
	t1.Mul(&p.y, &E)

	// X, Y, Z
	p.x.Mul(&L, &H)
	p.y.Sub(&G, &H).
		Mul(&p.y, &O).
		Sub(&p.y, &t1)
	p.z.Mul(&E, &p.z)

	t2.Mul(&L, &a.Y)
	J.Mul(&a.X, &O).
		Sub(&J, &t2)

	// Line evaluation
	evaluations.r0.Set(&L)
	evaluations.r1.Set(&J)
	evaluations.r2.Neg(&O)
}
