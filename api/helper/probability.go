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

func (r *Probability) NewGetRand() (key string) {
	r.sum()
	rNum := (rand.Intn(100000000) + 1) % r.proSum
	for k, v := range r.proArr {
		if rNum < v {
			key = k
			break
		} else {
			rNum -= v
		}
	}
	return
}
