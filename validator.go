// Package validator 实现了一个支持场景/国际化/自定义错误/自定义验证规则的 map[string]interface{} 元素批量验证器，意在支持各种框架的 model 层实现自动验证，亦可单独使用
// 该验证器是逻辑验证器（float64(10)/int32(10) 均可被 intValidator 验证通过）而不是强类型验证器（如果需要强类型验证，可以自定义验证规则，或直接断言）
// 考虑到经过 encoding/json 解析后的数字类型均被解析为 float64，故强类型验证器不太可用，故如此设计，如需强类型验证，请使用 AddValidator 自行扩展
package validator

import (
	"fmt"
	"net"
	"errors"
	"regexp"
	"strings"
	"strconv"
	"reflect"
	"unicode/utf8"
	"github.com/goindow/validator/i18n"

	// "github.com/goindow/toolbox"
)

const (
	DEFAULT_LANG = "ZH_CN"

	PATTERN_ZIPCODE = `^[1-9]\d{5}$`
	PATTERN_TEL = `^(0\d{2,3}(\-)?)?\d{7,8}$`
	PATTERN_MOBILE = `^((\+86)|(86))?(1(([35][0-9])|[8][0-9]|[7][01356789]|[4][579]))\d{8}$`
	PATTERN_EMAIL = `^[\w!#$%&'*+/=?^_` + "`" + `{|}~-]+(?:\.[\w!#$%&'*+/=?^_` + "`" + `{|}~-]+)*@(?:[\w](?:[\w-]*[\w])?\.)+[a-zA-Z0-9](?:[\w-]*[\w])?$`
)

// 验证规则
type Rule struct {
	// 必须，待验证属性，单个属性 string，多个属性 []string，其他类型或未定义将 panic
	Attr	  interface{}
	// 必须，验证规则，即验证器，不存在的验证器或未定义将 panic
	Rule      string
	// 可选，自定义错误信息
	Message	  string
	// 可选，可空限制，作用于除 requiredValidator 外的所有验证器
	// false(默认) - 有值验证/无值跳过(如果同时设置了 requiredValidator ，则报 required 错误，此时错误是由 requiredValidator 报出的)
	// true - 有值验证/无值报 required 错误
	Required  bool
	// 可选，符号限制，作用于 numberValidator、intValidator、floatValidator
	// 0(默认) - 正/负数，>0 - 正数(不包含0)，<0 - 负数(不包含0)
	Symbol	  int64
	// 可选，最大限制，作用于 stringValidator、numberValidator、intValidator、floatValidator
	// Max < Min 将 Panic
	// 作用于 stringValidator、intValidator 只能是整数，其他将 panic
	// 作用于 numberValidator、floatValidator 只能是整数或浮点数，其他将Panic
	Max       interface{}
	// 可选，最小限制，同 Rule.Max
	Min		  interface{}
	// 必选（inValidator），枚举限制，作用于 inValidator
	Enum	  []string
	// 必选（regexValidator），正则匹配模式，作用于 regexValidator
	Pattern   string
}

// 场景
type Scence string

// 验证规则集 - 单个场景
type ScenceRules []Rule

// 验证规则集 - 所有场景
//
// var rules = Rules{
// 	"create": {
// 		{Attr: []string{"Username", "Password"}, Rule: "required", Message: "用户名或密码不能为空"},
// 		{Attr: "Password", Rule: "string", Max: 18, Min: 6},
// 		{Attr: "Age", Rule: "number", Required: true},
// 	},
// 	"get": {
// 		{Attr: "Username", Rule: "required", Message: "{Label}不能为空"},
// 	},
// }
type Rules map[Scence]ScenceRules

type M map[string]interface{}

// 验证错误类型
type E map[string]error

// 验证器函数类型
type F func(string, Rule, M) E

