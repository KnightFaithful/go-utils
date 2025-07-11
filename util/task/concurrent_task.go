package task

import (
	"example.com/m/util/iters"
	"example.com/m/util/paginator"
	"example.com/m/util/utilerror"
	"fmt"
	"math"
	"sync"
)

func ConcurrentQueryTaskRunnerWithRoutineCount(totalList interface{}, oneTaskCount int64, task func(interface{}) (interface{}, *utilerror.UtilError),
	routineCount int, stopWhenErr ...bool) (map[interface{}]interface{}, *utilerror.UtilError) {

	w := &sync.WaitGroup{}
	ch := make(chan bool, routineCount)

	var errList []*utilerror.UtilError
	errListLock := sync.Mutex{}

	totalCount := iters.From(totalList).Count()
	totalpages := int(math.Ceil(float64(totalCount) / float64(oneTaskCount)))
	// 安全map
	newConnManger := NewConnManger()

	for index := 0; index < totalpages; index++ {
		w.Add(1)
		ch <- true
		// run
		innerIndex := int64(index)
		innerList := paginator.PagingBaseType(totalList, innerIndex+1, oneTaskCount)
		go func() {
			defer func() {
				<-ch
				w.Done()
			}()
			defer func() {
				if panicErr := recover(); panicErr != nil {
					fmt.Println(panicErr)
				}
			}()
			//默认遇到err时停止
			if len(errList) != 0 && (len(stopWhenErr) == 0 || stopWhenErr[0]) {
				return
			}
			result, err := task(innerList)
			if err != nil {
				errListLock.Lock()
				errList = append(errList, err)
				errListLock.Unlock()
			}
			newConnManger.Add(innerIndex, result)
		}()
	}
	w.Wait()
	// 后置处理
	var returnErr *utilerror.UtilError
	for _, err := range errList {
		if err != nil {
			if returnErr == nil {
				returnErr = err
				continue
			}
			returnErr = returnErr.AddError(err.Message())
		}
	}
	return newConnManger.Map, returnErr
}

func SyncExecuteTaskWithLimitCount(totalList interface{}, oneTaskCount int64, task func(interface{}) (interface{}, *utilerror.UtilError), stopWhenErr bool) (map[interface{}]interface{}, *utilerror.UtilError) {

	var errList []*utilerror.UtilError

	totalCount := iters.From(totalList).Count()
	totalpages := int(math.Ceil(float64(totalCount) / float64(oneTaskCount)))
	// 安全map
	newConnManger := NewConnManger()

	for index := 0; index < totalpages; index++ {
		// run
		innerIndex := int64(index)
		innerList := paginator.PagingBaseType(totalList, innerIndex+1, oneTaskCount)
		result, err := task(innerList)
		newConnManger.Add(innerIndex, result)
		if err != nil {
			errList = append(errList, err)
			if stopWhenErr {
				break
			}
		}
	}
	// 后置处理
	var returnErr *utilerror.UtilError
	for _, err := range errList {
		if err != nil {
			if returnErr == nil {
				returnErr = err
				continue
			}
			returnErr = returnErr.AddError(err.Message())
		}
	}
	return newConnManger.Map, returnErr
}
