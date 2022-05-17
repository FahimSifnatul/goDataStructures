package Queue

import (
	"errors"
	"fmt"
	"reflect"
)

// Queue a global function which creates, initializes and returns a queue instance
func Queue() *queueStruct {
	return &queueStruct{
		queue: make([]interface{}, 0),
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

// queueStruct where queue data are stored
type queueStruct struct {
	queue         []interface{}
	queueDataKind reflect.Kind
}

type queueMethods interface {
	// global methods

	// Push adds one or more elements to an existing queue.
	// returns error if data types mismatched and also doesn't push any value to the queue
	Push(elem ...interface{}) error

	// Pop removes the earliest inserted element from the caller queue
	Pop() error

	// Pops can pop multiple elements from the top of the queue
	// all you have to do is to provide an int type number i.e. how many element to pop
	// returns error if the popCount number is bigger than the size of the queue
	Pops(popCount int) error

	// RemoveAll it removes all elements from the caller queue
	// but doesn't remove the data type
	// suppose, data type of the caller queue is int
	// now caller queue calls this function then
	// it will remove all elements from the queue but
	// data type of the queue remain as int meaning
	// no data can be inserted except int for this queue
	RemoveAll()

	// Clear it removes all elements from the caller queue
	// and also removes the data type
	// suppose, data type of the caller queue is int
	// now caller queue calls this function then
	// it will remove all elements from the queue and
	// any data except invalidKind types can be inserted for this queue
	Clear()

	// Front returns the front element i.e. first inserted element from the queue
	// and error (if queue is empty)
	Front() (interface{}, error)

	// Fronts returns the earliest inserted elements equal to frontCount (stored in a slice)
	// and error (if any)
	Fronts(frontCount int) ([]interface{}, error)

	// FrontAndPop it retrieves the Front() element from the queue
	// returns the front element and also Pop() from the queue
	// also returns error (if any)
	FrontAndPop() (interface{}, error)

	// FrontsAndPops returns the earliest inserted elements equal to count (stored in a slice)
	// and also pop those elements from the queue
	// and error (if any)
	FrontsAndPops(count int) ([]interface{}, error)

	// Size returns the size of an existing queue
	Size() int

	// Empty checks whether the queue is empty or not
	// returns true if empty else false
	Empty() bool

	// Search finds the parametric element in the queue
	// if the element is found then returns the position from the Front else -1 (not found)
	// N.B. Front() is taken as position 1
	Search(elem interface{}) int

	// Display prints the stack value as slice on console screen
	// the values in slice are arranged from left to right
	// meaning that the left most data is the first inserted value
	// and the right most data is the last inserted value
	Display()

	// ToSlice returns the queue as slice
	ToSlice() []interface{}

	// private methods (for internal use only)

	// checkDataKind checks the data kind of the elements of a queue
	// when adding an element to a queue, at first the data kind is checked by this function
	// the queue data kind is of type builtin reflect.Kind
	// a queue must contain elements having same data kind
	checkDataKind(value interface{}) error
}

func (q *queueStruct) Push(elem ...interface{}) error {
	for _, e := range elem {
		if err := q.checkDataKind(e); err != nil {
			return err
		}
	}

	for _, e := range elem {
		q.queue = append(q.queue, e)
	}
	return nil
}

func (q *queueStruct) Pop() error {
	if q.Empty() {
		return errors.New("invalid operation as queue is empty")
	}

	q.queue = q.queue[1:]
	return nil
}

func (q *queueStruct) Pops(popCount int) error {
	queueSize := q.Size()
	if popCount > queueSize {
		errMsg := "invalid operation as pop count (%d) is greater than queue size(%d)"
		return fmt.Errorf(errMsg, popCount, queueSize)
	}

	q.queue = q.queue[popCount:]
	return nil
}

func (q *queueStruct) RemoveAll() {
	tempQueue := Queue()
	q.queue = tempQueue.queue
}

func (q *queueStruct) Clear() {
	tempQueue := Queue()
	q.queue = tempQueue.queue
	q.queueDataKind = tempQueue.queueDataKind
}

func (q *queueStruct) Front() (interface{}, error) {
	queueSize := q.Size()
	if queueSize == 0 {
		return nil, errors.New("invalid operation as queue is empty")
	}

	return q.queue[0], nil
}

func (q *queueStruct) Fronts(frontCount int) ([]interface{}, error) {
	queueSize := q.Size()
	if frontCount > queueSize {
		errMsg := "invalid operation as front count (%d) is greater than the queue size(%d)"
		return nil, fmt.Errorf(errMsg, frontCount, queueSize)
	}

	return q.queue[:frontCount], nil
}

func (q *queueStruct) FrontAndPop() (interface{}, error) {
	elem, err := q.Front()
	if err != nil {
		return nil, err
	}
	if err := q.Pop(); err != nil {
		return nil, err
	}
	return elem, nil
}

func (q *queueStruct) FrontsAndPops(count int) ([]interface{}, error) {
	elemSlice, err := q.Fronts(count)
	if err != nil {
		return nil, err
	}
	if err := q.Pops(count); err != nil {
		return nil, err
	}
	return elemSlice, nil
}

func (q *queueStruct) Size() int {
	return len(q.queue)
}

func (q *queueStruct) Empty() bool {
	if q.Size() == 0 {
		return true
	}
	return false
}

func (q *queueStruct) Search(elem interface{}) int {
	queueSize := q.Size()
	for i := 0; i < queueSize; i++ {
		if q.queue[i] == elem {
			return i + 1
		}
	}
	return -1
}

func (q *queueStruct) Display() {
	fmt.Println(q.queue)
}

func (q *queueStruct) ToSlice() []interface{} {
	return q.queue
}

func (q *queueStruct) checkDataKind(val interface{}) error {
	valKind := reflect.TypeOf(val).Kind()

	if q.queueDataKind != reflect.Invalid {
		if q.queueDataKind != valKind {
			return errors.New("invalid value type")
		}
		return nil
	}

	for _, kind := range invalidKind {
		if valKind == kind {
			return fmt.Errorf("%v is not supported type for queue", valKind)
		}
	}

	q.queueDataKind = valKind
	return nil
}
