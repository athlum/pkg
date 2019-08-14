package utils

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func GetBody(request *http.Request) ([]byte, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(request.Body)
	if err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}

func ToBody(context string) io.Reader {
	buf := new(bytes.Buffer)
	buf.WriteString(context)
	return buf
}

func ToBodyByte(context []byte) io.Reader {
	buf := new(bytes.Buffer)
	buf.Write(context)
	return buf
}

func GetResp(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func TooManyRequests() (int, []byte) {
	return http.StatusTooManyRequests, []byte(fmt.Sprintf("{\"success\": false}"))
}

func InternalServerError(err error) (int, []byte) {
	return http.StatusInternalServerError, []byte(fmt.Sprintf("{\"success\": false, \"message\": \"Internal server error: %v\"}", err.Error()))
}

func BadRequest(err error) (int, []byte) {
	return http.StatusBadRequest, []byte(fmt.Sprintf("{\"success\": false, \"message\": \"%v\"}", err.Error()))
}

func NotFound(err error) (int, []byte) {
	return http.StatusNotFound, []byte(fmt.Sprintf("{\"success\": false, \"message\": \"%v\"}", err.Error()))
}

func Ok(msg []byte) (int, []byte) {
	return http.StatusOK, msg
}

func GetTemplate(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}
