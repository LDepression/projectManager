/**
 * @Author: lenovo
 * @Description:
 * @File:  errs
 * @Version: 1.0.0
 * @Date: 2023/07/17 11:40
 */

package errs

import "fmt"

type ErrorCode int

type BError struct {
	Code ErrorCode
	Msg  string
}

func (e BError) Error() string {
	return fmt.Sprintf("code:%v,msg:%s", e.Code, e.Msg)
}

func NewError(code ErrorCode, msg string) *BError {
	return &BError{Code: code, Msg: msg}
}
