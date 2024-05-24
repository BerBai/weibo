package weibo

import (
	"os"
	"testing"
)

const (
	uid    = "6874180501"
	Cookie = "XSRF-TOKEN=JRdTtx4jrsyK5tWAHWvM8mpJ; SUB=_2AkMTKKlJf8NxqwFRmP8RzWLkbY10zwrEieKldFiSJRMxHRl-yT9kqlM8tRB6OKiHpmrIgcUy6YQdWlF4Q9LVcDAvvpWG;"
	//Cookie = "XSRF-TOKEN=JRdTtx4jrsyK5tWAHWvM8mpJ; SUB=_2A25LVGWpDeRhGeBP41cW-C7EzzmIHXVoKOdhrDV8PUJbkNAbLUPBkW1NRTRuIoa9saz5_WBV5HyIPpx9dy3wpwu6;"
)

func TestWeibo(t *testing.T) {
	c := Client{Cookie: Cookie, Proxy: os.Getenv("Proxy"), Check: DefaultCheck()}

	if mblogs, err := c.GetMblogs(uid, 1, true); err != nil {
		t.Error()
	} else {
		for _, mblog := range mblogs {
			t.Log(mblog)
			break
			if mblog.Retweeted != nil {
				t.Log(mblog.Retweeted.PicUrls())
			}
		}
	}

	activate, err := c.CheckCookie()
	t.Log(activate)
	if err != nil {
		t.Log(err.Error())
	}
}
