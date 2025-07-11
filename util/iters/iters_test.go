package iters_test

import (
	"example.com/m/util/iters"
	"fmt"
	"github.com/ahmetb/go-linq/v3"
	"testing"
)

type Cell struct {
	CellCode int
	CellName string
	Right    bool
	Key      []int
}

//func TestIters_Count(t *testing.T) {
//	c := &Cell{}
//
//	fmt.Println(reflect.TypeOf(c))
//
//	a, _ := reflect.TypeOf(c).Elem().FieldByName("Key")
//	fmt.Printf("%+v", a.Type.Elem().Kind())
//}

func TestIters_Where(t *testing.T) {
	l := []*Cell{
		&Cell{
			CellCode: 10,
			CellName: "11",
		},
		&Cell{
			CellCode: 11,
			CellName: "12",
		},
		&Cell{
			CellCode: 10,
			CellName: "99",
		},
	}

	var result []*Cell
	iters.From(l).Where(func(o interface{}) bool {
		o2 := o.(*Cell)
		if o2.CellCode == 11 {
			return true
		} else {
			return false
		}
	}).ToSlice(&result)

	if len(result) != 1 || result[0].CellCode != 11 {
		t.Error("where fail")
	}
}

func TestIters_Distinct(t *testing.T) {
	l := []*Cell{
		&Cell{
			CellCode: 10,
			CellName: "11",
		},
		&Cell{
			CellCode: 11,
			CellName: "12",
		},
		&Cell{
			CellCode: 10,
			CellName: "99",
		},
	}

	var l2 []int
	iters.From(l).Select(func(i interface{}) interface{} {
		o := i.(*Cell)
		return o.CellCode
	}).Distinct().ToSlice(&l2)

	for _, o := range l2 {
		if o != 10 && o != 11 {
			t.Errorf("TestIters_Distinct fail")
		}
	}
}

func TestIters_Append(t *testing.T) {
	l := []*Cell{
		&Cell{
			CellCode: 10,
			CellName: "11",
		},
	}

	var l2 []*Cell
	iters.From(l).Append(&Cell{
		CellCode: 13,
		CellName: "12",
	}).ToSlice(&l2)

	for i, o := range l2 {
		if i == 0 {
			if o.CellCode != 10 || o.CellName != "11" {
				t.Errorf("TestIters_Append fail")
			}
		}
		if i == 1 {
			if o.CellCode != 13 || o.CellName != "12" {
				t.Errorf("TestIters_Append fail")
			}
		}
	}
}

func TestIters_Prepend(t *testing.T) {
	l := []*Cell{
		&Cell{
			CellCode: 10,
			CellName: "11",
		},
	}

	var l2 []*Cell
	iters.From(l).Prepend(&Cell{
		CellCode: 13,
		CellName: "12",
	}).ToSlice(&l2)

	for i, o := range l2 {
		if i == 0 {
			if o.CellCode != 13 || o.CellName != "12" {
				t.Errorf("TestIters_Prepend fail")
			}
		}
		if i == 1 {
			if o.CellCode != 10 || o.CellName != "11" {
				t.Errorf("TestIters_Prepend fail")
			}
		}
	}
}

func TestIters_All(t *testing.T) {
	l1 := []bool{true, false}

	b := iters.From(l1).All(func(i interface{}) bool {
		return i.(bool)
	})

	if b {
		t.Errorf("%+v TestIters_All is %+v", l1, b)
	}

	l2 := []bool{true, true}

	b = iters.From(l1).All(func(i interface{}) bool {
		return i.(bool)
	})

	if b {
		t.Errorf("%+v TestIters_All is %+v", l2, b)
	}
}

func TestIters_AnyWith(t *testing.T) {
	l := []*Cell{
		&Cell{
			CellCode: 10,
			CellName: "11",
		},
		&Cell{
			CellCode: 11,
			CellName: "12",
			Right:    true,
		},
		&Cell{
			CellCode: 10,
			CellName: "99",
		},
	}

	b := iters.From(l).AnyWith(func(i interface{}) bool {
		return i.(*Cell).Right
	})

	if !b {
		t.Errorf("TestIters_AnyWith fail")
	}
}

