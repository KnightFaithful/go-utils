package httputil

import (
	"encoding/json"
	"example.com/m/util/stringutil"
	"example.com/m/util/utilerror"
	"fmt"

	"io/ioutil"
	"net/http"
	"strings"
)

type HttpHelper struct {
	Host    string
	Cookie  string
	NeedLog bool
}

func (h *HttpHelper) Post(url string, header map[string]string, request interface{}, response interface{}) (interface{}, *utilerror.UtilError) {
	return h.doRequest(http.MethodPost, url, header, request, response)
}

func (h *HttpHelper) Get(url string, header map[string]string, request interface{}, response interface{}) (interface{}, *utilerror.UtilError) {
	return h.doRequest(http.MethodGet, url, header, request, response)
}

func (h *HttpHelper) doRequest(method, url string, header map[string]string, request interface{}, response interface{}) (interface{}, *utilerror.UtilError) {
	absUrl := h.Host + url

	postRequest, err2 := http.NewRequest(method, absUrl, strings.NewReader(stringutil.Object2String(request)))
	if err2 != nil {
		return nil, utilerror.NewError(err2.Error())
	}

	if header == nil {
		header = make(map[string]string)
		header["cookie"] = h.Cookie
	}
	for key, val := range header {
		postRequest.Header.Set(key, val)
	}
	if h.NeedLog {
		fmt.Println(fmt.Sprintf("doRequest header=%v, request=%v", stringutil.Object2String(header), stringutil.Object2String(request)))
	}
	client := &http.Client{}
	res, err := client.Do(postRequest)
	if err != nil {
		return nil, utilerror.NewError(err.Error())
	}

	respBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, utilerror.NewError(err.Error())
	}
	if h.NeedLog {
		fmt.Println("doRequest response=" + string(respBytes))
	}
	err = json.Unmarshal(respBytes, response)
	if err != nil {
		return nil, utilerror.NewError(err.Error())
	}
	return response, nil
}
