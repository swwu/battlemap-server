package scripting

import (
	"github.com/swwu/v8.go"

	"github.com/swwu/battlemap-server/logging"
)

// utility functions for dealing with V8 objects
func StringFromV8Object(obj *v8.Object, key string, defaultVal string) string {
	if obj.HasProperty(key) {
		val := obj.GetProperty(key)
		if val.IsString() {
			return val.ToString()
		} else {
			logging.Warning.Println(
				"Tried to extract string value from non-string field:", key,
				"value is", val.ToString())
			return defaultVal
		}
	} else {
		logging.Warning.Println(
			"Tried to extract string value from empty field:", key)
		return defaultVal
	}
}

func NumberFromV8Object(obj *v8.Object, key string, defaultVal float64) float64 {
	if obj.HasProperty(key) {
		val := obj.GetProperty(key)
		if val.IsNumber() {
			return val.ToNumber()
		} else {
			logging.Warning.Println(
				"Tried to extract number value from non-number field:", key,
				"value is", val.ToString())
			return defaultVal
		}
	} else {
		logging.Warning.Println(
			"Tried to extract number value from empty field:", key)
		return defaultVal
	}
}

func BooleanFromV8Object(obj *v8.Object, key string, defaultVal bool) bool {
	if obj.HasProperty(key) {
		val := obj.GetProperty(key)
		if val.IsBoolean() {
			return val.ToBoolean()
		} else {
			logging.Warning.Println(
				"Tried to extract boolean value from non-boolean field:", key,
				"value is", val.ToString())
			return defaultVal
		}
	} else {
		logging.Warning.Println(
			"Tried to extract boolean value from empty field:", key)
		return defaultVal
	}
}

func FnFromV8Object(obj *v8.Object, key string, defaultVal *v8.Function) *v8.Function {
	if obj.HasProperty(key) {
		val := obj.GetProperty(key)
		if val.IsFunction() {
			return val.ToFunction()
		} else {
			logging.Warning.Println(
				"Tried to extract function value from non-function field:", key,
				"value is", val.ToString())
			return defaultVal
		}
	} else {
		logging.Warning.Println(
			"Tried to extract function value from empty field:", key)
		return defaultVal
	}
}

func StringArrFromV8Object(obj *v8.Object, key string, defaultVal []string) []string {
	if obj.HasProperty(key) {
		val := obj.GetProperty(key)
		if val.IsArray() {
			// define one of each type because *v8.Array doesn't implement
			// *v8.Object methods
			arrObj := val.ToObject()
			arrArr := val.ToArray()
			strings := []string{}

			for i := 0; i < arrArr.Length(); i++ {
				el := arrObj.GetElement(i)
				if !el.IsString() {
					logging.Warning.Println(
						"Tried to extract non-string value from array element:", i,
						"value is", el.ToString())
					continue
				}
				strings = append(strings, el.ToString())
			}
			return strings
		} else {
			logging.Warning.Println(
				"Tried to extract array value from non-array field:", key,
				"value is", val.ToString())
			return defaultVal
		}
	} else {
		//logging.Warning.Println(
		//	"Tried to extract array value from empty field:", key)
		return defaultVal
	}
}

func NumberFromV8Value(val *v8.Value, defaultVal float64) float64 {
	if val.IsNumber() {
		return val.ToNumber()
	} else {
		logging.Warning.Println(
			"Tried to extract number value from non-number value:", val.ToString())
		return defaultVal
	}
}

func BooleanFromV8Value(val *v8.Value, defaultVal bool) bool {
	if val.IsBoolean() {
		return val.ToBoolean()
	} else {
		logging.Warning.Println(
			"Tried to extract boolean value from non-boolean value:", val.ToString())
		return defaultVal
	}
}
