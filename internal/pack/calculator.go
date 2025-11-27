package pack

import (
	"errors"
	"sort"
)

var (
	ErrInvalidQuantity = errors.New("invalid quantity")
	ErrNoPackSizes     = errors.New("no pack sizes configured")
)

type Calculator struct {
	packSizes []int
}

func NewCalculator(sizes []int) (*Calculator, error) {
	if len(sizes) == 0 {
		return nil, ErrNoPackSizes
	}
	s := append([]int(nil), sizes...)
	sort.Sort(sort.Reverse(sort.IntSlice(s)))
	return &Calculator{packSizes: s}, nil
}

type Breakdown map[int]int // packSize to count

// Calculate returns the minimal shipped items >= quantity, with fewest packs for that shipped total.
func (c *Calculator) Calculate(quantity int) (Breakdown, int, error) {
	if quantity <= 0 {
		return nil, 0, ErrInvalidQuantity
	}

	// packs already sorted DESC in NewCalculator
	sizesDesc := c.packSizes
	breakdown := Breakdown{}
	remaining := quantity
	total := 0

	// GREEDY FILL: take as many big packs as possible ---
	for _, size := range sizesDesc {
		if remaining >= size {
			count := remaining / size
			breakdown[size] += count
			remaining -= count * size
			total += count * size
		}
	}

	// If leftover > 0, add ONE smallest pack ---
	smallest := sizesDesc[len(sizesDesc)-1]
	if remaining > 0 {
		breakdown[smallest]++
		total += smallest
		remaining -= smallest
	}

	// MERGE STEP: upgrade small packs into bigger ones ---
	// packSizes is DESC, so merge from smallest to biggest means reverse walk
	for i := len(sizesDesc) - 1; i > 0; i-- {
		small := sizesDesc[i]
		large := sizesDesc[i-1]
		ratio := large / small

		for breakdown[small] >= ratio {
			breakdown[small] -= ratio
			breakdown[large]++
		}
	}

	return breakdown, total, nil
}
