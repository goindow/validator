package i18n

const ZH_CN = "ZH_CN"

func init() {
    Errors[ZH_CN] = errors{
        // common
        "equal": "必须是 %v",
        // requiredValidator
        "required": "不能为空",
        // inValidator
        "in": "只能是 %v 中的一个",
        "inValid": "必须是字符串、数字、布尔值中的一种",
        // stringValidator
        "string": "必须是字符串",
        "stringLengthMax": "长度不能超过 %v",
        "stringLengthMin": "长度不能小于 %v",
        "stringLengthRange": "长度必须在 %v 到 %v 之间",
        "stringLengthEqual": "长度必须是 %v",
        // intValidator
        "integer": "必须是整数",
        "integerMax": "必须是不大于 %v 的整数",
        "integerMin": "必须是不小于 %v 的整数",
        "integerRange": "必须是介于 %v 到 %v 的整数",
        "integerPositive": "必须是正整数",
        "integerPositiveMax": "必须是不大于 %v 的正整数",
        "integerPositiveMin": "必须是不小于 %v 的正整数",
        "integerPositiveRange": "必须是介于 %v 到 %v 的正整数",
        "integerNegative": "必须是负整数",
        "integerNegativeMax": "必须是不大于 %v 的负整数",
        "integerNegativeMin": "必须是不小于 %v 的负整数",
        "integerNegativeRange": "必须是介于 %v 到 %v 的负整数",
        // floatValidator
        "decimal": "必须是小数",
        "decimalMax": "必须是不大于 %v 的小数",
        "decimalMin": "必须是不小于 %v 的小数",
        "decimalRange": "必须是介于 %v 到 %v 的小数",
        "decimalPositive": "必须是正小数",
        "decimalPositiveMax": "必须是不大于 %v 的正小数",
        "decimalPositiveMin": "必须是不小于 %v 的正小数",
        "decimalPositiveRange": "必须是介于 %v 到 %v 的正小数",
        "decimalNegative": "必须是负小数",
        "decimalNegativeMax": "必须是不大于 %v 的负小数",
        "decimalNegativeMin": "必须是不小于 %v 的负小数",
        "decimalNegativeRange": "必须是介于 %v 到 %v 的负小数",
        // numberValidator
        "number": "必须是数字",
        "numberMax": "必须是不大于 %v 的数",
        "numberMin": "必须是不小于 %v 的数",
        "numberRange": "必须是介于 %v 到 %v 的数",
        "numberPositive": "必须是正数",
        "numberPositiveMax": "必须是不大于 %v 的正数",
        "numberPositiveMin": "必须是不小于 %v 的正数",
        "numberPositiveRange": "必须是介于 %v 到 %v 的正数",
        "numberNegative": "必须是负数",
        "numberNegativeMax": "必须是不大于 %v 的负数",
        "numberNegativeMin": "必须是不小于 %v 的负数",
        "numberNegativeRange": "必须是介于 %v 到 %v 的负数",
        // booleanValidator
        "boolean": "必须是布尔值或布尔字符串",
        // regexValidator
        "regex": "格式不正确",
        // ipValidator
        "ip": "无效的 ip",
        // emailValidator
        "email": "无效的 email",
        // telValidator
        "tel": "无效的座机号",
        // mobileValidator
        "mobile": "无效的手机号",
        // zipcodeValidator
        "zipcode": "无效的邮编",
    }
}