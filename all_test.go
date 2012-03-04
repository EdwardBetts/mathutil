// Copyright (c) 2011 CZ.NIC z.s.p.o. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// blame: jnml, labs.nic.cz

package mathutil

import (
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"runtime"
	"sort"
	"testing"
)

func r32() *FC32 {
	r, err := NewFC32(math.MinInt32, math.MaxInt32, true)
	if err != nil {
		panic(err)
	}

	return r
}

var (
	r64lo = big.NewInt(math.MinInt64)
	r64hi = big.NewInt(math.MaxInt64)
)

func r64() *FCBig {
	r, err := NewFCBig(r64lo, r64hi, true)
	if err != nil {
		panic(err)
	}

	return r
}

func benchmark1eN(b *testing.B, r *FC32) {
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		r.Next()
	}
}

func BenchmarkFC1e3(b *testing.B) {
	b.StopTimer()
	r, _ := NewFC32(0, 1e3, false)
	benchmark1eN(b, r)
}

func BenchmarkFC1e6(b *testing.B) {
	b.StopTimer()
	r, _ := NewFC32(0, 1e6, false)
	benchmark1eN(b, r)
}

func BenchmarkFC1e9(b *testing.B) {
	b.StopTimer()
	r, _ := NewFC32(0, 1e9, false)
	benchmark1eN(b, r)
}

func Test0(t *testing.T) {
	const N = 10000
	for n := 1; n < N; n++ {
		lo, hi := 0, n-1
		period := int64(hi) - int64(lo) + 1
		r, err := NewFC32(lo, hi, false)
		if err != nil {
			t.Fatal(err)
		}
		if r.Cycle()-period > period {
			t.Fatalf("Cycle exceeds 2 * period")
		}
	}
	for n := 1; n < N; n++ {
		lo, hi := 0, n-1
		period := int64(hi) - int64(lo) + 1
		r, err := NewFC32(lo, hi, true)
		if err != nil {
			t.Fatal(err)
		}
		if r.Cycle()-2*period > period {
			t.Fatalf("Cycle exceeds 3 * period")
		}
	}
}

func Test1(t *testing.T) {
	const (
		N = 360
		S = 3
	)
	for hq := 0; hq <= 1; hq++ {
		for n := 1; n < N; n++ {
			for seed := 0; seed < S; seed++ {
				lo, hi := -n, 2*n
				period := int64(hi - lo + 1)
				r, err := NewFC32(lo, hi, hq == 1)
				if err != nil {
					t.Fatal(err)
				}
				r.Seed(int64(seed))
				m := map[int]bool{}
				v := make([]int, period, period)
				p := make([]int64, period, period)
				for i := lo; i <= hi; i++ {
					x := r.Next()
					p[i-lo] = r.Pos()
					if x < lo || x > hi {
						t.Fatal("t1.0")
					}
					if m[x] {
						t.Fatal("t1.1")
					}
					m[x] = true
					v[i-lo] = x
				}
				for i := lo; i <= hi; i++ {
					x := r.Next()
					if x < lo || x > hi {
						t.Fatal("t1.2")
					}
					if !m[x] {
						t.Fatal("t1.3")
					}
					if x != v[i-lo] {
						t.Fatal("t1.4")
					}
					if r.Pos() != p[i-lo] {
						t.Fatal("t1.5")
					}
					m[x] = false
				}
				for i := lo; i <= hi; i++ {
					r.Seek(p[i-lo] + 1)
					x := r.Prev()
					if x < lo || x > hi {
						t.Fatal("t1.6")
					}
					if x != v[i-lo] {
						t.Fatal("t1.7")
					}
				}
			}
		}
	}
}

func Test2(t *testing.T) {
	const (
		N = 370
		S = 3
	)
	for hq := 0; hq <= 1; hq++ {
		for n := 1; n < N; n++ {
			for seed := 0; seed < S; seed++ {
				lo, hi := -n, 2*n
				period := int64(hi - lo + 1)
				r, err := NewFC32(lo, hi, hq == 1)
				if err != nil {
					t.Fatal(err)
				}
				r.Seed(int64(seed))
				m := map[int]bool{}
				v := make([]int, period, period)
				p := make([]int64, period, period)
				for i := lo; i <= hi; i++ {
					x := r.Prev()
					p[i-lo] = r.Pos()
					if x < lo || x > hi {
						t.Fatal("t2.0")
					}
					if m[x] {
						t.Fatal("t2.1")
					}
					m[x] = true
					v[i-lo] = x
				}
				for i := lo; i <= hi; i++ {
					x := r.Prev()
					if x < lo || x > hi {
						t.Fatal("t2.2")
					}
					if !m[x] {
						t.Fatal("t2.3")
					}
					if x != v[i-lo] {
						t.Fatal("t2.4")
					}
					if r.Pos() != p[i-lo] {
						t.Fatal("t2.5")
					}
					m[x] = false
				}
				for i := lo; i <= hi; i++ {
					s := p[i-lo] - 1
					if s < 0 {
						s = r.Cycle() - 1
					}
					r.Seek(s)
					x := r.Next()
					if x < lo || x > hi {
						t.Fatal("t2.6")
					}
					if x != v[i-lo] {
						t.Fatal("t2.7")
					}
				}
			}
		}
	}
}

func benchmarkBig1eN(b *testing.B, r *FCBig) {
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		r.Next()
	}
}

func BenchmarkFCBig1e3(b *testing.B) {
	b.StopTimer()
	hi := big.NewInt(0).SetInt64(1e3)
	r, _ := NewFCBig(big0, hi, false)
	benchmarkBig1eN(b, r)
}

