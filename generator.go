package snowflakeid

import (
	"errors"
	"fmt"
	"sync/atomic"
	"time"
)

type GeneratorOptions struct {
	nSeqBits      int
	nShardBits    int
	timestampUnit time.Duration
}

type GeneratorOption func(options *GeneratorOptions)

func WithSequenceBitsLen(n int) GeneratorOption {
	return func(options *GeneratorOptions) {
		options.nSeqBits = n
	}
}

func WithShardBitsLen(n int) GeneratorOption {
	return func(options *GeneratorOptions) {
		options.nShardBits = n
	}
}

func WithTimestampUnit(t time.Duration) GeneratorOption {
	return func(options *GeneratorOptions) {
		options.timestampUnit = t
	}
}

type Generator[T ~int64] struct {
	options    GeneratorOptions
	seqCounter int64
	shard      int64
	epoch      time.Time
	maxSeq     int64
}

func NewGenerator[T ~int64](shard int, epoch time.Time, opts ...GeneratorOption) (*Generator[T], error) {
	options := &GeneratorOptions{
		nSeqBits:      8, // 256 sequences each millisecond by default
		nShardBits:    6, // 64 machines by default
		timestampUnit: time.Millisecond,
	}

	for _, opt := range opts {
		opt(options)
	}

	if options.nSeqBits < 1 || options.nSeqBits > 16 {
		return nil, errors.New("invalid options: SequenceBitsLen must be in range [1,16]")
	}

	if options.timestampUnit < time.Microsecond || options.timestampUnit > time.Second {
		return nil, fmt.Errorf("timestamp unit must between [%v, %v]", time.Microsecond, time.Second)
	}

	if options.nShardBits > 8 {
		return nil, errors.New("invalid options: ShardBitsLen must be in range [0,8]")
	}

	if shard < 0 {
		return nil, errors.New("shard cannot be negative")
	}

	if options.nShardBits+options.nSeqBits >= 20 {
		return nil, fmt.Errorf("invalid options: ShardBitsLen %d + SequenceBitsLen %d must be less than 20", options.nShardBits, options.nSeqBits)
	}

	return &Generator[T]{
		options: *options,
		epoch:   epoch,
		shard:   int64(shard) % (1 << options.nShardBits),
		maxSeq:  1 << options.nSeqBits,
	}, nil
}

func (g *Generator[T]) Next() T {
	elapsed := time.Since(g.epoch) / g.options.timestampUnit
	id := int64(elapsed) << (g.options.nSeqBits + g.options.nShardBits)
	if g.options.nShardBits > 0 {
		id |= g.shard << g.options.nSeqBits
	}
	seq := atomic.AddInt64(&g.seqCounter, 1)
	id |= seq % g.maxSeq
	return T(id)
}
