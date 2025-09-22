package validation

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/aasoft24/golara/wpkg/database"

	"github.com/aasoft24/golara/wpkg/gola"

	"gorm.io/gorm"
)

// Validator structure
type Validator struct {
	Data           map[string]interface{}
	Errors         map[string][]string
	DB             *gorm.DB
	CustomMessages map[string]string // Custom error messages
}

// Create new validator
func NewValidator(data map[string]interface{}, db *gorm.DB) *Validator {
	return &Validator{
		Data:   data,
		Errors: make(map[string][]string),
		DB:     db,
	}
}

// Validate function with Laravel-style rules

func (v *Validator) Validate(rules map[string]string, customMessages ...map[string]string) bool {
	// Set custom messages if provided
	if len(customMessages) > 0 {
		v.CustomMessages = customMessages[0]
	}

	for field, ruleStr := range rules {
		rulesArr := strings.Split(ruleStr, "|")
		value, exists := v.Data[field]

		// Handle []string values
		var values []interface{}
		switch val := value.(type) {
		case []string:
			for _, v := range val {
				values = append(values, v)
			}
		case []interface{}:
			values = val
		default:
			values = []interface{}{value}
		}

		// Validate each value (for multi or single)
		for _, singleVal := range values {
			for _, rule := range rulesArr {
				if rule == "" {
					continue
				}

				// Check for sometimes rule
				if rule == "sometimes" {
					if !exists {
						continue
					}
					continue
				}

				parts := strings.SplitN(rule, ":", 2)
				ruleName := parts[0]
				ruleValue := ""
				if len(parts) > 1 {
					ruleValue = parts[1]
				}

				if !exists && ruleName != "required" {
					continue
				}

				switch ruleName {
				case "required":
					if !v.validateRequired(singleVal) {
						v.addError(field, "The %s field is required", field)
					}
				case "email":
					if exists && !v.validateEmail(singleVal) {
						v.addError(field, "The %s must be a valid email address", field)
					}
				case "min":
					if exists && !v.validateMin(singleVal, ruleValue) {
						v.addError(field, "The %s must be at least %s characters", field, ruleValue)
					}
				case "max":
					if exists && !v.validateMax(singleVal, ruleValue) {
						v.addError(field, "The %s may not be greater than %s characters", field, ruleValue)
					}
				case "len":
					if exists && !v.validateLen(singleVal, ruleValue) {
						v.addError(field, "The %s must be %s characters", field, ruleValue)
					}
				case "numeric":
					if exists && !v.validateNumeric(singleVal) {
						v.addError(field, "The %s must be a number", field)
					}
				case "same":
					if exists && !v.validateSame(singleVal, v.Data[ruleValue]) {
						v.addError(field, "The %s and %s must match", field, ruleValue)
					}
				case "alpha":
					if exists && !v.validateAlpha(fmt.Sprintf("%v", singleVal)) {
						v.addError(field, "The %s may only contain letters", field)
					}
				case "alpha_num":
					if exists && !v.validateAlphaNum(fmt.Sprintf("%v", singleVal)) {
						v.addError(field, "The %s may only contain letters and numbers", field)
					}
				case "in":
					if exists && !v.validateIn(singleVal, ruleValue) {
						v.addError(field, "The selected %s is invalid", field)
					}
				case "not_in":
					if exists && !v.validateNotIn(singleVal, ruleValue) {
						v.addError(field, "The selected %s is invalid", field)
					}
				case "regex":
					if exists && !v.validateRegex(singleVal, ruleValue) {
						v.addError(field, "The %s format is invalid", field)
					}
				case "unique":
					if exists && !v.validateUnique(singleVal, ruleValue) {
						v.addError(field, "The %s has already been taken", field)
					}
				case "unique_multi":
					if exists && !v.validateUniqueMulti(singleVal, ruleValue) {
						v.addError(field, "The %s has already been taken", field)
					}
				case "unique_except":
					if exists && !v.validateUniqueExcept(singleVal, ruleValue) {
						v.addError(field, "The %s has already been taken", field)
					}
				case "unique_multi_except":
					if exists && !v.validateUniqueMultiExcept(singleVal, ruleValue) {
						v.addError(field, "The %s has already been taken", field)
					}
				}
			}
		}
	}

	return len(v.Errors) == 0
}

// func (v *Validator) Validate(rules map[string]string, customMessages ...map[string]string) bool {
// 	// Set custom messages if provided
// 	if len(customMessages) > 0 {
// 		v.CustomMessages = customMessages[0]
// 	}

// 	for field, ruleStr := range rules {
// 		rules := strings.Split(ruleStr, "|")
// 		value, exists := v.Data[field]

// 		for _, rule := range rules {
// 			if rule == "" {
// 				continue
// 			}

