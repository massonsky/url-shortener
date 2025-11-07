package utils

const (
	Alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Base     = 62
)

func Encode(num int64) string {

	if num == 0 {
		return string(Alphabet[0])
	}

	var buf []byte
	for num > 0 {
		rem := num % Base
		buf = append(buf, Alphabet[rem])
		num = num / Base
	}

	// reverse
	for i, j := 0, len(buf)-1; i < j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}
	return string(buf)
}
