package utils

import (
	"fmt"
	"testing"
)

func Test_decimalToBase64(t *testing.T) {
	for idx, rune := range CHARS() {

		number := idx + 1
		expected := string(rune)
		base64 := DecimalToBase64(number)
		if base64 != expected {
			redbase64 := "\033[1;31m" + base64 + "\033[0m"
			t.Errorf("Expected %v but got %v when converting %v", expected, redbase64, number)
		}
		for j, rune2 := range CHARS() {
			number := (idx+1)*64 + (j + 1)
			expected := string(rune) + string(rune2)
			base64 := DecimalToBase64(number)

			redbase64 := "\033[1;31m" + base64 + "\033[0m"
			if base64 != expected {
				t.Errorf("Expected %v but got %v when converting %v", expected, redbase64, number)
			}
		}
	}
}

func Test_base64toDecimal(t *testing.T) {
	for idx, rune := range CHARS() {
		base64 := string(rune)
		num := Base64ToDecimal(base64)
		expected := idx + 1
		if num != expected {
			redNum := "\033[1;31m" + fmt.Sprint(num) + "\033[0m"
			t.Errorf("Expected %d but got %v for char %v", expected, redNum, base64)
		}
		for j, rune2 := range CHARS() {

			base64 := string(rune) + string(rune2)
			num := Base64ToDecimal(base64)
			expected := (idx+1)*64 + j + 1
			if num != expected {
				redNum := "\033[1;31m" + fmt.Sprint(num) + "\033[0m"
				t.Errorf("Expected %d but got %v for char %v", expected, redNum, base64)
			}
		}
	}
}
