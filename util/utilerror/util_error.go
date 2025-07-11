package utilerror

import (
	"bytes"
	"fmt"
	"runtime"
)

type UtilError struct {
	Items []*UtilErrorItem
}

type UtilErrorItem struct {
	Message string
	File    string
	Line    int
	Params  []interface{}
}

func NewError(message string, params ...interface{}) *UtilError {
	_, file, line, _ := runtime.Caller(1)
	o := &UtilErrorItem{
		Message: message,
		File:    file,
		Line:    line,
		Params:  params,
	}
	return &UtilError{Items: []*UtilErrorItem{o}}
}

// AddError 不能将w作为params参数传递
func (w *UtilError) AddError(message string, params ...interface{}) *UtilError {
	_, file, line, _ := runtime.Caller(1)
	o := &UtilErrorItem{
		Message: message,
		File:    file,
		Line:    line,
		Params:  params,
	}
	w.Items = append(w.Items, o)
	return w
}

func (w *UtilError) Mark() *UtilError {
	_, file, line, _ := runtime.Caller(1)
	o := &UtilErrorItem{
		Message: w.Message(),
		File:    file,
		Line:    line,
		Params:  w.Params(),
	}
	w.Items = append(w.Items, o)
	return w
}

func (w *UtilError) Error() string {
	if w == nil {
		return ""
	}
	o := w.Items[len(w.Items)-1]
	return fmt.Sprintf(o.Message, o.Params...)
}

func (w *UtilError) DebugError() string {
	if w == nil {
		return "success"
	}
	buf := bytes.NewBufferString("")
	for i, o := range w.Items {
		s1 := fmt.Sprintf("%v %v |%v", o.File, o.Line, o.Message)
		// tips:解决用户使用的时候直接add自身出现的栈溢出问题
		var params []interface{}
		for _, c := range o.Params {
			if c != w {
				params = append(params, c)
			}
		}
		s2 := fmt.Sprintf(s1, params...)
		buf.WriteString(s2)
		if i != len(w.Items)-1 {
			buf.WriteString("\n")
		}
	}
	return buf.String()
}

func (w *UtilError) Message() string {
	if w == nil || len(w.Items) == 0 {
		return "success"
	}
	return w.Items[len(w.Items)-1].Message
}

func (w *UtilError) Params() []interface{} {
	if w == nil || len(w.Items) == 0 {
		return nil
	}
	if len(w.Items[len(w.Items)-1].Params) == 0 {
		return nil
	}
	return w.Items[len(w.Items)-1].Params
}

func (w *UtilError) Values() []interface{} {
	if w == nil || len(w.Items) == 0 {
		return []interface{}{}
	}
	return w.Items[len(w.Items)-1].Params
}

func FormatErrs(errs []*UtilError) *UtilError {
	var returnErr *UtilError
	for _, err := range errs {
		if err != nil {
			if returnErr == nil {
				returnErr = err
				continue
			}
			returnErr = returnErr.AddError(err.Message())
		}
	}
	return returnErr
}
