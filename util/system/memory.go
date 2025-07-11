package system

import (
	"fmt"
	"runtime"
)

func MemStat() {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	fmt.Printf("Alloc:%d(MB) HeapIdle:%d(MB) HeapReleased:%d(MB)\n", ms.Alloc/1024/1024, ms.HeapIdle/1024/1024, ms.HeapReleased/1024/1024)
}
