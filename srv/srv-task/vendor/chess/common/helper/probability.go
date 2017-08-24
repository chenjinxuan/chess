// 简易概率算法
package helper

import (
	"math/rand"
)

type Probability struct {
	proArr map[string]int
	proSum int
}

func NewProbability(proArr map[string]int) *Probability {
	return &Probability{
		proArr: proArr,
	}
}

func (r *Probability) sum() {
	r.proSum = 0
	for _, v := range r.proArr {
		r.proSum += v
	}
}

func (r *Probability) GetRand() (key string) {
	r.sum()

	for k, v := range r.proArr {
		rNum := rand.Intn(r.proSum) + 1
		if rNum <= v {
			key = k
			break
		} else {
			r.proSum -= v
		}
	}
	return
}
