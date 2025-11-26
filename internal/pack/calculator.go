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
	sort.Ints(s)
	return &Calculator{packSizes: s}, nil
}

type Breakdown map[int]int // packSize -> count

// Calculate returns the minimal shipped items >= quantity, with fewest packs for that shipped total.
func (c *Calculator) Calculate(quantity int) (Breakdown, int, error) {
	if quantity <= 0 {
		return nil, 0, ErrInvalidQuantity
	}
	maxSize := c.packSizes[len(c.packSizes)-1]
	// dynamic programming: dp[t] = min packs to make exactly t, and prev for reconstruction
	type node struct {
		packs int
		prev  int // previous total
		size  int // last pack size used
	}
	dp := map[int]node{0: {packs: 0, prev: -1, size: 0}}

	target := quantity
	found := -1
	limit := quantity + maxSize // initial search window

	for {
		for t := 1; t <= limit; t++ {
			if _, ok := dp[t]; ok {
				continue
			}
			bestPacks := int(^uint(0) >> 1) // max int
			bestPrev := -1
			bestSize := 0
			for _, sz := range c.packSizes {
				if t-sz < 0 {
					break
				}
				if prev, ok := dp[t-sz]; ok {
					if prev.packs+1 < bestPacks {
						bestPacks = prev.packs + 1
						bestPrev = t - sz
						bestSize = sz
					}
				}
			}
			if bestPrev != -1 {
				dp[t] = node{packs: bestPacks, prev: bestPrev, size: bestSize}
			}
		}
		// find minimal achievable total >= target
		for t := target; t <= limit; t++ {
			if _, ok := dp[t]; ok {
				found = t
				break
			}
		}
		if found != -1 {
			break
		}
		// expand window and continue
		limit += maxSize
		// safety bound to avoid infinite loops in pathological inputs
		if limit > quantity+maxSize*100 {
			return nil, 0, errors.New("unable to compute pack combination")
		}
	}

	// reconstruct breakdown from found
	res := Breakdown{}
	cur := found
	for cur > 0 {
		n := dp[cur]
		res[n.size]++
		cur = n.prev
	}
	return res, found, nil
}


