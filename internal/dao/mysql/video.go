package mysql

import (
	"tiktink/internal/model"
	"tiktink/pkg/tracer"

	"github.com/pkg/errors"
)

const (
	videoPageRows = 20 //分页一页的页数
)

type videoFunc interface {
	PublishVideo(video *model.Video) error
	QueryVideoExist(videoID string) (bool, error)
	QueryVideoByAuthorID(authorID string, pn int) ([]*model.VideoMSG, error)
}

type videoDealer struct{}

func NewVideoDealer() videoFunc {
	return &videoDealer{}
}
func (v *videoDealer) QueryVideoByAuthorID(authorID string, pn int) ([]*model.VideoMSG, error) {
	var videoMsgs []*model.VideoMSG
	offset := (pn - 1) * videoPageRows
	count := videoPageRows
	err := db.Raw("select `video_id`,`play_url`,`cover_url`,`favorite_count`,`comment_count`,`title`,"+
		"`user_id`,`user_name`,`follow_count`,`follower_count`"+
		"from `users` inner join `videos` on `videos`.`author_id`= `user_id`"+
		"where `videos`.author_id = ? limit ?,?", authorID, offset, count).Scan(&videoMsgs).Error
	if err != nil {
		return nil, errors.Wrap(err, tracer.FormatParam(authorID, pn))
	}
	return videoMsgs, nil
}

func (v *videoDealer) QueryVideoExist(videoID string) (bool, error) {
	cnt := new(int64)
	err := db.Raw("select 1 from videos where video_id = ?", videoID).Scan(cnt).Error
	if err != nil {
		return false, errors.Wrap(err, tracer.FormatParam(videoID))
	}
	return *cnt == 1, nil
}

func (v *videoDealer) PublishVideo(video *model.Video) error {
	err := db.Create(video).Error
	return errors.Wrap(err, tracer.FormatParam(video))
}
