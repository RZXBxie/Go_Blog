package other

import "encoding/json"

type Data struct {
	ID  *string         `json:"id"`
	Doc json.RawMessage `json:"doc"`
}

type ESIndexResponse struct {
	Data []Data `json:"data"`
}