func BenchmarkFCBig1e6(b *testing.B) {
	b.StopTimer()
	hi := big.NewInt(0).SetInt64(1e6)
	r, _ := NewFCBig(big0, hi, false)
	benchmarkBig1eN(b, r)
}

func BenchmarkFCBig1e9(b *testing.B) {
	b.StopTimer()
	hi := big.NewInt(0).SetInt64(1e9)
	r, _ := NewFCBig(big0, hi, false)
	benchmarkBig1eN(b, r)
}

func BenchmarkFCBig1e12(b *testing.B) {
	b.StopTimer()
	hi := big.NewInt(0).SetInt64(1e12)
	r, _ := NewFCBig(big0, hi, false)
	benchmarkBig1eN(b, r)
}

func BenchmarkFCBig1e15(b *testing.B) {
	b.StopTimer()
	hi := big.NewInt(0).SetInt64(1e15)
	r, _ := NewFCBig(big0, hi, false)
	benchmarkBig1eN(b, r)
}

func BenchmarkFCBig1e18(b *testing.B) {
	b.StopTimer()
	hi := big.NewInt(0).SetInt64(1e18)
	r, _ := NewFCBig(big0, hi, false)
	benchmarkBig1eN(b, r)
}

var (
	big0 = big.NewInt(0)
	big1 = big.NewInt(1)
)

func TestBig0(t *testing.T) {
	const N = 7400
	lo := big.NewInt(0)
	hi := big.NewInt(0)
	period := big.NewInt(0)
	c := big.NewInt(0)
	for n := int64(1); n < N; n++ {
		hi.SetInt64(n - 1)
		period.Set(hi)
		period.Sub(period, lo)
		period.Add(period, big1)
		r, err := NewFCBig(lo, hi, false)
		if err != nil {
			t.Fatal(err)
		}
		if r.cycle.Cmp(period) < 0 {
			t.Fatalf("Period exceeds cycle")
		}
		c.Set(r.Cycle())
		c.Sub(c, period)
		if c.Cmp(period) > 0 {
			t.Fatalf("Cycle exceeds 2 * period")
		}
	}
	for n := int64(1); n < N; n++ {
		hi.SetInt64(n - 1)
		period.Set(hi)
		period.Sub(period, lo)
		period.Add(period, big1)
		r, err := NewFCBig(lo, hi, true)
		if err != nil {
			t.Fatal(err)
		}
		if r.cycle.Cmp(period) < 0 {
			t.Fatalf("Period exceeds cycle")
		}
		c.Set(r.Cycle())
		c.Sub(c, period)
		c.Sub(c, period)
		if c.Cmp(period) > 0 {
			t.Fatalf("Cycle exceeds 3 * period")
		}
	}
}

func TestBig1(t *testing.T) {
	const (
		N = 120
		S = 3
	)
	lo := big.NewInt(0)
	hi := big.NewInt(0)
	seek := big.NewInt(0)
	for hq := 0; hq <= 1; hq++ {
		for n := int64(1); n < N; n++ {
			for seed := 0; seed < S; seed++ {
				lo64 := -n
				hi64 := 2 * n
				lo.SetInt64(lo64)
				hi.SetInt64(hi64)
				period := hi64 - lo64 + 1
				r, err := NewFCBig(lo, hi, hq == 1)
				if err != nil {
					t.Fatal(err)
				}
				r.Seed(int64(seed))
				m := map[int64]bool{}
				v := make([]int64, period, period)
				p := make([]int64, period, period)
				for i := lo64; i <= hi64; i++ {
					x := r.Next().Int64()
					p[i-lo64] = r.Pos().Int64()
					if x < lo64 || x > hi64 {
						t.Fatal("tb1.0")
					}
					if m[x] {
						t.Fatal("tb1.1")
					}
					m[x] = true
					v[i-lo64] = x
				}
				for i := lo64; i <= hi64; i++ {
					x := r.Next().Int64()
					if x < lo64 || x > hi64 {
						t.Fatal("tb1.2")
					}
					if !m[x] {
						t.Fatal("tb1.3")
					}
					if x != v[i-lo64] {
						t.Fatal("tb1.4")
					}
					if r.Pos().Int64() != p[i-lo64] {
						t.Fatal("tb1.5")
					}
					m[x] = false
				}
				for i := lo64; i <= hi64; i++ {
					r.Seek(seek.SetInt64(p[i-lo64] + 1))
					x := r.Prev().Int64()
					if x < lo64 || x > hi64 {
						t.Fatal("tb1.6")
					}
					if x != v[i-lo64] {
						t.Fatal("tb1.7")
					}
				}
			}
		}
	}
}

func TestBig2(t *testing.T) {
	const (
		N = 120
		S = 3
	)
	lo := big.NewInt(0)
	hi := big.NewInt(0)
	seek := big.NewInt(0)
	for hq := 0; hq <= 1; hq++ {
		for n := int64(1); n < N; n++ {
			for seed := 0; seed < S; seed++ {
				lo64, hi64 := -n, 2*n
				lo.SetInt64(lo64)
				hi.SetInt64(hi64)
				period := hi64 - lo64 + 1
				r, err := NewFCBig(lo, hi, hq == 1)
				if err != nil {
					t.Fatal(err)
				}
				r.Seed(int64(seed))
				m := map[int64]bool{}
				v := make([]int64, period, period)
				p := make([]int64, period, period)
				for i := lo64; i <= hi64; i++ {
					x := r.Prev().Int64()
					p[i-lo64] = r.Pos().Int64()
					if x < lo64 || x > hi64 {
						t.Fatal("tb2.0")
					}
					if m[x] {
						t.Fatal("tb2.1")
					}
					m[x] = true
					v[i-lo64] = x
				}
				for i := lo64; i <= hi64; i++ {
					x := r.Prev().Int64()
					if x < lo64 || x > hi64 {
						t.Fatal("tb2.2")
					}
					if !m[x] {
						t.Fatal("tb2.3")
					}
					if x != v[i-lo64] {
						t.Fatal("tb2.4")
					}
					if r.Pos().Int64() != p[i-lo64] {
						t.Fatal("tb2.5")
					}
					m[x] = false
				}
				for i := lo64; i <= hi64; i++ {
					s := p[i-lo64] - 1
					if s < 0 {
						s = r.Cycle().Int64() - 1
					}
					r.Seek(seek.SetInt64(s))
					x := r.Next().Int64()
					if x < lo64 || x > hi64 {
						t.Fatal("tb2.6")
					}
					if x != v[i-lo64] {
						t.Fatal("tb2.7")
					}
				}
			}
		}
	}
}

