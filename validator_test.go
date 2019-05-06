package validator

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	// "github.com/goindow/toolbox"
)

const (
	ZH_CN = "zh_cn"
	EN_US = "EN_US"
)

var (
	v          *validator
	rulesEmpty = Rules{}
	objEmpty   = map[string]interface{}{}
)

func init() {
	v = New()
	// toolbox.Dump(v)
}

func generator(message, new string, placeholder ...interface{}) string {
	e := strings.Replace(message, "{label}", new, -1)
	return fmt.Sprintf(e, placeholder...)
}

func fail(t *testing.T, s string) {
	echo(t, s, 1)
}

func ok(t *testing.T, s string) {
	echo(t, s, 2)
}

func echo(t *testing.T, s string, level uint) {
	switch level {
	case 1:
		t.Error("[fail] " + s)
	case 2:
		t.Log("[ok] " + s)
	}
}

/***** Lang() *****/

// 切换语言包
func Test_Lang(t *testing.T) {
	vEnglish := New().Lang(EN_US) // 重新实例化，避免语言对其他用例造成影响
	var rules = Rules{
		"create": {
			{Attr: "hobby", Rule: "required"},
		},
	}
	message := generator(vEnglish.default_errors["required"], "hobby")
	e := vEnglish.Validate(rules, objEmpty, "create")
	// toolbox.Dump(e) // [map[hobby:can not be empty]]
	if len(e) == 0 || e[0]["hobby"] != message {
		fail(t, "should print errors("+message+")")
	}
}

// 未定义语言包
func Test_Lang_Undefiend_Lang(t *testing.T) {
	LANG := "undefined_lang"
	defer func() {
		p := recover()
		// toolbox.Dump(p) // undefined_lang unsupport language
		if p == nil {
			fail(t, "should panic("+LANG+" unsupport language)")
		}
	}()
	v.Lang(LANG)
}

/***** Validator() *****/

// 未定义场景
func Test_Valiadte_Undefined_Scence(t *testing.T) {
	var SCENCE Scence = "undefinde_scence"
	defer func() {
		p := recover()
		// toolbox.Dump(p) // undefinde_scence scence undefined
		if p == nil {
			fail(t, "should panic("+string(SCENCE)+" scence undefined)")
		}
	}()
	v.Validate(rulesEmpty, objEmpty, SCENCE)
}

/***** dispatch() *****/

// 未定义验证器 Rule.Rule
func Test_Dispatch_Undefined_Rule(t *testing.T) {
	RULE := "undefined_rule"
	var rules = Rules{
		"create": {
			{Attr: "username", Rule: RULE},
		},
	}
	defer func() {
		p := recover()
		// toolbox.Dump(p) // undefined_rule validator undefined
		if p == nil {
			fail(t, "should panic("+RULE+" validator undefined)")
		}
	}()
	v.Validate(rules, objEmpty, "create")
}

// 未传参验证器 Rule.Rule
func Test_Dispatch_NotFound_Rule(t *testing.T) {
	var rules = Rules{
		"create": {
			{Attr: "username"},
		},
	}
	defer func() {
		p := recover()
		// toolbox.Dump(p) // {username   false <nil> <nil> [] 0} attribute 'Rule' not found
		if p == nil {
			fail(t, "should panic("+fmt.Sprint(rules["create"][0])+" attribute 'Rule' not found)")
		}
	}()
	v.Validate(rules, objEmpty, "create")
}

// 未传参须验证属性 Rule.Attr
func Test_Dispatch_NotFound_Attr(t *testing.T) {
	var rules = Rules{
		"create": {
			{Rule: "required"},
		},
	}
	defer func() {
		p := recover()
		// toolbox.Dump(p) // {<nil> required  false <nil> <nil> [] 0} attribute 'Attr' not found
		if p == nil {
			fail(t, "should panic("+fmt.Sprint(rules["create"][0])+" attribute 'Attr' not found)")
		}
	}()
	v.Validate(rules, objEmpty, "create")
}

// 参数类型错误 Rule.Attr
func Test_Dispatch_TypeErr_Attr(t *testing.T) {
	var rules = Rules{
		"create": {
			{Attr: 1234, Rule: "required"},
		},
	}
	defer func() {
		p := recover()
		// toolbox.Dump(p) // attribute 'Attr' should be 'string' or '[]string'
		if p == nil {
			fail(t, "should panic(attribute 'Attr' should be 'string' or '[]string')")
		}
	}()
	v.Validate(rules, objEmpty, "create")
}

/***** generator() *****/

// 内置错误信息
func Test_Generator_ErrorInfo(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "username", Rule: "required"},
		},
	}
	message := generator(v.default_errors["required"], "username")
	e := v.Validate(rules, objEmpty, "create")
	// toolbox.Dump(e) // [map[username:不能为空]]
	if len(e) == 0 || e[0]["username"] != message {
		fail(t, "should print error("+message+")")
	}
}

// 自定义错误信息
func Test_Generator_ErrorInfo_Custom(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "username", Rule: "required", Message: "用户名不能为空"},
		},
	}
	e := v.Validate(rules, objEmpty, "create")
	// toolbox.Dump(e) // [map[username:用户名不能为空]]
	if len(e) == 0 || e[0]["username"] != rules["create"][0].Message {
		fail(t, "should print error("+rules["create"][0].Message+")")
	}
}

/***** AddValidator() *****/

// 规则名已存在
func Test_AddValidator_Name_Already_Exists(t *testing.T) {
	defer func() {
		p := recover()
		// toolbox.Dump(p) // validator named 'mobile' already exists
		if p == nil {
			fail(t, "validator named 'mobile' already exists")
		}
	}()
	v.AddValidator("mobile", func(attr string, rule Rule, obj M) E {
		return nil
	})
}

// 自定义验证规则
func add() {
	name := "one"
	if _, ok := v.validators[name]; !ok {
		v.AddValidator(name, func(attr string, rule Rule, obj M) E {
			if _, ok := obj[attr]; !ok {
				if !rule.Required { // 允许为空
					return nil
				}
				return v.generator("required", attr, rule)
			}
			if obj[attr] != 1 {
				return E{attr: rule.Message}
			}
			return nil
		})
	}
}

