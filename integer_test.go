package validator

import (
	"testing"
	"strconv"
	"fmt"

	// "github.com/goindow/toolbox"
)

// 有值，string - 浮点数字符串
func Test_Rule_IntegerValidator_String(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "integer"},
		},
	}
	obj := map[string]interface{}{ "field": "3.14" }
	message := generator(v.default_errors["integer"], "field")
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:field 必须是整数]]
	if len(e) == 0 || e[0]["field"].Error() != message {
		fail(t, "should print error(" + message + ")")
	}
}

// 有值，string - 整数字符串
func Test_Rule_IntegerValidator_OK_String(t *testing.T) { // 整数字符串
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "integer"},
		},
	}
	obj := map[string]interface{}{ "field": "3" }
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 无值 Rule.Required == flase
func Test_Rule_IntegerValidator_Required_False_Empty(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "integer"},
		},
	}
	e := v.Validate(rules, objEmpty, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 无值 Rule.Required == true
func Test_Rule_IntegerValidator_Required_True_Empty(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "integer", Required: true},
		},
	}
	message := generator(v.default_errors["required"], "field")
	e := v.Validate(rules, objEmpty, "create")
	// toolbox.Dump(e) // [map[field:field 不能为空]]
	if len(e) == 0 || e[0]["field"].Error() != message {
		fail(t, "should print error(" + message + ")")
	}
}

// 参数错误，Max 非 int
func Test_Rule_IntegerValidator_Max_TypeErr(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "integer", Max: "10"},
		},
	}
	obj := map[string]interface{}{ "field": 28 }
	defer func() {
		p := recover()
		// toolbox.Dump(p) // {field int  false 10 <nil> [] 0} attribute 'Max' should be int
		if p == nil {
			fail(t, "should panic(" + fmt.Sprint(rules["create"][0]) + " attribute 'Max' should be int)")
		}
	}()
	v.Validate(rules, obj, "create")
}

// 参数错误，Min 非 int
func Test_Rule_IntegerValidator_Min_TypeErr(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "integer", Min: "10"},
		},
	}
	obj := map[string]interface{}{ "field": 28 }
	defer func() {
		p := recover()
		// toolbox.Dump(p) // {field int  false <nil> 10 [] 0} attribute 'Min' should be int
		if p == nil {
			fail(t, "should panic(" + fmt.Sprint(rules["create"][0]) + " attribute 'Min' should be int)")
		}
	}()
	v.Validate(rules, obj, "create")
}

// 参数错误，Max < Min
func Test_Rule_IntegerValidator_Max_Min_LogicErr(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "integer", Max: 10, Min: 100},
		},
	}
	obj := map[string]interface{}{ "field": 28 }
	defer func() {
		p := recover()
		// toolbox.Dump(p) // {field int  false 10 100 [] 0} attribute 'Max' should greater than 'Min'
		if p == nil {
			fail(t, "should panic(" + fmt.Sprint(rules["create"][0]) + " attribute 'Max' should greater than 'Min')")
		}
	}()
	v.Validate(rules, obj, "create")
}

// 参数错误，Symbol < 0 && Min >= 0
func Test_Rule_IntegerValidator_Symbol_Min_LogicErr(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "integer", Symbol: -1, Min: 1},
		},
	}
	obj := map[string]interface{}{ "field": -10 }
	defer func() {
		p := recover()
		// toolbox.Dump(p) // {field int  false <nil> 1 [] -1} attribute 'Min' should less than 0 when 'Symbal' = -1
		if p == nil {
			fail(t, "should panic(" + fmt.Sprint(rules["create"][0]) + " attribute 'Min' should less than 0 when 'Symbal' = " + strconv.FormatInt(rules["create"][0].Symbol, 10))
		}
	}()
	v.Validate(rules, obj, "create")
}

