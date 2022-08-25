package model

import "tiktink/internal/code"

type FeedReq struct {
	LatestTime int64  `form:"latest_time"`
	Token      string `form:"token"`
}

type FeedResp struct {
	StatusCode code.ResCode `json:"status_code"`
	StatusMsg  string       `json:"status_msg"`
	NextTime   int64        `json:"next_time"`
	VideoList  []*VideoMSG  `json:"video_list"`
}