// 使用自定义验证规则
func Test_AddValidator(t *testing.T) {
	add()
	// toolbox.Dump(v.validators)
	message := "必须等于一"
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "one", Message: message},
		},
	}
	obj := map[string]interface{}{"field": 2}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:必须等于一]]
	if len(e) == 0 || e[0]["field"] != message {
		fail(t, "should print error("+message+")")
	}
}

func Test_AddValidator_OK(t *testing.T) {
	add()
	// toolbox.Dump(v.validators)
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "one", Message: "必须等于一"},
		},
	}
	obj := map[string]interface{}{"field": 1}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

/***** funcValidator *****/

var f F = func(attr string, rule Rule, obj M) E {
	if obj["password"] != obj["rpassword"] {
		return E{attr: "两次输入不一致"}
	}
	return nil
}

// 有值，验证失败
func Test_Rule_FuncValidator(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "password", Rule: "func", Func: f},
		},
	}
	obj := map[string]interface{}{"password": 123, "rpassword": 123456}
	message := "两次输入不一致"
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[password:两次输入不一致]]
	if len(e) == 0 || e[0]["password"] != message {
		fail(t, "should print error("+message+")")
	}
}

// 有值，验证通过
func Test_Rule_FuncValidator_OK(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "password", Rule: "func", Func: f},
		},
	}
	obj := map[string]interface{}{"password": 123, "rpassword": 123}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 无值 Rule.Required == false
func Test_Rule_FuncValidator_Required_False_Empty(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "func", Func: f},
		},
	}
	e := v.Validate(rules, objEmpty, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 无值 Rule.Required == true
func Test_Rule_FuncValidator_Required_True_Empty(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "func", Func: f, Required: true},
		},
	}
	message := generator(v.default_errors["required"], "field")
	e := v.Validate(rules, objEmpty, "create")
	// toolbox.Dump(e) // [map[field:不能为空]]
	if len(e) == 0 || e[0]["field"] != message {
		fail(t, "should print error("+message+")")
	}
}

// 未传参 Rule.Func
func Test_Rule_FuncValidator_NotFound_Func(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "func"},
		},
	}
	defer func() {
		p := recover()
		// toolbox.Dump(p) // {field func  false 0 <nil> <nil> []  <nil>} attribute 'Func' not found or empty
		if p == nil {
			fail(t, "should panic("+fmt.Sprint(rules["create"][0])+" attribute 'Func' not found or empty)")
		}
	}()
	v.Validate(rules, objEmpty, "create")
}

/***** requiredValidator *****/

// 无值
func Test_Rule_RequiredValidator(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "username", Rule: "required"},
		},
	}
	message := generator(v.default_errors["required"], "username")
	e := v.Validate(rules, objEmpty, "create")
	// toolbox.Dump(e) // [map[username:不能为空]]
	if len(e) == 0 || e[0]["username"] != message {
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

/***** inValidator *****/

// 有值，不在 Rule.Enum 内
func Test_Rule_InValidator_String(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "gender", Rule: "in", Enum: []string{"male", "female"}},
		},
	}
	obj := map[string]interface{}{"gender": "unknown"}
	message := generator(v.default_errors["in"], "gender", "["+strings.Join(rules["create"][0].Enum, "、")+"]")
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[gender:只能是 [male、female] 中的一个]]
	if len(e) == 0 || e[0]["gender"] != message {
		fail(t, "should print error("+message+")")
	}
}

// 有值，在 Rule.Enum 内
func Test_Rule_InValidator_OK_String(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "gender", Rule: "in", Enum: []string{"male", "female"}},
		},
	}
	obj := map[string]interface{}{"gender": "male"}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 无值 Rule.Required == false
func Test_Rule_InValidator_Required_False_Empty(t *testing.T) {
	rules := Rules{
		"create": {
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
		"create": {
			{Attr: "gender", Rule: "in", Enum: []string{"male", "female"}, Required: true},
		},
	}
	message := generator(v.default_errors["required"], "gender")
	e := v.Validate(rules, objEmpty, "create")
	// toolbox.Dump(e) // [map[gender:不能为空]]
	if len(e) == 0 || e[0]["gender"] != message {
		fail(t, "should print error("+message+")")
	}
}

// 未传参，Rule.Enum
func Test_Rule_InValidator_NotFound_Enmu(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "gender", Rule: "in"},
		},
	}
	defer func() {
		p := recover()
		// toolbox.Dump(p) // {gender in  false <nil> <nil> [] 0} attribute 'Enum' not found or empty
		if p == nil {
			fail(t, "should panic("+fmt.Sprint(rules["create"][0])+" attribute 'Enum' not found or empty)")
		}
	}()
	v.Validate(rules, objEmpty, "create")
}

// 有值，类型错误
func Test_Rule_InValidator_TypeErr(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "gender", Rule: "in", Enum: []string{"male", "female"}},
		},
	}
	obj := map[string]interface{}{"gender": nil}
	message := generator(v.default_errors["inValid"], "gender")
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[gender:必须是字符串、数字、布尔中的一种]]
	if len(e) == 0 || e[0]["gender"] != message {
		fail(t, "should print error("+message+")")
	}
}

// 有值，string，同 Test_Rule_InValidator_OK_String

// 有值，float64
func Test_Rule_InValidator_OK_Float64(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "in", Enum: []string{"3.1415926", "3.14", "true", "10", "a"}},
		},
	}
	obj := map[string]interface{}{"field": 3.1415926}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 有值，float32
func Test_Rule_InValidator_OK_Float32(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "in", Enum: []string{"3.1415926", "3.14", "true", "10", "a"}},
		},
	}
	obj := map[string]interface{}{"field": float32(3.14)}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 有值，int64
func Test_Rule_InValidator_OK_Int64(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "in", Enum: []string{"3.1415926", "3.14", "true", "10", "a"}},
		},
	}
	obj := map[string]interface{}{"field": int64(10)}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 有值，int32
func Test_Rule_InValidator_OK_Int32(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "in", Enum: []string{"3.1415926", "3.14", "true", "10", "a"}},
		},
	}
	obj := map[string]interface{}{"field": int32(10)}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 有值，int16
