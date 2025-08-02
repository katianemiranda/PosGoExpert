package tax

import "testing"

func CalculateTax(amount float64) float64 {
	if amount == 0 {
		return 0.0
	}

	if amount >= 1000 {
		return 10.0
	}
	return 5.0
}

func FuzzyCalculateTax(f *testing.F) {
	seed := []float64{-1, -2, -2.5, 500.0, 1000.0, 1501.0}
	for _, amount := range seed {
		f.Add(amount)
	}
	f.Fuzz(func(t *testing.T, amount float64) {
		result := CalculateTax(amount)
		if amount <= 0 && result != 0.0 {
			t.Errorf("Received %f, expected 0.0 ", result)
		}
	})
}