// 私有，使用构造器初始化，方便链式调用
type validator struct {
	// 默认语言
	lang string
	// 默认错误
	default_errors map[string]string
	// 验证器
	validators map[string]F
	// 错误收集器
	errors []E
}

// New 构造器，validator.New()
func New() *validator {
	this := &validator{
		lang: DEFAULT_LANG,
		default_errors: i18n.Errors[DEFAULT_LANG],
	}
	// 挂载内置验证器
	this.mount()
	return this
}

// Lang 设置默认的错误信息语言
func (this *validator) Lang(lang string) *validator {
	l := strings.ToUpper(lang)
	if _, ok := i18n.Errors[l]; !ok {
		panic(lang + " unsupport language")
	}
	this.lang = l
	this.default_errors = i18n.Errors[l]
	return this
}

// AddValidator 自定义验证器，新增验证规则，方便扩展
func (this *validator) AddValidator() {}

// Validate 场景验证
func (this *validator) Validate(rules Rules, obj M, scence Scence) []E {
	// 清空 errors
	this.errors = make([]E, 0)
	// 场景不存在
	scenceRules, ok := rules[scence]
	if !ok {
		panic(fmt.Sprint(scence) + " scence undefined")
	}
	// 验证
	for _, rule := range scenceRules {
		this.dispatch(rule, obj)
	}
	// toolbox.Dump(this)
	// toolbox.Dump(i18n.Errors)
	// toolbox.Dump(this.errors)
	return this.errors
}

// dispatch 验证调度器
func (this *validator) dispatch(rule Rule, obj M) {
	name := rule.Rule
	// Rule.Rule 未定义
	if name == "" {
		panic(errors.New(fmt.Sprint(rule) + " attribute 'Rule' not found"))
	}
	// 验证器不存在
	f, ok := this.validators[name]
	if !ok {
		panic(errors.New(name + " validator undefined"))
	}
	// Rule.Attr 未定义
	attr := rule.Attr
	if attr == nil {
		panic(errors.New(fmt.Sprint(rule) + " attribute 'Attr' not found"))
	}
	// Rule.Attr 类型错误
	switch attr.(type) {
		case string:
			this.adapter(f, rule, obj, false)
		case []string:
			this.adapter(f, rule, obj, true)
		default:
			panic("attribute 'Attr' should be 'string' or '[]string'")
	}
}

// adapter 多字段适配器
func (this *validator) adapter(f F, rule Rule, obj M, ismultiple bool) {
	// 多字段
	if ismultiple {
		for _, attr := range rule.Attr.([]string) {
			this.validate(f, attr, rule, obj)
		}
		return 
	}
	// 单字段
	this.validate(f, rule.Attr.(string), rule, obj)
}

// validate 验证
func (this *validator) validate(f F, attr string, rule Rule, obj M) {
	if e := f(attr, rule, obj); e != nil {
		this.errors = append(this.errors, e)
	}
}

// generator 错误信息生成器
func (this *validator) generator(name string, attr string, rule Rule, placeholder ...interface{}) E {
	ok := false
	e := rule.Message
	// 自定义错误信息
	if e != "" {
		return E{attr: errors.New(e)}
	}
	// 内置错误信息
	e, ok = this.default_errors[name]
	if ok {
		// 替换标签
		e = strings.Replace(e, "{label}", attr, -1)
		// 替换占位符
		e = fmt.Sprintf(e, placeholder...)
	} else { // 内置错误信息不存在
		e = "unknow error"
	}
	return E{attr: errors.New(e)}
}

// mount 挂载内置验证器
func (this *validator) mount() {
	this.validators = map[string]F {
		"required": this.requiredValidator,
		"in": this.inValidator,
		"string": this.stringValidator,
		"integer": this.integerValidator,
		"decimal": this.decimalValidator,
		"number": this.numberValidator,
		"boolean": this.booleanValidator,
		"ip": this.ipValidator,
		"regex": this.regexValidator,
		"email": this.emailValidator,
		"tel": this.telValidator,
		"mobile": this.mobileValidator,
		"zipcode": this.zipcodeValidator,
		// 别名
		"int": this.integerValidator, // integer
		"float": this.decimalValidator, // decimal
		"bool": this.booleanValidator, // boolean
		"phone": this.mobileValidator, // mobile
	}
}