func Test_Rule_InValidator_OK_Int16(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "in", Enum: []string{"3.1415926", "3.14", "true", "10", "a"}},
		},
	}
	obj := map[string]interface{}{"field": int16(10)}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 有值，int8
func Test_Rule_InValidator_OK_Int8(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "in", Enum: []string{"3.1415926", "3.14", "true", "10", "a"}},
		},
	}
	obj := map[string]interface{}{"field": int8(10)}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 有值，int
func Test_Rule_InValidator_OK_Int(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "in", Enum: []string{"3.1415926", "3.14", "true", "10", "a"}},
		},
	}
	obj := map[string]interface{}{"field": 10}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 有值，bool
func Test_Rule_InValidator_OK_Bool(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "in", Enum: []string{"3.1415926", "3.14", "true", "10", "a"}},
		},
	}
	obj := map[string]interface{}{"field": true}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

/***** stringValidator *****/

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
	// toolbox.Dump(e) // [map[username:必须是字符串]]
	if len(e) == 0 || e[0]["username"] != message {
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
	// toolbox.Dump(e) // [map[username:不能为空]]
	if len(e) == 0 || e[0]["username"] != message {
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
	// toolbox.Dump(e) // [map[username:长度必须在 6 到 10 之间]]
	if len(e) == 0 || e[0]["username"] != message {
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
	// toolbox.Dump(e) // [map[username:长度必须是 10]]
	if len(e) == 0 || e[0]["username"] != message {
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
	// toolbox.Dump(e) // [map[username:长度不能超过 10]]
	if len(e) == 0 || e[0]["username"] != message {
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
	// toolbox.Dump(e) // [map[username:长度不能小于 10]]
	if len(e) == 0 || e[0]["username"] != message {
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

/***** integerValidator *****/

// 有值，string - 浮点数字符串
func Test_Rule_IntegerValidator_String(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "integer"},
		},
	}
	obj := map[string]interface{}{"field": "3.14"}
	message := generator(v.default_errors["integer"], "field")
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:必须是整数]]
	if len(e) == 0 || e[0]["field"] != message {
		fail(t, "should print error("+message+")")
	}
}

// 有值，string - 整数字符串
func Test_Rule_IntegerValidator_OK_String(t *testing.T) { // 整数字符串
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "integer"},
		},
	}
	obj := map[string]interface{}{"field": "3"}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 无值 Rule.Required == flase
func Test_Rule_IntegerValidator_Required_False_Empty(t *testing.T) {
	rules := Rules{
		"create": {
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
		"create": {
			{Attr: "field", Rule: "integer", Required: true},
		},
	}
	message := generator(v.default_errors["required"], "field")
	e := v.Validate(rules, objEmpty, "create")
	// toolbox.Dump(e) // [map[field:不能为空]]
	if len(e) == 0 || e[0]["field"] != message {
		fail(t, "should print error("+message+")")
	}
}

// 参数错误，Max 非 int
func Test_Rule_IntegerValidator_Max_TypeErr(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "integer", Max: "10"},
		},
	}
	obj := map[string]interface{}{"field": 28}
	defer func() {
		p := recover()
		// toolbox.Dump(p) // {field int  false 10 <nil> [] 0} attribute 'Max' should be int
		if p == nil {
			fail(t, "should panic("+fmt.Sprint(rules["create"][0])+" attribute 'Max' should be int)")
		}
	}()
	v.Validate(rules, obj, "create")
}

// 参数错误，Min 非 int
func Test_Rule_IntegerValidator_Min_TypeErr(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "integer", Min: "10"},
		},
	}
	obj := map[string]interface{}{"field": 28}
	defer func() {
		p := recover()
		// toolbox.Dump(p) // {field int  false <nil> 10 [] 0} attribute 'Min' should be int
		if p == nil {
			fail(t, "should panic("+fmt.Sprint(rules["create"][0])+" attribute 'Min' should be int)")
		}
	}()
	v.Validate(rules, obj, "create")
}

// 参数错误，Max < Min
func Test_Rule_IntegerValidator_Max_Min_LogicErr(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "integer", Max: 10, Min: 100},
		},
	}
	obj := map[string]interface{}{"field": 28}
	defer func() {
		p := recover()
		// toolbox.Dump(p) // {field int  false 10 100 [] 0} attribute 'Max' should greater than 'Min'
		if p == nil {
			fail(t, "should panic("+fmt.Sprint(rules["create"][0])+" attribute 'Max' should greater than 'Min')")
		}
	}()
	v.Validate(rules, obj, "create")
}

// 参数错误，Symbol < 0 && Min >= 0
func Test_Rule_IntegerValidator_Symbol_Min_LogicErr(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "integer", Symbol: -1, Min: 1},
		},
	}
	obj := map[string]interface{}{"field": -10}
	defer func() {
		p := recover()
		// toolbox.Dump(p) // {field int  false <nil> 1 [] -1} attribute 'Min' should less than 0 when 'Symbal' = -1
		if p == nil {
			fail(t, "should panic("+fmt.Sprint(rules["create"][0])+" attribute 'Min' should less than 0 when 'Symbal' = "+strconv.FormatInt(rules["create"][0].Symbol, 10))
		}
	}()
	v.Validate(rules, obj, "create")
}

// 参数错误，Symbol > 0 && Max <= 0
func Test_Rule_IntegerValidator_Symbol_Max_LogicErr(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "integer", Symbol: 1, Max: -1},
		},
	}
	obj := map[string]interface{}{"field": 10}
	defer func() {
		p := recover()
		// toolbox.Dump(p) // {field int  false -1 <nil> [] 1} attribute 'Max' should greater than 0 when 'Symbal' = 1
		if p == nil {
			fail(t, "should panic("+fmt.Sprint(rules["create"][0])+" attribute 'Max' should greater than 0 when 'Symbal' = "+strconv.FormatInt(rules["create"][0].Symbol, 10))
		}
	}()
	v.Validate(rules, obj, "create")
}

// 最大值，Max
func Test_Rule_IntegerValidator_Max(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "integer", Max: 150},
		},
	}
	obj := map[string]interface{}{"field": 151}
	message := generator(v.default_errors["integerMax"], "field", 150)
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:必须是不大于 150 的整数]]
	if len(e) == 0 || e[0]["field"] != message {
		fail(t, "should print error("+message+")")
	}
}

