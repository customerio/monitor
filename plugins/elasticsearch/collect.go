package elasticsearch

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/bitly/go-simplejson"
)

func (e *Elasticsearch) collect() {
	defer func() {
		if r := recover(); r != nil {
			e.clear()
		}
	}()

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

			indexes := node.GetPath("indices", "indexing", "index_total").MustInt()
			gets := node.GetPath("indices", "get", "total").MustInt()
			searches := node.GetPath("indices", "search", "query_total").MustInt()

			e.stats["indexes"] = indexes - e.previousIndexes
			e.stats["gets"] = gets - e.previousGets
			e.stats["searches"] = searches - e.previousSearches

			e.previousIndexes = indexes
			e.previousGets = gets
			e.previousSearches = searches

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