// 			// Check for sometimes rule
// 			if rule == "sometimes" {
// 				if !exists {
// 					continue // Skip validation if field doesn't exist
// 				}
// 				continue
// 			}

// 			parts := strings.SplitN(rule, ":", 2)
// 			ruleName := parts[0]
// 			ruleValue := ""
// 			if len(parts) > 1 {
// 				ruleValue = parts[1]
// 			}

// 			if !exists && ruleName != "required" {
// 				continue
// 			}

// 			switch ruleName {
// 			case "required":
// 				if !v.validateRequired(value) {
// 					v.addError(field, "The %s field is required", field)
// 				}
// 			case "email":
// 				if exists && !v.validateEmail(value) {
// 					v.addError(field, "The %s must be a valid email address", field)
// 				}
// 			case "min":
// 				if exists && !v.validateMin(value, ruleValue) {
// 					v.addError(field, "The %s must be at least %s characters", field, ruleValue)
// 				}
// 			case "max":
// 				if exists && !v.validateMax(value, ruleValue) {
// 					v.addError(field, "The %s may not be greater than %s characters", field, ruleValue)
// 				}
// 			case "len":
// 				if exists && !v.validateLen(value, ruleValue) {
// 					v.addError(field, "The %s must be %s characters", field, ruleValue)
// 				}
// 			case "numeric":
// 				if exists && !v.validateNumeric(value) {
// 					v.addError(field, "The %s must be a number", field)
// 				}
// 			case "same":
// 				if exists && !v.validateSame(value, v.Data[ruleValue]) {
// 					v.addError(field, "The %s and %s must match", field, ruleValue)
// 				}
// 			case "alpha":
// 				if exists && !v.validateAlpha(value.(string)) {
// 					v.addError(field, "The %s may only contain letters", field)
// 				}
// 			case "alpha_num":
// 				if exists && !v.validateAlphaNum(value.(string)) {
// 					v.addError(field, "The %s may only contain letters and numbers", field)
// 				}
// 			case "in":
// 				if exists && !v.validateIn(value, ruleValue) {
// 					v.addError(field, "The selected %s is invalid", field)
// 				}
// 			case "not_in":
// 				if exists && !v.validateNotIn(value, ruleValue) {
// 					v.addError(field, "The selected %s is invalid", field)
// 				}
// 			case "regex":
// 				if exists && !v.validateRegex(value, ruleValue) {
// 					v.addError(field, "The %s format is invalid", field)
// 				}
// 			case "unique":
// 				if exists && !v.validateUnique(value, ruleValue) {
// 					v.addError(field, "The %s has already been taken", field)
// 				}
// 			case "unique_multi":
// 				if exists && !v.validateUniqueMulti(value, ruleValue) {
// 					v.addError(field, "The %s has already been taken", field)
// 				}
// 			case "unique_except":
// 				if exists && !v.validateUniqueExcept(value, ruleValue) {
// 					v.addError(field, "The %s has already been taken", field)
// 				}
// 			case "unique_multi_except":
// 				if exists && !v.validateUniqueMultiExcept(value, ruleValue) {
// 					v.addError(field, "The %s has already been taken", field)
// 				}
// 			}
// 		}
// 	}

// 	return len(v.Errors) == 0
// }

// === Validators ===
func (v *Validator) validateRequired(value interface{}) bool {
	if value == nil {
		return false
	}
	if str, ok := value.(string); ok {
		return strings.TrimSpace(str) != ""
	}
	return true
}

func (v *Validator) validateEmail(value interface{}) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(str)
}

func (v *Validator) validateMin(value interface{}, minStr string) bool {
	min := parseInt(minStr)
	switch val := value.(type) {
	case string:
		return utf8.RuneCountInString(val) >= min
	case int:
		return val >= min
	case float64:
		return val >= float64(min)
	default:
		return false
	}
}

func (v *Validator) validateMax(value interface{}, maxStr string) bool {
	max := parseInt(maxStr)
	switch val := value.(type) {
	case string:
		return utf8.RuneCountInString(val) <= max
	case int:
		return val <= max
	case float64:
		return val <= float64(max)
	default:
		return false
	}
}

func (v *Validator) validateLen(value interface{}, lenStr string) bool {
	length := parseInt(lenStr)
	switch val := value.(type) {
	case string:
		return utf8.RuneCountInString(val) == length
	case int:
		return val == length
	default:
		return false
	}
}

func (v *Validator) validateNumeric(value interface{}) bool {
	switch value.(type) {
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		return true
	case string:
		str := value.(string)
		matched, _ := regexp.MatchString(`^\-?\d+(\.\d+)?$`, str)
		return matched
	default:
		return false
	}
}