func Test_Rule_IntegerValidator_Max_OK(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "integer", Max: 150},
		},
	}
	obj := map[string]interface{}{"field": 28}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 最小值，Min
func Test_Rule_IntegerValidator_Min(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "integer", Min: 18},
		},
	}
	obj := map[string]interface{}{"field": 17}
	message := generator(v.default_errors["integerMin"], "field", 18)
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:必须是不小于 18 的整数]]
	if len(e) == 0 || e[0]["field"] != message {
		fail(t, "should print error("+message+")")
	}
}

func Test_Rule_IntegerValidator_Min_OK(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "integer", Min: 18},
		},
	}
	obj := map[string]interface{}{"field": 18}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 范围值，Min ~ Max
func Test_Rule_IntegerValidator_Max_Min_Range(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "integer", Max: 35, Min: 18},
		},
	}
	obj := map[string]interface{}{"field": 51}
	message := generator(v.default_errors["integerRange"], "field", 18, 35)
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:必须是介于 18 到 35 的整数]]
	if len(e) == 0 || e[0]["field"] != message {
		fail(t, "should print error("+message+")")
	}
}

func Test_Rule_IntegerValidator_Max_Min_Range_OK(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "integer", Max: 35, Min: 18},
		},
	}
	obj := map[string]interface{}{"field": 22}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 固定值，Max == Min
func Test_Rule_IntegerValidator_Max_Min_Eq(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "integer", Max: 20, Min: 20},
		},
	}
	obj := map[string]interface{}{"field": 21}
	message := generator(v.default_errors["equal"], "field", 20)
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:必须是 20]]
	if len(e) == 0 || e[0]["field"] != message {
		fail(t, "should print error("+message+")")
	}
}

func Test_Rule_IntegerValidator_Max_Min_Eq_OK(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "integer", Max: 20, Min: 20},
		},
	}
	obj := map[string]interface{}{"field": 20}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 正整数 Rule.Symbol > 0
func Test_Rule_IntegerValidator_Symbol_Positive(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "integer", Symbol: 1},
		},
	}
	obj := map[string]interface{}{"field": 0}
	message := generator(v.default_errors["integerPositive"], "field")
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:必须是正整数]]
	if len(e) == 0 || e[0]["field"] != message {
		fail(t, "should print error("+message+")")
	}
}

func Test_Rule_IntegerValidator_Symbol_Positive_OK(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "integer", Symbol: 1},
		},
	}
	obj := map[string]interface{}{"field": 20}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 负整数 Rule.Symbol < 0
func Test_Rule_IntegerValidator_Symbol_Negative(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "integer", Symbol: -1},
		},
	}
	obj := map[string]interface{}{"field": 18}
	message := generator(v.default_errors["integerNegative"], "field")
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:必须是负整数]]
	if len(e) == 0 || e[0]["field"] != message {
		fail(t, "should print error("+message+")")
	}
}

func Test_Rule_IntegerValidator_Symbol_Negative_OK(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "integer", Symbol: -1},
		},
	}
	obj := map[string]interface{}{"field": -10}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 正整数最大值，Symbol（positive） + Max
func Test_Rule_IntegerValidator_Symbol_Positive_Max(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "integer", Symbol: 1, Max: 35},
		},
	}
	obj := map[string]interface{}{"field": 36}
	message := generator(v.default_errors["integerPositiveMax"], "field", 35)
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:必须是不大于 35 的正整数]]
	if len(e) == 0 || e[0]["field"] != message {
		fail(t, "should print error("+message+")")
	}
}

func Test_Rule_IntegerValidator_Symbol_Positive_Max_OK(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "integer", Symbol: 1, Max: 35},
		},
	}
	obj := map[string]interface{}{"field": 35}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 正整数最小值，Symbol（positive） + Min
func Test_Rule_IntegerValidator_Symbol_Positive_Min(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "integer", Symbol: 1, Min: 18},
		},
	}
	obj := map[string]interface{}{"field": 17}
	message := generator(v.default_errors["integerPositiveMin"], "field", 18)
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:必须是不小于 18 的正整数]]
	if len(e) == 0 || e[0]["field"] != message {
		fail(t, "should print error("+message+")")
	}
}

func Test_Rule_IntegerValidator_Symbol_Positive_Min_OK(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "integer", Symbol: 1, Min: 18},
		},
	}
	obj := map[string]interface{}{"field": 18}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 正整数范围值，Symbol（positive） + Max + Min
func Test_Rule_IntegerValidator_Symbol_Positive_Max_Min_Range(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "integer", Symbol: 1, Min: 18, Max: 35},
		},
	}
	obj := map[string]interface{}{"field": 17}
	message := generator(v.default_errors["integerPositiveRange"], "field", 18, 35)
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:必须是介于 18 到 35 的正整数]]
	if len(e) == 0 || e[0]["field"] != message {
		fail(t, "should print error("+message+")")
	}
}

func Test_Rule_IntegerValidator_Symbol_Positive_Max_Min_Range_OK(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "integer", Symbol: 1, Min: 18, Max: 35},
		},
	}
	obj := map[string]interface{}{"field": 20}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 正整数固定值，Symbol（positive） + Max + Min, Max == Min
func Test_Rule_IntegerValidator_Symbol_Positive_Max_Min_Eq(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "integer", Symbol: 1, Min: 18, Max: 18},
		},
	}
	obj := map[string]interface{}{"field": 17}
	message := generator(v.default_errors["equal"], "field", 18)
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:必须是 18]]
	if len(e) == 0 || e[0]["field"] != message {
		fail(t, "should print error("+message+")")
	}
}

func Test_Rule_IntegerValidator_Symbol_Positive_Max_Min_Eq_OK(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "integer", Symbol: 1, Min: 18, Max: 18},
		},
	}
	obj := map[string]interface{}{"field": 18}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 负整数最大值，Symbol（negative） + Max
func Test_Rule_IntegerValidator_Symbol_Negative_Max(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "integer", Symbol: -1, Max: -100},
		},
	}
	obj := map[string]interface{}{"field": -99}
	message := generator(v.default_errors["integerNegativeMax"], "field", -100)
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:必须是不大于 -100 的负整数]]
	if len(e) == 0 || e[0]["field"] != message {
		fail(t, "should print error("+message+")")
	}
}

