package helper

import (
	"reflect"
	//"strconv"
	"fmt"
	"strings"
)

const (
	HighChartTagxAxisCate = "cate"
	HighChartTagSeries    = "series"
)

type HighChartXAxis struct {
	Category []string `json:"category"`
}

type HighChartSeries struct {
	Name string    `json:"name"`
	Data []float64 `json:"data"`
}

type HighChartOutput struct {
	XAxis  HighChartXAxis    `json:"xAxis"`
	Series []HighChartSeries `json:"series"`
}

func (o *HighChartOutput) AppendSeries(name string, val float64) {
	find := false
	for i, _ := range o.Series {
		if o.Series[i].Name == name {
			o.Series[i].Data = append(o.Series[i].Data, val)
			find = true
			return
		}
	}
	if !find {
		newSeries := HighChartSeries{}
		newSeries.Name = name
		newSeries.Data = append(newSeries.Data, val)
		o.Series = append(o.Series, newSeries)
	}
}

type HighChartInvalidError struct {
	Type reflect.Type
}

func (e *HighChartInvalidError) Error() string {
	if e.Type == nil {
		return "highchart: Unmarshal(nil)"
	}
	if e.Type.Kind() != reflect.Ptr {
		return "highchart: Unmarshal(non-pointer " + e.Type.String() + ")"
	}
	return "highchart: Unmarshal(nil " + e.Type.String() + ")"
}

type HighChartHelp struct {
}

var HighChartHelper = new(HighChartHelp)

func (h *HighChartHelp) Unmarshal(v interface{}) (HighChartOutput, error) {
	vv := reflect.ValueOf(v)
	if vv.Kind() != reflect.Ptr || vv.IsNil() {
		return HighChartOutput{}, &HighChartInvalidError{reflect.TypeOf(v)}
	}

	rv := vv.Elem()
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		return HighChartOutput{}, &HighChartInvalidError{reflect.TypeOf(rv)}
	}

	n := rv.Len()
	out := HighChartOutput{}

	for i := 0; i < n; i++ {
		ele := rv.Index(i)
		var ele_val reflect.Value
		if ele.Kind() == reflect.Ptr {
			ele_val = ele.Elem()
		} else {
			ele_val = ele
		}

		ele_rt := ele_val.Type()
		ele_n := ele_val.NumField()

		for j := 0; j < ele_n; j++ {
			vf := ele_val.Field(j)
			tf := ele_rt.Field(j)

			tag := tf.Tag.Get("highchart")
			if tag == "-" || tag == "" || tag == "omitempty" {
				continue
			}

			tagArr := strings.SplitN(tag, ":", 2)

			if tagArr[0] == HighChartTagxAxisCate {
				out.XAxis.Category = append(out.XAxis.Category, vf.String())
			}

			if tagArr[0] == HighChartTagSeries {
				seriesName := tagArr[1]
				seriesNameVal, _ := h.checkFloat(vf)
				outPtr := &out
				outPtr.AppendSeries(seriesName, seriesNameVal)
			}
		}
	}
	//res, err := json.Marshal(out)
	return out, nil
}

func (h *HighChartHelp) checkFloat(v reflect.Value) (res float64, err error) {
	if v.Kind() == reflect.Float64 {
		return v.Float(), nil
	}

	if v.Kind() == reflect.Int || v.Kind() == reflect.Int16 || v.Kind() == reflect.Int32 || v.Kind() == reflect.Int64 {
		return float64(v.Int()), nil
	}

	return float64(0), fmt.Errorf("%+v can't parse to float64", v)
}
