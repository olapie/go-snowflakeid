package snowflakeid

import (
	"testing"
	"time"
)

func TestGenerator_Next(t *testing.T) {
	epoch := time.Now().AddDate(-1, 0, 0)
	g1, err := NewGenerator[int64](0, epoch)
	if err != nil {
		t.Fatalf("NewGenerator: %v", err)
	}
	g2, err := NewGenerator[int64](0, epoch, WithTimestampUnit(time.Microsecond))
	if err != nil {
		t.Fatalf("NewGenerator: %v", err)
	}
	g3, err := NewGenerator[int64](0, epoch, WithTimestampUnit(time.Second))
	if err != nil {
		t.Fatalf("NewGenerator: %v", err)
	}

	ids := make(map[int64]bool)
	// default seq bits len
	for i := 0; i < 256; i++ {
		id1 := g1.Next()
		id2 := g2.Next()
		id3 := g3.Next()
		if ids[id1] {
			t.Fatalf("Duplicate id1 %d", id1)
		}
		ids[id1] = true

		if ids[id2] {
			t.Fatalf("Duplicate id2 %d", id1)
		}
		ids[id2] = true

		if ids[id3] {
			t.Fatalf("Duplicate id3 %d", id1)
		}
		ids[id3] = true

		if id1 < id3 {
			t.Fatalf("id1 %d must be greater than id3 %d", id1, id3)
		}
		if id2 < id1 {
			t.Fatalf("id1 %d must be less than id2 %d", id1, id2)
		}
		t.Logf("id1 %d id2 %d id3 %d\n", id1, id2, id3)
	}
}