// requiredValidator 必填验证器
func (this *validator) requiredValidator(attr string, rule Rule, obj M) E {
	if _, ok := obj[attr]; ok {
		return nil
	}
	return this.generator("required", attr, rule)
}

// inValidator 枚举验证器
// 支持类型 int64、int32、int16、int8、int、float64、float32、string、bool
// Rule.Required    bool        可选    false(默认) - 被验证字段有值验证/无值跳过，true - 被验证字段无值，验证失败，报 reqired 错误
// Rule.Enum        []string    必须    被验证字段必须在 Rule.Enum 中
func (this *validator) inValidator(attr string, rule Rule, obj M) E {
	enum := rule.Enum
	if len(enum) == 0 {
		panic(errors.New(fmt.Sprint(rule) + " attribute 'Enum' not found or empty"))
	}
	// 必填检测
	if _, ok := obj[attr]; !ok {
		if !rule.Required {	// 允许为空
			return nil
		}
		return this.generator("required", attr, rule)
	}
	// 枚举检测
	var field string
	switch v := obj[attr].(type) {
		case int64:
			field = strconv.Itoa(int(v))
		case int32:
			field = strconv.Itoa(int(v))
		case int16:
			field = strconv.Itoa(int(v))
		case int8:
			field = strconv.Itoa(int(v))
		case int:
			field = strconv.Itoa(v)
		case float64:
			field = strconv.FormatFloat(v, 'f', -1, 64)
		case float32:
			field = strconv.FormatFloat(float64(v), 'f', -1, 32)
		case string:
			field = v
		case bool:
			field = strconv.FormatBool(v)
		default:
			return this.generator("inValid", attr, rule)
	}
	in := false
	for _, v := range enum {
		if v == field {
			in = true
			break
		}
	}
	if !in {
		return this.generator("in", attr, rule, "[" +strings.Join(enum, "、") + "]")
	}
	return nil
}

// stringValidator 字符串验证器
// Rule.Required    bool    可选    false(默认) - 被验证字段有值验证/无值跳过，true - 被验证字段无值，验证失败，报 reqired 错误
// Rule.Max         int     可选    被验证字段长度不能大于 Rule.Max
// Rule.Min         in      可选    被验证字段长度不能小于 Rule.Min
func (this *validator) stringValidator(attr string, rule Rule, obj M) E {
	// 必填检测
	if _, ok := obj[attr]; !ok {
		if !rule.Required {	// 允许为空
			return nil
		}
		return this.generator("required", attr, rule)
	}
	// 类型检测
	if reflect.ValueOf(obj[attr]).Kind() != reflect.String {
		return this.generator("string", attr, rule)
	}
	// 长度检测
	max := rule.Max
	min := rule.Min
	if max != nil || min != nil {
		// 逻辑错误
		if max != nil && reflect.ValueOf(max).Kind() != reflect.Int{
			panic(errors.New(fmt.Sprint(rule) + " attribute 'Max' should be int"))
		}
		if min != nil && reflect.ValueOf(min).Kind() != reflect.Int{
			panic(errors.New(fmt.Sprint(rule) + " attribute 'Min' should be int"))
		}
		if max != nil && min != nil && min.(int) > max.(int) {
			panic(errors.New(fmt.Sprint(rule) + " attribute 'Max' should greater than 'Min'"))
		}
		// 比较
		length := int(utf8.RuneCountInString(obj[attr].(string)))
		if max != nil && min == nil && length > max.(int) { // only Max
			return this.generator("stringLengthMax", attr, rule, max)
		}
		if min != nil && max == nil && length < min.(int) { // only Min
			return this.generator("stringLengthMin", attr, rule, min)
		}
		if max != nil && min != nil && (length > max.(int) || length < min.(int)) { // both
			if max != min {
				return this.generator("stringLengthRange", attr, rule, min, max) // range
			}
			return this.generator("stringLengthEqual", attr, rule, max) // euqal
		}
	}
	return nil
}

