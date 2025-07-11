package collection

// 空结构体
var exists = struct{}{}

// Set is the main interface
type Set struct {
	// struct为结构体类型的变量
	m map[interface{}]struct{}
}

func NewSet(items ...interface{}) *Set {
	s := &Set{}
	// 声明map类型的数据结构
	s.m = make(map[interface{}]struct{})
	s.Add(items...)
	return s
}

func (s *Set) Add(items ...interface{}) {
	for _, item := range items {
		s.m[item] = exists
	}
}

// 添加已存在的元素时，返回false
func (s *Set) AddOne(item interface{}) bool {
	if _, ok := s.m[item]; ok {
		return false
	}
	s.m[item] = exists
	return true
}

func (s *Set) Remove(items ...interface{}) {
	for _, item := range items {
		delete(s.m, item)
	}
}

func (s *Set) Contains(item interface{}) bool {
	_, ok := s.m[item]
	return ok
}

func (s *Set) Size() int {
	return len(s.m)
}

func (s *Set) Clear() {
	s.m = make(map[interface{}]struct{})
}

func (s *Set) Equal(other *Set) bool {
	// 如果两者Size不相等，就不用比较了
	if s.Size() != other.Size() {
		return false
	}
	// 迭代查询遍历
	for key := range s.m {
		// 只要有一个不存在就返回false
		if !other.Contains(key) {
			return false
		}
	}
	return true
}

func (s *Set) IsSubset(other *Set) bool {
	if s.Size() > other.Size() {
		return false
	}
	// 迭代遍历
	for key := range s.m {
		if !other.Contains(key) {
			return false
		}
	}
	return true
}

func (s *Set) ToSlice() []interface{} {
	results := make([]interface{}, 0)
	for key := range s.m {
		results = append(results, key)
	}
	return results
}

func (s *Set) ToStringSlice() []string {
	results := make([]string, 0)
	for key := range s.m {
		results = append(results, key.(string))
	}
	return results
}

func (s *Set) IsEmpty() bool {
	return len(s.m) == 0
}

func (s *Set) Copy() *Set {
	newSet := NewSet()
	for key := range s.m {
		newSet.Add(key)
	}
	return newSet
}

func (s *Set) InterSet(other *Set) *Set { // 交集
	newSet := NewSet()
	for key := range s.m {
		if other.Contains(key) {
			newSet.Add(key)
		}
	}
	return newSet
}

func (s *Set) UnionSet(other *Set) *Set { // 并集
	newSet := other.Copy()
	newSet.Add(s.ToSlice()...)
	return newSet
}

func (s *Set) ToIntSlice() []int {
	results := make([]int, 0)
	for key := range s.m {
		results = append(results, int(key.(int64)))
	}
	return results
}

func (s *Set) ToInt64Slice() []int64 {
	results := make([]int64, 0)
	for key := range s.m {
		results = append(results, key.(int64))
	}
	return results
}
