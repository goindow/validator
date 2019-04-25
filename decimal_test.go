package validator

import (
	"fmt"
	"strconv"
	"testing"
	// "github.com/goindow/toolbox"
)

// 有值，string - 整数字符串
func Test_Rule_DecimalValidator_String(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "decimal"},
		},
	}
	obj := map[string]interface{}{"field": "3"}
	message := generator(v.default_errors["decimal"], "field")
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:field 必须是小数]]
	if len(e) == 0 || e[0]["field"].Error() != message {
		fail(t, "should print error("+message+")")
	}
}

// 有值，string - 浮点数字符串
func Test_Rule_DecimalValidator_OK_String(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "decimal"},
		},
	}
	obj := map[string]interface{}{"field": "3.14"}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 无值 Rule.Required == flase
func Test_Rule_DecimalValidator_Required_False_Empty(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "decimal"},
		},
	}
	e := v.Validate(rules, objEmpty, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 无值 Rule.Required == true
func Test_Rule_DecimalValidator_Required_True_Empty(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "decimal", Required: true},
		},
	}
	message := generator(v.default_errors["required"], "field")
	e := v.Validate(rules, objEmpty, "create")
	// toolbox.Dump(e) // [map[field:field 不能为空]]
	if len(e) == 0 || e[0]["field"].Error() != message {
		fail(t, "should print error("+message+")")
	}
}

// 参数错误，Max 非 float64
func Test_Rule_DecimalValidator_Max_TypeErr(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "decimal", Max: "abc"},
		},
	}
	obj := map[string]interface{}{"field": 3.14}
	defer func() {
		p := recover()
		// toolbox.Dump(p) // {field decimal  false abc <nil> [] 0} attribute 'Max' should be int or float64
		if p == nil {
			fail(t, "should panic("+fmt.Sprint(rules["create"][0])+" attribute 'Max' should be int or float64)")
		}
	}()
	v.Validate(rules, obj, "create")
}

// 参数错误，Min 非 float64
func Test_Rule_DecimalValidator_Min_TypeErr(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "decimal", Min: "abc"},
		},
	}
	obj := map[string]interface{}{"field": 3.14}
	defer func() {
		p := recover()
		// toolbox.Dump(p) // {field decimal  false abc <nil> [] 0} attribute 'Min' should be int or float64
		if p == nil {
			fail(t, "should panic("+fmt.Sprint(rules["create"][0])+" attribute 'Min' should be int or float64)")
		}
	}()
	v.Validate(rules, obj, "create")
}

// 参数错误，Max < Min
func Test_Rule_DecimalValidator_Max_Min_LogicErr(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "decimal", Max: 2.2, Min: 6},
		},
	}
	obj := map[string]interface{}{"field": 5.4}
	defer func() {
		p := recover()
		// toolbox.Dump(p) // {field decimal  false 10 100 [] 0} attribute 'Max' should greater than 'Min'
		if p == nil {
			fail(t, "should panic("+fmt.Sprint(rules["create"][0])+" attribute 'Max' should greater than 'Min')")
		}
	}()
	v.Validate(rules, obj, "create")
}

// 参数错误，Symbol < 0 && Min >= 0
func Test_Rule_DecimalValidator_Symbol_Min_LogicErr(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "decimal", Symbol: -1, Min: 3.14},
		},
	}
	obj := map[string]interface{}{"field": -3.32}
	defer func() {
		p := recover()
		// toolbox.Dump(p) // {field decimal  false <nil> 3.14 [] -1} attribute 'Min' should less than 0 when 'Symbal' = -1
		if p == nil {
			fail(t, "should panic("+fmt.Sprint(rules["create"][0])+" attribute 'Min' should less than 0 when 'Symbal' = "+strconv.FormatInt(rules["create"][0].Symbol, 10))
		}
	}()
	v.Validate(rules, obj, "create")
}