// intValidator 整数验证器（整数/无小数位的浮点数/整数字符串）
// 支持类型 int64、int32、int16、int8、int、float64、float32、string
// Rule.Required    bool     可选    false(默认) - 被验证字段有值验证/无值跳过，true - 被验证字段无值，验证失败，报 reqired 错误
// Rule.Symbol      int64    可选    0(默认) - 正/负数，>0 - 正数(不包含0)，<0 - 负数(不包含0)
// Rule.Max         int      可选    被验证字段大小不能大于 Rule.Max
// Rule.Min         int      可选    被验证字段大小不能小于 Rule.Min
func (this *validator) integerValidator(attr string, rule Rule, obj M) E {
	// 必填检测
	if _, ok := obj[attr]; !ok {
		if !rule.Required {	// 允许为空
			return nil
		}
		return this.generator("required", attr, rule)
	}
	// 类型检测
	var field float64
	switch v := obj[attr].(type) {
		case int64:
			field = float64(v)
		case int32:
			field = float64(v)
		case int16:
			field = float64(v)
		case int8:
			field = float64(v)
		case int:
			field = float64(v)
		case float64:
			if v - float64(int(v)) != 0 { // 带小数位
				return this.generator("integer", attr, rule)	
			}
			field = v
		case float32:
			if v - float32(int(v)) != 0 { // 带小数位
				return this.generator("integer", attr, rule)	
			}
			field = float64(v)
		case string:
			f, err := strconv.ParseFloat(v, 64)
			if err != nil { // 不能转换为 float64
				return this.generator("integer", attr, rule)
			}
			if f - float64(int(f)) != 0 { // 带小数位
				return this.generator("integer", attr, rule)	
			}
			field = f
		default:
			return this.generator("integer", attr, rule)
	}
	// 正负检测
	symbol := rule.Symbol
	if (symbol > 0 && field <= 0) || (symbol < 0 && field >= 0) {
		if symbol > 0 {
			return this.generator("integerPositive", attr, rule)
		}
		return this.generator("integerNegative", attr, rule)
	}
	// 大小检测
	max := rule.Max
	min := rule.Min
	if max != nil || min != nil {
		errPrefix := "integer"
		// 逻辑错误
		if max != nil && reflect.ValueOf(max).Kind() != reflect.Int {
			panic(errors.New(fmt.Sprint(rule) + " attribute 'Max' should be int"))
		}
		if min != nil && reflect.ValueOf(min).Kind() != reflect.Int {
			panic(errors.New(fmt.Sprint(rule) + " attribute 'Min' should be int"))
		}
		if max != nil && min != nil && min.(int) > max.(int) {
			panic(errors.New(fmt.Sprint(rule) + " attribute 'Max' should greater than 'Min'"))
		}
		if symbol > 0 { 
			errPrefix += "Positive"
			if max != nil && max.(int) <= 0 { // 要求被检测属性是正数，而最大值被设置成负数，panic
				panic(errors.New(fmt.Sprint(rule) + " attribute 'Max' should greater than 0 when 'Symbal' = " + strconv.FormatInt(symbol, 10)))
			}
		}
		if symbol < 0 {
			errPrefix += "Negative"
			if min != nil && min.(int) >= 0 { // 要求被检测属性是负数，而最小值被设置成正数，panic
				panic(errors.New(fmt.Sprint(rule) + " attribute 'Min' should less than 0 when 'Symbal' = " + strconv.FormatInt(symbol, 10)))
			}
		}
		// 比较
		if max != nil && min == nil && field > (float64)(max.(int)) { // only Max
			return this.generator(errPrefix + "Max", attr, rule, max)
		}
		if min != nil && max == nil && field < (float64)(min.(int)) { // only Min
			return this.generator(errPrefix + "Min", attr, rule, min)
		}
		if max != nil && min != nil && (field > (float64)(max.(int)) || field < (float64)(min.(int))) { // both
			if max != min {	// range
				return this.generator(errPrefix + "Range", attr, rule, min, max)
			}
			return this.generator("equal", attr, rule, max) // euqal
		}
	}
	return nil
}

