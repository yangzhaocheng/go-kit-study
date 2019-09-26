package impl

import "errors"

//ArithmeticService implement Service interface
type ArithmeticService struct {
}

// Add implement Add method
func (s ArithmeticService) Add(a, b int) int {
	return a + b
}

// Subtract implement Subtract method
func (s ArithmeticService) Subtract(a, b int) int {
	return a - b
}

// Multiply implement Multiply method
func (s ArithmeticService) Multiply(a, b int) int {
	return a * b
}

// Divide implement Divide method
func (s ArithmeticService) Divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("the dividend can not be zero!")
	}

	return a / b, nil
}
