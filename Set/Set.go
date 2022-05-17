package Set

import (
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"time"
)

// Set a global function which creates, initializes and returns a set instance
func Set() *setStruct {
	return &setStruct{
		set: make(map[interface{}]bool),
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

// setMethods stores interface declaration of all setStruct methods
type setMethods interface {
	// global methods

	// Add adds one or more elements to an existing set
	// returns error if data types mismatched and also doesn't push any value to the set
	Add(elem ...interface{}) error

	// Remove removes one or more elements from an existing set
	Remove(elems ...interface{})

	// RemoveAll it removes all elements from the caller set
	// but doesn't remove the data type
	// suppose, data type of the caller set is int
	// now caller set calls this function then
	// it will remove all elements from the set but
	// data type of the set remain as int meaning
	// no data can be inserted except int for this set
	RemoveAll()

	// Clear it removes all elements from the caller set
	// and also removes the data type
	// suppose, data type of the caller set is int
	// now caller set calls this function then
	// it will remove all elements from the set and
	// any data except invalidKind types can be inserted for this set
	Clear()

	// Copy copies the existing set to a new set and returns the new set
	Copy() *setStruct

	// Len returns the length of the existing set
	Len() int

	// Union performs the set union operation among the existing set and sets passed as params,
	// stores data in a new set and returns the new set
	Union(sets ...*setStruct) (*setStruct, error)

	// Intersection performs the set intersection operation among the existing set and sets passed as params,
	// stores data in a new set and returns the new set
	Intersection(sets ...*setStruct) (*setStruct, error)

	// Difference performs the set difference operation from the existing set and sets passed as params,
	// stores data in a new set and returns the new set
	// the set difference is found as follows
	// the set calling this method - parametric set1 - parametric set2 - parametric set3 -...
	Difference(sets ...*setStruct) (*setStruct, error)

	// MakeDisjoint makes the caller set and parametric set disjoint to each other.
	// suppose, the call is like x.MakeDisjoint(y)
	// then this function makes the sets x and y disjoint to each other
	MakeDisjoint(set *setStruct) error

	// MakeSubSet creates and returns a sub set of the caller set having randomized elements equal to passed parameter
	// suppose, the call is like x.MakeSubSet(y)
	// then the function creates a sub set of x having randomized elements equal to y and returns the sub set
	// y = 0 is valid value as it will return empty set
	// y < -1 or y > number of elements present in x is invalid choice
	MakeSubSet(elemNum int) (*setStruct, error)

	// Has checks whether the existing set has a specific element or not
	Has(elem interface{}) bool

	// IsDisjoint checks whether two sets are disjoint to each other or not
	// suppose, the call is like x.IsDisjoint(y)
	// then this function determines whether x and y are disjoint to each other or not
	// and returns boolean value (true, false) and error (if any)
	IsDisjoint(sets *setStruct) (bool, error)

	// IsSubSet checks whether the caller set is a sub set of the parametric set
	// suppose, the call is like x.IsSubSet(y)
	// then the functions checks whether x is a sub set of y or not
	// and returns boolean value (true, false) and error (if any)
	IsSubSet(set *setStruct) (bool, error)

	// IsSuperSet checks whether the caller set is a super set of the parametric set
	// suppose, the call is like x.IsSubSet(y)
	// then the functions checks whether x is the super set of y or not
	// and returns boolean value (true, false) and error (if any)
	IsSuperSet(set *setStruct) (bool, error)

	// ToSlice converts set to golang slice and return the slice
	ToSlice() []interface{}

	// Display converts set to a golang slice and
	// prints the converted set (slice) on console screen
	Display()

	// private methods (for internal use only)

	// checkDataKind checks the data kind of the elements of a set
	// when adding an element to a set, at first the data kind is checked by this function
	// the set data kind is of type builtin reflect.Kind
	// a set must contain elements having same data kind
	checkDataKind(value interface{}) error
}

// setStruct where set data are stored
type setStruct struct {
	set         map[interface{}]bool
	setDataKind reflect.Kind
}

func (s *setStruct) Add(elem ...interface{}) error {
	for _, e := range elem {
		if err := s.checkDataKind(e); err != nil {
			return err
		}
	}

	for _, e := range elem {
		s.set[e] = true
	}
	return nil
}

func (s *setStruct) Remove(elem ...interface{}) {
	for _, e := range elem {
		delete(s.set, e)
	}
}

func (s *setStruct) RemoveAll() {
	tempSet := Set()
	s.set = tempSet.set
}

func (s *setStruct) Clear() {
	tempSet := Set()
	s.set = tempSet.set
	s.setDataKind = tempSet.setDataKind
}

func (s *setStruct) Copy() *setStruct {
	return &setStruct{
		set:         s.set,
		setDataKind: s.setDataKind,
	}
}

func (s *setStruct) Len() int {
	return len(s.set)
}

func (s *setStruct) Union(sets ...*setStruct) (*setStruct, error) {
	unionSet := s.Copy()

	for _, set := range sets {
		if unionSet.setDataKind == reflect.Invalid && set.setDataKind != reflect.Invalid {
			unionSet.setDataKind = set.setDataKind
		}

		if set.setDataKind != reflect.Invalid {
			if unionSet.setDataKind != set.setDataKind {
				return nil, errors.New("mismatched data types among sets")
			}
			for key := range set.set {
				unionSet.set[key] = true
			}
		}
	}

	return unionSet, nil
}

func (s *setStruct) Intersection(sets ...*setStruct) (*setStruct, error) {
	intersectionSet := Set()
	totalSetCount := len(sets) + 1 // +1 for s
	elemFreqCount := make(map[interface{}]int)

	if s.setDataKind != reflect.Invalid {
		intersectionSet.setDataKind = s.setDataKind
		for key := range s.set {
			elemFreqCount[key] += 1
		}
	}

	for _, set := range sets {
		if intersectionSet.setDataKind == reflect.Invalid && set.setDataKind != reflect.Invalid {
			intersectionSet.setDataKind = set.setDataKind
		}

		if set.setDataKind != reflect.Invalid {
			if intersectionSet.setDataKind != set.setDataKind {
				return nil, errors.New("mismatched data types among sets")
			}
			for key := range set.set {
				elemFreqCount[key] += 1
			}
		}
	}

	for elem, freq := range elemFreqCount {
		if freq == totalSetCount {
			intersectionSet.set[elem] = true
		}
	}

	return intersectionSet, nil
}

func (s *setStruct) Difference(sets ...*setStruct) (*setStruct, error) {
	diffSet := s.Copy()
	unionSet, err := Set().Union(sets...)
	if err != nil {
		return nil, err
	}

	if diffSet.setDataKind != reflect.Invalid && unionSet.setDataKind != reflect.Invalid && diffSet.setDataKind != unionSet.setDataKind {
		return nil, errors.New("mismatched data types among sets")
	}

	for elem := range unionSet.set {
		if diffSet.set[elem] {
			diffSet.Remove(elem)
		}
	}

	return diffSet, nil
}

func (s *setStruct) MakeDisjoint(set *setStruct) error {
	if s.setDataKind != reflect.Invalid && set.setDataKind != reflect.Invalid && s.setDataKind != set.setDataKind {
		return errors.New("mismatched data types among sets")
	}

	for elem := range set.set {
		if s.set[elem] {
			s.Remove(elem)
			set.Remove(elem)
		}
	}

	return nil
}

func (s *setStruct) MakeSubSet(elemNum int) (*setStruct, error) {
	setSlice := s.ToSlice()
	setSliceLen := len(setSlice)

	subSet := Set()
	if elemNum < 0 || elemNum > setSliceLen {
		return subSet, errors.New("invalid element number provided to make sub set")
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(setSliceLen, func(i, j int) { setSlice[i], setSlice[j] = setSlice[j], setSlice[i] })

	subSet.setDataKind = s.setDataKind
	for _, elem := range setSlice[:elemNum] {
		subSet.set[elem] = true
	}
	return subSet, nil
}

func (s *setStruct) Has(elem interface{}) bool {
	if _, has := s.set[elem]; !has {
		return false
	}
	return true
}

func (s *setStruct) IsDisjoint(set *setStruct) (bool, error) {
	disjointSet, err := s.Intersection(set)
	if disjointSet.Len() == 0 || err != nil {
		return false, err
	}
	return true, nil
}

func (s *setStruct) IsSubSet(set *setStruct) (bool, error) {
	if s.setDataKind != reflect.Invalid && set.setDataKind != reflect.Invalid && s.setDataKind != set.setDataKind {
		return false, errors.New("mismatched data types among sets")
	}

	for elem := range s.set {
		if !set.set[elem] {
			return false, nil
		}
	}
	return true, nil
}

func (s *setStruct) IsSuperSet(set *setStruct) (bool, error) {
	if s.setDataKind != reflect.Invalid && set.setDataKind != reflect.Invalid && s.setDataKind != set.setDataKind {
		return false, errors.New("mismatched data types among sets")
	}

	for elem := range set.set {
		if !s.set[elem] {
			return false, nil
		}
	}
	return true, nil
}

func (s *setStruct) ToSlice() []interface{} {
	setSlice := make([]interface{}, 0)
	for elem := range s.set {
		setSlice = append(setSlice, elem)
	}
	return setSlice
}

func (s *setStruct) Display() {
	setSlice := s.ToSlice()
	fmt.Println(setSlice)
}

func (s *setStruct) checkDataKind(val interface{}) error {
	valKind := reflect.TypeOf(val).Kind()

	if s.setDataKind != reflect.Invalid && s.setDataKind != valKind {
		return errors.New("invalid value type")
	}

	for _, kind := range invalidKind {
		if valKind == kind {
			return fmt.Errorf("%v is not supported type for set", valKind)
		}
	}

	s.setDataKind = valKind
	return nil
}
