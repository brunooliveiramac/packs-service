package pack

import "testing"

func TestCalculator(t *testing.T) {
	calc, err := NewCalculator([]int{250, 500, 1000, 2000, 5000})
	if err != nil { t.Fatalf("new calc: %v", err) }

	tests := []struct{
		qty int
		wantShipped int
		want map[int]int
	}{
		{1, 250, map[int]int{250:1}},
		{250, 250, map[int]int{250:1}},
		{251, 500, map[int]int{500:1}},
		{501, 750, map[int]int{500:1, 250:1}},
		{12001, 12250, map[int]int{5000:2, 2000:1, 250:1}},
	}

	for _, tc := range tests {
		brk, shipped, err := calc.Calculate(tc.qty)
		if err != nil { t.Fatalf("qty=%d err=%v", tc.qty, err) }
		if shipped != tc.wantShipped {
			t.Fatalf("qty=%d shipped=%d want=%d", tc.qty, shipped, tc.wantShipped)
		}
		for size, count := range tc.want {
			if brk[size] != count {
				t.Fatalf("qty=%d size=%d got=%d want=%d", tc.qty, size, brk[size], count)
			}
		}
	}
}