func TestPermutations(t *testing.T) {
	data := sort.IntSlice{3, 2, 1}
	check := [][]int{
		{1, 2, 3},
		{1, 3, 2},
		{2, 1, 3},
		{2, 3, 1},
		{3, 1, 2},
		{3, 2, 1},
	}
	i := 0
	for PermutationFirst(data); ; i++ {
		if i >= len(check) {
			t.Fatalf("too much permutations generated: %d > %d", i+1, len(check))
		}

		for j, v := range check[i] {
			got := data[j]
			if got != v {
				t.Fatalf("permutation %d:\ndata: %v\ncheck: %v\nexpected data[%d] == %d, got %d", i, data, check[i], j, v, got)
			}
		}

		if !PermutationNext(data) {
			if i != len(check)-1 {
				t.Fatal("permutations generated", i, "expected", len(check))
			}
			break
		}
	}
}

func TestIsPrime(t *testing.T) {
	const p4M = 283146 // # of primes < 4e6
	n := 0
	for i := uint32(0); i <= 4e6; i++ {
		if IsPrime(i) {
			n++
		}
	}
	t.Log(n)
	if n != p4M {
		t.Fatal(n)
	}
}

func BenchmarkIsPrime(b *testing.B) {
	b.StopTimer()
	n := make([]uint32, b.N)
	rng := rand.New(rand.NewSource(1))
	for i := 0; i < b.N; i++ {
		n[i] = rng.Uint32()
	}
	b.StartTimer()
	for _, n := range n {
		IsPrime(n)
	}
}

func BenchmarkNextPrime(b *testing.B) {
	b.StopTimer()
	n := make([]uint32, b.N)
	rng := rand.New(rand.NewSource(1))
	for i := 0; i < b.N; i++ {
		n[i] = rng.Uint32()
	}
	b.StartTimer()
	for _, n := range n {
		NextPrime(n)
	}
}

func TestNextPrime(t *testing.T) {
	const p4M = 283146 // # of primes < 4e6
	n := 0
	var p uint32
	for {
		p, _ = NextPrime(p)
		if p >= 4e6 {
			break
		}
		n++
	}
	t.Log(n)
	if n != p4M {
		t.Fatal(n)
	}
}

func TestNextPrime2(t *testing.T) {
	type data struct {
		x  uint32
		y  uint32
		ok bool
	}
	tests := []data{
		{0, 2, true},
		{1, 2, true},
		{2, 3, true},
		{3, 5, true},
		{math.MaxUint32, 0, false},
		{math.MaxUint32 - 1, 0, false},
		{math.MaxUint32 - 2, 0, false},
		{math.MaxUint32 - 3, 0, false},
		{math.MaxUint32 - 4, 0, false},
		{math.MaxUint32 - 5, math.MaxUint32 - 4, true},
	}

	for _, test := range tests {
		y, ok := NextPrime(test.x)
		if ok != test.ok || ok && y != test.y {
			t.Fatalf("x %d, got y %d ok %t, expected y %d ok %t", test.x, y, ok, test.y, test.ok)
		}
	}
}

func TestISqrt(t *testing.T) {
	for n := int64(0); n < 5e6; n++ {
		x := int64(ISqrt(uint32(n)))
		if x2 := x * x; x2 > n {
			t.Fatalf("got ISqrt(%d) == %d, too big", n, x)
		}
		if x2 := x*x + 2*x + 1; x2 < n {
			t.Fatalf("got ISqrt(%d) == %d, too low", n, x)
		}
	}
	for n := int64(math.MaxUint32); n > math.MaxUint32-5e6; n-- {
		x := int64(ISqrt(uint32(n)))
		if x2 := x * x; x2 > n {
			t.Fatalf("got ISqrt(%d) == %d, too big", n, x)
		}
		if x2 := x*x + 2*x + 1; x2 < n {
			t.Fatalf("got ISqrt(%d) == %d, too low", n, x)
		}
	}
	rng := rand.New(rand.NewSource(1))
	for i := 0; i < 5e6; i++ {
		n := int64(rng.Uint32())
		x := int64(ISqrt(uint32(n)))
		if x2 := x * x; x2 > n {
			t.Fatalf("got ISqrt(%d) == %d, too big", n, x)
		}
		if x2 := x*x + 2*x + 1; x2 < n {
			t.Fatalf("got ISqrt(%d) == %d, too low", n, x)
		}
	}
}

