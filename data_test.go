package bank

import (
	"testing"
)

func Benchmark_loadData(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		loadData()
	}

}
