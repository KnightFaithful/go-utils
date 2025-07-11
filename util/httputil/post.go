package httputil

import (
	"bytes"
	"example.com/m/util/convert"
	"example.com/m/util/utilerror"

	"io/ioutil"
	"mime/multipart"
	"net/http"
	"unsafe"
)

type PostWithFormDataRequest struct {
	Url          string
	Header       map[string]string
	FormDataList []*FormDataItem
}

type FormDataItem struct {
	FieldName string
	FileName  string
	FileValue []byte
}

type PostWithFormDataResponse struct {
	StatusCode   int
	ResponseBody string
}

func PostWithFormData(req *PostWithFormDataRequest) (*PostWithFormDataResponse, *utilerror.UtilError) {
	if req == nil {
		return nil, utilerror.ErrorRequestParamCanNotBeNull
	}

	reqBody := new(bytes.Buffer)
	w := multipart.NewWriter(reqBody)
	//构建form-data
	for _, formData := range req.FormDataList {
		createFormFile, err := w.CreateFormFile(formData.FieldName, formData.FileName)
		if err != nil {
			return nil, utilerror.NewError(err.Error())
		}
		_, err = createFormFile.Write(formData.FileValue)
		if err != nil {
			return nil, utilerror.NewError(err.Error())
		}
	}
	//close才能写入reqBody
	err := w.Close()
	if err != nil {
		return nil, utilerror.NewError(err.Error())
	}

	postRequest, err2 := http.NewRequest(http.MethodPost, req.Url, reqBody)
	if err2 != nil {
		return nil, utilerror.NewError(err2.Error())
	}

	postRequest.Header.Set("Content-Type", w.FormDataContentType())
	for key, val := range req.Header {
		postRequest.Header.Set(key, val)
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

	str := (*string)(unsafe.Pointer(&respBytes))
	return &PostWithFormDataResponse{
		StatusCode:   res.StatusCode,
		ResponseBody: convert.StringValue(str),
	}, nil
}
