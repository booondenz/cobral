package helper

// FormatNumber 把没有分隔符的数字串加上千分位逗号。
// 输入 [123456789] -> 输出 [123,456,789]。
// 这个函数已经写好了，作为「字符串按字节处理」的例子读一读。
func FormatNumber(s string) string {
	length := len(s)
	if length == 0 {
		return "<nil>"
	}
	newLength := length + (length-1)/3
	newNumber := make([]rune, newLength)
	commas := 0
	for i, c := range s {
		if c < '0' || c > '9' {
			return "<nil>"
		}
		if (length-i)%3 == 0 && i != 0 {
			newNumber[i+commas] = ','
			commas += 1
		}
		newNumber[i+commas] = c
	}
	return string(newNumber)
}
