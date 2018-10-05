package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/lagoon-platform/api/storage"
	"github.com/stretchr/testify/assert"
)

//************************************************
// EMPTY STORAGE WITHOUT ANY KEYS - START
//************************************************
func TestGetNoContent(t *testing.T) {
	usedStorage = storage.GetMockStorage()
	defer usedStorage.Clean()

	logger = *log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	initLog(false, false)

	handler := http.HandlerFunc(getValue)
	server := httptest.NewServer(handler)
	defer server.Close()

	resp, err := http.Get(server.URL + "/" + "dummy_id")
	if err != nil {
		t.Fatal(err)
	}

	checkResponseCode(t, http.StatusNotFound, resp)
}

func TestDeleteNoContent(t *testing.T) {

	usedStorage = storage.GetMockStorage()
	defer usedStorage.Clean()

	logger = *log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	initLog(false, false)

	handler := http.HandlerFunc(deleteValue)
	server := httptest.NewServer(handler)
	defer server.Close()

	req := httptest.NewRequest(http.MethodDelete, server.URL+"/"+"dummy_id", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()

	checkResponseCode(t, http.StatusNotFound, resp)
}

func TestGetKeysNoContent(t *testing.T) {
	usedStorage = storage.GetMockStorage()
	defer usedStorage.Clean()

	logger = *log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	initLog(false, false)

	defer usedStorage.Clean()

	handler := http.HandlerFunc(getKeys)
	server := httptest.NewServer(handler)
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	keys := make([]string, 0)
	err = json.Unmarshal(body, &keys)
	if err != nil {
		t.Fatal(err)
	}

	assert.Nil(t, err)
	assert.Equal(t, len(keys), 0)

	checkResponseCode(t, http.StatusOK, resp)
}

//************************************************
// EMPTY STORAGE WITHOUT ANY KEYS - END
//************************************************

func TestSaveValue(t *testing.T) {

	strKey := "test_key"
	strValue := "test_value"

	usedStorage = storage.GetMockStorage()
	defer usedStorage.Clean()

	logger = *log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	initLog(false, false)

	handler := http.HandlerFunc(saveValue)
	server := httptest.NewServer(handler)
	defer server.Close()

	body := StorePostRequest{
		Key:   strKey,
		Value: strValue,
	}

	jsonStr, err := json.Marshal(body)
	assert.Nil(t, err)

	req := httptest.NewRequest(http.MethodPost, server.URL, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()

	checkResponseCode(t, http.StatusCreated, resp)

	b, err := usedStorage.Contains(strKey)
	assert.Nil(t, err)
	assert.True(t, b)
	b, val, err := usedStorage.Get(strKey)
	assert.Nil(t, err)
	assert.True(t, b)
	assert.Equal(t, string(val), strValue)

}

func TestGetKeys(t *testing.T) {
	usedStorage = storage.GetMockStorage()
	defer usedStorage.Clean()

	strKey1 := "test_key1"
	strValue1 := "test_value1"
	strKey2 := "test_key2"
	strValue2 := "test_value2"
	strKey3 := "test_key3"
	strValue3 := "test_value3"

	usedStorage.StoreString(strKey1, strValue1)
	usedStorage.StoreString(strKey2, strValue2)
	usedStorage.StoreString(strKey3, strValue3)

	logger = *log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	initLog(false, false)

	handler := http.HandlerFunc(getKeys)
	server := httptest.NewServer(handler)
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	checkResponseCode(t, http.StatusOK, resp)

	keys := make([]string, 0)
	err = json.Unmarshal(body, &keys)
	if err != nil {
		t.Fatal(err)
	}

	assert.Nil(t, err)
	assert.Equal(t, len(keys), 3)
	assert.Contains(t, keys, strKey1, strKey2, strKey3)
}

func TestGetValue(t *testing.T) {

	strKey := "test_key"
	strValue := "test_value"

	usedStorage = storage.GetMockStorage()
	usedStorage.StoreString(strKey, strValue)
	defer usedStorage.Clean()

	logger = *log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	initLog(false, false)

	application = App{}
	application.initialize()

	req, _ := http.NewRequest("GET", "/storage/"+strKey, nil)
	respRecorder := executeRequest(req)

	checkCode(t, http.StatusOK, respRecorder.Code)

	sPr := StorePostRequest{}
	err := json.Unmarshal(respRecorder.Body.Bytes(), &sPr)
	assert.Nil(t, err)
	assert.Equal(t, sPr.Key, strKey)
	assert.Equal(t, sPr.Value, strValue)

}

func TestDeleteValue(t *testing.T) {

	strKey := "test_key"
	strValue := "test_value"

	usedStorage = storage.GetMockStorage()
	usedStorage.StoreString(strKey, strValue)
	defer usedStorage.Clean()

	logger = *log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	initLog(false, false)

	application = App{}
	application.initialize()

	req, _ := http.NewRequest("DELETE", "/storage/"+strKey, nil)
	respRecorder := executeRequest(req)
	checkCode(t, http.StatusOK, respRecorder.Code)

	b, err := usedStorage.Contains(strKey)

	assert.Nil(t, err)
	assert.False(t, b)

}