package reflect

import (
	"fmt"
	"reflect"
	"strings"
)

func TransMap2ExplicitType(m map[string]interface{}, s interface{}) error {
	//sType := reflect.TypeOf(s)
	sValue := reflect.ValueOf(&s)
	// switch sType.Kind() {
	// case reflect.Struct:
	// 	for k, v := range m {
	// 		nValue := sValue.FieldByNameFunc(func(fnm string) bool {
	// 			return strings.EqualFold(fnm, k)
	// 		})
	// 		fmt.Printf("IsValid() : %t\n", nValue.IsValid())
	// 		fmt.Printf("CanSet() : %t\n", nValue.CanSet())
	// 		if nValue.IsValid() && nValue.CanSet() {
	// 			nValue.Set(reflect.ValueOf(v))
	// 		}
	// 	}
	// }

	for k, v := range m {
		nValue := sValue.FieldByNameFunc(func(fnm string) bool {
			b := strings.EqualFold(fnm, k)
			fmt.Printf("%v  %t\n", fnm, b)
			return b
		})
		nValue.Set(reflect.ValueOf(&v))
		// fmt.Printf("nValue  IsValid() : %t\n", nValue.IsValid())
		// fmt.Printf("nValue  CanSet() : %t\n", nValue.CanSet())
		// if nValue.IsValid() && nValue.CanSet() {
		// 	nValue.Set(reflect.ValueOf(v))
		// }
	}

	return nil
}