func Test_Rule_IntegerValidator_Symbol_Negative_Max_OK(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "integer", Symbol: -1, Max: -100},
		},
	}
	obj := map[string]interface{}{"field": -101}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 负整数最小值，Symbol（negative） + Min
func Test_Rule_IntegerValidator_Symbol_Negative_Min(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "integer", Symbol: -1, Min: -100},
		},
	}
	obj := map[string]interface{}{"field": -101}
	message := generator(v.default_errors["integerNegativeMin"], "field", -100)
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:必须是不小于 -100 的负整数]]
	if len(e) == 0 || e[0]["field"] != message {
		fail(t, "should print error("+message+")")
	}
}

func Test_Rule_IntegerValidator_Symbol_Negative_Min_OK(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "integer", Symbol: -1, Min: -100},
		},
	}
	obj := map[string]interface{}{"field": -99}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 负整数范围值，Symbol（negative） + Max + Min
func Test_Rule_IntegerValidator_Symbol_Negative_Max_Min_Range(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "integer", Symbol: -1, Min: -100, Max: -1},
		},
	}
	obj := map[string]interface{}{"field": -101}
	message := generator(v.default_errors["integerNegativeRange"], "field", -100, -1)
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:必须是介于 -100 到 -1 的负整数]]
	if len(e) == 0 || e[0]["field"] != message {
		fail(t, "should print error("+message+")")
	}
}

func Test_Rule_IntegerValidator_Symbol_Negative_Max_Min_Range_OK(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "integer", Symbol: -1, Min: -100, Max: -1},
		},
	}
	obj := map[string]interface{}{"field": -99}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 负整数固定值，Symbol（negative） + Max + Min, Max == Min
func Test_Rule_IntegerValidator_Symbol_Negative_Max_Min_Eq(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "integer", Symbol: -1, Min: -100, Max: -100},
		},
	}
	obj := map[string]interface{}{"field": -101}
	message := generator(v.default_errors["equal"], "field", -100)
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:必须是 -100]]
	if len(e) == 0 || e[0]["field"] != message {
		fail(t, "should print error("+message+")")
	}
}

func Test_Rule_IntegerValidator_Symbol_Negative_Max_Min_Eq_OK(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "integer", Symbol: -1, Min: -100, Max: -100},
		},
	}
	obj := map[string]interface{}{"field": -100}
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
		"create": {
			{Attr: "field", Rule: "integer"},
		},
	}
	obj := map[string]interface{}{"field": float64(3.14)}
	message := generator(v.default_errors["integer"], "field")
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:必须是整数]]
	if len(e) == 0 || e[0]["field"] != message {
		fail(t, "should print error("+message+")")
	}
}

// 有值，float64 - 无小数点浮点数
func Test_Rule_IntegerValidator_OK_Float64(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "integer"},
		},
	}
	obj := map[string]interface{}{"field": float64(3)}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 有值，float32 - 有小数点浮点数
func Test_Rule_IntegerValidator_Float32(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "integer"},
		},
	}
	obj := map[string]interface{}{"field": float32(3.14)}
	message := generator(v.default_errors["integer"], "field")
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:必须是整数]]
	if len(e) == 0 || e[0]["field"] != message {
		fail(t, "should print error("+message+")")
	}
}

// 有值，float32 - 无小数点浮点数
func Test_Rule_IntegerValidator_OK_Float32(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "integer"},
		},
	}
	obj := map[string]interface{}{"field": float32(3)}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 有值，int64
func Test_Rule_IntegerValidator_OK_Int64(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "integer"},
		},
	}
	obj := map[string]interface{}{"field": int64(3)}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 有值，int32
func Test_Rule_IntegerValidator_OK_Int32(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "integer"},
		},
	}
	obj := map[string]interface{}{"field": int32(3)}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 有值，int16
func Test_Rule_IntegerValidator_OK_Int16(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "integer"},
		},
	}
	obj := map[string]interface{}{"field": int16(3)}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 有值，int8
func Test_Rule_IntegerValidator_OK_Int8(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "integer"},
		},
	}
	obj := map[string]interface{}{"field": int8(3)}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 有值，int
func Test_Rule_IntegerValidator_OK_Int(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "integer"},
		},
	}
	obj := map[string]interface{}{"field": 3}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

