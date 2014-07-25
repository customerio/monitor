package elasticsearch

import (
	"encoding/json"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"net/http"
)

func (e *Elasticsearch) collect() {
	cluster, err := http.Get("http://" + e.server + "/_cluster/stats")
	if err != nil {
		panic(err)
	}
	defer cluster.Body.Close()

	node, err := http.Get("http://" + e.server + "/_nodes/_local/stats")
	if err != nil {
		panic(err)
	}
	defer node.Body.Close()

	cbody, err := ioutil.ReadAll(cluster.Body)
	if err != nil {
		panic(err)
	}

	cjson, err := simplejson.NewJson(cbody)
	if err != nil {
		panic(err)
	}

	nbody, err := ioutil.ReadAll(node.Body)
	if err != nil {
		panic(err)
	}

	njson, err := simplejson.NewJson(nbody)
	if err != nil {
		panic(err)
	}

	e.stats = make(map[string]int)

	if nodes := njson.Get("nodes").MustMap(); len(nodes) > 0 {
		for _, data := range nodes {
			j, _ := json.Marshal(data)
			node, _ := simplejson.NewJson(j)

			e.stats["cpu"] = node.GetPath("process", "cpu", "percent").MustInt()
			e.stats["memory"] = node.GetPath("jvm", "mem", "heap_used_in_bytes").MustInt()
			e.stats["indexes"] = node.GetPath("indices", "indexing", "index_current").MustInt()
			e.stats["gets"] = node.GetPath("indices", "get", "current").MustInt()
			e.stats["searches"] = node.GetPath("indices", "search", "query_current").MustInt()
			break
		}
	}

	e.stats["nodes"] = cjson.GetPath("nodes", "count", "total").MustInt()
	e.stats["docs"] = cjson.GetPath("indices", "docs", "count").MustInt()

	status := cjson.Get("status").MustString()

	if status == "green" {
		e.stats["status"] = GREEN
	} else if status == "yellow" {
		e.stats["status"] = YELLOW
	} else {
		e.stats["status"] = RED
	}
}
