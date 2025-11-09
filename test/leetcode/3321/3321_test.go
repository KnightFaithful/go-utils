package _321

import (
	"cmp"
	"fmt"
	"github.com/emirpasic/gods/v2/trees/redblacktree"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func findXSum(nums []int, k int, x int) []int64 {
	var res []int64
	sw := NewSlidingWindow(x)
	for i := 0; i < k; i++ {
		sw.Push(int64(nums[i]))
	}
	for i := k; i < len(nums); i++ {
		res = append(res, sw.Sum())
		sw.Push(int64(nums[i]))
		sw.Pop(int64(nums[i-k]))
	}
	res = append(res, sw.Sum())
	return res
}

type Item struct {
	Value int64 //值
	Num   int64 //频率
}

func less(p, q Item) int {
	return int(cmp.Or(p.Num-q.Num, p.Value-q.Value))
}

type SlidingWindow struct {
	master, backup *redblacktree.Tree[Item, struct{}]
	sum            int64
	x              int
	count          map[int64]int64
}

func NewSlidingWindow(x int) *SlidingWindow {
	sw := &SlidingWindow{
		master: redblacktree.NewWith[Item, struct{}](less),
		backup: redblacktree.NewWith[Item, struct{}](less),
		sum:    0,
		x:      x,
		count:  make(map[int64]int64),
	}
	return sw
}

func (sw *SlidingWindow) Push(val int64) {
	if sw.count[val] == 0 && sw.master.Size() < sw.x {
		sw.operateMaster(val, 1)
		return
	}
	temp := Item{Value: val, Num: sw.count[val]}
	if _, ok := sw.master.Get(temp); ok {
		sw.operateMaster(val, 1)
		return
	}
	sw.operateBackup(val, 1)
}

func (sw *SlidingWindow) Pop(val int64) {
	temp := Item{Value: val, Num: sw.count[val]}
	if _, ok := sw.master.Get(temp); ok {
		sw.operateMaster(val, -1)
		return
	}
	sw.operateBackup(val, -1)
}

func (sw *SlidingWindow) operateMaster(val int64, num int64) {
	temp := Item{Value: val, Num: sw.count[val]}
	sw.count[val] += num
	sw.master.Remove(temp)
	temp.Num = sw.count[val]
	sw.master.Put(temp, struct{}{})
	sw.sum += num * val
	if num < 0 {
		sw.resort()
	}
}

func (sw *SlidingWindow) operateBackup(val int64, num int64) {
	temp := Item{Value: val, Num: sw.count[val]}
	sw.count[val] += num
	sw.backup.Remove(temp)
	temp.Num = sw.count[val]
	if temp.Num == 0 {
		return
	}
	sw.backup.Put(temp, struct{}{})
	if num > 0 {
		sw.resort()
	}
}

func (sw *SlidingWindow) resort() {
	masterMin := sw.master.Left().Key
	if !sw.backup.Empty() {
		backupMax := sw.backup.Right().Key
		if less(masterMin, backupMax) < 0 {
			sw.master.Remove(masterMin)
			sw.master.Put(backupMax, struct{}{})
			sw.backup.Remove(backupMax)
			sw.backup.Put(masterMin, struct{}{})
			sw.sum += backupMax.Value*backupMax.Num - masterMin.Value*masterMin.Num
		}
	}
}

func (sw *SlidingWindow) Sum() int64 {
	return sw.sum
}

func Test_findXSum(t *testing.T) {
	Convey("findXSum", t, func() {
		Convey("findXSum test1", func() {
			res := findXSum([]int{1, 1, 2, 2, 3, 4, 2, 3}, 6, 2)
			fmt.Println(res)
		})
	})
}