// 参数错误，Symbol > 0 && Max <= 0
func Test_Rule_IntegerValidator_Symbol_Max_LogicErr(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "integer", Symbol: 1, Max: -1},
		},
	}
	obj := map[string]interface{}{ "field": 10 }
	defer func() {
		p := recover()
		// toolbox.Dump(p) // {field int  false -1 <nil> [] 1} attribute 'Max' should greater than 0 when 'Symbal' = 1
		if p == nil {
			fail(t, "should panic(" + fmt.Sprint(rules["create"][0]) + " attribute 'Max' should greater than 0 when 'Symbal' = " + strconv.FormatInt(rules["create"][0].Symbol, 10))
		}
	}()
	v.Validate(rules, obj, "create")
}

// 最大值，Max
func Test_Rule_IntegerValidator_Max(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "integer", Max: 150},
		},
	}
	obj := map[string]interface{}{ "field": 151 }
	message := generator(v.default_errors["integerMax"], "field", 150)
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:field 必须是不大于 150 的整数]]
	if len(e) == 0 || e[0]["field"].Error() != message {
		fail(t, "should print error(" + message + ")")
	}
}

func Test_Rule_IntegerValidator_Max_OK(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "integer", Max: 150},
		},
	}
	obj := map[string]interface{}{ "field": 28 }
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 最小值，Min
func Test_Rule_IntegerValidator_Min(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "integer", Min: 18},
		},
	}
	obj := map[string]interface{}{ "field": 17 }
	message := generator(v.default_errors["integerMin"], "field", 18)
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:field 必须是不小于 18 的整数]]
	if len(e) == 0 || e[0]["field"].Error() != message {
		fail(t, "should print error(" + message + ")")
	}
}

func Test_Rule_IntegerValidator_Min_OK(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "integer", Min: 18},
		},
	}
	obj := map[string]interface{}{ "field": 18 }
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 范围值，Min ~ Max
func Test_Rule_IntegerValidator_Max_Min_Range(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "integer", Max: 35, Min: 18},
		},
	}
	obj := map[string]interface{}{ "field": 51 }
	message := generator(v.default_errors["integerRange"], "field", 18, 35)
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:field 必须是介于 18 到 35 的整数]]
	if len(e) == 0 || e[0]["field"].Error() != message {
		fail(t, "should print error(" + message + ")")
	}
}

func Test_Rule_IntegerValidator_Max_Min_Range_OK(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "integer", Max: 35, Min: 18},
		},
	}
	obj := map[string]interface{}{ "field": 22 }
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 固定值，Max == Min
func Test_Rule_IntegerValidator_Max_Min_Eq(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "integer", Max: 20, Min: 20},
		},
	}
	obj := map[string]interface{}{ "field": 21 }
	message := generator(v.default_errors["equal"], "field", 20)
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:field 必须是 20]]
	if len(e) == 0 || e[0]["field"].Error() != message {
		fail(t, "should print error(" + message + ")")
	}
}

func Test_Rule_IntegerValidator_Max_Min_Eq_OK(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "integer", Max: 20, Min: 20},
		},
	}
	obj := map[string]interface{}{ "field": 20 }
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 正整数 Rule.Symbol > 0
func Test_Rule_IntegerValidator_Symbol_Positive(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "integer", Symbol: 1},
		},
	}
	obj := map[string]interface{}{ "field": 0 }
	message := generator(v.default_errors["integerPositive"], "field")
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:field 必须是正整数]]
	if len(e) == 0 || e[0]["field"].Error() != message {
		fail(t, "should print error(" + message + ")")
	}
}

func Test_Rule_IntegerValidator_Symbol_Positive_OK(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "integer", Symbol: 1},
		},
	}
	obj := map[string]interface{}{ "field": 20 }
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 负整数 Rule.Symbol < 0
func Test_Rule_IntegerValidator_Symbol_Negative(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "integer", Symbol: -1},
		},
	}
	obj := map[string]interface{}{ "field": 18 }
	message := generator(v.default_errors["integerNegative"], "field")
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:field 必须是负整数]]
	if len(e) == 0 || e[0]["field"].Error() != message {
		fail(t, "should print error(" + message + ")")
	}
}

