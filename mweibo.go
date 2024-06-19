package weibo

import (
	"fmt"
	strip "github.com/grokify/html-strip-tags-go"
	"strings"
)

type CMblogBody struct {
	Data struct {
		Cards []*Card `json:"cards"`
	} `json:"data"`
	Ok  int    `json:"ok"`
	Msg string `json:"msg,omitempty"`
}

type Card struct {
	CardType       int8         `json:"card_type"`
	ShowType       int8         `json:"show_Type"`
	SkipGroupTitle bool         `json:"skip_group_title,omitempty"`
	CardGroup      []*CardGroup `json:"card_group,omitempty"`
	Mblog          CMblog       `json:"mblog,omitempty"`
}

type CardGroup struct {
	CardType int8   `json:"card_type"`
	ShowType int8   `json:"show_Type"`
	Mblog    CMblog `json:"mblog"`
}

type CMblog struct {
	CreatedAt   string      `json:"created_at"`
	ID          string      `json:"id"`
	Text        string      `json:"text"`
	PicIds      []string    `json:"pic_ids"`
	User        *User       `json:"user"`
	IsLongText  bool        `json:"isLongText"`
	ActionInfo  *ActionInfo `json:"action_info"`
	PicNum      int8        `json:"pic_num"`
	MblogID     string      `json:"bid"`
	Pics        []*Pics     `json:"pics,omitempty"`
	Retweeted   *CMblog     `json:"retweeted_status,omitempty"`
	LongTextRaw string
}

type ActionInfo struct {
	Comment *Comment `json:"comment"`
}

type Comment struct {
	List  []*CommentBlog `json:"list"`
	Count int64          `json:"count"`
}

type CommentBlog struct {
	CreatedAt string `json:"created_at"`
	ID        int64  `json:"id"`
	Text      string `json:"text"`
	User      *User  `json:"user"`
}

type Pics struct {
	Pid   string `json:"pid"`
	Url   string `json:"url"`
	Large *Large `json:"large"`
}

type Large struct {
	Url string `json:"url"`
}

func (m *CMblog) TheText() string {
	if m.LongTextRaw != "" {
		return m.LongTextRaw
	}
	text := strings.ReplaceAll(m.Text, "<br />", "\n")
	text = strip.StripTags(text)

	return text
}

func (c *Client) FetchCMblogLongText(mblog *CMblog) error {
	if mblog.IsLongText {
		if longtext, err := c.GetMblogLongText(mblog.MblogID); err != nil {
			if err == BadRequest {
				return nil
			}
			return err
		} else {
			mblog.LongTextRaw = longtext
			return nil
		}
	} else {
		return nil
	}
}

// 获取互动信息，需要cookie
func (c *Client) GetCMblogs(userid string, page int, longtext bool) ([]*CMblog, error) {
	blogUrl := fmt.Sprintf("https://m.weibo.cn/api/container/getIndex?containerid=230869%s_-_comment&page_type=03&page=%d", userid, page)
	body := &CMblogBody{}
	if err := c.getJSON(blogUrl, body); err != nil {
		return nil, err
	} else if body.Ok != 1 {
		return nil, fmt.Errorf("body not ok")
	}
	var mblogs []*CMblog
	for _, card := range body.Data.Cards {
		//if card.ShowType == 0 {
		//	continue
		//}
		if card.CardType == 11 {
			if len(card.CardGroup) == 0 || card.SkipGroupTitle {
				continue
			}
			if longtext {
				if err := c.FetchCMblogLongText(&card.CardGroup[0].Mblog); err != nil {
					return nil, err
				}
				if card.CardGroup[0].Mblog.Retweeted != nil {
					if err := c.FetchCMblogLongText(card.CardGroup[0].Mblog.Retweeted); err != nil {
						return nil, err
					}
				}
			}
			mblogs = append(mblogs, &card.CardGroup[0].Mblog)
		} else if card.CardType == 9 {
			if longtext {
				if err := c.FetchCMblogLongText(&card.Mblog); err != nil {
					return nil, err
				}
				if card.Mblog.Retweeted != nil {
					if err := c.FetchCMblogLongText(card.Mblog.Retweeted); err != nil {
						return nil, err
					}
				}
			}
			mblogs = append(mblogs, &card.Mblog)
		}
	}
	return mblogs, nil
}

// 手机端api，获取全部微博，需要cookie
func (c *Client) GetMMblogs(userid string, page int, longtext bool) ([]*CMblog, error) {
	blogUrl := fmt.Sprintf("https://m.weibo.cn/api/container/getIndex?containerid=230413%s_-_WEIBO_SECOND_PROFILE_WEIBO&page_type=01&page=%d", userid, page)
	body := &CMblogBody{}
	if err := c.getJSON(blogUrl, body); err != nil {
		return nil, err
	} else if body.Ok != 1 {
		return nil, fmt.Errorf("body not ok")
	}
	var mblogs []*CMblog
	for _, card := range body.Data.Cards {
		//if card.ShowType == 0 {
		//	continue
		//}
		if card.CardType == 11 {
			if len(card.CardGroup) == 0 || card.SkipGroupTitle {
				continue
			}
			if longtext {
				if err := c.FetchCMblogLongText(&card.CardGroup[0].Mblog); err != nil {
					return nil, err
				}
				if card.CardGroup[0].Mblog.Retweeted != nil {
					if err := c.FetchCMblogLongText(card.CardGroup[0].Mblog.Retweeted); err != nil {
						return nil, err
					}
				}
			}
			mblogs = append(mblogs, &card.CardGroup[0].Mblog)
		} else if card.CardType == 9 {
			if longtext {
				if err := c.FetchCMblogLongText(&card.Mblog); err != nil {
					return nil, err
				}
				if card.Mblog.Retweeted != nil {
					if err := c.FetchCMblogLongText(card.Mblog.Retweeted); err != nil {
						return nil, err
					}
				}
			}
			mblogs = append(mblogs, &card.Mblog)
		}
	}
	return mblogs, nil
}
