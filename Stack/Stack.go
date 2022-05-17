package Stack

import (
	"errors"
	"fmt"
	"reflect"
)

// Stack a global function which creates, initializes and returns a stack instance
func Stack() *stackStruct {
	return &stackStruct{
		stack: make([]interface{}, 0),
	}
}

// not supported data kinds are stored here
var (
	invalidKind = []reflect.Kind{
		reflect.Array,
		reflect.Chan,
		reflect.Func,
		reflect.Interface,
		reflect.Map,
		reflect.Ptr,
		reflect.Slice,
		reflect.Struct,
		reflect.UnsafePointer,
	}
)

// stackStruct where stack data are stored
type stackStruct struct {
	stack         []interface{}
	stackDataKind reflect.Kind
}

type stackMethods interface {
	// global methods

	// Push adds one or more elements to an existing stack.
	// returns error if data types mismatched and also doesn't push any value to the stack
	Push(elem ...interface{}) error

	// Pop removes the top element i.e. last inserted element from the stack
	Pop() error

	// Pops can pop multiple elements from the top of the stack
	// all you have to do is to provide an int type number i.e. how many element to pop
	// returns error if the popCount number is bigger than the size of the stack
	Pops(popCount int) error

	// RemoveAll it removes all elements from the caller stack
	// but doesn't remove the data type
	// suppose, data type of the caller stack is int
	// now caller stack calls this function then
	// it will remove all elements from the stack but
	// data type of the stack remain as int meaning
	// no data can be inserted except int for this stack
	RemoveAll()

	// Clear it removes all elements from the caller stack
	// and also removes the data type
	// suppose, data type of the caller stack is int
	// now caller stack calls this function then
	// it will remove all elements from the stack and
	// any data except invalidKind types can be inserted for this stack
	Clear()

	// Top returns the top element i.e. last inserted element from the stack
	// and error (if stack is empty)
	Top() (interface{}, error)

	// Tops returns top elements i.e. the latest elements equal to topCount (stored in a slice)
	// and error (if stack is empty)
	Tops(topCount int) ([]interface{}, error)

	// TopAndPop it retrieves the Top() element from the stack
	// returns the top element and also Pop() from the stack
	// also returns error (if any)
	TopAndPop() (interface{}, error)

	// TopsAndPops returns top elements i.e. the latest elements equal to count (stored in a slice)
	// and also pop those elements from the stack
	// and error (if any)
	TopsAndPops(count int) ([]interface{}, error)

	// Size returns the size of an existing stack
	Size() int

	// Empty checks whether the stack is empty or not
	// returns true if empty else false
	Empty() bool

	// Search finds the parametric element in the stack
	// if the element is found then returns the position from the Top() else -1 (not found)
	// N.B. Top() is taken as position 1
	Search(elem interface{}) int

	// Display prints the stack value as slice on console screen
	// the values in slice are arranged from left to right
	// meaning that the left most data is the first inserted value
	// and the right most data is the last inserted value
	Display()

	// ToSlice returns the stack as slice
	ToSlice() []interface{}

	// private methods (for internal use only)

	// checkDataKind checks the data kind of the elements of a stack
	// when adding an element to a stack, at first the data kind is checked by this function
	// the stack data kind is of type builtin reflect.Kind
	// a stack must contain elements having same data kind
	checkDataKind(value interface{}) error
}

func (st *stackStruct) Push(elem ...interface{}) error {
	for _, e := range elem {
		if err := st.checkDataKind(e); err != nil {
			return err
		}
	}

	for _, e := range elem {
		st.stack = append(st.stack, e)
	}
	return nil
}

func (st *stackStruct) Pop() error {
	stackSize := st.Size()
	if stackSize == 0 {
		return errors.New("invalid operation as stack is empty")
	}

	st.stack = st.stack[:stackSize-1]
	return nil
}

func (st *stackStruct) Pops(popCount int) error {
	stackSize := st.Size()
	if popCount > stackSize {
		errMsg := "invalid operation as pop count (%d) is greater than the stack size(%d)"
		return fmt.Errorf(errMsg, popCount, stackSize)
	}

	st.stack = st.stack[:stackSize-popCount]
	return nil
}

func (st *stackStruct) RemoveAll() {
	tempStack := Stack()
	st.stack = tempStack.stack
}

func (st *stackStruct) Clear() {
	tempStack := Stack()
	st.stack = tempStack.stack
	st.stackDataKind = tempStack.stackDataKind
}

func (st *stackStruct) Top() (interface{}, error) {
	stackSize := st.Size()
	if stackSize == 0 {
		return nil, errors.New("invalid operation as stack is empty")
	}

	return st.stack[stackSize-1], nil
}

func (st *stackStruct) Tops(topCount int) ([]interface{}, error) {
	stackSize := st.Size()
	if topCount > stackSize {
		errMsg := "invalid operation as top count (%d) is greater than the stack size(%d)"
		return nil, fmt.Errorf(errMsg, topCount, stackSize)
	}

	return st.stack[stackSize-topCount : stackSize], nil
}

func (st *stackStruct) TopAndPop() (interface{}, error) {
	elem, err := st.Top()
	if err != nil {
		return nil, err
	}
	if err := st.Pop(); err != nil {
		return nil, err
	}
	return elem, nil
}

func (st *stackStruct) TopsAndPops(count int) ([]interface{}, error) {
	elemSlice, err := st.Tops(count)
	if err != nil {
		return nil, err
	}
	if err := st.Pops(count); err != nil {
		return nil, err
	}
	return elemSlice, nil
}

func (st *stackStruct) Size() int {
	return len(st.stack)
}

func (st *stackStruct) Empty() bool {
	if st.Size() == 0 {
		return true
	}
	return false
}

func (st *stackStruct) Search(elem interface{}) int {
	stackSize := st.Size()
	for i := stackSize - 1; i >= 0; i-- {
		if st.stack[i] == elem {
			return stackSize - i
		}
	}
	return -1
}

func (st *stackStruct) Display() {
	fmt.Println(st.stack)
}

func (st *stackStruct) ToSlice() []interface{} {
	return st.stack
}

func (st *stackStruct) checkDataKind(val interface{}) error {
	valKind := reflect.TypeOf(val).Kind()

	if st.stackDataKind != reflect.Invalid && st.stackDataKind != valKind {
		return errors.New("invalid value type")
	}

	for _, kind := range invalidKind {
		if valKind == kind {
			return fmt.Errorf("%v is not supported type for stack", valKind)
		}
	}

	st.stackDataKind = valKind
	return nil
}