func (v *Validator) validateSame(value interface{}, other interface{}) bool {
	return fmt.Sprintf("%v", value) == fmt.Sprintf("%v", other)
}

// validateAlpha: শুধুমাত্র a-zA-Z অনুমোদন করবে
func (v *Validator) validateAlpha(value string) bool {
	for _, c := range value {
		if (c < 'a' || c > 'z') && (c < 'A' || c > 'Z') {
			return false
		}
	}
	return true
}

// validateAlphaNum: শুধুমাত্র a-zA-Z0-9 অনুমোদন করবে
func (v *Validator) validateAlphaNum(value string) bool {
	for _, c := range value {
		if (c < 'a' || c > 'z') && (c < 'A' || c > 'Z') && (c < '0' || c > '9') {
			return false
		}
	}
	return true
}

// validateIn rule
func (v *Validator) validateIn(value interface{}, ruleValue string) bool {
	allowedValues := strings.Split(ruleValue, ",")
	strValue := fmt.Sprintf("%v", value)

	for _, allowed := range allowedValues {
		if strValue == allowed {
			return true
		}
	}
	return false
}

// validateNotIn rule
func (v *Validator) validateNotIn(value interface{}, ruleValue string) bool {
	return !v.validateIn(value, ruleValue)
}

// validateRegex rule
func (v *Validator) validateRegex(value interface{}, ruleValue string) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}

	matched, err := regexp.MatchString(ruleValue, str)
	return err == nil && matched
}

// Unique DB check
func (v *Validator) validateUnique(value interface{}, ruleValue string) bool {
	if v.DB == nil {
		return true
	}
	parts := strings.Split(ruleValue, ",")
	if len(parts) != 2 {
		return true
	}
	table := parts[0]
	column := parts[1]

	var count int64
	v.DB.Table(table).Where(fmt.Sprintf("%s = ?", column), value).Count(&count)
	return count == 0
}

// Validator struct e add korun
func (v *Validator) validateUniqueMulti(value interface{}, ruleValue string) bool {
	if v.DB == nil {
		return true
	}

	tables := strings.Split(ruleValue, ",")
	if len(tables) < 1 {
		return true
	}

	for _, tableColumn := range tables {
		parts := strings.Split(tableColumn, ".")
		if len(parts) != 2 {
			continue
		}
		table := parts[0]
		column := parts[1]

		var count int64
		v.DB.Table(table).Where(fmt.Sprintf("%s = ?", column), value).Count(&count)
		if count > 0 {
			return false
		}
	}
	return true
}

// Validator struct e add korun
func (v *Validator) validateUniqueExcept(value interface{}, ruleValue string) bool {
	if v.DB == nil {
		return true
	}

	// Format: table,column,id,except_id
	parts := strings.Split(ruleValue, ",")
	if len(parts) < 4 {
		return true
	}

	table := parts[0]
	column := parts[1]
	idColumn := parts[2]
	exceptID := parts[3]

	var count int64
	v.DB.Table(table).
		Where(fmt.Sprintf("%s = ?", column), value).
		Where(fmt.Sprintf("%s != ?", idColumn), exceptID).
		Count(&count)

	return count == 0
}

// validateUniqueMultiExcept for multiple tables
func (v *Validator) validateUniqueMultiExcept(value interface{}, ruleValue string) bool {
	if v.DB == nil {
		return true
	}

	// Format: table1.column1:id_column1:except_id1,table2.column2:id_column2:except_id2
	tablesRules := strings.Split(ruleValue, ",")

	for _, tableRule := range tablesRules {
		parts := strings.Split(tableRule, ":")
		if len(parts) < 3 {
			continue
		}

		tableColumn := strings.Split(parts[0], ".")
		if len(tableColumn) != 2 {
			continue
		}

		table := tableColumn[0]
		column := tableColumn[1]
		idColumn := parts[1]
		exceptID := parts[2]

		var count int64
		v.DB.Table(table).
			Where(fmt.Sprintf("%s = ?", column), value).
			Where(fmt.Sprintf("%s != ?", idColumn), exceptID).
			Count(&count)

		if count > 0 {
			return false
		}
	}
	return true
}

// Add error
func (v *Validator) addError(field, format string, args ...interface{}) {
	// Custom message check korun
	key := field + "." + strings.Split(format, " ")[0] // field.required, field.email, etc.
	if customMsg, exists := v.CustomMessages[key]; exists {
		message := fmt.Sprintf(customMsg, args...)
		v.Errors[field] = append(v.Errors[field], message)
	} else {
		message := fmt.Sprintf(format, args...)
		v.Errors[field] = append(v.Errors[field], message)
	}
}
func (v *Validator) GetErrors() map[string][]string {
	return v.Errors
}

