package rest

import (
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/ekara-platform/api/storage"
)

func TestGetInfo(t *testing.T) {

	usedStorage = storage.GetMockStorage()
	// TODO HERE ADD THE REAL STUFF WHICH IS SUPPOSED TO BE RETURNED
	defer usedStorage.Clean(storage.STORAGE_PREFIX)

	logger = *log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	initLog(false, false)

	application = Api{}
	application.initialize()

	req, _ := http.NewRequest(http.MethodGet, "/infos/", nil)
	respRecorder := executeRequest(req)
	resp := respRecorder.Result()

	checkResponseCode(t, http.StatusOK, resp)
	if body := getNotEmptyBody(t, resp); body != "" {
		if b, result := checkJsonRoundTpip(t, body, &ApiInfo{}); !b {
			t.Errorf(`Not the same received "%s" converted to "%s"`, body, result)
		}
	}
}
