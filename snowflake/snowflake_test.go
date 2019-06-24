package snowflake

import (
	"fmt"
	"testing"
)

func TestSnowFlake(t *testing.T) {
	s, _ := New(1, 1)
	fmt.Println(s.NextId())
	fmt.Println(s.NextBase62Id())
}
