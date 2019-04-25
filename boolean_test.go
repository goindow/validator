package validator

import (
	"testing"

	// "github.com/goindow/toolbox"
)

// 有值，类型错误
func Test_Rule_BooleanValidator(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "boolean"},
		},
	}
	obj := map[string]interface{}{ "field": 1 }
	message := generator(v.default_errors["boolean"], "field")
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:field 必须是布尔值或布尔字符串]]
	if len(e) == 0 || e[0]["field"].Error() != message {
		fail(t, "should print error(" + message + ")")
	}
}

// 有值，bool
func Test_Rule_BooleanValidator_OK(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "boolean"},
		},
	}
	obj := map[string]interface{}{ "field": false }
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 无值 Rule.Required == false
func Test_Rule_BooleanValidator_Required_False_Empty(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "boolean"},
		},
	}
	e := v.Validate(rules, objEmpty, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 无值 Rule.Required == true
func Test_Rule_BooleanValidator_Required_True_Empty(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "boolean", Required: true},
		},
	}
	message := generator(v.default_errors["required"], "field")
	e := v.Validate(rules, objEmpty, "create")
	// toolbox.Dump(e) // [map[field:field 不能为空]]
	if len(e) == 0 || e[0]["field"].Error() != message {
		fail(t, "should print error(" + message + ")")
	}
}

// 有值，bool，同 Test_Rule_BooleanValidator_OK

// 有值，string，非布尔字符串
func Test_Rule_BooleanValidator_String(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "boolean"},
		},
	}
	obj := map[string]interface{}{ "field": "abc" }
	message := generator(v.default_errors["boolean"], "field")
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:field 必须是布尔值或布尔字符串]]
	if len(e) == 0 || e[0]["field"].Error() != message {
		fail(t, "should print error(" + message + ")")
	}
}

// 有值，string，布尔字符串
func Test_Rule_BooleanValidator_OK_String(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "boolean"},
		},
	}
	obj := map[string]interface{}{ "field": "true" }
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}