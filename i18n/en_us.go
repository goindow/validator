package i18n

const EN_US = "EN_US"

func init() {
    Errors[EN_US] = errors{
        "required": "missing field {label}",
        "string": "need string type {label}",
    }
}