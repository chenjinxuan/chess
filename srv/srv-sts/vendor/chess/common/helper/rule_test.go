package helper

import (
	"testing"
)

func TestRule(t *testing.T) {
	ctx := `{"from": ["android","android-ky"],"clients_type": "black","clients": ["android_300","android-ky_300"],"time_type": "long","start": 0,"end": 0}`
	rule := NewRule(ctx)
	t.Log(rule.IsMeet("android", 300))

	ctx = `{"from": ["android","android-ky"],"clients_type": "gte","clients": ["android_300","android-ky_300"],"time_type": "long","start": 0,"end": 0}`
	rule = NewRule(ctx)
	t.Log(rule.IsMeet("android", 299))

}
