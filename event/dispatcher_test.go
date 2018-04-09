package event

import "testing"

func testfun0(e Event) { print("hello-0\n") }
func testfun1(e Event) { print("hello-1\n") }
func testfun2(e Event) { print("hello-2\n") }
func TestDispatch(t *testing.T) {

	listener1 := Listener{Callable: testfun0, Priority: 1}
	listener2 := Listener{Callable: testfun1, Priority: 0}
	listener3 := Listener{Callable: testfun2, Priority: 2}

	SharedDispatcher().AddListener(EventTxEnterPool, listener1)
	SharedDispatcher().AddListener(EventTxEnterPool, listener2)
	SharedDispatcher().AddListenerExecOnce(EventTxEnterPool, listener3)
	event := NewParamsEvent(EventTxEnterPool)
	print("dispatch event test 1\n")
	SharedDispatcher().Dispatch(event)
	print("dispatch event test 2\n")
	SharedDispatcher().Dispatch(event)

	SharedDispatcher().RemoveAll(EventTxEnterPool)

	print("dispatch event test 3\n")
	SharedDispatcher().Dispatch(event)
}
