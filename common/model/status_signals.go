package model

type StopSignal struct {
	WorkerId int64 `json:"worker_id"`
}

type StartSignal struct {
	WorkerId int64 `json:"worker_id"`
}