func (v *Validator) HasErrors() bool {
	return len(v.Errors) > 0
}

func parseInt(s string) int {
	var n int
	fmt.Sscanf(s, "%d", &n)
	return n
}

// === ValidateRequest helper ===
// validation/validate_request.go
func notwork(c *gola.Context, db *gorm.DB, rules map[string]string) (map[string][]string, map[string]string) {
	data := make(map[string]interface{})
	old := make(map[string]string)

	for field := range rules {
		value := c.Input(field)
		data[field] = value
		old[field] = value
	}

	v := NewValidator(data, db)

	if !v.Validate(rules) {
		return v.GetErrors(), old // full map[string][]string
	}

	return nil, old
}

func ValidateRequest(c *gola.Context, rules map[string]string) (map[string]string, map[string]string) {
	data := make(map[string]interface{})
	old := make(map[string]string)

	for field := range rules {
		// 🔹 Check if field has multiple values
		if arr := c.PostFormArray(field); len(arr) > 0 {
			data[field] = arr
			old[field] = strings.Join(arr, ",") // old value as comma-separated string
		} else {
			// fallback single value
			value := c.Input(field)
			data[field] = value
			old[field] = value
		}
	}

	db := database.DB
	v := NewValidator(data, db)
	errors := make(map[string]string)

	if !v.Validate(rules) {
		for field, msgs := range v.GetErrors() {
			if len(msgs) > 0 {
				errors[field] = msgs[0]
			}
		}
	}

	return errors, old
}

// In validation package
// func ValidateRequest(c *gola.Context, rules map[string]string) (map[string]string, map[string]string) {
// 	data := make(map[string]interface{})
// 	old := make(map[string]string)

// 	for field := range rules {
// 		value := c.Input(field)
// 		data[field] = value
// 		old[field] = value
// 	}

// 	db := database.DB

// 	v := NewValidator(data, db)
// 	errors := make(map[string]string)

// 	if !v.Validate(rules) {
// 		for field, msgs := range v.GetErrors() {
// 			if len(msgs) > 0 {
// 				errors[field] = msgs[0]
// 			}
// 		}
// 	}

// 	return errors, old
// }

func ValidateCustom(c *gola.Context, rules map[string]string, customMessages map[string]string) (map[string]string, map[string]string) {
	data := make(map[string]interface{})
	old := make(map[string]string)

	for field := range rules {
		value := c.Input(field)
		data[field] = value
		old[field] = value
	}

	db := database.DB

	v := NewValidator(data, db)
	errors := make(map[string]string)

	if !v.Validate(rules, customMessages) {
		for field, msgs := range v.GetErrors() {
			if len(msgs) > 0 {
				errors[field] = msgs[0]
			}
		}
	}

	return errors, old
}

func ValidateRequestJSON(c *gola.Context, rules map[string]string) bool {
	data := make(map[string]interface{})
	for field := range rules {
		// প্রথমে check করো multi value
		if arr := c.PostFormArray(field); len(arr) > 0 {
			data[field] = arr
		} else {
			// fallback single value
			data[field] = c.Input(field)
		}
	}

	db := database.DB
	v := NewValidator(data, db)
	if !v.Validate(rules) {
		c.JSON(422, map[string]interface{}{
			"errors":  v.GetErrors(),
			"message": "Validation failed",
		})
		return false
	}
	return true
}

// func ValidateRequestJSON(c *gola.Context, rules map[string]string) bool {
// 	data := make(map[string]interface{})
// 	for field := range rules {
// 		data[field] = c.Input(field)
// 	}

// 	db := database.DB

// 	v := NewValidator(data, db)
// 	if !v.Validate(rules) {
// 		c.JSON(422, map[string]interface{}{
// 			"errors":  v.GetErrors(),
// 			"message": "Validation failed",
// 		})
// 		return false
// 	}
// 	return true
// }

// func ValidateRequest(c *gola.Context, db *gorm.DB, rules map[string]string, responseType string) bool {
// 	data := make(map[string]interface{})
// 	for field := range rules {
// 		data[field] = c.Input(field)
// 	}

// 	v := NewValidator(data, db)
// 	if !v.Validate(rules) {
// 		errors := v.GetErrors() // map[string][]string
// 		if responseType == "json" {
// 			// সব field error JSON এ পাঠানো
// 			c.JSON(422, map[string]map[string][]string{"errors": errors})
// 		} else {
// 			// HTML flash: প্রথম error set করা
// 			for field, msgs := range errors {
// 				if len(msgs) > 0 {
// 					c.SetFlash("error", fmt.Sprintf("%s: %s", field, msgs[0]))
// 					c.RedirectBack()
// 					break
// 				}
// 			}
// 		}
// 		return false
// 	}
// 	return true
// }
