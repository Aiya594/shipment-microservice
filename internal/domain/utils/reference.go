package utils

import (
	"fmt"
	"sync"
	"time"
)

var (
	counter int
	mu      sync.Mutex
)

func GenerateReferenceNum() string {
	//example:REF-21032026-0001
	mu.Lock()
	defer mu.Unlock()
	counter++
	ref := fmt.Sprintf("REF-%s-%04d", time.Now().Format("02012006"), counter)
	return ref
}
