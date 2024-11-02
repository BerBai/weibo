package weibo

import "fmt"

type CommentBody struct {
	Ok          int           `json:"ok"`
	Data        []*Comments   `json:"data"`
	RootComment []interface{} `json:"rootComment"`
	TotalNumber int           `json:"total_number"`
	TipMsg      string        `json:"tip_msg"`
	MaxId       int64         `json:"max_id"`
	TrendsText  string        `json:"trendsText"`
}

type Comments struct {
	User                 *User         `json:"user"`
	CreatedAt            string        `json:"created_at"`
	Id                   int64         `json:"id"`
	Rootid               int64         `json:"rootid"`
	Rootidstr            string        `json:"rootidstr"`
	FloorNumber          int           `json:"floor_number"`
	Text                 string        `json:"text"`
	DisableReply         int           `json:"disable_reply"`
	RestrictOperate      int           `json:"restrictOperate"`
	SourceAllowclick     int           `json:"source_allowclick"`
	SourceType           int           `json:"source_type"`
	Source               string        `json:"source"`
	Mid                  string        `json:"mid"`
	Idstr                string        `json:"idstr"`
	UrlObjects           []interface{} `json:"url_objects"`
	Liked                bool          `json:"liked"`
	Readtimetype         string        `json:"readtimetype"`
	AnalysisExtra        string        `json:"analysis_extra"`
	SafeTags             int64         `json:"safe_tags"`
	MarkType             int           `json:"mark_type,omitempty"`
	MatchAiPlayPicture   bool          `json:"match_ai_play_picture"`
	Rid                  string        `json:"rid"`
	AllowFollow          bool          `json:"allow_follow"`
	ItemCategory         string        `json:"item_category"`
	Comments             []*Comments   `json:"comments"`
	ReplyComment         *Comments     `json:"reply_comment,omitempty"`
	MaxId                int64         `json:"max_id"`
	TotalNumber          int           `json:"total_number"`
	IsLikedByMblogAuthor bool          `json:"isLikedByMblogAuthor"`
	LikeCounts           int           `json:"like_counts"`
	MoreInfo             struct {
		Scheme        string `json:"scheme"`
		Text          string `json:"text"`
		HighlightText string `json:"highlight_text"`
	} `json:"more_info"`
	TextRaw string      `json:"text_raw"`
	Urls    []UrlStruct `json:"url_struct,omitempty"`
}

type UrlStruct struct {
	UrlTitle    string                 `json:"url_title"`
	UrlTypePic  string                 `json:"url_type_pic"`
	OriUrl      string                 `json:"ori_url"`
	ShortUrl    string                 `json:"short_url"`
	LongUrl     string                 `json:"long_url"`
	UrlType     int                    `json:"url_type"`
	Result      bool                   `json:"result"`
	StorageType string                 `json:"storage_type"`
	Hide        int                    `json:"hide"`
	ObjectType  string                 `json:"object_type"`
	Position    int                    `json:"position"`
	PicInfos    map[string]interface{} `json:"pic_infos"`
	PicIds      []string               `json:"pic_ids"`
	GifName     string                 `json:"gif_name"`
	H5TargetUrl string                 `json:"h5_target_url"`
	NeedSaveObj int                    `json:"need_save_obj"`
}

// GetComments
// - flow：0-按热度排序；1-按时间排序
// - mid：博文的id；评论id
// - userid：用户id
// - isMix：0-首页；1-后续页，需要max_id参数（上一页返回JSON的max_id)
// - fetchLevel：0-博文下的评论；1-评论下的评论
// - type：feed-简要
func (c *Client) GetComments(flow int, mid int64, userid string, isMax int, maxId int64, fetchLevel int, longtext bool) (*CommentBody, error) {
	blogUrl := fmt.Sprintf("https://weibo.com/ajax/statuses/buildComments?flow=%d&is_reload=1&id=%d&is_show_bulletin=2&is_mix=%d&max_id=%d&count=20&type=1&uid=%s&fetch_level=%d&locale=zh-CN", flow, mid, isMax, maxId, userid, fetchLevel)

	body := &CommentBody{}
	if err := c.getJSON(blogUrl, body); err != nil {
		return nil, err
	} else if body.Ok != 1 {
		return nil, fmt.Errorf("body not ok")
	}
	//var mblogs []*Mblog
	//for _, v := range body.Data.List {
	//	if longtext {
	//		if err := c.FetchMblogLongText(v); err != nil {
	//			return nil, err
	//		}
	//		if v.Retweeted != nil {
	//			if err := c.FetchMblogLongText(v.Retweeted); err != nil {
	//				return nil, err
	//			}
	//		}
	//	}
	//	mblogs = append(mblogs, v)
	//}
	return body, nil
}
