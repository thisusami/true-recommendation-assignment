package util

func IsNotEmptyString(s string) bool {
	return s != ""
}
func MinMaxInteger(i int, min int, max int) bool {
	if min == -1 {
		return i <= max
	} else if max == -1 {
		return i >= min
	}
	return i >= min && i <= max
}
