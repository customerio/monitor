package etcd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"

	etcd "github.com/coreos/etcd/client"

	"github.com/customerio/monitor/metrics"
	"github.com/customerio/monitor/notifiers/slack"
)

type Etcd struct {
	hostname    string
	client      etcd.Client
	last        string
	slackClient *slack.Client
}

func New(slackUrl, u, hostname string, skipNotification bool) *Etcd {
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
		hostname: hostname,
		client:   c,
		slackClient: slack.New(&slack.Config{
			URL:      slackUrl,
			Username: "etcd plugin",
			Icon:     ":etcd_plugin:",
			Enabled:  !skipNotification,
		}),
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
		if health {
			c.slackClient.Resolve(c.hostname, msg)
		} else {
			c.slackClient.Trigger(c.hostname, msg)
		}
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
