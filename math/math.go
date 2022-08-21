package math

// 把数n分成m份
// 返回每一份大小和最后一份的大小
func Split(n, m uint) (uint, uint) {
	per := (n + m - 1) / m
	last := n % per
	if last == 0 {
		last = per
	}
	return per, last
}
