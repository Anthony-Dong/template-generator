package utils

func LowerCaseFiledFirst(str string) string {
	if str == "" {
		return str
	}
	arr := []byte(str)
	if arr[0] >= 'A' && arr[0] <= 'Z' {
		arr[0] = arr[0] + 32
	}
	return string(arr)
}