func TestFactorInt(t *testing.T) {
	chk := func(n uint64, f []FactorTerm) bool {
		if n < 2 {
			return len(f) == 0
		}

		for i := 1; i < len(f); i++ { // verify ordering
			if t, u := f[i-1], f[i]; t.Prime >= u.Prime {
				return false
			}
		}

		x := uint64(1)
		for _, v := range f {
			if p := v.Prime; p < 0 || !IsPrime(uint32(v.Prime)) {
				return false
			}

			for i := uint32(0); i < v.Power; i++ {
				x *= uint64(v.Prime)
				if x > math.MaxUint32 {
					return false
				}
			}
		}
		return x == n
	}

	for n := uint64(0); n < 3e5; n++ {
		f := FactorInt(uint32(n))
		if !chk(n, f) {
			t.Fatalf("bad FactorInt(%d): %v", n, f)
		}
	}
	for n := uint64(math.MaxUint32); n > math.MaxUint32-12e4; n-- {
		f := FactorInt(uint32(n))
		if !chk(n, f) {
			t.Fatalf("bad FactorInt(%d): %v", n, f)
		}
	}
	rng := rand.New(rand.NewSource(1))
	for i := 0; i < 13e4; i++ {
		n := rng.Uint32()
		f := FactorInt(n)
		if !chk(uint64(n), f) {
			t.Fatalf("bad FactorInt(%d): %v", n, f)
		}
	}
}

func TestFactorIntB(t *testing.T) {
	const N = 3e5 // must be < math.MaxInt32
	factors := make([][]FactorTerm, N+1)
	// set up the divisors
	for prime := uint32(2); prime <= N; prime, _ = NextPrime(prime) {
		for n := int(prime); n <= N; n += int(prime) {
			factors[n] = append(factors[n], FactorTerm{prime, 0})
		}
	}
	// set up the powers
	for n := 2; n <= N; n++ {
		f := factors[n]
		m := uint32(n)
		for i, v := range f {
			for m%v.Prime == 0 {
				m /= v.Prime
				v.Power++
			}
			f[i] = v
		}
		factors[n] = f
	}
	// check equal
	for n, e := range factors {
		g := FactorInt(uint32(n))
		if len(e) != len(g) {
			t.Fatal(n, "len", g, "!=", e)
		}

		for i, ev := range e {
			gv := g[i]
			if ev.Prime != gv.Prime {
				t.Fatal(n, "prime", gv, ev)
			}

			if ev.Power != gv.Power {
				t.Fatal(n, "power", gv, ev)
			}
		}
	}
}

func BenchmarkISqrt(b *testing.B) {
	b.StopTimer()
	n := make([]uint32, b.N)
	rng := rand.New(rand.NewSource(1))
	for i := 0; i < b.N; i++ {
		n[i] = rng.Uint32()
	}
	b.StartTimer()
	for _, n := range n {
		ISqrt(n)
	}
}

func BenchmarkFactorInt(b *testing.B) {
	b.StopTimer()
	n := make([]uint32, b.N)
	rng := rand.New(rand.NewSource(1))
	for i := 0; i < b.N; i++ {
		n[i] = rng.Uint32()
	}
	b.StartTimer()
	for _, n := range n {
		FactorInt(n)
	}
}

func TestIsPrimeUint16(t *testing.T) {
	for n := 0; n <= math.MaxUint16; n++ {
		if IsPrimeUint16(uint16(n)) != IsPrime(uint32(n)) {
			t.Fatal(n)
		}
	}
}

func BenchmarkIsPrimeUint16(b *testing.B) {
	b.StopTimer()
	n := make([]uint16, b.N)
	rng := rand.New(rand.NewSource(1))
	for i := 0; i < b.N; i++ {
		n[i] = uint16(rng.Uint32())
	}
	b.StartTimer()
	for _, n := range n {
		IsPrimeUint16(n)
	}
}

func TestNextPrimeUint16(t *testing.T) {
	for n := 0; n <= math.MaxUint16; n++ {
		p, ok := NextPrimeUint16(uint16(n))
		p2, ok2 := NextPrime(uint32(n))
		switch {
		case ok:
			if !ok2 || uint32(p) != p2 {
				t.Fatal(n, p, ok)
			}
		case !ok && ok2:
			if p2 < 65536 {
				t.Fatal(n, p, ok)
			}
		}
	}
}

func BenchmarkNextPrimeUint16(b *testing.B) {
	b.StopTimer()
	n := make([]uint16, b.N)
	rng := rand.New(rand.NewSource(1))
	for i := 0; i < b.N; i++ {
		n[i] = uint16(rng.Uint32())
	}
	b.StartTimer()
	for _, n := range n {
		NextPrimeUint16(n)
	}
}

/*

From: http://graphics.stanford.edu/~seander/bithacks.html#CountBitsSetKernighan

Counting bits set, Brian Kernighan's way

unsigned int v; // count the number of bits set in v
unsigned int c; // c accumulates the total bits set in v
for (c = 0; v; c++)
{
  v &= v - 1; // clear the least significant bit set
}

Brian Kernighan's method goes through as many iterations as there are set bits.
So if we have a 32-bit word with only the high bit set, then it will only go
once through the loop.

Published in 1988, the C Programming Language 2nd Ed. (by Brian W. Kernighan
and Dennis M. Ritchie) mentions this in exercise 2-9. On April 19, 2006 Don
Knuth pointed out to me that this method "was first published by Peter Wegner
in CACM 3 (1960), 322. (Also discovered independently by Derrick Lehmer and
published in 1964 in a book edited by Beckenbach.)"
*/
func bcnt(v uint64) (c int) {
	for ; v != 0; c++ {
		v &= v - 1
	}
	return
}