// 参数错误，Symbol > 0 && Max <= 0
func Test_Rule_DecimalValidator_Symbol_Max_LogicErr(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "decimal", Symbol: 1, Max: -3},
		},
	}
	obj := map[string]interface{}{"field": 2.2}
	defer func() {
		p := recover()
		// toolbox.Dump(p) // {field decimal  false -3 <nil> [] 1} attribute 'Max' should greater than 0 when 'Symbal' = 1
		if p == nil {
			fail(t, "should panic("+fmt.Sprint(rules["create"][0])+" attribute 'Min' should less than 0 when 'Symbal' = "+strconv.FormatInt(rules["create"][0].Symbol, 10))
		}
	}()
	v.Validate(rules, obj, "create")
}

// 最大值，Max
func Test_Rule_DecimalValidator_Max(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "decimal", Max: 10},
		},
	}
	obj := map[string]interface{}{"field": 22.2}
	message := generator(v.default_errors["decimalMax"], "field", 10)
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:field 必须是不大于 10 的小数]]
	if len(e) == 0 || e[0]["field"].Error() != message {
		fail(t, "should print error("+message+")")
	}
}

func Test_Rule_DecimalValidator_Max_OK_Int(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "decimal", Max: 10},
		},
	}
	obj := map[string]interface{}{"field": 9.98}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

func Test_Rule_DecimalValidator_Max_OK_Float64(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "decimal", Max: 9.99},
		},
	}
	obj := map[string]interface{}{"field": 9.98}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 最小值，Min
func Test_Rule_DecimalValidator_Min(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "decimal", Min: 10},
		},
	}
	obj := map[string]interface{}{"field": 3.14}
	message := generator(v.default_errors["decimalMin"], "field", 10)
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:field 必须是不小于 10 的小数]]
	if len(e) == 0 || e[0]["field"].Error() != message {
		fail(t, "should print error("+message+")")
	}
}

func Test_Rule_DecimalValidator_Min_OK_Int(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "decimal", Min: 10},
		},
	}
	obj := map[string]interface{}{"field": 10.35}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

func Test_Rule_DecimalValidator_Min_OK_Float64(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "decimal", Min: 10.35},
		},
	}
	obj := map[string]interface{}{"field": 12.3}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 范围值，Min ~ Max
func Test_Rule_DecimalValidator_Max_Min_Range(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "decimal", Min: 10, Max: 11},
		},
	}
	obj := map[string]interface{}{"field": 11.11}
	message := generator(v.default_errors["decimalRange"], "field", 10, 11)
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:field 必须是介于 10 到 11 的小数]]
	if len(e) == 0 || e[0]["field"].Error() != message {
		fail(t, "should print error("+message+")")
	}
}

func Test_Rule_DecimalValidator_Max_Min_Range_OK(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "decimal", Min: 10, Max: 11},
		},
	}
	obj := map[string]interface{}{"field": 10.88}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 固定值，Max == Min
func Test_Rule_DecimalValidator_Max_Min_Eq(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "decimal", Min: 10.88, Max: 10.88},
		},
	}
	obj := map[string]interface{}{"field": 11.11}
	message := generator(v.default_errors["equal"], "field", 10.88)
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:field 必须是 10.88]]
	if len(e) == 0 || e[0]["field"].Error() != message {
		fail(t, "should print error("+message+")")
	}
}

func Test_Rule_DecimalValidator_Max_Min_Eq_OK(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "decimal", Min: 10.88, Max: 10.88},
		},
	}
	obj := map[string]interface{}{"field": 10.88}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 正浮点数 Rule.Symbol > 0
func Test_Rule_DecimalValidator_Symbol_Positive(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "decimal", Symbol: 1},
		},
	}
	obj := map[string]interface{}{"field": -3.14}
	message := generator(v.default_errors["decimalPositive"], "field")
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:field 必须是正小数]]
	if len(e) == 0 || e[0]["field"].Error() != message {
		fail(t, "should print error("+message+")")
	}
}

func Test_Rule_DecimalValidator_Symbol_Positive_OK(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "decimal", Symbol: 1},
		},
	}
	obj := map[string]interface{}{"field": 3.14}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 负浮点数 Rule.Symbol < 0
