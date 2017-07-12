package helper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"sort"
	"runtime"
	"github.com/davecgh/go-spew/spew"
	"chess/common/log"
)

// A data structure to hold key/value pairs
type Pair struct {
	Key   int
	Value float32
}

// A slice of pairs that implements sort.Interface to sort by values
type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

func SortMapByValue(m map[int]float32) PairList {
	pl := make(PairList, len(m))
	i := 0
	for k, v := range m {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

// 获取流量值 M -> G
func GenFlowValue(num int) string {
	var val string

	if num < 1024 {
		val = fmt.Sprintf("%dM", num)
	} else {
		val = fmt.Sprintf("%dG", num/1024)
	}

	return val
}

func Round(f float64, n int) float64 {
	pow10_n := math.Pow10(n)
	return math.Trunc((f+0.5/pow10_n)*pow10_n) / pow10_n
}

func LoadFileConfig(path string, c interface{}) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, c)
}

func TypeConvert(from interface{}, to interface{}) bool {
	fromByte, err := json.Marshal(from)
	if err != nil {
		return false
	}
	err = json.Unmarshal(fromByte, to)
	if err != nil {
		return false
	}
	return true
}

// Exists reports whether the named file or directory exists.
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// 产生panic时的调用栈打印
func PrintPanicStack(extras ...interface{}) {
	if x := recover(); x != nil {
		log.Error(x)
		i := 0
		funcName, file, line, ok := runtime.Caller(i)
		for ok {
			log.Errorf("frame %v:[func:%v,file:%v,line:%v]\n", i, runtime.FuncForPC(funcName).Name(), file, line)
			i++
			funcName, file, line, ok = runtime.Caller(i)
		}

		for k := range extras {
			log.Errorf("EXRAS#%v DATA:%v\n", k, spew.Sdump(extras[k]))
		}
	}
}