func TestPopCount(t *testing.T) {
	const N = 4e5
	maxUint64 := big.NewInt(0)
	maxUint64.SetBit(maxUint64, 64, 1)
	maxUint64.Sub(maxUint64, big.NewInt(1))
	rng := r64()
	for i := 0; i < N; i++ {
		n := uint64(rng.Next().Int64())
		if g, e := PopCountByte(byte(n)), bcnt(uint64(byte(n))); g != e {
			t.Fatal(n, g, e)
		}

		if g, e := PopCountUint16(uint16(n)), bcnt(uint64(uint16(n))); g != e {
			t.Fatal(n, g, e)
		}

		if g, e := PopCountUint32(uint32(n)), bcnt(uint64(uint32(n))); g != e {
			t.Fatal(n, g, e)
		}

		if g, e := PopCount(int(n)), bcnt(uint64(uint(n))); g != e {
			t.Fatal(n, g, e)
		}

		if g, e := PopCountUint(uint(n)), bcnt(uint64(uint(n))); g != e {
			t.Fatal(n, g, e)
		}

		if g, e := PopCountUint64(n), bcnt(n); g != e {
			t.Fatal(n, g, e)
		}

		if g, e := PopCountUintptr(uintptr(n)), bcnt(uint64(n)); g != e {
			t.Fatal(n, g, e)
		}
	}
}

var gcds = []struct{ a, b, gcd uint64 }{
	{8, 12, 4},
	{12, 18, 6},
	{42, 56, 14},
	{54, 24, 6},
	{252, 105, 21},
	{1989, 867, 51},
	{1071, 462, 21},
	{2 * 3 * 5 * 7 * 11, 5 * 7 * 11 * 13 * 17, 5 * 7 * 11},
	{2 * 3 * 5 * 7 * 7 * 11, 5 * 7 * 7 * 11 * 13 * 17, 5 * 7 * 7 * 11},
	{2 * 3 * 5 * 7 * 7 * 11, 5 * 7 * 7 * 13 * 17, 5 * 7 * 7},
	{2 * 3 * 5 * 7 * 11, 13 * 17 * 19, 1},
}

func TestGCD(t *testing.T) {
	for i, v := range gcds {
		if v.a <= math.MaxUint16 && v.b <= math.MaxUint16 {
			if g, e := uint64(GCDUint16(uint16(v.a), uint16(v.b))), v.gcd; g != e {
				t.Errorf("%d: got gcd(%d, %d) %d, exp %d", i, v.a, v.b, g, e)
			}
			if g, e := uint64(GCDUint16(uint16(v.b), uint16(v.a))), v.gcd; g != e {
				t.Errorf("%d: got gcd(%d, %d) %d, exp %d", i, v.b, v.a, g, e)
			}
		}
		if v.a <= math.MaxUint32 && v.b <= math.MaxUint32 {
			if g, e := uint64(GCDUint32(uint32(v.a), uint32(v.b))), v.gcd; g != e {
				t.Errorf("%d: got gcd(%d, %d) %d, exp %d", i, v.a, v.b, g, e)
			}
			if g, e := uint64(GCDUint32(uint32(v.b), uint32(v.a))), v.gcd; g != e {
				t.Errorf("%d: got gcd(%d, %d) %d, exp %d", i, v.b, v.a, g, e)
			}
		}
		if g, e := GCDUint64(v.a, v.b), v.gcd; g != e {
			t.Errorf("%d: got gcd(%d, %d) %d, exp %d", i, v.a, v.b, g, e)
		}
		if g, e := GCDUint64(v.b, v.a), v.gcd; g != e {
			t.Errorf("%d: got gcd(%d, %d) %d, exp %d", i, v.b, v.a, g, e)
		}
	}
}

func lg2(n uint64) (lg int) {
	if n == 0 {
		return -1
	}

	for n >>= 1; n != 0; n >>= 1 {
		lg++
	}
	return
}

func TestLog2(t *testing.T) {
	if g, e := Log2Byte(0), -1; g != e {
		t.Error(g, e)
	}
	if g, e := Log2Uint16(0), -1; g != e {
		t.Error(g, e)
	}
	if g, e := Log2Uint32(0), -1; g != e {
		t.Error(g, e)
	}
	if g, e := Log2Uint64(0), -1; g != e {
		t.Error(g, e)
	}
	const N = 1e6
	rng := r64()
	for i := 0; i < N; i++ {
		n := uint64(rng.Next().Int64())
		if g, e := Log2Uint64(n), lg2(n); g != e {
			t.Fatalf("%b %d %d", n, g, e)
		}
		if g, e := Log2Uint32(uint32(n)), lg2(n&0xffffffff); g != e {
			t.Fatalf("%b %d %d", n, g, e)
		}
		if g, e := Log2Uint16(uint16(n)), lg2(n&0xffff); g != e {
			t.Fatalf("%b %d %d", n, g, e)
		}
		if g, e := Log2Byte(byte(n)), lg2(n&0xff); g != e {
			t.Fatalf("%b %d %d", n, g, e)
		}
	}
}

func TestBitLen(t *testing.T) {
	if g, e := BitLenByte(0), 0; g != e {
		t.Error(g, e)
	}
	if g, e := BitLenUint16(0), 0; g != e {
		t.Error(g, e)
	}
	if g, e := BitLenUint32(0), 0; g != e {
		t.Error(g, e)
	}
	if g, e := BitLenUint64(0), 0; g != e {
		t.Error(g, e)
	}
	if g, e := BitLenUintptr(0), 0; g != e {
		t.Error(g, e)
	}
	const N = 1e6
	rng := r64()
	for i := 0; i < N; i++ {
		n := uint64(rng.Next().Int64())
		if g, e := BitLenUintptr(uintptr(n)), lg2(uint64(n))+1; g != e {
			t.Fatalf("%b %d %d", n, g, e)
		}
		if g, e := BitLenUint64(n), lg2(n)+1; g != e {
			t.Fatalf("%b %d %d", n, g, e)
		}
		if g, e := BitLenUint32(uint32(n)), lg2(n&0xffffffff)+1; g != e {
			t.Fatalf("%b %d %d", n, g, e)
		}
		if g, e := BitLen(int(n)), lg2(uint64(uint(n)))+1; g != e {
			t.Fatalf("%b %d %d", n, g, e)
		}
		if g, e := BitLenUint(uint(n)), lg2(uint64(uint(n)))+1; g != e {
			t.Fatalf("%b %d %d", n, g, e)
		}
		if g, e := BitLenUint16(uint16(n)), lg2(n&0xffff)+1; g != e {
			t.Fatalf("%b %d %d", n, g, e)
		}
		if g, e := BitLenByte(byte(n)), lg2(n&0xff)+1; g != e {
			t.Fatalf("%b %d %d", n, g, e)
		}
	}
}

