# validator
Package validator 实现了一个支持场景/国际化/自定义错误/自定义验证规则的 map[string]interface{} 元素批量验证器，意在支持各种框架的 model 层实现自动验证，亦可单独使用

## 说明
- 该验证器是 ***逻辑验证器***（float64(10)/int32(10)/"10" 均可被 intValidator 验证通过）而不是 ***强类型验证器***
- 考虑到经过 encoding/json 解析后的数字类型均被解析为 float64，强类型验证器不太可用，故此设计
- 如果需要强类型验证，可以使用 AddValidator 自行扩展

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
- ***validator.Rule*** struct 验证规则
	- ***Attr***        interface{}    **必选**，待验证属性，单个属性 string，多个属性 []string，其他类型或未定义将 panic
	- ***Rule***        string         **必选**，验证规则，即验证器，不存在的验证器或未定义将 panic
	- ***Message***     string         **可选**，自定义错误信息
	- ***Required***    bool           **可选**，可空限制，作用于除 requiredValidator 外的所有验证器，false(默认) - 有值验证/无值跳过，true - 有值验证/无值报错
	- ***Symbol***      int64          **可选**，符号限制，作用于 numberValidator、integerValidator、decimalValidator，0(默认) - 正/负数，>0 - 正数(不包含0)，<0 - 负数(不包含0)
	- ***Max***         interface{}    **可选**，最大限制，作用于 stringValidator、numberValidator、integerValidator、decimalValidator
	- ***Min***         interface{}    **可选**，最小限制，同 Max
	- ***Enum***        []string       **必选**（inValidator）**，枚举限制，作用于 inValidator
	- ***Pattern***     string         **必选**（regexValidator）**，正则匹配模式，作用于 regexValidator
- ***validator.Scence*** string 场景
- ***validator.ScenceRules*** []validator.Rule 验证规则集 - 单一场景
- ***validator.Rules map[Scence]ScenceRules*** 验证规则集 - 所有场景
```go
rules := validator.Rules{ // validator.Rules
	// validator.ScenceRules
	"create": { // validator.Scence
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
- 在 i18n 下，新建错误信息对应的语言文件，格式参考已有文件，包本身自带两种语言(zh_cn、en_us)，默认语言为 zh_cn
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
	v := validator.New()

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
- [inValidator](#inValidator)
- [stringValidator](#stringValidator)
- [integerValidator](#integerValidator)
- [decimalValidator](#decimalValidator)
- [numberValidator](#numberValidator)
- [booleanValidator](#booleanValidator)
- [ipValidator](#ipValidator)
- [regexValidator](#regexValidator)
- [emailValidator](#emailValidator)
- [teldValidator](#teldValidator)
- [mobileValidator](#mobileValidator)
- [zipcodeValidator](#zipcodeValidator)


### requiredValidator
- 必填
- Rule.Rule        string      必选    required
```go
rule := {Attr: []string{"username", "password"}, Rule: "required"}
```

### inValidator
- 枚举，被验证字段支持类型 int64、int32、int16、int8、int、float64、float32、string、bool
- Rule.Rule        string      必选    in
- Rule.Required    bool        可选    false(默认) - 被验证字段有值验证/无值跳过，true - 被验证字段无值，验证失败，报 reqired 错误
- Rule.Enum        []string    必选    被验证字段必须在 Rule.Enum 中
```go
rule := {Attr: "gender", Rule: "in", Enum: {"male", "female", "unknown"}}
rule := {Attr: "gender", Rule: "in", Enum: {"male", "female", "unknown"}, Required: true} // 有值验证，无值跳过
```

### stringValidator
- 字符串
- Rule.Rule        string    必选    string
- Rule.Required    bool      可选    false(默认) - 被验证字段有值验证/无值跳过，true - 被验证字段无值，验证失败，报 reqired 错误
- Rule.Max         int       可选    被验证字段长度不能大于 Rule.Max
- Rule.Min         in        可选    被验证字段长度不能小于 Rule.Min
```go
rule := {Attr: "name", Rule: "string"}
rule := {Attr: "name", Rule: "string", Min: 6} // utf8 字符数，即字符串长度，兼容中文
rule := {Attr: "name", Rule: "string", Min: 6, Max: 18}
rule := {Attr: "name", Rule: "string", Min: 6, Max: 18, required: true}
```

### integerValidator
- 整数，被验证字段支持类型 int64、int32、int16、int8、int、float64、float32、string
- Rule.Rule        string    必选    integer/int
- Rule.Required    bool      可选    false(默认) - 被验证字段有值验证/无值跳过，true - 被验证字段无值，验证失败，报 reqired 错误
- Rule.Symbol      int64     可选    0(默认) - 正/负数，>0 - 正数(不包含0)，<0 - 负数(不包含0)
- Rule.Max         int       可选    被验证字段大小不能大于 Rule.Max
- Rule.Min         int       可选    被验证字段大小不能小于 Rule.Min
```go
// int 为 integer 的别名，都指向 integerValidator 验证器
rule := {Attr: "age", Rule: "int"}
rule := {Attr: "age", Rule: "integer", Symobl: 1} // 正整数
rule := {Attr: "age", Rule: "integer", Min: 18}
rule := {Attr: "age", Rule: "integer", Min: 18, Max: 18} // == 18
rule := {Attr: "age", Rule: "integer", Min: 18, Max: 35, required: true}

