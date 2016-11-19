package dice

import (
	"math/rand"
	"time"
)

type Die struct {
	Min int
	Max int
}

var (
	D2  = Die{1, 2}
	D4  = Die{1, 4}
	D6  = Die{1, 4}
	D8  = Die{1, 8}
	D20 = Die{1, 20}
)

func Even(number int) bool {
	return number%2 == 0
}

func MakeDie(min int, max int) *Die {
	return &Die{min, max}
}

func (d *Die) RollEven() int {
	var num int
	for {
		num = d.Roll()
		if Even(num) == true {
			break
		}
	}
	return num
}

func (d *Die) RollOdd() int {
	var num int
	for {
		num = d.Roll()
		if Even(num) == false {
			break
		}
	}
	return num
}

func (d *Die) Roll() int {
	rand.Seed(time.Now().UTC().UnixNano() + time.Now().UTC().UnixNano())
	return d.Min + rand.Intn(d.Max-d.Min)
}

func RollMultiple(count int, highest int) []int {
	list := make([]int, 0)
	d := &Die{1, highest}
	for i := 0; i < count; i++ {
		list = append(list, d.Roll())
	}
	return list
}

func Init() {
	rand.Seed(time.Now().UTC().UnixNano() + time.Now().UTC().UnixNano())
}
