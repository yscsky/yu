package yu

// Serial 序列生成器
type Serial struct {
	serial chan int
	count  int
	reset  int
}

// NewSerial 创建Serial
func NewSerial(reset int) *Serial {
	return &Serial{
		serial: make(chan int),
		count:  0,
		reset:  reset,
	}
}

// Start 开始序列生成器
func (s *Serial) Start() {
	go func() {
		for {
			s.serial <- s.count
			s.count++
			if s.count > s.reset {
				s.count = 0
			}
		}
	}()
}

// Get 获取序列值
func (s *Serial) Get() int {
	return <-s.serial
}

// Reset 重置计数
func (s *Serial) Reset(c int) {
	s.count = c
}