func TestIters_Concat(t *testing.T) {
	l1 := []*Cell{
		&Cell{
			CellCode: 10,
			CellName: "11",
		},
	}

	l2 := []*Cell{
		&Cell{
			CellCode: 11,
			CellName: "12",
			Right:    true,
		},
	}

	var l []*Cell
	iters.From(l1).Concat(iters.From(l2)).ToSlice(&l)

	for i, o := range l {
		if i == 0 {
			if o.CellCode != 10 || o.CellName != "11" {
				t.Errorf("TestIters_AnyWith fail")
			}
		}
		if i == 1 {
			if o.CellCode != 11 || o.CellName != "12" {
				t.Errorf("TestIters_AnyWith fail")
			}
		}
	}
}

func TestIters_Sort(t *testing.T) {
	l1 := []*Cell{
		&Cell{
			CellCode: 12,
			CellName: "11",
		},
		&Cell{
			CellCode: 11,
			CellName: "12",
			Right:    true,
		},
	}

	var l []*Cell
	iters.From(l1).Sort(func(i, j interface{}) bool {
		return i.(*Cell).CellCode < j.(*Cell).CellCode
	}).ToSlice(&l)

	for i, o := range l {
		if i == 0 {
			if o.CellCode != 11 || o.CellName != "12" {
				t.Errorf("TestIters_AnyWith fail")
			}
		}
		if i == 1 {
			if o.CellCode != 12 || o.CellName != "11" {
				t.Errorf("TestIters_AnyWith fail")
			}
		}
	}
}

func TestIters_Max(t *testing.T) {
	l1 := []int{
		11, 123, 88,
	}

	a := iters.From(l1).Max().(int)

	if a != 123 {
		t.Error("TestIters_Max fail")
	}
}

func TestIters_Except(t *testing.T) {
	l1 := []int{1, 2, 3, 4}
	l2 := []int{2, 3}
	var a []int
	iters.From(l1).Except(iters.From(l2)).ToSlice(&a)
	fmt.Println(a)
	if len(a) != 2 {
		t.Error("TestIters_Except fail")
	}
	for _, o := range a {
		if o != 1 && o != 4 {
			t.Error("TestIters_Except fail")
		}
	}
}

func TestIters_ToMapByKey(t *testing.T) {
	l1 := []*Cell{
		&Cell{
			CellCode: 12,
			CellName: "11",
		},
		&Cell{
			CellCode: 11,
			CellName: "12",
			Right:    true,
		},
	}

	m1 := make(map[int]*Cell)
	iters.From(l1).ToMapByKey(&m1, func(i interface{}) interface{} {
		return i.(*Cell).CellCode
	})
	if m1[11] != l1[1] || m1[12] != l1[0] {
		t.Errorf("TestIters_ToMapByKey fail")
	}
}

func TestIters_ToMapBy(t *testing.T) {
	l1 := []*Cell{
		&Cell{
			CellCode: 12,
			CellName: "11",
		},
		&Cell{
			CellCode: 11,
			CellName: "12",
			Right:    true,
		},
	}

	m1 := make(map[int]string)
	iters.From(l1).ToMapBy(&m1,
		func(i interface{}) interface{} {
			return i.(*Cell).CellCode
		},
		func(j interface{}) interface{} {
			return j.(*Cell).CellName
		},
	)
	if m1[11] != "12" || m1[12] != "11" {
		t.Errorf("TestIters_ToMapBy fail")
	}
}

func TestIters_SumInts(t *testing.T) {
	l1 := []*Cell{
		&Cell{
			CellCode: 12,
			CellName: "11",
		},
		&Cell{
			CellCode: 11,
			CellName: "12",
			Right:    true,
		},
	}

	totalCode := iters.From(l1).Select(func(i interface{}) interface{} {
		return i.(*Cell).CellCode
	}).SumInts()
	if totalCode != 23 {
		t.Errorf("TestIters_SumInts fail")
	}
}

func TestIters_SumInts2(t *testing.T) {
	b := 1
	a := &b
	func() {
		c := 3
		a = &c
	}()
	fmt.Println(*a)
	fmt.Println(b)
}

func TestIters_MapValueSelect(t *testing.T) {
	m := map[string]*Cell{
		"1": &Cell{
			CellCode: 10,
			CellName: "11",
		},
		"2": &Cell{
			CellCode: 11,
			CellName: "12",
		},
		"3": &Cell{
			CellCode: 10,
			CellName: "99",
		},
	}

	var result []*Cell
	iters.From(m).Select(func(o interface{}) interface{} {
		o2 := o.(linq.KeyValue)
		return o2.Value
	}).ToSlice(&result)

	if len(result) != 3 {
		t.Error("map value select fail")
	}
}
