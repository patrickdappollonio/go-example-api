package app

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func errorJSON(w http.ResponseWriter, format string, args ...interface{}) {
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": fmt.Sprintf(format, args...),
	})
}
