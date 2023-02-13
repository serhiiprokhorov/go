package main

import (
	"sync"
	"fmt"
)

type MessagesControl struct {
	messages [1000]string
	messagesIdx int
	messagesLock sync.Mutex
	waiting MessageWaitControl
}

func (mc *MessagesControl) unsafe_updateMessageIdx() int {
	fmt.Println("unsafe_updateMessageIdx")
	mc.messagesIdx = (mc.messagesIdx+1) % len(mc.messages)
	return mc.messagesIdx
}

func (mc *MessagesControl) addMessage(m string) {
	withLocked(&mc.messagesLock, func() int {
		fmt.Println("addMessage", m)
		mc.messages[mc.messagesIdx]=m
		return mc.unsafe_updateMessageIdx()
	})

	mc.waiting.broadcastIfWaiting()
}

func (mc *MessagesControl) lastMessageIdx() int {
	return withLocked(&mc.messagesLock, func() int {
		ret := mc.messagesIdx
		 fmt.Println("lastMessageIdx", ret)
		return mc.messagesIdx
	})
}

func (mc *MessagesControl) copyMessages(behind int, actual int, op func (x string) ) {
	doCopy := func (toCopy []string) {
		for x := range toCopy {
			op(toCopy[x])
		}
	}

	if(behind < actual) {
		fmt.Println("copyMessage ", behind, actual)
		doCopy(mc.messages[behind:actual])
	} else {
		fmt.Println("copyMessage ", behind, ":")
		doCopy(mc.messages[behind:])
		fmt.Println("copyMessage ", ":", behind)
		doCopy(mc.messages[:behind])
	}
}
