package number

import (
	"math/rand"
)

func Int4() int {
	return 1000 + rand.Intn(8999)
}

// 6位int
func Int6() int {
	return 100000 + rand.Intn(899999)
}

// 10位int
func Int10() int64 {
	return 1000000000 + rand.Int63n(8999999999)
}

// 7位int
func Int7() int64 {
	return 1000000 + rand.Int63n(8999999)
}