func Test_Rule_DecimalValidator_Symbol_Negative(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "decimal", Symbol: -1},
		},
	}
	obj := map[string]interface{}{"field": 3.14}
	message := generator(v.default_errors["decimalNegative"], "field")
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:field 必须是负小数]]
	if len(e) == 0 || e[0]["field"].Error() != message {
		fail(t, "should print error("+message+")")
	}
}

func Test_Rule_DecimalValidator_Symbol_Negative_OK(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "decimal", Symbol: -1},
		},
	}
	obj := map[string]interface{}{"field": -3.14}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 正浮点数最大值，Symbol（positive） + Max
func Test_Rule_DecimalValidator_Symbol_Positive_Max(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "decimal", Symbol: 1, Max: 3.5},
		},
	}
	obj := map[string]interface{}{"field": 3.77}
	message := generator(v.default_errors["decimalPositiveMax"], "field", 3.5)
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:field 必须是不大于 3.5 的正小数]]
	if len(e) == 0 || e[0]["field"].Error() != message {
		fail(t, "should print error("+message+")")
	}
}

func Test_Rule_DecimalValidator_Symbol_Positive_Max_OK(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "decimal", Symbol: 1, Max: 3.5},
		},
	}
	obj := map[string]interface{}{"field": 3.14}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 正浮点数最小值，Symbol（positive） + Min
func Test_Rule_DecimalValidator_Symbol_Positive_Min(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "decimal", Symbol: 1, Min: 3.5},
		},
	}
	obj := map[string]interface{}{"field": 3.14}
	message := generator(v.default_errors["decimalPositiveMin"], "field", 3.5)
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:field 必须是不小于 3.5 的正小数]]
	if len(e) == 0 || e[0]["field"].Error() != message {
		fail(t, "should print error("+message+")")
	}
}

func Test_Rule_DecimalValidator_Symbol_Positive_Min_OK(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "decimal", Symbol: 1, Min: 3.5},
		},
	}
	obj := map[string]interface{}{"field": 3.66}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 正浮点数范围值，Symbol（positive） + Max + Min
func Test_Rule_DecimalValidator_Symbol_Positive_Max_Min_Range(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "decimal", Symbol: 1, Min: 3.5, Max: 4},
		},
	}
	obj := map[string]interface{}{"field": 3.33}
	message := generator(v.default_errors["decimalPositiveRange"], "field", 3.5, 4)
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:field 必须是介于 3.5 到 4 的正小数]]
	if len(e) == 0 || e[0]["field"].Error() != message {
		fail(t, "should print error("+message+")")
	}
}

func Test_Rule_DecimalValidator_Symbol_Positive_Max_Min_Range_OK(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "decimal", Symbol: 1, Min: 3.5, Max: 4},
		},
	}
	obj := map[string]interface{}{"field": 3.66}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 正浮点数固定值，Symbol（positive） + Max + Min, Max == Min
func Test_Rule_DecimalValidator_Symbol_Positive_Max_Min_Eq(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "decimal", Symbol: 1, Min: 3.5, Max: 3.5},
		},
	}
	obj := map[string]interface{}{"field": 3.4}
	message := generator(v.default_errors["equal"], "field", 3.5)
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:field 必须是 3.5]]
	if len(e) == 0 || e[0]["field"].Error() != message {
		fail(t, "should print error("+message+")")
	}
}

func Test_Rule_DecimalValidator_Symbol_Positive_Max_Min_Eq_OK(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "decimal", Symbol: 1, Min: 3.5, Max: 3.5},
		},
	}
	obj := map[string]interface{}{"field": 3.5}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 负浮点数最大值，Symbol（negative） + Max
func Test_Rule_DecimalValidator_Symbol_Negative_Max(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "decimal", Symbol: -1, Max: -3.14},
		},
	}
	obj := map[string]interface{}{"field": -1.1}
	message := generator(v.default_errors["decimalNegativeMax"], "field", -3.14)
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:field 必须是不大于 -3.14 的负小数]]
	if len(e) == 0 || e[0]["field"].Error() != message {
		fail(t, "should print error("+message+")")
	}
}

