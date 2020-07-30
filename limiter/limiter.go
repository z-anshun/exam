package limiter

type Limiter struct {
	Conn chan int
}

var L = NewLimiter()

//设置1000个 每10分钟只允许1000个弹幕 超过1000条就忽略 向b站学习
func NewLimiter() *Limiter {
	return &Limiter{Conn: make(chan int, 1000)}
}

func (l *Limiter) SetLimiter() {
	l.Conn <- 1
}

//释放
func (l *Limiter) ReSetLimiter() {
	for {
		<-l.Conn
	}
}
