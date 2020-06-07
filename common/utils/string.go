package utils

import "strings"

//SplitString 拆分字符串
func SplitString(line, tag string) []string {
	fields := strings.Split(line, tag)
	for idx, field := range fields {
		fields[idx] = strings.TrimSpace(field)
	}
	return fields
}

func SplitString2(line, tag string) []string {
	ret := strings.Split(line, tag)
	fields := make([]string, 0, len(ret))
	for _, field := range ret {
		if strings.TrimSpace(field) != "" {
			fields = append(fields, field)
		}
	}
	return fields
}

func ParseFlight(s string) (letters, numbers string) {
	var l, n []rune
	for _, r := range s {
		switch {
		case r >= 'A' && r <= 'Z':
			l = append(l, r)
		case r >= 'a' && r <= 'z':
			l = append(l, r)
		case r >= '0' && r <= '9':
			n = append(n, r)
		}
	}
	return strings.ToUpper(string(l)), string(n)
}

func ResetZZContract(s string) string {
	var l, n []rune
	for _, r := range s {
		switch {
		case r >= 'A' && r <= 'Z':
			l = append(l, r)
		case r >= 'a' && r <= 'z':
			l = append(l, r)
		case r >= '0' && r <= '9':
			n = append(n, r)
		}
	}
	if n[0] == '0' {
		l = append(l, '2')
	} else {
		l = append(l, '2') //TODO:2020年开始的都是2开头
	}
	l = append(l, n...)
	return strings.ToUpper(string(l))
}

func GetVariety(s string) string {
	var l []rune
	for _, r := range s {
		switch {
		case r >= 'A' && r <= 'Z':
			l = append(l, r)
		case r >= 'a' && r <= 'z':
			l = append(l, r)
		}
	}
	return strings.ToUpper(string(l))
}
