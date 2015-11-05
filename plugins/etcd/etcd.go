package etcd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"

	etcd "github.com/coreos/etcd/client"

	"github.com/customerio/monitor/metrics"
)

type Etcd struct {
	slackURL string
	client   etcd.Client
	last     string
}

func (c *Etcd) postSlack(msg string) {
	if len(c.slackURL) == 0 {
		fmt.Printf("%s", msg)
		return
	}

	client := &http.Client{Timeout: time.Second * 10}

	//ip := cfg.hostIP
	var ip string
	if len(ip) == 0 {
		var err error
		ip, err = os.Hostname()
		if err != nil {
			ip = "unknown"
		}
	}
	type message struct {
		Text string `json:"text"`
	}
	m := message{Text: fmt.Sprintf("report from host %s\n%s", ip, msg)}

	body, err := json.Marshal(&m)
	if err != nil {
		fmt.Printf("etcd: post stack notification: %v: %s\n", err, msg)
		return
	}

	v := url.Values{}
	v.Set("payload", string(body))
	resp, err := client.PostForm(c.slackURL, v)
	if err != nil {
		fmt.Printf("etcd: post stack notification: %v: %s\n", err, msg)
		return
	}

	defer resp.Body.Close()
	io.Copy(ioutil.Discard, resp.Body)
}

func New(slack string, u string) *Etcd {
	urls := strings.Split(u, ",")
	c, err := etcd.New(etcd.Config{
		Endpoints:               urls,
		Transport:               etcd.DefaultTransport,
		HeaderTimeoutPerRequest: 10 * time.Second,
	})
	if err != nil {
		panic(err)
	}

	return &Etcd{
		slackURL: slack,
		client:   c,
	}
}

func (c *Etcd) Collect(batch *metrics.Batch) {
	mi := etcd.NewMembersAPI(c.client)
	ms, err := mi.List(context.TODO())

	b := &bytes.Buffer{}
	healthy := 0
	health := false
	defer func() {
		msg := string(b.Bytes())
		// If the messages are the same we're done.
		if msg == c.last {
			return
		}

		// We don't want to log healthy state initially.
		if healthy == len(ms) && len(c.last) == 0 {
			c.last = msg
			return
		}

		c.last = msg
		c.postSlack(msg)
	}()

	if err != nil {
		fmt.Fprintf(b, "etcd: cluster may be unhealthy: failed to list members\n")
		return
	}
	hc := &http.Client{Timeout: time.Second * 5}

	for _, m := range ms {
		if len(m.ClientURLs) == 0 {
			fmt.Fprintf(b, "etcd: member %s is unreachable: no available published client urls\n", m.ID)
			continue
		}

		checked := false
		for _, url := range m.ClientURLs {
			resp, err := hc.Get(url + "/health")
			if err != nil {
				fmt.Fprintf(b, "failed to check the health of member %s on %s: %v\n", m.ID, url, err)
				continue
			}

			result := struct{ Health string }{}
			nresult := struct{ Health bool }{}
			bytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Fprintf(b, "failed to check the health of member %s on %s: %v\n", m.ID, url, err)
				continue
			}
			resp.Body.Close()

			err = json.Unmarshal(bytes, &result)
			if err != nil {
				err = json.Unmarshal(bytes, &nresult)
			}
			if err != nil {
				fmt.Fprintf(b, "failed to check the health of member %s on %s: %v\n", m.ID, url, err)
				continue
			}

			checked = true
			if result.Health == "true" || nresult.Health == true {
				healthy++
				health = true
				fmt.Fprintf(b, "member %s is healthy: got healthy result from %s\n", m.ID, url)
			} else {
				fmt.Fprintf(b, "member %s is unhealthy: got unhealthy result from %s\n", m.ID, url)
			}
			break
		}
		if !checked {
			fmt.Fprintf(b, "member %s is unreachable: %v are all unreachable\n", m.ID, m.ClientURLs)
		}
	}

	if health {
		fmt.Fprintln(b, "cluster is healthy")
	} else {
		fmt.Fprintln(b, "cluster is unhealthy")
	}
}

func (c *Etcd) clear() {
}
