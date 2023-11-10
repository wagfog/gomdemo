package common

import "encoding/json"

func SwapTo(request, product interface{}) error {
	Byteas, err := json.Marshal(request)
	if err != nil {
		return err
	}
	err = json.Unmarshal(Byteas, product)
	return err
}
