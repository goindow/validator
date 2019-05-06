package i18n

const EN_US = "EN_US"

func init() {
    Errors[EN_US] = errors{
         // common
        "equal": "must be euqal to %v",
        // requiredValidator
        "required": "can not be empty",
        // inValidator
        "in": "must be in %v",
        "inValid": "must be one of string, number, boolean",
        // stringValidator
        "string": "must be a string",
        "stringLengthMax": "{label}'s maximum length is %v",
        "stringLengthMin": "{label}'s minimum length is %v",
        "stringLengthRange": "{label}'s length is %v to %v",
        "stringLengthEqual": "{label}'s length must be equal to %v",
        // intValidator
        "integer": "must be an integer",
        "integerMax": "must be an integer with a maximum value of %v",
        "integerMin": "must be an integer with a minimum value of %v",
        "integerRange": "must be an integer of %v to %v",
        "integerPositive": "must be a positive integer",
        "integerPositiveMax": "must be a positive integer with a maximum value of %v",
        "integerPositiveMin": "must be a positive integer with a minimum value of %v",
        "integerPositiveRange": "must be a positive integer of %v to %v",
        "integerNegative": "must be a negative integer",
        "integerNegativeMax": "must be a negative integer with a maximum value of %v",
        "integerNegativeMin": "must be a negative integer with a minimum value of %v",
        "integerNegativeRange": "must be a negative integer of %v to %v",
        // floatValidator
        "decimal": "must be a decimal",
        "decimalMax": "must be a decimal with a maximum value of %v",
        "decimalMin": "must be a decimal with a minimum value of %v",
        "decimalRange": "must be a decimal of %v to %v",
        "decimalPositive": "must be a positive decimal",
        "decimalPositiveMax": "must be a positive decimal with a maximum value of %v",
        "decimalPositiveMin": "must be a positive decimal with a minimum value of %v",
        "decimalPositiveRange": "must be a positive decimal of %v to %v",
        "decimalNegative": "must be a negative decimal",
        "decimalNegativeMax": "must be a negative decimal with a maximum value of %v",
        "decimalNegativeMin": "must be a negative decimal with a minimum value of %v",
        "decimalNegativeRange": "must be a negative decimal of %v to %v",
        // numberValidator
        "number": "must be a number",
        "numberMax": "must be a number with a maximum value of %v",
        "numberMin": "must be a number with a minimum value of %v",
        "numberRange": "must be a number of %v to %v",
        "numberPositive": "must be a positive number",
        "numberPositiveMax": "must be a positive number with a maximum value of %v",
        "numberPositiveMin": "must be a positive number with a minimum value of %v",
        "numberPositiveRange": "must be a positive number of %v to %v",
        "numberNegative": "must be a negative number",
        "numberNegativeMax": "must be a negative number with a maximum value of %v",
        "numberNegativeMin": "must be a negative number with a minimum value of %v",
        "numberNegativeRange": "must be a negative number of %v to %v",
        // booleanValidator
        "boolean": "must be a boolean or string",
        // regexValidator
        "regex": "must be in a valid format",
        // ipValidator
        "ip": "must be a valid ip address",
        // emailValidator
        "email": "must be a valid email address",
        // telValidator
        "tel": "must be a valid telephone number",
        // mobileValidator
        "mobile": "must be a valid telephone or mobile phone number",
        // zipcodeValidator
        "zipcode": "must be a valid zipcode",
    }
}