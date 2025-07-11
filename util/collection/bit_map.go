package collection

import "math"

type Bitmap struct {
	words  []uint64
	length int
}

func NewBitmap(capacities ...int) *Bitmap {
	capacity := 0
	if len(capacities) > 0 && capacities[0] > 0 {
		capacity = capacities[0]/64 + 1
	}
	return &Bitmap{
		length: 0,
		words:  make([]uint64, capacity),
	}
}
func (m *Bitmap) Has(num int) bool {
	word, bit := m.getWordAndBit(num)
	return word < len(m.words) && (m.words[word]&(1<<bit)) != 0
}

func (m *Bitmap) Add(num int) {
	word, bit := m.getWordAndBit(num)
	for word >= len(m.words) {
		m.words = append(m.words, 0)
	}
	// 判断num是否已经存在bitmap中
	if m.words[word]&(1<<bit) == 0 {
		m.words[word] |= 1 << bit
		m.length++
	}
}

func (m *Bitmap) Copy() *Bitmap {
	res := NewBitmap()
	for _, word := range m.words {
		res.words = append(res.words, word)
	}
	return res
}

func (m *Bitmap) AddAndReturn(num int) *Bitmap {
	res := m.Copy()
	res.Add(num)
	return res
}

func (m *Bitmap) Len() int {
	return m.length
}

func (m *Bitmap) Capacity() int {
	return len(m.words) * 64
}

func (m *Bitmap) WordLength() int {
	return len(m.words)
}

func (m *Bitmap) getWordAndBit(num int) (int, uint) {
	return num / 64, uint(num % 64)
}

func (m *Bitmap) Union(other *Bitmap) *Bitmap {
	res := NewBitmap()
	for i := 0; i < m.WordLength() && i < other.WordLength(); i++ {
		res.words = append(res.words, m.words[i]|other.words[i])
	}
	for i := res.WordLength(); i < m.WordLength(); i++ {
		res.words = append(res.words, m.words[i])
	}
	for i := res.WordLength(); i < other.WordLength(); i++ {
		res.words = append(res.words, other.words[i])
	}
	return res
}

func (m *Bitmap) Diff(other *Bitmap) *Bitmap {
	res := NewBitmap()
	for i := 0; i < m.Capacity() && i < other.Capacity(); i++ {
		if m.Has(i) && !other.Has(i) {
			res.Add(i)
		}
	}
	return res
}

func (m *Bitmap) Inter(other *Bitmap) *Bitmap {
	res := NewBitmap()
	for i := 0; i < m.Capacity() && i < other.Capacity(); i++ {
		if m.Has(i) && other.Has(i) {
			res.Add(i)
		}
	}
	return res
}

func (m *Bitmap) Not() *Bitmap {
	res := NewBitmap(m.Capacity())
	for i := 0; i < m.WordLength(); i++ {
		res.words[i] = math.MaxUint64 ^ m.words[i]
	}
	return res
}
