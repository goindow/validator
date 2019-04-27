package i18n

const ZH_CN = "ZH_CN"

func init() {
    Errors[ZH_CN] = errors{
        // common
        "equal": "{label} 必须是 %v",
        // requiredValidator
        "required": "{label} 不能为空",
        // inValidator
        "in": "{label} 只能是 %v 中的一个",
        "inValid": "{label} 必须是字符串、数字、布尔值中的一种",
        // stringValidator
        "string": "{label} 必须是字符串",
        "stringLengthMax": "{label} 长度不能超过 %v",
        "stringLengthMin": "{label} 长度不能小于 %v",
        "stringLengthRange": "{label} 长度必须在 %v 到 %v 之间",
        "stringLengthEqual": "{label} 长度必须是 %v",
        // intValidator
        "integer": "{label} 必须是整数",
        "integerMax": "{label} 必须是不大于 %v 的整数",
        "integerMin": "{label} 必须是不小于 %v 的整数",
        "integerRange": "{label} 必须是介于 %v 到 %v 的整数",
        "integerPositive": "{label} 必须是正整数",
        "integerPositiveMax": "{label} 必须是不大于 %v 的正整数",
        "integerPositiveMin": "{label} 必须是不小于 %v 的正整数",
        "integerPositiveRange": "{label} 必须是介于 %v 到 %v 的正整数",
        "integerNegative": "{label} 必须是负整数",
        "integerNegativeMax": "{label} 必须是不大于 %v 的负整数",
        "integerNegativeMin": "{label} 必须是不小于 %v 的负整数",
        "integerNegativeRange": "{label} 必须是介于 %v 到 %v 的负整数",
        // floatValidator
        "decimal": "{label} 必须是小数",
        "decimalMax": "{label} 必须是不大于 %v 的小数",
        "decimalMin": "{label} 必须是不小于 %v 的小数",
        "decimalRange": "{label} 必须是介于 %v 到 %v 的小数",
        "decimalPositive": "{label} 必须是正小数",
        "decimalPositiveMax": "{label} 必须是不大于 %v 的正小数",
        "decimalPositiveMin": "{label} 必须是不小于 %v 的正小数",
        "decimalPositiveRange": "{label} 必须是介于 %v 到 %v 的正小数",
        "decimalNegative": "{label} 必须是负小数",
        "decimalNegativeMax": "{label} 必须是不大于 %v 的负小数",
        "decimalNegativeMin": "{label} 必须是不小于 %v 的负小数",
        "decimalNegativeRange": "{label} 必须是介于 %v 到 %v 的负小数",
        // numberValidator
        "number": "{label} 必须是数字",
        "numberMax": "{label} 必须是不大于 %v 的数",
        "numberMin": "{label} 必须是不小于 %v 的数",
        "numberRange": "{label} 必须是介于 %v 到 %v 的数",
        "numberPositive": "{label} 必须是正数",
        "numberPositiveMax": "{label} 必须是不大于 %v 的正数",
        "numberPositiveMin": "{label} 必须是不小于 %v 的正数",
        "numberPositiveRange": "{label} 必须是介于 %v 到 %v 的正数",
        "numberNegative": "{label} 必须是负数",
        "numberNegativeMax": "{label} 必须是不大于 %v 的负数",
        "numberNegativeMin": "{label} 必须是不小于 %v 的负数",
        "numberNegativeRange": "{label} 必须是介于 %v 到 %v 的负数",
        // booleanValidator
        "boolean": "{label} 必须是布尔值或布尔字符串",
        // regexValidator
        "regex": "{label} 格式不正确",
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