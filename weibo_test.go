package weibo

import (
	"os"
	"testing"
)

const (
	uid    = "1223178222"
	Cookie = "SUB=_2AkMTKKlJf8NxqwFRmP8RzWLkbY10zwrEieKldFiSJRMxHRl-yT9kqlM8tRB6OKiHpmrIgcUy6YQdWlF4Q9LVcDAvvpWG;"
)

func TestWeibo(t *testing.T) {
	c := Client{Cookie: Cookie, Proxy: os.Getenv("Proxy")}
	if mblogs, err := c.GetMblogs(uid, 1, true); err != nil {
		t.Error()
	} else {
		for _, mblog := range mblogs {
			t.Log(mblog)
		}
	}
}
