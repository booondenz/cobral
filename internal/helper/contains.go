package helper

// ContainsInt 判断整数切片里有没有指定的那个整数。
// 有就返回 true，没有返回 false。
func ContainsInt(ints []int, specifiedInt int) bool {
	for _, i := range ints {
		if i == specifiedInt {
			return true
		}
	}
	return false
}
