package snowflakeid

import "time"

var defaultGenerator *Generator[int64]

func init() {
	epoch := time.Date(2023, time.August, 27, 15, 4, 5, 0, time.UTC)
	var err error
	defaultGenerator, err = NewGenerator[int64](0, epoch, WithSequenceBitsLen(6))
	if err != nil {
		panic(err)
	}
}

func Next() int64 {
	return defaultGenerator.Next()
}

func SetDefault(g *Generator[int64]) {
	defaultGenerator = g
}
