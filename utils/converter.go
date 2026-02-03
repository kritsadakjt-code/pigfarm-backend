package utils

import "strconv"

func UintToString(id uint) string {
	if id == 0 {
		return ""
	}
	return strconv.FormatUint(uint64(id), 10)
}

func StringToUint(s string) uint {
	if s == "" {
		return 0
	}

	id, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0
	}
	return uint(id)
}