// decimalValidator 小数验证器（有小数位的浮点数/有小数位的浮点数字符串）
// 支持类型 float64、float32、string
// Rule.Required    bool           可选    false(默认) - 被验证字段有值验证/无值跳过，true - 被验证字段无值，验证失败，报 reqired 错误
// Rule.Symbol      int64          可选    0(默认) - 正/负数，>0 - 正数(不包含0)，<0 - 负数(不包含0)
// Rule.Max         int|float64    可选    被验证字段大小不能大于 Rule.Max
// Rule.Min         int|float64    可选    被验证字段大小不能小于 Rule.Min
func (this *validator) decimalValidator(attr string, rule Rule, obj M) E {
	// 必填检测
	if _, ok := obj[attr]; !ok {
		if !rule.Required {	// 允许为空
			return nil
		}
		return this.generator("required", attr, rule)
	}
	// 类型检测
	var field float64
	switch v := obj[attr].(type) {
		case float64:
			field = v
		case float32:
			field = float64(v)
		case string:
			f, err := strconv.ParseFloat(v, 64)
			if err != nil { // 不能转换为 float64
				return this.generator("decimal", attr, rule)
			}
			field = f
		default:
			return this.generator("decimal", attr, rule)
	}
	// 不带小数位
	if field - float64(int(field)) == 0 {
		return this.generator("decimal", attr, rule)	
	}
	// 正负检测
	symbol := rule.Symbol
	if (symbol > 0 && field <= 0) || (symbol < 0 && field >= 0) {
		if symbol > 0 {
			return this.generator("decimalPositive", attr, rule)
		}
		return this.generator("decimalNegative", attr, rule)
	}
	// 大小检测
	max := rule.Max
	min := rule.Min
	var fmax, fmin float64
	if max != nil || min != nil {
		errPrefix := "decimal"
		// 逻辑错误
		if max != nil {
			switch v := max.(type) {
				case float64:
					fmax = v
				case int:
					fmax = float64(v)
				default:
					panic(errors.New(fmt.Sprint(rule) + " attribute 'Max' should be int or float64"))
			}
		}
		if min != nil {
			switch v := min.(type) {
				case float64:
					fmin = v
				case int:
					fmin = float64(v)
				default:
					panic(errors.New(fmt.Sprint(rule) + " attribute 'Min' should be int or float64"))
			}
		}
		if max != nil && min != nil && fmin > fmax {
			panic(errors.New(fmt.Sprint(rule) + " attribute 'Max' should greater than 'Min'"))
		}
		if symbol > 0 { 
			errPrefix += "Positive"
			if max != nil && fmax <= 0 { // 要求被检测属性是正数，而最大值被设置成负数，panic
				panic(errors.New(fmt.Sprint(rule) + " attribute 'Max' should greater than 0 when 'Symbal' = " + strconv.FormatInt(symbol, 10)))
			}
		}
		if symbol < 0 {
			errPrefix += "Negative"
			if min != nil && fmin >= 0 { // 要求被检测属性是负数，而最小值被设置成正数，panic
				panic(errors.New(fmt.Sprint(rule) + " attribute 'Min' should less than 0 when 'Symbal' = " + strconv.FormatInt(symbol, 10)))
			}
		}
		// 比较
		if max != nil && min == nil && field > fmax{ // only Max
			return this.generator(errPrefix + "Max", attr, rule, max)
		}
		if min != nil && max == nil && field < fmin { // only Min
			return this.generator(errPrefix + "Min", attr, rule, min)
		}
		if max != nil && min != nil && (field > fmax || field < fmin) { // both
			if max != min {	// range
				return this.generator(errPrefix + "Range", attr, rule, min, max)
			}
			return this.generator("equal", attr, rule, max) // euqal
		}
	}
	return nil
}

