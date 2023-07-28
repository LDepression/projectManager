/**
 * @Author: lenovo
 * @Description:
 * @File:  tms
 * @Version: 1.0.0
 * @Date: 2023/07/26 17:00
 */

package tms

import "time"

func Format(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}
func FormatYMD(t time.Time) string {
	return t.Format("2006-01-02")
}
func FormatByMill(t int64) string {
	return time.UnixMilli(t).Format("2006-01-02 15:04:05")
}