func BenchmarkGCDByte(b *testing.B) {
	const N = 1 << 16
	type t byte
	type u struct{ a, b t }
	b.StopTimer()
	rng := r32()
	a := make([]u, N)
	for i := range a {
		a[i] = u{t(rng.Next()), t(rng.Next())}
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		v := a[i&(N-1)]
		GCDByte(byte(v.a), byte(v.b))
	}
}

func BenchmarkGCDUint16(b *testing.B) {
	const N = 1 << 16
	type t uint16
	type u struct{ a, b t }
	b.StopTimer()
	rng := r32()
	a := make([]u, N)
	for i := range a {
		a[i] = u{t(rng.Next()), t(rng.Next())}
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		v := a[i&(N-1)]
		GCDUint16(uint16(v.a), uint16(v.b))
	}
}

func BenchmarkGCDUint32(b *testing.B) {
	const N = 1 << 16
	type t uint32
	type u struct{ a, b t }
	b.StopTimer()
	rng := r32()
	a := make([]u, N)
	for i := range a {
		a[i] = u{t(rng.Next()), t(rng.Next())}
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		v := a[i&(N-1)]
		GCDUint32(uint32(v.a), uint32(v.b))
	}
}

func BenchmarkGCDUint64(b *testing.B) {
	const N = 1 << 16
	type t uint64
	type u struct{ a, b t }
	b.StopTimer()
	rng := r64()
	a := make([]u, N)
	for i := range a {
		a[i] = u{t(rng.Next().Int64()), t(rng.Next().Int64())}
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		v := a[i&(N-1)]
		GCDUint64(uint64(v.a), uint64(v.b))
	}
}

func BenchmarkLog2Byte(b *testing.B) {
	const N = 1 << 16
	type t byte
	b.StopTimer()
	rng := r32()
	a := make([]t, N)
	for i := range a {
		a[i] = t(rng.Next())
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		Log2Byte(byte(a[i&(N-1)]))
	}
}

func BenchmarkLog2Uint16(b *testing.B) {
	const N = 1 << 16
	type t uint16
	b.StopTimer()
	rng := r32()
	a := make([]t, N)
	for i := range a {
		a[i] = t(rng.Next())
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		Log2Uint16(uint16(a[i&(N-1)]))
	}
}

func BenchmarkLog2Uint32(b *testing.B) {
	const N = 1 << 16
	type t uint32
	b.StopTimer()
	rng := r32()
	a := make([]t, N)
	for i := range a {
		a[i] = t(rng.Next())
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		Log2Uint32(uint32(a[i&(N-1)]))
	}
}

func BenchmarkLog2Uint64(b *testing.B) {
	const N = 1 << 16
	type t uint64
	b.StopTimer()
	rng := r64()
	a := make([]t, N)
	for i := range a {
		a[i] = t(rng.Next().Int64())
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		Log2Uint64(uint64(a[i&(N-1)]))
	}
}
func BenchmarkBitLenByte(b *testing.B) {
	const N = 1 << 16
	type t byte
	b.StopTimer()
	rng := r32()
	a := make([]t, N)
	for i := range a {
		a[i] = t(rng.Next())
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		BitLenByte(byte(a[i&(N-1)]))
	}
}

func BenchmarkBitLenUint16(b *testing.B) {
	const N = 1 << 16
	type t uint16
	b.StopTimer()
	rng := r32()
	a := make([]t, N)
	for i := range a {
		a[i] = t(rng.Next())
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		BitLenUint16(uint16(a[i&(N-1)]))
	}
}

func BenchmarkBitLenUint32(b *testing.B) {
	const N = 1 << 16
	type t uint32
	b.StopTimer()
	rng := r32()
	a := make([]t, N)
	for i := range a {
		a[i] = t(rng.Next())
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		BitLenUint32(uint32(a[i&(N-1)]))
	}
}

func BenchmarkBitLen(b *testing.B) {
	const N = 1 << 16
	type t int
	b.StopTimer()
	rng := r64()
	a := make([]t, N)
	for i := range a {
		a[i] = t(rng.Next().Int64())
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		BitLen(int(a[i&(N-1)]))
	}
}

func BenchmarkBitLenUint(b *testing.B) {
	const N = 1 << 16
	type t uint
	b.StopTimer()
	rng := r64()
	a := make([]t, N)
	for i := range a {
		a[i] = t(rng.Next().Int64())
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		BitLenUint(uint(a[i&(N-1)]))
	}
}

func BenchmarkBitLenUintptr(b *testing.B) {
	const N = 1 << 16
	type t uintptr
	b.StopTimer()
	rng := r64()
	a := make([]t, N)
	for i := range a {
		a[i] = t(rng.Next().Int64())
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		BitLenUintptr(uintptr(a[i&(N-1)]))
	}
}