func Test_Rule_IntegerValidator_Symbol_Negative_OK(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "integer", Symbol: -1},
		},
	}
	obj := map[string]interface{}{ "field": -10 }
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 正整数最大值，Symbol（positive） + Max
func Test_Rule_IntegerValidator_Symbol_Positive_Max(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "integer", Symbol: 1, Max: 35},
		},
	}
	obj := map[string]interface{}{ "field": 36 }
	message := generator(v.default_errors["integerPositiveMax"], "field", 35)
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:field 必须是不大于 35 的正整数]]
	if len(e) == 0 || e[0]["field"].Error() != message {
		fail(t, "should print error(" + message + ")")
	}
}

func Test_Rule_IntegerValidator_Symbol_Positive_Max_OK(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "integer", Symbol: 1, Max: 35},
		},
	}
	obj := map[string]interface{}{ "field": 35 }
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 正整数最小值，Symbol（positive） + Min
func Test_Rule_IntegerValidator_Symbol_Positive_Min(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "integer", Symbol: 1, Min: 18},
		},
	}
	obj := map[string]interface{}{ "field": 17 }
	message := generator(v.default_errors["integerPositiveMin"], "field", 18)
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:field 必须是不小于 18 的正整数]]
	if len(e) == 0 || e[0]["field"].Error() != message {
		fail(t, "should print error(" + message + ")")
	}
}

func Test_Rule_IntegerValidator_Symbol_Positive_Min_OK(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "integer", Symbol: 1, Min: 18},
		},
	}
	obj := map[string]interface{}{ "field": 18 }
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 正整数范围值，Symbol（positive） + Max + Min
func Test_Rule_IntegerValidator_Symbol_Positive_Max_Min_Range(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "integer", Symbol: 1, Min: 18, Max: 35},
		},
	}
	obj := map[string]interface{}{ "field": 17 }
	message := generator(v.default_errors["integerPositiveRange"], "field", 18, 35)
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:field 必须是介于 18 到 35 的正整数]]
	if len(e) == 0 || e[0]["field"].Error() != message {
		fail(t, "should print error(" + message + ")")
	}
}

func Test_Rule_IntegerValidator_Symbol_Positive_Max_Min_Range_OK(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "integer", Symbol: 1, Min: 18, Max: 35},
		},
	}
	obj := map[string]interface{}{ "field": 20 }
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 正整数固定值，Symbol（positive） + Max + Min, Max == Min
func Test_Rule_IntegerValidator_Symbol_Positive_Max_Min_Eq(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "integer", Symbol: 1, Min: 18, Max: 18},
		},
	}
	obj := map[string]interface{}{ "field": 17 }
	message := generator(v.default_errors["equal"], "field", 18)
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:field 必须是 18]]
	if len(e) == 0 || e[0]["field"].Error() != message {
		fail(t, "should print error(" + message + ")")
	}
}

func Test_Rule_IntegerValidator_Symbol_Positive_Max_Min_Eq_OK(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "integer", Symbol: 1, Min: 18, Max: 18},
		},
	}
	obj := map[string]interface{}{ "field": 18 }
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 负整数最大值，Symbol（negative） + Max
func Test_Rule_IntegerValidator_Symbol_Negative_Max(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "integer", Symbol: -1, Max: -100},
		},
	}
	obj := map[string]interface{}{ "field": -99 }
	message := generator(v.default_errors["integerNegativeMax"], "field", -100)
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:field 必须是不大于 -100 的负整数]]
	if len(e) == 0 || e[0]["field"].Error() != message {
		fail(t, "should print error(" + message + ")")
	}
}

func Test_Rule_IntegerValidator_Symbol_Negative_Max_OK(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "integer", Symbol: -1, Max: -100},
		},
	}
	obj := map[string]interface{}{ "field": -101 }
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 负整数最小值，Symbol（negative） + Min
func Test_Rule_IntegerValidator_Symbol_Negative_Min(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "integer", Symbol: -1, Min: -100},
		},
	}
	obj := map[string]interface{}{ "field": -101 }
	message := generator(v.default_errors["integerNegativeMin"], "field", -100)
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:field 必须是不小于 -100 的负整数]]
	if len(e) == 0 || e[0]["field"].Error() != message {
		fail(t, "should print error(" + message + ")")
	}
}

