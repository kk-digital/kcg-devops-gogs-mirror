package pretty

import "encoding/json"

func JSON(v interface{}) string {
	data, _ := json.MarshalIndent(v, "", "\t")
	return string(data)
}