// numberValidator 数字验证器（整数/浮点数/数字字符串）
// 支持类型 int64、int32、int16、int8、int、float64、float32、string
// Rule.Required    bool           可选    false(默认) - 被验证字段有值验证/无值跳过，true - 被验证字段无值，验证失败，报 reqired 错误
// Rule.Symbol      int64          可选    0(默认) - 正/负数，>0 - 正数(不包含0)，<0 - 负数(不包含0)
// Rule.Max         int|float64    可选    被验证字段大小不能大于 Rule.Max
// Rule.Min         int|float64    可选    被验证字段大小不能小于 Rule.Min
func (this *validator) numberValidator(attr string, rule Rule, obj M) E {
	// 必填检测
	if _, ok := obj[attr]; !ok {
		if !rule.Required {	// 允许为空
			return nil
		}
		return this.generator("required", attr, rule)
	}
	// 类型检测
	var field float64
	switch v := obj[attr].(type) {
		case int64:
			field = float64(v) 
		case int32:
			field = float64(v)
		case int16:
			field = float64(v)
		case int8:
			field = float64(v)
		case int:
			field = float64(v)
		case float64:
			field = v
		case float32:
			field = float64(v)
		case string:
			f, err := strconv.ParseFloat(v, 64)
			if err != nil { // 不能转换为 float64
				return this.generator("number", attr, rule)
			}
			field = f
		default:
			return this.generator("number", attr, rule)
	}
	// 正负检测
	symbol := rule.Symbol
	if (symbol > 0 && field <= 0) || (symbol < 0 && field >= 0) {
		if symbol > 0 {
			return this.generator("numberPositive", attr, rule)
		}
		return this.generator("numberNegative", attr, rule)
	}
	// 大小检测
	max := rule.Max
	min := rule.Min
	var fmax, fmin float64
	if max != nil || min != nil {
		errPrefix := "number"
		// 逻辑错误
		if max != nil {
			switch v := max.(type) {
				case float64:
					fmax = v
				case int:
					fmax = float64(v)
				default:
					panic(errors.New(fmt.Sprint(rule) + " attribute 'Max' should be int or float64"))
			}
		}
		if min != nil {
			switch v := min.(type) {
				case float64:
					fmin = v
				case int:
					fmin = float64(v)
				default:
					panic(errors.New(fmt.Sprint(rule) + " attribute 'Min' should be int or float64"))
			}
		}
		if max != nil && min != nil && fmin > fmax {
			panic(errors.New(fmt.Sprint(rule) + " attribute 'Max' should greater than 'Min'"))
		}
		if symbol > 0 { 
			errPrefix += "Positive"
			if max != nil && fmax <= 0 { // 要求被检测属性是正数，而最大值被设置成负数，panic
				panic(errors.New(fmt.Sprint(rule) + " attribute 'Max' should greater than 0 when 'Symbal' = " + strconv.FormatInt(symbol, 10)))
			}
		}
		if symbol < 0 {
			errPrefix += "Negative"
			if min != nil && fmin >= 0 { // 要求被检测属性是负数，而最小值被设置成正数，panic
				panic(errors.New(fmt.Sprint(rule) + " attribute 'Min' should less than 0 when 'Symbal' = " + strconv.FormatInt(symbol, 10)))
			}
		}
		// 比较
		if max != nil && min == nil && field > fmax{ // only Max
			return this.generator(errPrefix + "Max", attr, rule, max)
		}
		if min != nil && max == nil && field < fmin { // only Min
			return this.generator(errPrefix + "Min", attr, rule, min)
		}
		if max != nil && min != nil && (field > fmax || field < fmin) { // both
			if max != min {	// range
				return this.generator(errPrefix + "Range", attr, rule, min, max)
			}
			return this.generator("equal", attr, rule, max) // euqal
		}
	}
 	return nil
}

