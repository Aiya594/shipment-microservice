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
	mu.Unlock()
	date := time.Now().Format("21012006")
	counter++
	return fmt.Sprintf("REF-%s-%04d", date, counter)
}
