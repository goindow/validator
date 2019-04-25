package validator

import (
	"testing"
	"strings"
	"fmt"

	// "github.com/goindow/toolbox"
)

// 有值，不在 Rule.Enum 内
func Test_Rule_InValidator_String(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "gender", Rule: "in", Enum: []string{"male", "female"}},
		},
	}
	obj := map[string]interface{}{ "gender": "unknown" }
	message := generator(v.default_errors["in"], "gender", "[" +strings.Join(rules["create"][0].Enum, "、") + "]")
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[gender:gender 只能是 [male、female] 中的一个]]
	if len(e) == 0 || e[0]["gender"].Error() != message {
		fail(t, "should print error(" + message + ")")
	}
}

// 有值，在 Rule.Enum 内
func Test_Rule_InValidator_OK_String(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "gender", Rule: "in", Enum: []string{"male", "female"}},
		},
	}
	obj := map[string]interface{}{ "gender": "male" }
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 无值 Rule.Required == false
func Test_Rule_InValidator_Required_False_Empty(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "gender", Rule: "in", Enum: []string{"male", "female"}},
		},
	}
	e := v.Validate(rules, objEmpty, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 无值 Rule.Required == true
func Test_Rule_InValidator_Required_True_Empty(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "gender", Rule: "in", Enum: []string{"male", "female"}, Required: true},
		},
	}
	message := generator(v.default_errors["required"], "gender")
	e := v.Validate(rules, objEmpty, "create")
	// toolbox.Dump(e) // [map[gender:gender 不能为空]]
	if len(e) == 0 || e[0]["gender"].Error() != message {
		fail(t, "should print error(" + message + ")")
	}
}

// 未传参，Rule.Enum
func Test_Rule_InValidator_NotFound_Enmu(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "gender", Rule: "in"},
		},
	}
	defer func() {
		p := recover()
		// toolbox.Dump(p) // {gender in  false <nil> <nil> [] 0} attribute 'Enum' not found or empty
		if p == nil {
			fail(t, "should panic(" + fmt.Sprint(rules["create"][0]) + " attribute 'Enum' not found or empty)")
		}
	}()
	v.Validate(rules, objEmpty, "create")
}

// 有值，类型错误
func Test_Rule_InValidator_TypeErr(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "gender", Rule: "in", Enum: []string{"male", "female"}},
		},
	}
	obj := map[string]interface{}{ "gender": nil }
	message := generator(v.default_errors["inValid"], "gender")
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[gender:gender 必须是字符串、数字、布尔中的一种]]
	if len(e) == 0 || e[0]["gender"].Error() != message {
		fail(t, "should print error(" + message + ")")
	}
}

// 有值，string，同 Test_Rule_InValidator_OK_String

// 有值，float64
func Test_Rule_InValidator_OK_Float64(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "in", Enum: []string{"3.1415926", "3.14", "true", "10", "a"}},
		},
	}
	obj := map[string]interface{}{ "field": 3.1415926 }
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 有值，float32
func Test_Rule_InValidator_OK_Float32(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "in", Enum: []string{"3.1415926", "3.14", "true", "10", "a"}},
		},
	}
	obj := map[string]interface{}{ "field": float32(3.14) }
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 有值，int64
func Test_Rule_InValidator_OK_Int64(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "in", Enum: []string{"3.1415926", "3.14", "true", "10", "a"}},
		},
	}
	obj := map[string]interface{}{ "field": int64(10) }
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 有值，int32
func Test_Rule_InValidator_OK_Int32(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "in", Enum: []string{"3.1415926", "3.14", "true", "10", "a"}},
		},
	}
	obj := map[string]interface{}{ "field": int32(10) }
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 有值，int16
func Test_Rule_InValidator_OK_Int16(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "in", Enum: []string{"3.1415926", "3.14", "true", "10", "a"}},
		},
	}
	obj := map[string]interface{}{ "field": int16(10) }
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 有值，int8
func Test_Rule_InValidator_OK_Int8(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "in", Enum: []string{"3.1415926", "3.14", "true", "10", "a"}},
		},
	}
	obj := map[string]interface{}{ "field": int8(10) }
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 有值，int
func Test_Rule_InValidator_OK_Int(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "in", Enum: []string{"3.1415926", "3.14", "true", "10", "a"}},
		},
	}
	obj := map[string]interface{}{ "field": 10 }
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 有值，bool
func Test_Rule_InValidator_OK_Bool(t *testing.T) {
	rules := Rules{
		"create" : {
			{Attr: "field", Rule: "in", Enum: []string{"3.1415926", "3.14", "true", "10", "a"}},
		},
	}
	obj := map[string]interface{}{ "field": true }
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}