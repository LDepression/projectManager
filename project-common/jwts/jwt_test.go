/**
 * @Author: lenovo
 * @Description:
 * @File:  jwt_test
 * @Version: 1.0.0
 * @Date: 2023/07/20 20:49
 */

package jwts

import "testing"

func TestParseToken(t *testing.T) {
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTA0NjEzNTYsInRva2VuIjoiMTAwNiJ9.h3FhkFQ4gAnzRGV8f7g1ySxEYpTLDc7ssyT3FFotc1M"
	ParseToken(tokenString, "msproject")
}
