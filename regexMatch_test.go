package validator

import (
	"fmt"
	"testing"
	// "github.com/goindow/toolbox"
)

/***** regexValidator *****/

// 有值，不匹配
func Test_Rule_RegexValidator(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "regex", Pattern: "^[1][3,4,5,6,7,8,9][0-9]{9}$"},
		},
	}
	obj := map[string]interface{}{"field": "https://"}
	message := generator(v.default_errors["regex"], "field")
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:field 格式不对]]
	if len(e) == 0 || e[0]["field"].Error() != message {
		fail(t, "should print error("+message+")")
	}
}

// 有值，匹配
func Test_Rule_RegexValidator_OK(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "regex", Pattern: "^[1][3,4,5,6,7,8,9][0-9]{9}$"},
		},
	}
	obj := map[string]interface{}{"field": "15990573367"}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 有值，类型错误
func Test_Rule_RegexValidator_TypeErr(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "regex", Pattern: "^[1][3,4,5,6,7,8,9][0-9]{9}$"},
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
func Test_Rule_RegexValidator_Required_False_Empty(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "regex", Pattern: "^[1][3,4,5,6,7,8,9][0-9]{9}$"},
		},
	}
	e := v.Validate(rules, objEmpty, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 无值 Rule.Required == true
func Test_Rule_RegexValidator_Required_True_Empty(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "regex", Pattern: "^[1][3,4,5,6,7,8,9][0-9]{9}$", Required: true},
		},
	}
	message := generator(v.default_errors["required"], "field")
	e := v.Validate(rules, objEmpty, "create")
	// toolbox.Dump(e) // [map[field:field 不能为空]]
	if len(e) == 0 || e[0]["field"].Error() != message {
		fail(t, "should print error("+message+")")
	}
}

// 未传参，Rule.Pattern
func Test_Rule_RegexValidator_NotFound_Pattern(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "regex"},
		},
	}
	defer func() {
		p := recover()
		// toolbox.Dump(p) // {field regex  false 0 <nil> <nil> [] } attribute 'Pattern' not found or empty
		if p == nil {
			fail(t, "should panic("+fmt.Sprint(rules["create"][0])+" attribute 'Pattern' not found or empty)")
		}
	}()
	v.Validate(rules, objEmpty, "create")
}

/***** emailValidator *****/

// 有值，不匹配
func Test_Rule_EmailValidator(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "email"},
		},
	}
	obj := map[string]interface{}{"field": "76788424"}
	message := generator(v.default_errors["email"], "field")
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:无效的 email]]
	if len(e) == 0 || e[0]["field"].Error() != message {
		fail(t, "should print error("+message+")")
	}
}

// 有值，匹配
func Test_Rule_EmailValidator_OK(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "email"},
		},
	}
	obj := map[string]interface{}{"field": "hyb76788424@163.com"}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 有值，类型错误
func Test_Rule_EmailValidator_TypeErr(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "email"},
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
func Test_Rule_EmailValidator_Required_False_Empty(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "email"},
		},
	}
	e := v.Validate(rules, objEmpty, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 无值 Rule.Required == true
func Test_Rule_EmailValidator_Required_True_Empty(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "email", Required: true},
		},
	}
	message := generator(v.default_errors["required"], "field")
	e := v.Validate(rules, objEmpty, "create")
	// toolbox.Dump(e) // [map[field:field 不能为空]]
	if len(e) == 0 || e[0]["field"].Error() != message {
		fail(t, "should print error("+message+")")
	}
}

/***** telValidator *****/

// 有值，不匹配
func Test_Rule_TelValidator(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "tel"},
		},
	}
	obj := map[string]interface{}{"field": "7678842412"}
	message := generator(v.default_errors["tel"], "field")
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:无效的座机号]]
	if len(e) == 0 || e[0]["field"].Error() != message {
		fail(t, "should print error("+message+")")
	}
}

// 有值，匹配
func Test_Rule_TelValidator_OK(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "tel"},
		},
	}
	obj := map[string]interface{}{"field": "8518523"}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 有值，类型错误
func Test_Rule_TelValidator_TypeErr(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "tel"},
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
func Test_Rule_TelValidator_Required_False_Empty(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "tel"},
		},
	}
	e := v.Validate(rules, objEmpty, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 无值 Rule.Required == true
func Test_Rule_TelValidator_Required_True_Empty(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "tel", Required: true},
		},
	}
	message := generator(v.default_errors["required"], "field")
	e := v.Validate(rules, objEmpty, "create")
	// toolbox.Dump(e) // [map[field:field 不能为空]]
	if len(e) == 0 || e[0]["field"].Error() != message {
		fail(t, "should print error("+message+")")
	}
}

/***** mobileValidator *****/

// 有值，不匹配
func Test_Rule_MobileValidator(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "mobile"},
		},
	}
	obj := map[string]interface{}{"field": "159905733"}
	message := generator(v.default_errors["mobile"], "field")
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:field 手机号格式不对]]
	if len(e) == 0 || e[0]["field"].Error() != message {
		fail(t, "should print error("+message+")")
	}
}

// 有值，匹配
func Test_Rule_MobileValidator_OK(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "mobile"},
		},
	}
	obj := map[string]interface{}{"field": "15990573367"}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 有值，类型错误
func Test_Rule_MobileValidator_TypeErr(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "mobile"},
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
func Test_Rule_MobileValidator_Required_False_Empty(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "mobile"},
		},
	}
	e := v.Validate(rules, objEmpty, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 无值 Rule.Required == true
func Test_Rule_MobileValidator_Required_True_Empty(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "mobile", Required: true},
		},
	}
	message := generator(v.default_errors["required"], "field")
	e := v.Validate(rules, objEmpty, "create")
	// toolbox.Dump(e) // [map[field:field 不能为空]]
	if len(e) == 0 || e[0]["field"].Error() != message {
		fail(t, "should print error("+message+")")
	}
}

/***** zipcodeValidator *****/

// 有值，不匹配
func Test_Rule_ZipcodeValidator(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "zipcode"},
		},
	}
	obj := map[string]interface{}{"field": "159905733"}
	message := generator(v.default_errors["zipcode"], "field")
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:无效的邮编]]
	if len(e) == 0 || e[0]["field"].Error() != message {
		fail(t, "should print error("+message+")")
	}
}

// 有值，匹配
func Test_Rule_ZipcodeValidator_OK(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "zipcode"},
		},
	}
	obj := map[string]interface{}{"field": "333000"}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 有值，类型错误
func Test_Rule_ZipcodeValidator_TypeErr(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "zipcode"},
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
func Test_Rule_ZipcodeValidator_Required_False_Empty(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "zipcode"},
		},
	}
	e := v.Validate(rules, objEmpty, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 无值 Rule.Required == true
func Test_Rule_ZipcodeValidator_Required_True_Empty(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "zipcode", Required: true},
		},
	}
	message := generator(v.default_errors["required"], "field")
	e := v.Validate(rules, objEmpty, "create")
	// toolbox.Dump(e) // [map[field:field 不能为空]]
	if len(e) == 0 || e[0]["field"].Error() != message {
		fail(t, "should print error("+message+")")
	}
}
