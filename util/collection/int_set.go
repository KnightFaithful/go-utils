package collection

// 空结构体
var IntExists = struct{}{}

// IntSet is the main interface
type IntSet struct {
	// struct为结构体类型的变量
	M map[int64]struct{}
}

func NewIntSet(items ...int64) *IntSet {
	// 获取IntSet的地址
	s := &IntSet{}
	// 声明map类型的数据结构
	s.M = make(map[int64]struct{})
	s.Add(items...)
	return s
}

func (s *IntSet) Add(items ...int64) {
	for _, item := range items {
		s.M[item] = IntExists
	}
}

func (s *IntSet) Remove(items ...int64) {
	for _, item := range items {
		delete(s.M, item)
	}
}

func (s *IntSet) Contains(item int64) bool {
	_, ok := s.M[item]
	return ok
}

func (s *IntSet) Size() int {
	return len(s.M)
}

func (s *IntSet) Clear() {
	s.M = make(map[int64]struct{})
}

func (s *IntSet) Equal(other *IntSet) bool {
	// 如果两者Size不相等，就不用比较了
	if s.Size() != other.Size() {
		return false
	}
	// 迭代查询遍历
	for key := range s.M {
		// 只要有一个不存在就返回false
		if !other.Contains(key) {
			return false
		}
	}
	return true
}

func (s *IntSet) IsSubset(other *IntSet) bool {
	if s.Size() > other.Size() {
		return false
	}
	// 迭代遍历
	for key := range s.M {
		if !other.Contains(key) {
			return false
		}
	}
	return true
}

func (s *IntSet) ToSlice() []int64 {
	results := make([]int64, 0)
	for key := range s.M {
		results = append(results, key)
	}
	return results
}

func (s *IntSet) IsEmpty() bool {
	return len(s.M) == 0
}

func (s *IntSet) Copy() *IntSet {
	newIntSet := NewIntSet()
	for key := range s.M {
		newIntSet.Add(key)
	}
	return newIntSet
}

func (s *IntSet) InterSet(other *IntSet) *IntSet { // 交集
	newSet := NewIntSet()
	for key := range s.M {
		if other.Contains(key) {
			newSet.Add(key)
		}
	}
	return newSet
}

func (s *IntSet) HasIntersection(other *IntSet) bool { // 交集
	for key := range s.M {
		if other.Contains(key) {
			return true
		}
	}
	return false
}

func (s *IntSet) UnionSet(other *IntSet) *IntSet { // 并集
	newSet := other.Copy()
	newSet.Add(s.ToSlice()...)
	return newSet
}

func (s *IntSet) DiffSet(other *IntSet) *IntSet { // 差集
	newSet := s.Copy()
	newSet.Remove(other.ToSlice()...)
	return newSet
}
