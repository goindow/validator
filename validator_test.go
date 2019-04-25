package validator

import (
	"testing"
	"fmt"
	"errors"
	"strings"

	// "github.com/goindow/toolbox"
)

const (
	ZH_CN = "zh_cn"
	EN_US = "EN_US"
)

var (
	v *validator
	rulesEmpty = Rules{}
	objEmpty = map[string]interface{}{}
)

func init() {
	v = New()
	// toolbox.Dump(v)
}

/***** Lang() *****/

// 切换语言包
func Test_Lang(t *testing.T) {
	vEnglish := New().Lang(EN_US)	// 重新实例化，避免语言对其他用例造成影响
	var rules = Rules{
		"create": {
			{Attr: "hobby", Rule: "required"},
		},
	}
	message := generator(vEnglish.default_errors["required"], "hobby")
	e := vEnglish.Validate(rules, objEmpty, "create")
	// toolbox.Dump(e) // [map[hobby:missing field hobby]]
	if len(e) == 0 || e[0]["hobby"].Error() != message {
		fail(t, "should print errors(" + message + ")")
	}	
}

// 未定义语言包
func Test_Lang_Undefiend_Lang(t *testing.T) {
	LANG := "undefined_lang"
	defer func() {
		p := recover()
		// toolbox.Dump(p) // undefined_lang unsupport language
		if p == nil {
			fail(t, "should panic(" + LANG + " unsupport language)")
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
			fail(t, "should panic(" + string(SCENCE) + " scence undefined)")
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
			fail(t, "should panic(" + RULE + " validator undefined)")
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
			fail(t, "should panic(" + fmt.Sprint(rules["create"][0]) + " attribute 'Rule' not found)")
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
			fail(t, "should panic(" + fmt.Sprint(rules["create"][0]) + " attribute 'Attr' not found)")
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
	// toolbox.Dump(e) // [map[username:username 不能为空]]
	if len(e) == 0 || e[0]["username"].Error() != message {
		fail(t, "should print error(" + message + ")")
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
	if len(e) == 0 || e[0]["username"].Error() != rules["create"][0].Message {
		fail(t, "should print error(" + rules["create"][0].Message + ")")
	}
}

/***** AddValidator() *****/

// 自定义验证规则，规则名已存在
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
				if !rule.Required {	// 允许为空
					return nil
				}
				return v.generator("required", attr, rule)
			}
			if obj[attr] != 1 {
				return E{attr: errors.New(rule.Message)}
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
	if len(e) == 0 || e[0]["field"].Error() != message {
		fail(t, "should print error(" + message + ")")
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

/***** utils *****/

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