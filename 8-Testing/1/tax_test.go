package tax

import "testing"

func TestCalculateTax(t *testing.T) { 
	amount := 500.0
	expected := 5.0
	result := CalculateTax(amount)

	if result != expected {
		t.Errorf("Expected %f but got %f", expected, result)
	}
}