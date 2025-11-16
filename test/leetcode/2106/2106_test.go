package _106

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"math"
	"testing"
)

func maxTotalFruits(fruits [][]int, startPos int, k int) int {
	//----------------------------预处理----------------------------
	//把startPos和无穷大作为2个特殊点放入数组
	//前缀和，避免区间重复求和
	var indexList, sum []int
	pre := 0
	insertStartPos := false
	startIndex := -1
	for _, n := range fruits {
		if !insertStartPos && n[0] >= startPos {
			insertStartPos = true
			startIndex = len(indexList)
			if n[0] != startPos {
				indexList = append(indexList, startPos)
				sum = append(sum, pre)
			}
		}
		indexList = append(indexList, n[0])
		sum = append(sum, pre+n[1])
		pre += n[1]
	}
	if !insertStartPos {
		startIndex = len(indexList)
		indexList = append(indexList, startPos)
		sum = append(sum, pre)
	}
	indexList = append(indexList, math.MaxInt)
	sum = append(sum, pre)

	//----------------------------找到第一个需要开始遍历的点----------------------------
	cur := BinarySearchInt(indexList, startPos-k)
	if cur < 0 {
		cur = -cur - 1
		if cur == 0 && abs(startPos-indexList[cur]) > k {
			//到不了左端点
			return 0
		}
		if cur == len(indexList)-1 && abs(indexList[cur]-startPos) > k {
			//到不了右端点
			return 0
		}
	}

	//----------------------------开始遍历可能的点----------------------------
	res := 0
	for ; indexList[cur] <= startPos+k; cur++ {
		leftIndex := 0
		rightIndex := 0
		if indexList[cur] <= startPos {
			//先去左边
			//查看最远能到哪
			rightIndex = BinarySearchInt(indexList, indexList[cur]+(k-abs(startPos-indexList[cur])))
			if rightIndex < 0 {
				//二分查找如果数字不存在，-rightIndex - 1表示第一个比这个数大的元素所在位置，先去左边，则-rightIndex - 1这个地方到不了
				rightIndex = -rightIndex - 2
			}
			leftIndex = cur
			rightIndex = max(rightIndex, startIndex)
		} else {
			//先去右边
			//查看最远能到哪
			rightIndex = BinarySearchInt(indexList, indexList[cur]-(k-abs(startPos-indexList[cur])))
			if rightIndex < 0 {
				//二分查找如果数字不存在，-rightIndex - 1表示第一个比这个数大的元素所在位置，先去右边，则-rightIndex - 1这个地方到得了
				rightIndex = -rightIndex - 1
			}
			leftIndex = min(rightIndex, startIndex)
			rightIndex = cur
		}
		//区间求和
		diff := 0
		if leftIndex > 0 {
			diff = sum[leftIndex-1]
		}
		res = max(res, sum[rightIndex]-diff)
	}
	return res
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func BinarySearchInt(arr []int, target int) int {
	left := 0
	right := len(arr) - 1
	for left <= right {
		mid := (left + right) / 2
		if arr[mid] == target {
			return mid
		} else if arr[mid] > target {
			right = mid - 1
		} else {
			left = mid + 1
		}
	}
	return -left - 1
}
func Test_maxTotalFruits(t *testing.T) {
	Convey("maxTotalFruits", t, func() {
		Convey("maxTotalFruits test1", func() {
			res := maxTotalFruits([][]int{{29, 9}, {31, 1}}, 31, 2)
			fmt.Println(res)
		})
	})
}