func Test_Rule_DecimalValidator_Symbol_Negative_Max_OK(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "decimal", Symbol: -1, Max: -3.14},
		},
	}
	obj := map[string]interface{}{"field": -3.34}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 负浮点数最小值，Symbol（negative） + Min
func Test_Rule_DecimalValidator_Symbol_Negative_Min(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "decimal", Symbol: -1, Min: -3.14},
		},
	}
	obj := map[string]interface{}{"field": -3.5}
	message := generator(v.default_errors["decimalNegativeMin"], "field", -3.14)
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:field 必须是不小于 -3.14 的负小数]]
	if len(e) == 0 || e[0]["field"].Error() != message {
		fail(t, "should print error("+message+")")
	}
}

func Test_Rule_DecimalValidator_Symbol_Negative_Min_OK(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "decimal", Symbol: -1, Min: -3.14},
		},
	}
	obj := map[string]interface{}{"field": -1.23}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 负浮点数范围值，Symbol（negative） + Max + Min
func Test_Rule_DecimalValidator_Symbol_Negative_Max_Min_Range(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "decimal", Symbol: -1, Min: -3.14, Max: -1},
		},
	}
	obj := map[string]interface{}{"field": -3.5}
	message := generator(v.default_errors["decimalNegativeRange"], "field", -3.14, -1)
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:field 必须是介于 -3.14 到 -1 的负小数]]
	if len(e) == 0 || e[0]["field"].Error() != message {
		fail(t, "should print error("+message+")")
	}
}

func Test_Rule_DecimalValidator_Symbol_Negative_Max_Min_Range_OK(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "decimal", Symbol: -1, Min: -3.14, Max: -1},
		},
	}
	obj := map[string]interface{}{"field": -2.22}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 负浮点数固定值，Symbol（negative） + Max + Min, Max == Min
func Test_Rule_DecimalValidator_Symbol_Negative_Max_Min_Eq(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "decimal", Symbol: -1, Min: -3.14, Max: -3.14},
		},
	}
	obj := map[string]interface{}{"field": -3.5}
	message := generator(v.default_errors["equal"], "field", -3.14)
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:field 必须是 -3.14]]
	if len(e) == 0 || e[0]["field"].Error() != message {
		fail(t, "should print error("+message+")")
	}
}

func Test_Rule_DecimalValidator_Symbol_Negative_Max_Min_Eq_OK(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "decimal", Symbol: -1, Min: -3.14, Max: -3.14},
		},
	}
	obj := map[string]interface{}{"field": -3.14}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}

}

// 有值，string - 整数字符串，同 Test_Rule_DecimalValidator_String
// 有值，string - 浮点数字付串，同 Test_Rule_DecimalValidator_OK_String

// 有值，float64 - 无小数点浮点数
func Test_Rule_DecimalValidator_Float64(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "decimal"},
		},
	}
	obj := map[string]interface{}{"field": float64(33)}
	message := generator(v.default_errors["decimal"], "field")
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:field 必须是小数]]
	if len(e) == 0 || e[0]["field"].Error() != message {
		fail(t, "should print error("+message+")")
	}
}

// 有值，float64 - 有小数点浮点数
func Test_Rule_DecimalValidator_OK_Float64(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "decimal"},
		},
	}
	obj := map[string]interface{}{"field": 3.14}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 有值，float32 - 无小数点浮点数
func Test_Rule_DecimalValidator_Float32(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "decimal"},
		},
	}
	obj := map[string]interface{}{"field": float32(33)}
	message := generator(v.default_errors["decimal"], "field")
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:field 必须是小数]]
	if len(e) == 0 || e[0]["field"].Error() != message {
		fail(t, "should print error("+message+")")
	}
}

// 有值，float32 - 有小数点浮点数
func Test_Rule_DecimalValidator_OK_Float32(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "decimal"},
		},
	}
	obj := map[string]interface{}{"field": float32(3.14)}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}
