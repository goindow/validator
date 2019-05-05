package i18n

const EN_US = "EN_US"

func init() {
    Errors[EN_US] = errors{
         // common
        "equal": "{label} must be euqal to %v",
        // requiredValidator
        "required": "{label} can not be empty",
        // inValidator
        "in": "{label} must be in %v",
        "inValid": "{label} must be one of string, number, boolean",
        // stringValidator
        "string": "{label} must be a string",
        "stringLengthMax": "{label}'s maximum length is %v",
        "stringLengthMin": "{label}'s minimum length is %v",
        "stringLengthRange": "{label}'s length is %v to %v",
        "stringLengthEqual": "{label}'s length must be equal to %v",
        // intValidator
        "integer": "{label} must be an integer",
        "integerMax": "{label} must be an integer with a maximum value of %v",
        "integerMin": "{label} must be an integer with a minimum value of %v",
        "integerRange": "{label} must be an integer of %v to %v",
        "integerPositive": "{label} must be a positive integer",
        "integerPositiveMax": "{label} must be a positive integer with a maximum value of %v",
        "integerPositiveMin": "{label} must be a positive integer with a minimum value of %v",
        "integerPositiveRange": "{label} must be a positive integer of %v to %v",
        "integerNegative": "{label} must be a negative integer",
        "integerNegativeMax": "{label} must be a negative integer with a maximum value of %v",
        "integerNegativeMin": "{label} must be a negative integer with a minimum value of %v",
        "integerNegativeRange": "{label} must be a negative integer of %v to %v",
        // floatValidator
        "decimal": "{label} must be a decimal",
        "decimalMax": "{label} must be a decimal with a maximum value of %v",
        "decimalMin": "{label} must be a decimal with a minimum value of %v",
        "decimalRange": "{label} must be a decimal of %v to %v",
        "decimalPositive": "{label} must be a positive decimal",
        "decimalPositiveMax": "{label} must be a positive decimal with a maximum value of %v",
        "decimalPositiveMin": "{label} must be a positive decimal with a minimum value of %v",
        "decimalPositiveRange": "{label} must be a positive decimal of %v to %v",
        "decimalNegative": "{label} must be a negative decimal",
        "decimalNegativeMax": "{label} must be a negative decimal with a maximum value of %v",
        "decimalNegativeMin": "{label} must be a negative decimal with a minimum value of %v",
        "decimalNegativeRange": "{label} must be a negative decimal of %v to %v",
        // numberValidator
        "number": "{label} must be a number",
        "numberMax": "{label} must be a number with a maximum value of %v",
        "numberMin": "{label} must be a number with a minimum value of %v",
        "numberRange": "{label} must be a number of %v to %v",
        "numberPositive": "{label} must be a positive number",
        "numberPositiveMax": "{label} must be a positive number with a maximum value of %v",
        "numberPositiveMin": "{label} must be a positive number with a minimum value of %v",
        "numberPositiveRange": "{label} must be a positive number of %v to %v",
        "numberNegative": "{label} must be a negative number",
        "numberNegativeMax": "{label} must be a negative number with a maximum value of %v",
        "numberNegativeMin": "{label} must be a negative number with a minimum value of %v",
        "numberNegativeRange": "{label} must be a negative number of %v to %v",
        // booleanValidator
        "boolean": "{label} must be a boolean or string",
        // regexValidator
        "regex": "{label} must be in a valid format",
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