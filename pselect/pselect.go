package pselect

// 优先级select ch1 的任务先执行完毕后才会执行 ch2 里面的任务
func Select[T1, T2 any](ch1 <-chan T1, f1 func(T1), ch2 <-chan T2, f2 func(T2)) {
	for {
		select {
		case a := <-ch1:
			f1(a)
		case b := <-ch2:
		priority:
			for {
				select {
				case a := <-ch1:
					f1(a)
				default:
					break priority
				}
			}
			f2(b)
		}
	}
}