/***** decimalValidator *****/

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
	// toolbox.Dump(e) // [map[field:必须是小数]]
	if len(e) == 0 || e[0]["field"] != message {
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
	// toolbox.Dump(e) // [map[field:不能为空]]
	if len(e) == 0 || e[0]["field"] != message {
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
	// toolbox.Dump(e) // [map[field:必须是不大于 10 的小数]]
	if len(e) == 0 || e[0]["field"] != message {
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
	// toolbox.Dump(e) // [map[field:必须是不小于 10 的小数]]
	if len(e) == 0 || e[0]["field"] != message {
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
	// toolbox.Dump(e) // [map[field:必须是介于 10 到 11 的小数]]
	if len(e) == 0 || e[0]["field"] != message {
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
	// toolbox.Dump(e) // [map[field:必须是 10.88]]
	if len(e) == 0 || e[0]["field"] != message {
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
	// toolbox.Dump(e) // [map[field:必须是正小数]]
	if len(e) == 0 || e[0]["field"] != message {
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
	// toolbox.Dump(e) // [map[field:必须是负小数]]
	if len(e) == 0 || e[0]["field"] != message {
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
	// toolbox.Dump(e) // [map[field:必须是不大于 3.5 的正小数]]
	if len(e) == 0 || e[0]["field"] != message {
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
	// toolbox.Dump(e) // [map[field:必须是不小于 3.5 的正小数]]
	if len(e) == 0 || e[0]["field"] != message {
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
	// toolbox.Dump(e) // [map[field:必须是介于 3.5 到 4 的正小数]]
	if len(e) == 0 || e[0]["field"] != message {
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
	// toolbox.Dump(e) // [map[field:必须是 3.5]]
	if len(e) == 0 || e[0]["field"] != message {
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
	// toolbox.Dump(e) // [map[field:必须是不大于 -3.14 的负小数]]
	if len(e) == 0 || e[0]["field"] != message {
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
	// toolbox.Dump(e) // [map[field:必须是不小于 -3.14 的负小数]]
	if len(e) == 0 || e[0]["field"] != message {
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
	// toolbox.Dump(e) // [map[field:必须是介于 -3.14 到 -1 的负小数]]
	if len(e) == 0 || e[0]["field"] != message {
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
	// toolbox.Dump(e) // [map[field:必须是 -3.14]]
	if len(e) == 0 || e[0]["field"] != message {
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
	// toolbox.Dump(e) // [map[field:必须是小数]]
	if len(e) == 0 || e[0]["field"] != message {
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
	// toolbox.Dump(e) // [map[field:必须是小数]]
	if len(e) == 0 || e[0]["field"] != message {
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

/***** numberValidator *****/

// 有值，string - 非数字字符串
func Test_Rule_NumberValidator_String(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "number"},
		},
	}
	obj := map[string]interface{}{"field": "hyb"}
	message := generator(v.default_errors["number"], "field")
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:必须是数字]]
	if len(e) == 0 || e[0]["field"] != message {
		fail(t, "should print error("+message+")")
	}
}

// 有值，string - 数字字符串
func Test_Rule_NumberValidator_OK_String(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "number"},
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
func Test_Rule_NumberValidator_Required_False_Empty(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "number"},
		},
	}
	e := v.Validate(rules, objEmpty, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 无值 Rule.Required == true
func Test_Rule_NumberValidator_Required_True_Empty(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "number", Required: true},
		},
	}
	message := generator(v.default_errors["required"], "field")
	e := v.Validate(rules, objEmpty, "create")
	// toolbox.Dump(e) // [map[field:不能为空]]
	if len(e) == 0 || e[0]["field"] != message {
		fail(t, "should print error("+message+")")
	}
}

// 参数错误，Max 非 float64
func Test_Rule_NumberValidator_Max_TypeErr(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "number", Max: "abc"},
		},
	}
	obj := map[string]interface{}{"field": 3.14}
	defer func() {
		p := recover()
		// toolbox.Dump(p) // {field number  false abc <nil> [] 0} attribute 'Max' should be int or float64
		if p == nil {
			fail(t, "should panic("+fmt.Sprint(rules["create"][0])+" attribute 'Max' should be int or float64)")
		}
	}()
	v.Validate(rules, obj, "create")
}

// 参数错误，Min 非 float64
func Test_Rule_NumberValidator_Min_TypeErr(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "number", Min: "abc"},
		},
	}
	obj := map[string]interface{}{"field": 3.14}
	defer func() {
		p := recover()
		// toolbox.Dump(p) // {field number  false abc <nil> [] 0} attribute 'Min' should be int or float64
		if p == nil {
			fail(t, "should panic("+fmt.Sprint(rules["create"][0])+" attribute 'Min' should be int or float64)")
		}
	}()
	v.Validate(rules, obj, "create")
}

// 参数错误，Max < Min
func Test_Rule_NumberValidator_Max_Min_LogicErr(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "number", Max: 2.2, Min: 6},
		},
	}
	obj := map[string]interface{}{"field": 5.4}
	defer func() {
		p := recover()
		// toolbox.Dump(p) // {field number  false 10 100 [] 0} attribute 'Max' should greater than 'Min'
		if p == nil {
			fail(t, "should panic("+fmt.Sprint(rules["create"][0])+" attribute 'Max' should greater than 'Min')")
		}
	}()
	v.Validate(rules, obj, "create")
}

// 参数错误，Symbol < 0 && Min >= 0
func Test_Rule_NumberValidator_Symbol_Min_LogicErr(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "number", Symbol: -1, Min: 3.14},
		},
	}
	obj := map[string]interface{}{"field": -3.32}
	defer func() {
		p := recover()
		// toolbox.Dump(p) // {field number  false <nil> 3.14 [] -1} attribute 'Min' should less than 0 when 'Symbal' = -1
		if p == nil {
			fail(t, "should panic("+fmt.Sprint(rules["create"][0])+" attribute 'Min' should less than 0 when 'Symbal' = "+strconv.FormatInt(rules["create"][0].Symbol, 10))
		}
	}()
	v.Validate(rules, obj, "create")
}

// 参数错误，Symbol > 0 && Max <= 0
func Test_Rule_NumberValidator_Symbol_Max_LogicErr(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "number", Symbol: 1, Max: -3},
		},
	}
	obj := map[string]interface{}{"field": 2.2}
	defer func() {
		p := recover()
		// toolbox.Dump(p) // {field number  false -3 <nil> [] 1} attribute 'Max' should greater than 0 when 'Symbal' = 1
		if p == nil {
			fail(t, "should panic("+fmt.Sprint(rules["create"][0])+" attribute 'Min' should less than 0 when 'Symbal' = "+strconv.FormatInt(rules["create"][0].Symbol, 10))
		}
	}()
	v.Validate(rules, obj, "create")
}

// 最大值，Max
func Test_Rule_NumberValidator_Max(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "number", Max: 10},
		},
	}
	obj := map[string]interface{}{"field": 22.2}
	message := generator(v.default_errors["numberMax"], "field", 10)
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:必须是不大于 10 的数]]
	if len(e) == 0 || e[0]["field"] != message {
		fail(t, "should print error("+message+")")
	}
}

func Test_Rule_NumberValidator_Max_OK_Int(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "number", Max: 10},
		},
	}
	obj := map[string]interface{}{"field": 9.98}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

func Test_Rule_NumberValidator_Max_OK_Float64(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "number", Max: 9.99},
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
func Test_Rule_NumberValidator_Min(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "number", Min: 10},
		},
	}
	obj := map[string]interface{}{"field": 3.14}
	message := generator(v.default_errors["numberMin"], "field", 10)
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:必须是不小于 10 的数]]
	if len(e) == 0 || e[0]["field"] != message {
		fail(t, "should print error("+message+")")
	}
}

func Test_Rule_NumberValidator_Min_OK_Int(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "number", Min: 10},
		},
	}
	obj := map[string]interface{}{"field": 10.35}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

