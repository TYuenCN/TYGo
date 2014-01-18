package YSMD5Util_test

import (
	//"fmt"
	"strings"
	"testing"
	"yuensoft.com/utils/crypto/YSMD5Util"
)

func TestMD5FromString(t *testing.T) {
	if ok := strings.EqualFold("8b1a9953c4611296a827abf8c47804d7", YSMD5Util.MD5FromString("Hello")); ok {
		t.Errorf("sting for \"Hello\" MD5 is : %s\nFunction return is : %s", "8b1a9953c4611296a827abf8c47804d7", YSMD5Util.MD5FromString("Hello"))
	}
}
