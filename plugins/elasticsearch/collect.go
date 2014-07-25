package elasticsearch

import (
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

	node, err := http.Get("http://" + e.server + "/_stats")
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

	e.stats = map[string]int{
		"nodes":   cjson.GetPath("nodes", "count", "total").MustInt(),
		"cpu":     cjson.GetPath("nodes", "process", "cpu", "percent").MustInt(),
		"memory":  cjson.GetPath("nodes", "jvm", "mem", "heap_used_in_bytes").MustInt(),
		"docs":    njson.GetPath("_all", "total", "docs", "count").MustInt(),
		"indexes": njson.GetPath("_all", "total", "indexing", "index_current").MustInt(),
		"gets":    njson.GetPath("_all", "total", "get", "current").MustInt(),
	}

	status := cjson.Get("status").MustString()

	if status == "green" {
		e.stats["status"] = GREEN
	} else if status == "yellow" {
		e.stats["status"] = YELLOW
	} else {
		e.stats["status"] = RED
	}
}
