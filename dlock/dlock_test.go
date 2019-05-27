package dlock

import (
	"log"
	"sync"
	"testing"
	"time"
)

const consulAddress = "172.77.0.66:8500"

func TestLock(t *testing.T) {
	var wg sync.WaitGroup

	key := "dLockTestKey"

	wg.Add(1000)
	for i := 0; i < 1000; i++ {
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

			<-time.After(15 * time.Second)
		}(i)
	}
	wg.Wait()
}

//func TestLock2(t *testing.T) {
//	dLock, err := NewDLock(consulAddress, "testtLockorg", 15000)
//	if err != nil {
//		t.Error(err)
//	}
//
//	locked, err := dLock.Lock()
//	if err != nil {
//		t.Error(err)
//	}
//
//	if locked {
//		for {
//			if dLock.IsHeld() {
//				fmt.Println(true)
//			} else {
//				fmt.Println(false)
//				break
//			}
//
//			<-time.After(time.Second * 1)
//		}
//	}
//}
