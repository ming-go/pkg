package dlock

import (
	"log"
	"sync"
	"testing"
	"time"
)

const worker = 1000
const consulAddress = "172.77.0.66:8500"

func TestLock(t *testing.T) {
	var wg sync.WaitGroup

	key := "dLockTestKey1234567"

	wg.Add(worker)
	for i := 0; i < worker; i++ {
		go func(id int) {
			defer wg.Done()
			cfg := NewDefaultConfig(consulAddress, key)
			dLock, err := NewDLock(cfg)
			if err != nil {
				log.Println(err)
				return
			}

			_, err = dLock.Lock()
			if err != nil {
				log.Println(err)
				return
			}

			if !dLock.IsHeld() {
				log.Println("not hold lock")
				return
			}

			defer dLock.Unlock()

			log.Println(id, "Get lock")

			<-time.After(1000 * time.Second)
		}(i)
	}
	wg.Wait()
}
