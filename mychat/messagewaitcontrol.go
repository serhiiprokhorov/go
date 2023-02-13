package main

import (
	"sync"
	"sync/atomic"
	"fmt"
)

type MessageWaitControl struct {
	messageCond sync.Cond
	waitingCC atomic.Int32
}

func (mwc *MessageWaitControl) broadcastIfWaiting() {
	w := mwc.waitingCC.Load()
	fmt.Println("broadcastIfWaiting ", w)
	if w != 0 {
		fmt.Println("do broadcast")
		mwc.messageCond.Broadcast()	
	}
}

func (mwc *MessageWaitControl) waitWhile( op func () bool ) {
	fmt.Println("waitWhile started")
	fmt.Println("waitWhile acq lock")
	mwc.messageCond.L.Lock()
	fmt.Println("waitWhile lock acq")
	fmt.Println("waitWhile inc waitingCC")
	mwc.waitingCC.Add(1)
	for op() {
		fmt.Println("waitWhile wait")
		mwc.messageCond.Wait()
		fmt.Println("waitWhile recheck")
	}
	fmt.Println("waitWhile dec waitingCC")
	mwc.waitingCC.Add(-1)
	fmt.Println("waitWhile unlock")
	mwc.messageCond.L.Unlock()
	fmt.Println("waitWhile done")
}
