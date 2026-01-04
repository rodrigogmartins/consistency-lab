package api

import "consistency-lab/internal/store"

type PutBody struct {
	Value string `json:"value"`
}

type PutResp struct {
	Item store.Item `json:"item"`
	Mode string     `json:"mode"`
}

type ErrResp struct {
	Error string `json:"error"`
}
