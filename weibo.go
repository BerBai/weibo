package weibo

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var BadRequest = errors.New("BadRequest")

type User struct {
	ID   int64  `json:"id"`
	Name string `json:"screen_name"`
}

type Mblog struct {
	User        *User  `json:"user"`
	CreatedAt   string `json:"created_at"`
	ID          int64  `json:"id"`
	MblogID     string `json:"mblogid"`
	TextRaw     string `json:"text_raw"`
	Text        string `json:"text"`
	IsLongText  bool   `json:"isLongText"`
	Retweeted   *Mblog `json:"retweeted_status,omitempty"`
	LongTextRaw string
}

func (m *Mblog) TheText() string {
	if m.LongTextRaw != "" {
		return m.LongTextRaw
	}
	return m.TextRaw
}

func (mblog *Mblog) String() string {
	text := strings.ReplaceAll(mblog.TextRaw, "\n", "\\n")
	if len([]rune(text)) > 50 {
		text = string([]rune(text)[0:50]) + "..."
	}
	return fmt.Sprintf("%d | %s | %v | %v | %s", mblog.ID, mblog.MblogID, mblog.IsLongText, len(mblog.LongTextRaw) > 0, text)
}

type MymblogBody struct {
	Data struct {
		List []*Mblog `json:"list"`
	} `json:"data"`
	Ok int `json:"ok"`
}

type LongtextBody struct {
	Data struct {
		LongTextContent string `json:"longTextContent"`
	} `json:"data"`
	Ok int `json:"ok"`
}

type Client struct {
	Cookie string
	Proxy  string
}

func (c *Client) getJSON(_url string, body any) error {
	client := &http.Client{}
	if c.Proxy != "" {
		if proxyUrl, err := url.Parse(c.Proxy); err == nil {
			client.Transport = &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			}
		}
	}

	req, err := http.NewRequest("GET", _url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:107.0) Gecko/20100101 Firefox/107.0")
	req.Header.Set("Host", "weibo.com")
	req.Header.Set("Cookie", c.Cookie)
	req.Header.Set("Accept", "*/*")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusBadRequest {
		return BadRequest
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, body); err != nil {
		return err
	}
	return nil
}

func (c *Client) GetMblogs(userid string, page int, longtext bool) ([]*Mblog, error) {
	url := fmt.Sprintf("https://weibo.com/ajax/statuses/mymblog?uid=%s&page=%d&feature=0", userid, page)
	body := &MymblogBody{}
	if err := c.getJSON(url, body); err != nil {
		return nil, err
	} else if body.Ok != 1 {
		return nil, fmt.Errorf("body not ok")
	}
	var mblogs []*Mblog
	for _, v := range body.Data.List {
		if longtext {
			if err := c.FetchMblogLongText(v); err != nil {
				return nil, err
			}
			if v.Retweeted != nil {
				if err := c.FetchMblogLongText(v.Retweeted); err != nil {
					return nil, err
				}
			}
		}
		mblogs = append(mblogs, v)
	}
	return mblogs, nil
}

func (c *Client) GetMblogLongText(mblogid string) (longtext string, err error) {
	url := fmt.Sprintf("https://weibo.com/ajax/statuses/longtext?id=%s", mblogid)
	body := &LongtextBody{}
	if err = c.getJSON(url, body); err != nil {
		return
	}
	if body.Ok != 1 {
		err = fmt.Errorf("body not ok")
		return
	}
	longtext = body.Data.LongTextContent
	return
}

func (c *Client) FetchMblogLongText(mblog *Mblog) error {
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

type Database struct {
	DN  string
	DSN string
	db  *sql.DB
}

func (database *Database) getdb() (*sql.DB, error) {
	if database.db == nil {
		db, err := sql.Open(database.DN, database.DSN)
		if err != nil {
			return nil, err
		}
		database.db = db
	}
	return database.db, nil
}

func (database *Database) Close() {
	if database.db != nil {
		database.db.Close()
		database.db = nil
	}
}

func (database *Database) Migrate() error {
	db, err := database.getdb()
	if err != nil {
		return err
	}

	if _, err = db.Exec("CREATE TABLE IF NOT EXISTS mblog (UID BIGINT NOT NULL, ID BIGINT NOT NULL, MblogID VARCHAR(64) NOT NULL, TheText TEXT, CreatedAt CHAR(32), RetweetedUID BIGINT NOT NULL, RetweetedID BIGINT NOT NULL, RetweetedMblogID VARCHAR(64) NOT NULL, RetweetedTheText TEXT, RetweetedCreatedAt CHAR(32), PRIMARY KEY (UID,ID,MblogID))"); err != nil {
		return err
	}
	return nil
}

func (database *Database) HasMblog(mblog *Mblog) (bool, error) {
	db, err := database.getdb()
	if err != nil {
		return false, err
	}

	rows, err := db.Query("SELECT UID, ID, MblogID FROM mblog WHERE UID = ? AND ID = ? AND MblogID = ?", mblog.User.ID, mblog.ID, mblog.MblogID)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		return true, nil
	}
	return false, nil
}

func (database *Database) AddMblog(mblog *Mblog) error {
	db, err := database.getdb()
	if err != nil {
		return err
	}

	var uid, id int64
	var mblogID, theText, createdAt string
	if mblog.Retweeted != nil {
		uid = mblog.Retweeted.User.ID
		id = mblog.Retweeted.ID
		mblogID = mblog.Retweeted.MblogID
		theText = mblog.Retweeted.TheText()
		createdAt = mblog.Retweeted.CreatedAt
	}
	if _, err := db.Exec("INSERT INTO mblog(UID, ID, MblogID, TheText, CreatedAt, RetweetedUID, RetweetedID, RetweetedMblogID, RetweetedTheText, RetweetedCreatedAt) VALUES(?,?,?,?,?,?,?,?,?,?)",
		mblog.User.ID, mblog.ID, mblog.MblogID, mblog.TheText(), mblog.CreatedAt, uid, id, mblogID, theText, createdAt); err != nil {
		return err
	}
	return nil
}