func Test_Rule_IntegerValidator_Symbol_Negative_Min_OK(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "integer", Symbol: -1, Min: -100},
		},
	}
	obj := map[string]interface{}{ "field": -99 }
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 负整数范围值，Symbol（negative） + Max + Min
func Test_Rule_IntegerValidator_Symbol_Negative_Max_Min_Range(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "integer", Symbol: -1, Min: -100, Max: -1},
		},
	}
	obj := map[string]interface{}{ "field": -101 }
	message := generator(v.default_errors["integerNegativeRange"], "field", -100, -1)
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:field 必须是介于 -100 到 -1 的负整数]]
	if len(e) == 0 || e[0]["field"].Error() != message {
		fail(t, "should print error(" + message + ")")
	}
}

func Test_Rule_IntegerValidator_Symbol_Negative_Max_Min_Range_OK(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "integer", Symbol: -1, Min: -100, Max: -1},
		},
	}
	obj := map[string]interface{}{ "field": -99 }
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 负整数固定值，Symbol（negative） + Max + Min, Max == Min
func Test_Rule_IntegerValidator_Symbol_Negative_Max_Min_Eq(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "integer", Symbol: -1, Min: -100, Max: -100},
		},
	}
	obj := map[string]interface{}{ "field": -101 }
	message := generator(v.default_errors["equal"], "field", -100)
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:field 必须是 -100]]
	if len(e) == 0 || e[0]["field"].Error() != message {
		fail(t, "should print error(" + message + ")")
	}
}

func Test_Rule_IntegerValidator_Symbol_Negative_Max_Min_Eq_OK(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "integer", Symbol: -1, Min: -100, Max: -100},
		},
	}
	obj := map[string]interface{}{ "field": -100 }
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 有值，string - 浮点数字符串，同 Test_Rule_IntegerValidator_String
// 有值，string - 整数字符串，同 Test_Rule_IntegerValidator_OK_String

// 有值，float64 - 有小数点浮点数
func Test_Rule_IntegerValidator_Float64(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "integer"},
		},
	}
	obj := map[string]interface{}{ "field": float64(3.14) }
	message := generator(v.default_errors["integer"], "field")
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:field 必须是整数]]
	if len(e) == 0 || e[0]["field"].Error() != message {
		fail(t, "should print error(" + message + ")")
	}
}

// 有值，float64 - 无小数点浮点数
func Test_Rule_IntegerValidator_OK_Float64(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "integer"},
		},
	}
	obj := map[string]interface{}{ "field": float64(3) }
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 有值，float32 - 有小数点浮点数
func Test_Rule_IntegerValidator_Float32(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "integer"},
		},
	}
	obj := map[string]interface{}{ "field": float32(3.14) }
	message := generator(v.default_errors["integer"], "field")
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:field 必须是整数]]
	if len(e) == 0 || e[0]["field"].Error() != message {
		fail(t, "should print error(" + message + ")")
	}
}

// 有值，float32 - 无小数点浮点数
func Test_Rule_IntegerValidator_OK_Float32(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "integer"},
		},
	}
	obj := map[string]interface{}{ "field": float32(3) }
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}	
}

// 有值，int64
func Test_Rule_IntegerValidator_OK_Int64(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "integer"},
		},
	}
	obj := map[string]interface{}{ "field": int64(3) }
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 有值，int32
func Test_Rule_IntegerValidator_OK_Int32(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "integer"},
		},
	}
	obj := map[string]interface{}{ "field": int32(3) }
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 有值，int16
func Test_Rule_IntegerValidator_OK_Int16(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "integer"},
		},
	}
	obj := map[string]interface{}{ "field": int16(3) }
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 有值，int8
func Test_Rule_IntegerValidator_OK_Int8(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "integer"},
		},
	}
	obj := map[string]interface{}{ "field": int8(3) }
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 有值，int
func Test_Rule_IntegerValidator_OK_Int(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "integer"},
		},
	}
	obj := map[string]interface{}{ "field": 3 }
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}