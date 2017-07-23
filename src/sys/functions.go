package sys

import "strconv"

func StringToInt32(n int64) string {
	return strconv.FormatInt(int64(n), 10)
}