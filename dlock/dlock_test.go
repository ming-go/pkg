package dlock

import (
	"fmt"
	"testing"
	"time"
)

const consulAddress = "172.77.0.22:8500"

func TestLock(t *testing.T) {
	dLock, err := NewDLock(consulAddress, "testtLockorg123", 10000)
	if err != nil {
		t.Error(err)
	}

	_, err = dLock.Lock()
	if err != nil {
		t.Error(err)
	}

	go func() {
		for {
			fmt.Println("lock status", dLock.IsHeld())
			<-time.After(1 * time.Second)
		}
	}()
	<-time.After(5 * time.Second)

	fmt.Println(dLock.Unlock())

	<-time.After(30 * time.Second)

	//if locked {
	//	for {
	//		if dLock.IsHeld() {
	//			fmt.Println(true)
	//		} else {
	//			fmt.Println(false)
	//			break
	//		}

	//		<-time.After(time.Second * 1)
	//	}
	//}

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
