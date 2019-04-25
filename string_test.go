package validator

import (
	"fmt"
	"testing"
	// "github.com/goindow/toolbox"
)

// 有值，类型错误
func Test_Rule_StringValidator(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "username", Rule: "string"},
		},
	}
	obj := map[string]interface{}{"username": 123}
	message := generator(v.default_errors["string"], "username")
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[username:username 必须是字符串]]
	if len(e) == 0 || e[0]["username"].Error() != message {
		fail(t, "should print error("+message+")")
	}
}

// 有值，类型正确
func Test_Rule_StringValidator_OK(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "username", Rule: "string"},
		},
	}
	obj := map[string]interface{}{"username": "hyb"}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 无值 Rule.Required == false
func Test_Rule_StringValidator_Required_False_Empty(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "username", Rule: "string"},
		},
	}
	e := v.Validate(rules, objEmpty, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 无值 Rule.Required == true
func Test_Rule_StringValidator_Required_True_Empty(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "username", Rule: "string", Required: true},
		},
	}
	message := generator(v.default_errors["required"], "username")
	e := v.Validate(rules, objEmpty, "create")
	// toolbox.Dump(e) // [map[username:username 不能为空]]
	if len(e) == 0 || e[0]["username"].Error() != message {
		fail(t, "should print error("+message+")")
	}
}

// 参数错误，Max 非 int
func Test_Rule_StringValidator_Max_TypeErr(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "username", Rule: "string", Max: "10"},
		},
	}
	obj := map[string]interface{}{"username": "hyb"}
	defer func() {
		p := recover()
		// toolbox.Dump(p) // {username string  false 10 <nil> [] 0} attribute 'Max' should be int
		if p == nil {
			fail(t, "should panic("+fmt.Sprint(rules["create"][0])+" attribute 'Max' should be int)")
		}
	}()
	v.Validate(rules, obj, "create")
}

// 参数错误，Min 非 int
func Test_Rule_StringValidator_Min_TypeErr(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "username", Rule: "string", Min: "10"},
		},
	}
	obj := map[string]interface{}{"username": "hyb"}
	defer func() {
		p := recover()
		// toolbox.Dump(p) // {username string  false 10 <nil> [] 0} attribute 'Min' should be int
		if p == nil {
			fail(t, "should panic("+fmt.Sprint(rules["create"][0])+" attribute 'Min' should be int)")
		}
	}()
	v.Validate(rules, obj, "create")
}

// 参数错误，Max < Min
func Test_Rule_StringValidator_Max_Min_LogicErr(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "username", Rule: "string", Max: 10, Min: 100},
		},
	}
	obj := map[string]interface{}{"username": "hyb"}
	defer func() {
		p := recover()
		// toolbox.Dump(p) // {username string  false 10 100 [] 0} attribute 'Max' should greater than 'Min'
		if p == nil {
			fail(t, "should panic("+fmt.Sprint(rules["create"][0])+" attribute 'Max' should greater than 'Min')")
		}
	}()
	v.Validate(rules, obj, "create")
}

// 范围长度，Min ~ Max
func Test_Rule_StringValidator_Max_Min_Range(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "username", Rule: "string", Max: 10, Min: 6},
		},
	}
	obj := map[string]interface{}{"username": "hyb"}
	message := generator(v.default_errors["stringLengthRange"], "username", 6, 10)
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[username:username 长度必须在 6 到 10 之间]]
	if len(e) == 0 || e[0]["username"].Error() != message {
		fail(t, "should print error("+message+")")
	}
}

func Test_Rule_StringValidator_Max_Min_Range_OK(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "username", Rule: "string", Max: 10, Min: 6},
		},
	}
	obj := map[string]interface{}{"username": "1234567"}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 固定长度，Max == Min
func Test_Rule_StringValidator_Max_Min_Eq(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "username", Rule: "string", Max: 10, Min: 10},
		},
	}
	obj := map[string]interface{}{"username": "hyb"}
	message := generator(v.default_errors["stringLengthEqual"], "username", 10)
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[username:username 长度必须是 10]]
	if len(e) == 0 || e[0]["username"].Error() != message {
		fail(t, "should print error("+message+")")
	}
}

func Test_Rule_StringValidator_Max_Min_Eq_OK(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "username", Rule: "string", Max: 10, Min: 10},
		},
	}
	obj := map[string]interface{}{"username": "1234567890"}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 最大长度，Max
func Test_Rule_StringValidator_Max(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "username", Rule: "string", Max: 10},
		},
	}
	obj := map[string]interface{}{"username": "12345678901"}
	message := generator(v.default_errors["stringLengthMax"], "username", 10)
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[username:username 长度不能超过 10]]
	if len(e) == 0 || e[0]["username"].Error() != message {
		fail(t, "should print error("+message+")")
	}
}

func Test_Rule_StringValidator_Max_OK(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "username", Rule: "string", Max: 10},
		},
	}
	obj := map[string]interface{}{"username": "1234"}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 最小长度，Min
func Test_Rule_StringValidator_Min(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "username", Rule: "string", Min: 10},
		},
	}
	obj := map[string]interface{}{"username": "hyb"}
	message := generator(v.default_errors["stringLengthMin"], "username", 10)
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[username:username 长度不能小于 10]]
	if len(e) == 0 || e[0]["username"].Error() != message {
		fail(t, "should print error("+message+")")
	}
}

func Test_Rule_StringValidator_Min_OK(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "username", Rule: "string", Min: 10},
		},
	}
	obj := map[string]interface{}{"username": "12345678901"}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}
