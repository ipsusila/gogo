package data

import (
	"fmt"
	"log"
	"reflect"
)

func flatMap(dst map[string]interface{}, arg reflect.Value, baseKey, sep string) {
	if arg.Kind() != reflect.Map {
		if arg.CanInterface() {
			dst[baseKey] = arg.Interface()
		} else {
			//Debug
			log.Printf("flatMap: arg !CanInterface(). ignored: %v", arg)
		}
		return
	}

	//get all map keys
	keys := arg.MapKeys()
	for _, key := range keys {
		//If key is not convertible to interface,
		//ignore key:value entry.
		if !key.CanInterface() {
			log.Printf("flatMap: map key !CanInterface(). ignored: %v", key)
			continue
		}
		str := fmt.Sprintf("%v", key.Interface())
		if len(baseKey) != 0 {
			str = baseKey + sep + str
		}

		//Get map value. If map value is an interface,
		//retrieve the underlying interface value.
		val := arg.MapIndex(key)
		if val.Kind() == reflect.Interface {
			val = reflect.ValueOf(val.Interface())
		}
		flatMap(dst, val, str, sep)
	}
}

//FlatMap converts map to flat map.
func FlatMap(src interface{}, sep string) map[string]interface{} {
	dest := make(map[string]interface{})

	//User reflection to retrieve underlying value.
	vmap := reflect.Indirect(reflect.ValueOf(src))
	if vmap.Kind() == reflect.Map {
		flatMap(dest, vmap, "", sep)
	}

	return dest
}

//IsFlatMap returns true if given argument is a flat map,
//i.e. a map in which it's value is not a map.
func IsFlatMap(m interface{}) bool {
	vmap := reflect.Indirect(reflect.ValueOf(m))
	if vmap.Kind() != reflect.Map {
		return false
	}

	//Check map value.
	keys := vmap.MapKeys()
	for _, key := range keys {
		val := vmap.MapIndex(key)
		if val.Kind() == reflect.Interface {
			val = reflect.ValueOf(val.Interface())
		}

		//If value type is a map, return false
		if val.Kind() == reflect.Map {
			return false
		}
	}

	return true
}
