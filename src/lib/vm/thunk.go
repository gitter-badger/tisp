package vm

import (
	"../types"
	"sync"
	"sync/atomic"
)

type State uint32

const (
	ILLEGAL State = iota
	VALUE
	LOCKED
	APP
)

type Thunk struct {
	Result    types.Object
	function  *Thunk
	args      *Thunk
	state     State
	blackHole sync.WaitGroup
}

func NewValueThunk(v types.Object) *Thunk {
	return &Thunk{Result: v, state: VALUE}
}

func NewAppThunk(f *Thunk, args *Thunk) *Thunk {
	t := &Thunk{function: f, args: args, state: APP}
	t.blackHole.Add(1)
	return t
}

func (t *Thunk) Eval() { // into WHNF
	if !t.compareAndSwapState(APP, LOCKED) {
		// Some goroutine is evaluating this thunk or it has been evaluated already.
		return
	}

	go t.function.Eval()
	go t.args.Eval()

	t.function.Wait()
	t.args.Wait()

	f, ok := t.function.Result.(types.Callable)

	if !ok {
		panic("Something not callable was called.")
	}

	args, ok := t.args.Result.(types.Dictionary)

	if !ok {
		panic("Something which is not a dictionary was used as arguments.")
	}

	t.Result = f.Call(args)
	t.function = nil
	t.args = nil

	t.storeState(VALUE)

	t.blackHole.Done()
}

func (t *Thunk) Wait() {
	t.blackHole.Wait()
}

func (t *Thunk) compareAndSwapState(old, new State) bool {
	return atomic.CompareAndSwapUint32(
		(*uint32)(&t.state),
		uint32(old),
		uint32(new))
}

func (t *Thunk) storeState(new State) {
	atomic.StoreUint32((*uint32)(&t.state), uint32(new))
}