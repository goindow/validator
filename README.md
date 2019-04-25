# validator
Package validator 实现了一个支持场景/国际化/自定义错误/自定义验证规则的 map[string]interface{} 元素批量验证器，意在支持各种框架的 model 层实现自动验证，亦可单独使用

## 说明
- 该验证器是逻辑验证器（float64(10)/int32(10) 均可被 intValidator 验证通过）而不是强类型验证器（如果需要强类型验证，可以自定义验证规则，或直接断言），
- 考虑到经过 encoding/json 解析后的数字类型均被解析为 float64，强类型验证器不太可用，故此设计，如需强类型验证，请使用 AddValidator 自行扩展

## 特性
- 支持场景
- 支持国际化
- 支持批量验证
- 支持自定义验证器
- 支持自定义错误信息

## 安装
```bash
go get github.com/goindow/validator
```

## 示例
```go
package main

import (
	"github.com/goindow/validator"
	"fmt"
)

func main() {
	user := map[string]interface{}{
		// "username": "hyb",
		"password": "******",
		"gender": "male",
		"age": 17,
		"weight": "53kg",
		"email": "hyb76788424#163.com",
	}

	rules := validator.Rules{
		"create": {
			{ Attr: []string{"username", "password"}, Rule: "required" },
			{ Attr: "password", Rule: "regex", Pattern: `[A-Z]{1}\w{5,}`, Message: "密码必须由大写字母开头"},
			{ Attr: "gender", Rule: "in", Enum: []string{"0", "1"} },
			{ Attr: "age", Rule: "int", Min: 18 },
			{ Attr: "weight", Rule: "number", Symbol: 1 },
			{ Attr: "email", Rule: "email" },
		},
		"read": {
			{ Attr: "id", Rule: "int", Symbol: 1 },
		},
	}

	if e := validator.New().Validate(rules, user, "create"); e != nil {
		// todo: handle errors
		for _, i := range e {
			for k, v := range i {
				fmt.Printf("%v => %v\n", k, v)
			}
		}
		// username => username 不能为空
		// password => 密码必须由大写字母开头
		// gender => gender 只能是 [0、1] 中的一个
		// age => age 必须是不小于 18 的整数
		// weight => weight 必须是数字
		// email => 无效的 email
	}
	// todo: do something
}
```

## 如何定义验证规则
- validator.Scence string 场景
- validator.Rule struct 规则
	-- Attr string
	-- Rule string
	-- Message string
	-- Required bool
	-- Symbol int64
	-- Max interface{}
	-- Min interface{}
	-- Enum []string
	-- Pattern string
- validator.ScenceRules []validator.Rule 验证规则集 - 单一场景
- validator.Rules map[Scence]ScenceRules 验证规则集 - 所有场景
```go
rules := validator.Rules{ // validator.Rules
	"create": { // validators.ScenceRules
		{ Attr: []string{"username", "password"}, Rule: "required" }, // validator.Rule
		{ Attr: "password", Rule: "regex", Pattern: `[A-Z]{1}\w{5,}`, Message: "密码必须由大写字母开头"},
		{ Attr: "gender", Rule: "in", Enum: []string{"0", "1"} },
		{ Attr: "age", Rule: "int", Min: 18 },
		{ Attr: "weight", Rule: "number", Symbol: 1 },
		{ Attr: "email", Rule: "email" },
	},
	"read": {
		{ Attr: "id", Rule: "int", Symbol: 1 },
	},
}
```

## 国际化
- Lang(lang string) *validator
- 在 i18n 下，新建错误信息对应的语言文件，格式参考已有文件
- 包本身自带两种语言(zh_cn、en_us)，默认语言为 zh_cn
```go
// touch ./i18n/en_us.go

v := validator.New().Lang("en_us")
```

## 自定义错误信息
- Rule.Message string
```go
rules := validator.Rules{
	"create": {
		{ Attr: "password", Rule: "regex", Pattern: `[A-Z]{1}\w{5,}`, Message: "密码必须由大写字母开头"},
	}
}
```

## 自定义验证器
- AddValidator(name string, customValidator F)
```go
package main

import (
	"github.com/goindow/validator"
	"errors"
	"fmt"
)

func main() {

	// 自定义错误处理函数，函数类型为 validator.F
	var oneValidator validator.F = func(attr string, rule validator.Rule, obj validator.M) validator.E {
		if _, ok := obj[attr]; !ok {
			return validator.E{attr: errors.New("not found")}
		}
		if obj[attr] != 1 {
			e := rule.Message
			if e == "" {
				e = "必须等于一"
			}
			return validator.E{attr: errors.New(e)}	
		}
		return nil
	}

	v := validator.New()

	// 挂载
	v.AddValidator("one", oneValidator)

	// 使用
	user := map[string]interface{}{
		"name": "hyb",
	}
	rules := validator.Rules{
		"someone": {
			{Attr: "name", Rule: "one"},
		},
	}
	e := v.Validate(rules, user, "someone")
	fmt.Println(e)
	// [map[name:必须等于一]]
}
```

## 内置验证器
- [requiredValidator](#requiredValidator)


### requiredValidator
- 必填验证器
```go
rule := {Attr: "field", Rule: "Required"}
```