func Test_Rule_NumberValidator_Min_OK_Float64(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "number", Min: 10.35},
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
func Test_Rule_NumberValidator_Max_Min_Range(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "number", Min: 10, Max: 11},
		},
	}
	obj := map[string]interface{}{"field": 11.11}
	message := generator(v.default_errors["numberRange"], "field", 10, 11)
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:必须是介于 10 到 11 的数]]
	if len(e) == 0 || e[0]["field"] != message {
		fail(t, "should print error("+message+")")
	}
}

func Test_Rule_NumberValidator_Max_Min_Range_OK(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "number", Min: 10, Max: 11},
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
func Test_Rule_NumberValidator_Max_Min_Eq(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "number", Min: 10.88, Max: 10.88},
		},
	}
	obj := map[string]interface{}{"field": 11.11}
	message := generator(v.default_errors["equal"], "field", 10.88)
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:必须是 10.88]]
	if len(e) == 0 || e[0]["field"] != message {
		fail(t, "should print error("+message+")")
	}
}

func Test_Rule_NumberValidator_Max_Min_Eq_OK(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "number", Min: 10.88, Max: 10.88},
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
func Test_Rule_NumberValidator_Symbol_Positive(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "number", Symbol: 1},
		},
	}
	obj := map[string]interface{}{"field": -3.14}
	message := generator(v.default_errors["numberPositive"], "field")
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:必须是正数]]
	if len(e) == 0 || e[0]["field"] != message {
		fail(t, "should print error("+message+")")
	}
}

func Test_Rule_NumberValidator_Symbol_Positive_OK(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "number", Symbol: 1},
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
func Test_Rule_NumberValidator_Symbol_Negative(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "number", Symbol: -1},
		},
	}
	obj := map[string]interface{}{"field": 3.14}
	message := generator(v.default_errors["numberNegative"], "field")
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:必须是负数]]
	if len(e) == 0 || e[0]["field"] != message {
		fail(t, "should print error("+message+")")
	}
}

func Test_Rule_NumberValidator_Symbol_Negative_OK(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "number", Symbol: -1},
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
func Test_Rule_NumberValidator_Symbol_Positive_Max(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "number", Symbol: 1, Max: 3.5},
		},
	}
	obj := map[string]interface{}{"field": 3.77}
	message := generator(v.default_errors["numberPositiveMax"], "field", 3.5)
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:必须是不大于 3.5 的正数]]
	if len(e) == 0 || e[0]["field"] != message {
		fail(t, "should print error("+message+")")
	}
}

func Test_Rule_NumberValidator_Symbol_Positive_Max_OK(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "number", Symbol: 1, Max: 3.5},
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
func Test_Rule_NumberValidator_Symbol_Positive_Min(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "number", Symbol: 1, Min: 3.5},
		},
	}
	obj := map[string]interface{}{"field": 3.14}
	message := generator(v.default_errors["numberPositiveMin"], "field", 3.5)
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:必须是不小于 3.5 的正数]]
	if len(e) == 0 || e[0]["field"] != message {
		fail(t, "should print error("+message+")")
	}
}

func Test_Rule_NumberValidator_Symbol_Positive_Min_OK(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "number", Symbol: 1, Min: 3.5},
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
func Test_Rule_NumberValidator_Symbol_Positive_Max_Min_Range(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "number", Symbol: 1, Min: 3.5, Max: 4},
		},
	}
	obj := map[string]interface{}{"field": 3.33}
	message := generator(v.default_errors["numberPositiveRange"], "field", 3.5, 4)
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:必须是介于 3.5 到 4 的正数]]
	if len(e) == 0 || e[0]["field"] != message {
		fail(t, "should print error("+message+")")
	}
}

func Test_Rule_NumberValidator_Symbol_Positive_Max_Min_Range_OK(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "number", Symbol: 1, Min: 3.5, Max: 4},
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
func Test_Rule_NumberValidator_Symbol_Positive_Max_Min_Eq(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "number", Symbol: 1, Min: 3.5, Max: 3.5},
		},
	}
	obj := map[string]interface{}{"field": 3.4}
	message := generator(v.default_errors["equal"], "field", 3.5)
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:必须是 3.5]]
	if len(e) == 0 || e[0]["field"] != message {
		fail(t, "should print error("+message+")")
	}
}

func Test_Rule_NumberValidator_Symbol_Positive_Max_Min_Eq_OK(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "number", Symbol: 1, Min: 3.5, Max: 3.5},
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
func Test_Rule_NumberValidator_Symbol_Negative_Max(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "number", Symbol: -1, Max: -3.14},
		},
	}
	obj := map[string]interface{}{"field": -1.1}
	message := generator(v.default_errors["numberNegativeMax"], "field", -3.14)
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:必须是不大于 -3.14 的负数]]
	if len(e) == 0 || e[0]["field"] != message {
		fail(t, "should print error("+message+")")
	}
}

func Test_Rule_NumberValidator_Symbol_Negative_Max_OK(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "number", Symbol: -1, Max: -3.14},
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
func Test_Rule_NumberValidator_Symbol_Negative_Min(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "number", Symbol: -1, Min: -3.14},
		},
	}
	obj := map[string]interface{}{"field": -3.5}
	message := generator(v.default_errors["numberNegativeMin"], "field", -3.14)
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:必须是不小于 -3.14 的负数]]
	if len(e) == 0 || e[0]["field"] != message {
		fail(t, "should print error("+message+")")
	}
}

func Test_Rule_NumberValidator_Symbol_Negative_Min_OK(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "number", Symbol: -1, Min: -3.14},
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
func Test_Rule_NumberValidator_Symbol_Negative_Max_Min_Range(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "number", Symbol: -1, Min: -3.14, Max: -1},
		},
	}
	obj := map[string]interface{}{"field": -3.5}
	message := generator(v.default_errors["numberNegativeRange"], "field", -3.14, -1)
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:必须是介于 -3.14 到 -1 的负数]]
	if len(e) == 0 || e[0]["field"] != message {
		fail(t, "should print error("+message+")")
	}
}

