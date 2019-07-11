package etcder

import (
	"fmt"
	"testing"
)

func TestGetRandomStrs(t *testing.T) {
	fmt.Println(GetRandomKv(10, 3, 5, 8, 10))
}