// float64(18)、float32(18)、"18" 都会被认为是整数
```

### decimalValidator
- 小数，被验证字段支持类型 float64、float32、string
- Rule.Rule        string         必选    decimal/float
- Rule.Required    bool           可选    false(默认) - 被验证字段有值验证/无值跳过，true - 被验证字段无值，验证失败，报 reqired 错误
- Rule.Symbol      int64          可选    0(默认) - 正/负数，>0 - 正数(不包含0)，<0 - 负数(不包含0)
- Rule.Max         int|float64    可选    被验证字段大小不能大于 Rule.Max
- Rule.Min         int|float64    可选    被验证字段大小不能小于 Rule.Min
```go
// float 为 decimal 的别名，都指向 decimalValidator 验证器
rule := {Attr: "field", Rule: "float"}
rule := {Attr: "field", Rule: "decimal", Symobl: -1} // 负小数
rule := {Attr: "field", Rule: "decimal", Min: 2}
rule := {Attr: "field", Rule: "decimal", Min: 3.14, Max: 3.14} // == 3.14
rule := {Attr: "field", Rule: "decimal", Min: 3, Max: 3.14, required: true}

// float64(18)、float32(18)、"18" 没有小数位会验证失败
```

### numberValidator
- 数字，被验证字段支持类型 int64、int32、int16、int8、int、float64、float32、string
- Rule.Rule        string         必选    number
- Rule.Required    bool           可选    false(默认) - 被验证字段有值验证/无值跳过，true - 被验证字段无值，验证失败，报 reqired 错误
- Rule.Symbol      int64          可选    0(默认) - 正/负数，>0 - 正数(不包含0)，<0 - 负数(不包含0)
- Rule.Max         int|float64    可选    被验证字段大小不能大于 Rule.Max
- Rule.Min         int|float64    可选    被验证字段大小不能小于 Rule.Min
```go
rule := {Attr: "weight", Rule: "number"}
rule := {Attr: "weight", Rule: "number", Symobl: 1}
rule := {Attr: "weight", Rule: "number", Min: 45}
rule := {Attr: "weight", Rule: "number", Min: 45, Max: 45} // == 45
rule := {Attr: "weight", Rule: "number", Min: 45, Max: 49.9, required: true}
```

### booleanValidator
- 布尔，被验证字段支持类型 bool、string
- Rule.Rule        string    必选    boolean/bool
- Rule.Required    bool      可选    false(默认) - 被验证字段有值验证/无值跳过，true - 被验证字段无值，验证失败，报 reqired 错误
```go
// bool 为 boolean 的别名，都指向 booleanValidator 验证器
rule := {Attr: "admin", Rule: "bool"}
rule := {Attr: "admin", Rule: "boolean"}

// 布尔值[true、false]、字符串表示的布尔值["1"、"0"、"t"、"f"、"true"、"false"(忽略大小写)] 都会被认为是布尔
```

### ipValidator
- ipv4/ipv6，被验证字段支持类型 string
- Rule.Rule        string    必选    ip
- Rule.Required    bool      可选    false(默认) - 被验证字段有值验证/无值跳过，true - 被验证字段无值，验证失败，报 reqired 错误
```go
rule := {Attr: "ip", Rule: "ip"}
```

### regexValidator
- 正则，被验证字段支持类型 string
- Rule.Rule        string    必选    regex
- Rule.Required    bool      可选    false(默认) - 被验证字段有值验证/无值跳过，true - 被验证字段无值，验证失败，报 reqired 错误
- Rule.Pattern     string    必选    正则模式字符串
```go
rule := {Attr: "password", Rule: "regex", Pattern: `[A-Z]{1}\w{5,}`},
```

### emailValidator
- email，被验证字段支持类型 string
- Rule.Rule        string    必选    email
- Rule.Required    bool      可选    false(默认) - 被验证字段有值验证/无值跳过，true - 被验证字段无值，验证失败，报 reqired 错误
```go
rule := {Attr: "email", Rule: "email"}

// pattern = `^[\w!#$%&'*+/=?^_` + "`" + `{|}~-]+(?:\.[\w!#$%&'*+/=?^_` + "`" + `{|}~-]+)*@(?:[\w](?:[\w-]*[\w])?\.)+[a-zA-Z0-9](?:[\w-]*[\w])?$`
```

### telValidator
- 中国大陆座机号，被验证字段支持类型 string
- Rule.Rule        string    必选    tel
- Rule.Required    bool      可选    false(默认) - 被验证字段有值验证/无值跳过，true - 被验证字段无值，验证失败，报 reqired 错误
```go
rule := {Attr: "tel", Rule: "tel"}

// pattern = `^(0\d{2,3}(\-)?)?\d{7,8}$`
```

### mobileValidator
- 中国大陆手机号，被验证字段支持类型 string
- Rule.Rule        string    必选    mobile
- Rule.Required    bool      可选    false(默认) - 被验证字段有值验证/无值跳过，true - 被验证字段无值，验证失败，报 reqired 错误
```go
rule := {Attr: "mobile", Rule: "mobile"}

// pattern = `^((\+86)|(86))?(1(([35][0-9])|[8][0-9]|[7][01356789]|[4][579]))\d{8}$`
```

### zipcodeValidator
- 中国大陆邮编，被验证字段支持类型 string
- Rule.Rule        string    必选    zipcode
- Rule.Required    bool      可选    false(默认) - 被验证字段有值验证/无值跳过，true - 被验证字段无值，验证失败，报 reqired 错误
```go
rule := {Attr: "zipcode", Rule: "zipcode"}

// pattern = `^[1-9]\d{5}$`
```