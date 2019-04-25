package validator

import (
	"testing"
	// "github.com/goindow/toolbox"
)

// 有值，不匹配
func Test_Rule_IpValidator(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "ip"},
		},
	}
	obj := map[string]interface{}{"field": "111.222.333.444"}
	message := generator(v.default_errors["ip"], "field")
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:无效的 ip]]
	if len(e) == 0 || e[0]["field"].Error() != message {
		fail(t, "should print error("+message+")")
	}
}

// 有值，匹配
func Test_Rule_IpValidator_OK_V4(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "ip"},
		},
	}
	obj := map[string]interface{}{"field": "127.0.0.1"}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

func Test_Rule_IpValidator_OK_V6(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "ip"},
		},
	}
	obj := map[string]interface{}{"field": "2001:4860:0:2001::68"}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 有值，类型错误
func Test_Rule_IpValidator_TypeErr(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "ip"},
		},
	}
	obj := map[string]interface{}{"field": 15990573367}
	message := generator(v.default_errors["string"], "field")
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:field 必须是字符串]]
	if len(e) == 0 || e[0]["field"].Error() != message {
		fail(t, "should print error("+message+")")
	}
}

// 无值 Rule.Required == false
func Test_Rule_IpValidator_Required_False_Empty(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "ip"},
		},
	}
	e := v.Validate(rules, objEmpty, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 无值 Rule.Required == true
func Test_Rule_IpValidator_Required_True_Empty(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "ip", Required: true},
		},
	}
	message := generator(v.default_errors["required"], "field")
	e := v.Validate(rules, objEmpty, "create")
	// toolbox.Dump(e) // [map[field:field 不能为空]]
	if len(e) == 0 || e[0]["field"].Error() != message {
		fail(t, "should print error("+message+")")
	}
}