func BenchmarkBitLenUint64(b *testing.B) {
	const N = 1 << 16
	type t uint64
	b.StopTimer()
	rng := r64()
	a := make([]t, N)
	for i := range a {
		a[i] = t(rng.Next().Int64())
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		BitLenUint64(uint64(a[i&(N-1)]))
	}
}

func BenchmarkPopCountByte(b *testing.B) {
	const N = 1 << 16
	type t byte
	b.StopTimer()
	rng := r32()
	a := make([]t, N)
	for i := range a {
		a[i] = t(rng.Next())
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		PopCountByte(byte(a[i&(N-1)]))
	}
}

func BenchmarkPopCountUint16(b *testing.B) {
	const N = 1 << 16
	type t uint16
	b.StopTimer()
	rng := r32()
	a := make([]t, N)
	for i := range a {
		a[i] = t(rng.Next())
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		PopCountUint16(uint16(a[i&(N-1)]))
	}
}

func BenchmarkPopCountUint32(b *testing.B) {
	const N = 1 << 16
	type t uint32
	b.StopTimer()
	rng := r32()
	a := make([]t, N)
	for i := range a {
		a[i] = t(rng.Next())
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		PopCountUint32(uint32(a[i&(N-1)]))
	}
}

func BenchmarkPopCount(b *testing.B) {
	const N = 1 << 16
	type t int
	b.StopTimer()
	rng := r64()
	a := make([]t, N)
	for i := range a {
		a[i] = t(rng.Next().Int64())
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		PopCount(int(a[i&(N-1)]))
	}
}

func BenchmarkPopCountUint(b *testing.B) {
	const N = 1 << 16
	type t uint
	b.StopTimer()
	rng := r64()
	a := make([]t, N)
	for i := range a {
		a[i] = t(rng.Next().Int64())
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		PopCountUint(uint(a[i&(N-1)]))
	}
}

func BenchmarkPopCountUintptr(b *testing.B) {
	const N = 1 << 16
	type t uintptr
	b.StopTimer()
	rng := r64()
	a := make([]t, N)
	for i := range a {
		a[i] = t(rng.Next().Int64())
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		PopCountUintptr(uintptr(a[i&(N-1)]))
	}
}

func BenchmarkPopCountUint64(b *testing.B) {
	const N = 1 << 16
	type t uint64
	b.StopTimer()
	rng := r64()
	a := make([]t, N)
	for i := range a {
		a[i] = t(rng.Next().Int64())
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		PopCountUint64(uint64(a[i&(N-1)]))
	}
}

func TestUintptrBits(t *testing.T) {
	switch g := UintptrBits(); g {
	case 32, 64:
		// ok
		t.Log(g)
	default:
		t.Fatalf("got %d, expected 32 or 64", g)
	}
}

func BenchmarkUintptrBits(b *testing.B) {
	for i := 0; i < b.N; i++ {
		UintptrBits()
	}
}

func TestUint64ToBigInt(t *testing.T) {
	const N = 2e5
	data := []uint64{0, 1, math.MaxInt64 - 1, math.MaxInt64, math.MaxInt64 + 1, math.MaxUint64 - 1, math.MaxUint64}

	var e big.Int
	f := func(n uint64) {
		g := Uint64ToBigInt(n)
		e.SetString(fmt.Sprintf("%d", n), 10)
		if g.Cmp(&e) != 0 {
			t.Errorf("got %s(0x%x), exp %d(0x%x)", g, g, n, n)
		}
	}

	for _, v := range data {
		f(v)
	}

	r := r64()
	for i := 0; i < N; i++ {
		f(uint64(r.Next().Int64()))
	}
}

func BenchmarkUint64ToBigInt(b *testing.B) {
	const N = 1 << 16
	b.StopTimer()
	a := make([]uint64, N)
	r := r64()
	for i := range a {
		a[i] = uint64(r.Next().Int64())
	}
	runtime.GC()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		Uint64ToBigInt(a[i&(N-1)])
	}
}

func TestUint64FromBigInt(t *testing.T) {
	const N = 2e5
	data := []struct {
		s  string
		e  uint64
		ok bool
	}{
		{"-2", 0, false},
		{"-1", 0, false},
		{"0", 0, true},
		{"1", 1, true},
		{"2", 2, true},

		{"4294967294", 4294967294, true},
		{"4294967295", 4294967295, true},
		{"4294967296", 4294967296, true},
		{"4294967297", 4294967297, true},
		{"4294967298", 4294967298, true},

		{"18446744073709551613", 18446744073709551613, true},
		{"18446744073709551614", 18446744073709551614, true},
		{"18446744073709551615", 18446744073709551615, true},
		{"18446744073709551616", 0, false},
		{"18446744073709551617", 0, false},
		{"18446744073709551618", 0, false},
	}

	var x big.Int
	f := func(s string, e uint64, ok bool) {
		x.SetString(s, 10)
		switch g, gok := Uint64FromBigInt(&x); {
		case gok != ok:
			t.Errorf("%s: got %t, exp %t", s, gok, ok)
		case ok && g != e:
			t.Errorf("%s: got %d, exp %d", s, g, s)
		}

	}

	for _, v := range data {
		f(v.s, v.e, v.ok)
	}
	r := r64()
	for i := 0; i < N; i++ {
		n := uint64(r.Next().Int64())
		f(fmt.Sprintf("%d", n), n, true)
	}
}

func BenchmarkUint64FromBigInt(b *testing.B) {
	const N = 1 << 16
	b.StopTimer()
	a := make([]*big.Int, N)
	r := r64()
	for i := range a {
		a[i] = Uint64ToBigInt(uint64(r.Next().Int64()))
	}
	runtime.GC()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		Uint64FromBigInt(a[i&(N-1)])
	}
}

