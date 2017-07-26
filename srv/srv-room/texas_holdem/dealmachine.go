package texas_holdem

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type DealMachine struct {
	cards        [SUITSIZE * CARDRANK]*Card //52张牌
	topCardIndex int                        //当前牌顶
	initilized   bool                       //是否已经初始化
}

func NewDealMachine() *DealMachine {
	d := new(DealMachine)
	d.initilized = false
	return d
}

/*
初始化牌组
对于花色：0代表黑桃、1代表红桃、2代表梅花、3代表方块，详见card包
对于值：0代表two，1代表three .. 12代表A
*/
func (dm *DealMachine) Init() {
	for i := 0; i < SUITSIZE; i++ {
		for j := 0; j < CARDRANK; j++ {
			dm.cards[i*CARDRANK+j] = new(Card)
			dm.cards[i*CARDRANK+j].Suit = i
			dm.cards[i*CARDRANK+j].Value = j
		}
	}
	dm.topCardIndex = 0
	dm.initilized = true
}

/*
洗牌！！游戏每次开始时候调用，允许多次调用。
随机序列生成的逻辑是这样的：
从后往前，N个数为例。
先生成一0~~N-1的随机数i，然后置换i和N之间的位置
同理处理N-1....
*/
func (dm *DealMachine) Shuffle() error {
	if dm.initilized == false {
		return errors.New("you must init DealMachine first")
	}
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	for i := SUITSIZE*CARDRANK - 1; i > 0; i-- {
		index := r.Int() % i
		dm.swapCard(i, index)
	}
	dm.topCardIndex = 0
	return nil
}

/*
调用此函数发一张牌
*/
func (dm *DealMachine) Deal() *Card {
	c := dm.cards[dm.topCardIndex]
	dm.topCardIndex++
	if dm.topCardIndex == SUITSIZE*CARDRANK {
		_ = dm.Shuffle()
	}
	return c
}

func (dm *DealMachine) swapCard(a int, b int) {
	tmp := dm.cards[a]
	dm.cards[a] = dm.cards[b]
	dm.cards[b] = tmp
}

func (dm *DealMachine) ShowCards() {
	for i := 0; i < SUITSIZE*CARDRANK; i++ {
		fmt.Printf("%s %d, ", SUITNAME[dm.cards[i].Suit], dm.cards[i].Value)
	}
	fmt.Println()
}