// boolValidator 布尔验证器（布尔值/字符串表示的布尔值[1、0、t、f、true、false(忽略大小写)]）
// 支持类型 bool、string
// Rule.Required    bool    可选    false(默认) - 被验证字段有值验证/无值跳过，true - 被验证字段无值，验证失败，报 reqired 错误
func (this *validator) booleanValidator(attr string, rule Rule, obj M) E {
	// 必填检测
	if _, ok := obj[attr]; !ok {
		if !rule.Required {	// 允许为空
			return nil
		}
		return this.generator("required", attr, rule)
	}
	// 类型检测
	switch v := obj[attr].(type) {
		case bool:
			return nil
		case string:
			if _, err := strconv.ParseBool(strings.ToLower(v)); err != nil {
				return this.generator("boolean", attr, rule)
			}
		default:
			return this.generator("boolean", attr, rule)
	}
	return nil
}

// ipValidator ip 验证器
func (this *validator) ipValidator(attr string, rule Rule, obj M) E {
	// 必填检测
	if _, ok := obj[attr]; !ok {
		if !rule.Required {	// 允许为空
			return nil
		}
		return this.generator("required", attr, rule)
	}
	// 字符串检测
	if reflect.ValueOf(obj[attr]).Kind() != reflect.String {
		return this.generator("string", attr, rule)
	}
	// ip 检测
	if ip := net.ParseIP(obj[attr].(string)); ip == nil {
		return this.generator("ip", attr, rule)
	}
	return nil
}

// regexMatch 正则匹配
func (this *validator) regexMatch(attr string, rule Rule, obj M, kind string) E {
	pattern := rule.Pattern
	if len(pattern) == 0 {
		panic(errors.New(fmt.Sprint(rule) + " attribute 'Pattern' not found or empty"))
	}
	// 必填检测
	if _, ok := obj[attr]; !ok {
		if !rule.Required {	// 允许为空
			return nil
		}
		return this.generator("required", attr, rule)
	}
	// 字符串检测
	if reflect.ValueOf(obj[attr]).Kind() != reflect.String {
		return this.generator("string", attr, rule)
	}
	// 正则检测
	regex := regexp.MustCompile(pattern)
	if !regex.MatchString(obj[attr].(string)) {
		return this.generator(kind, attr, rule)	
	}
	return nil
}

// regexValidator 正则验证器
func (this *validator) regexValidator(attr string, rule Rule, obj M) E {
	return this.regexMatch(attr, rule, obj, "regex")
}

// emailValidator 邮箱验证器
func (this *validator) emailValidator(attr string, rule Rule, obj M) E {
	rule.Pattern = PATTERN_EMAIL
	return this.regexMatch(attr, rule, obj, "email")
}

// mobileValidator 中国大陆座机号验证器
func (this *validator) telValidator(attr string, rule Rule, obj M) E {
	rule.Pattern = PATTERN_TEL
	return this.regexMatch(attr, rule, obj, "tel")
}

// mobileValidator 中国大陆手机号验证器
func (this *validator) mobileValidator(attr string, rule Rule, obj M) E {
	rule.Pattern = PATTERN_MOBILE
	return this.regexMatch(attr, rule, obj, "mobile")
}

// mobileValidator 中国大陆邮编验证器
func (this *validator) zipcodeValidator(attr string, rule Rule, obj M) E {
	rule.Pattern = PATTERN_ZIPCODE
	return this.regexMatch(attr, rule, obj, "zipcode")
}