func TestModPowByte(t *testing.T) {
	data := []struct{ b, e, m, r byte }{
		{2, 11, 23, 1}, // 23|M11
		{2, 11, 89, 1}, // 89|M11
		{2, 23, 47, 1}, // 47|M23
		{5, 3, 13, 8},
	}

	for _, v := range data {
		if g, e := ModPowByte(v.b, v.e, v.m), v.r; g != e {
			t.Errorf("b %d e %d m %d: got %d, exp %d", v.b, v.e, v.m, g, e)
		}
	}
}

func TestModPowUint16(t *testing.T) {
	data := []struct{ b, e, m, r uint16 }{
		{2, 11, 23, 1},     // 23|M11
		{2, 11, 89, 1},     // 89|M11
		{2, 23, 47, 1},     // 47|M23
		{2, 929, 13007, 1}, // 13007|M929
		{4, 13, 497, 445},
		{5, 3, 13, 8},
	}

	for _, v := range data {
		if g, e := ModPowUint16(v.b, v.e, v.m), v.r; g != e {
			t.Errorf("b %d e %d m %d: got %d, exp %d", v.b, v.e, v.m, g, e)
		}
	}
}

func TestModPowUint32(t *testing.T) {
	data := []struct{ b, e, m, r uint32 }{
		{2, 23, 47, 1},        // 47|M23
		{2, 67, 193707721, 1}, // 193707721|M67
		{2, 929, 13007, 1},    // 13007|M929
		{4, 13, 497, 445},
		{5, 3, 13, 8},
	}

	for _, v := range data {
		if g, e := ModPowUint32(v.b, v.e, v.m), v.r; g != e {
			t.Errorf("b %d e %d m %d: got %d, exp %d", v.b, v.e, v.m, g, e)
		}
	}
}

func TestModPowUint64(t *testing.T) {
	data := []struct{ b, e, m, r uint64 }{
		{2, 23, 47, 1},        // 47|M23
		{2, 67, 193707721, 1}, // 193707721|M67
		{2, 929, 13007, 1},    // 13007|M929
		{4, 13, 497, 445},
		{5, 3, 13, 8},
	}

	for _, v := range data {
		if g, e := ModPowUint64(v.b, v.e, v.m), v.r; g != e {
			t.Errorf("b %d e %d m %d: got %d, exp %d", v.b, v.e, v.m, g, e)
		}
	}
}

func TestModPowBigInt(t *testing.T) {
	data := []struct{ b, e, m, r int64 }{
		{2, 23, 47, 1},        // 47|M23
		{2, 67, 193707721, 1}, // 193707721|M67
		{2, 929, 13007, 1},    // 13007|M929
		{4, 13, 497, 445},
		{5, 3, 13, 8},
	}

	for _, v := range data {
		b, e, m, r := big.NewInt(v.b), big.NewInt(v.e), big.NewInt(v.m), big.NewInt(v.r)
		if g, e := ModPowBigInt(b, e, m), r; g.Cmp(e) != 0 {
			t.Errorf("b %s e %s m %s: got %s, exp %s", b, e, m, g, e)
		}
	}
}

func BenchmarkModPowByte(b *testing.B) {
	const N = 1 << 16
	b.StopTimer()
	type t struct{ b, e, m byte }
	a := make([]t, N)
	r := r32()
	for i := range a {
		a[i] = t{
			byte(r.Next()),
			byte(r.Next()),
			byte(r.Next() | 1),
		}
	}
	runtime.GC()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		v := a[i&(N-1)]
		ModPowByte(v.b, v.e, v.m)
	}
}

func BenchmarkModPowUint16(b *testing.B) {
	const N = 1 << 16
	b.StopTimer()
	type t struct{ b, e, m uint16 }
	a := make([]t, N)
	r := r32()
	for i := range a {
		a[i] = t{
			uint16(r.Next()),
			uint16(r.Next()),
			uint16(r.Next() | 1),
		}
	}
	runtime.GC()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		v := a[i&(N-1)]
		ModPowUint16(v.b, v.e, v.m)
	}
}

func BenchmarkModPowUint32(b *testing.B) {
	const N = 1 << 16
	b.StopTimer()
	type t struct{ b, e, m uint32 }
	a := make([]t, N)
	r := r32()
	for i := range a {
		a[i] = t{
			uint32(r.Next()),
			uint32(r.Next()),
			uint32(r.Next() | 1),
		}
	}
	runtime.GC()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		v := a[i&(N-1)]
		ModPowUint32(v.b, v.e, v.m)
	}
}

func BenchmarkModPowUint64(b *testing.B) {
	const N = 1 << 16
	b.StopTimer()
	type t struct{ b, e, m uint64 }
	a := make([]t, N)
	r := r64()
	for i := range a {
		a[i] = t{
			uint64(r.Next().Int64()),
			uint64(r.Next().Int64()),
			uint64(r.Next().Int64() | 1),
		}
	}
	runtime.GC()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		v := a[i&(N-1)]
		ModPowUint64(v.b, v.e, v.m)
	}
}

func BenchmarkModPowBigInt(b *testing.B) {
	const N = 1 << 16
	b.StopTimer()
	type t struct{ b, e, m *big.Int }
	a := make([]t, N)
	mx := big.NewInt(math.MaxInt64)
	mx.Mul(mx, mx)
	r, err := NewFCBig(big.NewInt(1), mx, true)
	if err != nil {
		b.Fatal(err)
	}
	for i := range a {
		a[i] = t{
			r.Next(),
			r.Next(),
			r.Next(),
		}
	}
	runtime.GC()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		v := a[i&(N-1)]
		ModPowBigInt(v.b, v.e, v.m)
	}
}
