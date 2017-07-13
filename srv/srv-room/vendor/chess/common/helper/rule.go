package helper

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"strconv"
	"strings"
	"time"
)

/**
* rule data
{
	'from': [],
	'clients_type': 'all',
	'clients': ['android_130','ios_131'],
	'time_type': 'short',
	'start': 14002321142,
	'end': 14002321142
}
*/

const (
	ClientsTypeAll   = "all"
	ClientsTypeWhite = "white"
	ClientsTypeBlack = "black"
	ClientsTypeGte   = "gte"
	ClientsTypeLte   = "lte"

	TimeTypeLong  = "long"
	TimeTypeShort = "short"
)

type Rule struct {
	From        []string `json:"from"`                            // [] 代表全部  or  ['ios', 'android']
	ClientsType string   `json:"clients_type" binding:"required"` // all,white,black,lte,gte
	Clients     []string `json:"clients"`                         // "from_ver"
	TimeType    string   `json:"time_type" binding:"required"`    // long,short
	Start       int64    `json:"start"`
	End         int64    `json:"end"`
}

func NewRule(ctr string) *Rule {
	rule := new(Rule)

	err := json.Unmarshal([]byte(ctr), rule)
	if err != nil {
		return nil
	}

	err = binding.Validator.ValidateStruct(rule)
	if err != nil {
		return nil
	}

	return rule
}

func (r *Rule) IsMeet(from string, ver int) bool {
	from = strings.ToLower(from)
	client := fmt.Sprintf("%s_%d", from, ver)

	// 校验from
	if len(r.From) > 0 && StringInSlice(r.From, from) == -1 {
		return false
	}

	// 校验 版本 黑白名单
	switch r.ClientsType {
	case ClientsTypeWhite:
		if StringInSlice(r.Clients, client) == -1 {
			return false
		}
	case ClientsTypeBlack:
		if StringInSlice(r.Clients, client) != -1 {
			return false
		}
	case ClientsTypeGte, ClientsTypeLte: // >= 某个版本 || <= 某个版本
		for _, cli := range r.Clients {
			tmp := strings.Split(cli, "_")
			if len(tmp) != 2 {
				return false
			}
			_from := tmp[0]
			_ver, err := strconv.Atoi(tmp[1])
			if err != nil {
				return false
			}

			if _from == from {
				if r.ClientsType == ClientsTypeGte && ver < _ver {
					return false
				}
				if r.ClientsType == ClientsTypeLte && ver > _ver {
					return false
				}

				break
			}
		}
	}

	// 校验时间
	if r.TimeType == TimeTypeShort {
		nowUnix := time.Now().Unix()
		if nowUnix > r.End || nowUnix < r.Start {
			return false
		}
	}

	return true
}
