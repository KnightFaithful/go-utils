package copier

import (
	"encoding/json"
	"fmt"
	"testing"
)

type PrintSetting struct {
	Size string `json:"size"`
	Qty  int64  `json:"qty"`
}

func TestSetStructFieldByMap(t *testing.T) {
	a := &PrintSetting{}
	d := map[string]interface{}{"size": "asd", "qty": "23"}

	r, _ := SetStructFieldByMap(a, d)

	fmt.Println("==============")
	fmt.Println(r)
	b, _ := json.Marshal(r)
	fmt.Println(string(b))

}
