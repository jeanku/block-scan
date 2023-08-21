package scan

import (
	"container/list"
	"fmt"
	"math/big"
)

type GasValue struct {
	GasUsed   *big.Int
	GasTipCap int64
	GasFeeCap int64
	GasBase   uint64
	GasPrice  int64
	GasLimit  uint64
}

type fifo struct {
	length     int
	totalLimit uint64
	ll         *list.List
}

func FifoNew(length int) *fifo {
	return &fifo{
		length: length,
		ll:     list.New(),
	}
}

func (f *fifo) Set(value GasValue) {
	f.totalLimit += value.GasLimit
	f.ll.PushBack(value)
	if f.Len() > f.length {
		f.Eliminate()
	}
}

func (f *fifo) Eliminate() {
	f.RemoveElement(f.ll.Front())
}

func (f *fifo) Avg() uint64 {
	length := f.Len()
	if length == 0 {
		return 0
	}
	return f.totalLimit / uint64(f.Len())
}

func (f *fifo) Len() int {
	return f.ll.Len()
}

func (f *fifo) RemoveElement(e *list.Element) {
	if e != nil {
		f.ll.Remove(e)
		f.totalLimit -= e.Value.(GasValue).GasLimit
	}
}

func (f *fifo) Gas() *big.Int {
	node := f.ll.Front()
	total := big.NewInt(0)
	for {
		if node == nil {
			break
		}
		value := node.Value.(GasValue)
		total = new(big.Int).Add(total, value.GasUsed)
		node = node.Next()
	}
	return new(big.Int).Div(total, big.NewInt(int64(f.ll.Len())))
}

func (f *fifo) Print() {
	node := f.ll.Front()
	for {
		if node == nil {
			break
		}
		value := node.Value.(GasValue)
		node = node.Next()
		fmt.Print(" node:", value.GasLimit, " ")
	}
	fmt.Println("total: ", f.totalLimit, f.Avg())
}
