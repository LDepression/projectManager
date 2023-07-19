/**
 * @Author: lenovo
 * @Description:
 * @File:  encrypts
 * @Version: 1.0.0
 * @Date: 2023/07/19 20:30
 */

package encrypts

import (
	"crypto/md5"
	"encoding/hex"
	"io"
)

func Md5(str string) string {
	hash := md5.New()
	_, _ = io.WriteString(hash, str)
	return hex.EncodeToString(hash.Sum(nil))
}