func Test_Rule_NumberValidator_Symbol_Negative_Max_Min_Range_OK(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "number", Symbol: -1, Min: -3.14, Max: -1},
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
func Test_Rule_NumberValidator_Symbol_Negative_Max_Min_Eq(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "number", Symbol: -1, Min: -3.14, Max: -3.14},
		},
	}
	obj := map[string]interface{}{"field": -3.5}
	message := generator(v.default_errors["equal"], "field", -3.14)
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:必须是 -3.14]]
	if len(e) == 0 || e[0]["field"] != message {
		fail(t, "should print error("+message+")")
	}
}

func Test_Rule_NumberValidator_Symbol_Negative_Max_Min_Eq_OK(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "number", Symbol: -1, Min: -3.14, Max: -3.14},
		},
	}
	obj := map[string]interface{}{"field": -3.14}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 有值，string - 非数字字符串，同 Test_Rule_NumberValidator_String
// 有值，string - 数字字符串，同 Test_Rule_NumberValidator_OK_String

// 有值，float64
func Test_Rule_NumberValidator_OK_Float64(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "number"},
		},
	}
	obj := map[string]interface{}{"field": 3.14}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 有值，float32
func Test_Rule_NumberValidator_OK_Float32(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "number"},
		},
	}
	obj := map[string]interface{}{"field": float32(3.14)}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 有值，int64
func Test_Rule_NumberValidator_OK_Int64(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "number"},
		},
	}
	obj := map[string]interface{}{"field": int64(3)}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 有值，int32
func Test_Rule_NumberValidator_OK_Int32(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "number"},
		},
	}
	obj := map[string]interface{}{"field": int32(3)}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 有值，int16
func Test_Rule_NumberValidator_OK_Int16(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "number"},
		},
	}
	obj := map[string]interface{}{"field": int16(3)}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 有值，int8
func Test_Rule_NumberValidator_OK_Int8(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "number"},
		},
	}
	obj := map[string]interface{}{"field": int8(3)}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 有值，int
func Test_Rule_NumberValidator_OK_Int(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "number"},
		},
	}
	obj := map[string]interface{}{"field": 3}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

/***** booleanValidator *****/

// 有值，类型错误
func Test_Rule_BooleanValidator(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "boolean"},
		},
	}
	obj := map[string]interface{}{"field": 1}
	message := generator(v.default_errors["boolean"], "field")
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:必须是布尔值或布尔字符串]]
	if len(e) == 0 || e[0]["field"] != message {
		fail(t, "should print error("+message+")")
	}
}

// 有值，bool
func Test_Rule_BooleanValidator_OK(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "boolean"},
		},
	}
	obj := map[string]interface{}{"field": false}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

// 无值 Rule.Required == false
func Test_Rule_BooleanValidator_Required_False_Empty(t *testing.T) {
	rules := Rules{
		"create": {
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
		"create": {
			{Attr: "field", Rule: "boolean", Required: true},
		},
	}
	message := generator(v.default_errors["required"], "field")
	e := v.Validate(rules, objEmpty, "create")
	// toolbox.Dump(e) // [map[field:不能为空]]
	if len(e) == 0 || e[0]["field"] != message {
		fail(t, "should print error("+message+")")
	}
}

// 有值，bool，同 Test_Rule_BooleanValidator_OK

// 有值，string，非布尔字符串
func Test_Rule_BooleanValidator_String(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "boolean"},
		},
	}
	obj := map[string]interface{}{"field": "abc"}
	message := generator(v.default_errors["boolean"], "field")
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // [map[field:必须是布尔值或布尔字符串]]
	if len(e) == 0 || e[0]["field"] != message {
		fail(t, "should print error("+message+")")
	}
}

// 有值，string，布尔字符串
func Test_Rule_BooleanValidator_OK_String(t *testing.T) {
	rules := Rules{
		"create": {
			{Attr: "field", Rule: "boolean"},
		},
	}
	obj := map[string]interface{}{"field": "true"}
	e := v.Validate(rules, obj, "create")
	// toolbox.Dump(e) // []
	if len(e) != 0 {
		fail(t, "should print nothing")
	}
}

/***** ipValidator *****/

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
	if len(e) == 0 || e[0]["field"] != message {
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
	// toolbox.Dump(e) // [map[field:必须是字符串]]
	if len(e) == 0 || e[0]["field"] != message {
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
	// toolbox.Dump(e) // [map[field:不能为空]]
	if len(e) == 0 || e[0]["field"] != message {
		fail(t, "should print error("+message+")")
	}
}

//***** regexValidator *****/

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
	// toolbox.Dump(e) // [map[field:格式不对]]
	if len(e) == 0 || e[0]["field"] != message {
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
	// toolbox.Dump(e) // [map[field:必须是字符串]]
	if len(e) == 0 || e[0]["field"] != message {
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
	// toolbox.Dump(e) // [map[field:不能为空]]
	if len(e) == 0 || e[0]["field"] != message {
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
	if len(e) == 0 || e[0]["field"] != message {
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
	// toolbox.Dump(e) // [map[field:必须是字符串]]
	if len(e) == 0 || e[0]["field"] != message {
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
	// toolbox.Dump(e) // [map[field:不能为空]]
	if len(e) == 0 || e[0]["field"] != message {
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
	if len(e) == 0 || e[0]["field"] != message {
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
	// toolbox.Dump(e) // [map[field:必须是字符串]]
	if len(e) == 0 || e[0]["field"] != message {
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
	// toolbox.Dump(e) // [map[field:不能为空]]
	if len(e) == 0 || e[0]["field"] != message {
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
	// toolbox.Dump(e) // [map[field:手机号格式不对]]
	if len(e) == 0 || e[0]["field"] != message {
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
	// toolbox.Dump(e) // [map[field:必须是字符串]]
	if len(e) == 0 || e[0]["field"] != message {
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
	// toolbox.Dump(e) // [map[field:不能为空]]
	if len(e) == 0 || e[0]["field"] != message {
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
	if len(e) == 0 || e[0]["field"] != message {
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
	// toolbox.Dump(e) // [map[field:必须是字符串]]
	if len(e) == 0 || e[0]["field"] != message {
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
	// toolbox.Dump(e) // [map[field:不能为空]]
	if len(e) == 0 || e[0]["field"] != message {
		fail(t, "should print error("+message+")")
	}
}
