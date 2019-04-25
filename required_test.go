package validator

import (
	"testing"
	// "github.com/goindow/toolbox"
)

// 无值
func Test_Rule_RequiredValidator(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "username", Rule: "required"},
		},
	}
	message := generator(v.default_errors["required"], "username")
	e := v.Validate(rules, objEmpty, "create")
	// toolbox.Dump(e) // [map[username:username 不能为空]]
	if len(e) == 0 || e[0]["username"].Error() != message {
		fail(t, "should print error("+message+")")
	}
}

// 有值
func Test_Rule_RequiredValidator_OK(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "username", Rule: "required"},
		},
	}
	obj := map[string]interface{}{"username": 123}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}
