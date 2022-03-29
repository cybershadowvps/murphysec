package api

import (
	"murphysec-cli-simple/utils/must"
)

func StartCheck(taskId string) error {
	must.True(taskId != "")
	httpReq := C.PostJson("/message/v2/access/client/start_check", map[string]interface{}{"task_info": taskId})
	return C.Do(httpReq, nil)
}
