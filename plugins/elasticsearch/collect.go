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

	if nodes := njson.Get("nodes").MustMap(); len(nodes) > 0 {
		for _, data := range nodes {
			j, _ := json.Marshal(data)
			node, _ := simplejson.NewJson(j)

			e.gauges["cpu"].Update(float64(node.GetPath("process", "cpu", "percent").MustInt()))
			e.gauges["memory"].Update(float64(node.GetPath("jvm", "mem", "heap_used_in_bytes").MustInt()))

			indexes := node.GetPath("indices", "indexing", "index_total").MustInt()
			gets := node.GetPath("indices", "get", "total").MustInt()
			searches := node.GetPath("indices", "search", "query_total").MustInt()

			e.gauges["indexes"].Update(float64(indexes - e.previousIndexes))
			e.gauges["gets"].Update(float64(gets - e.previousGets))
			e.gauges["searches"].Update(float64(searches - e.previousSearches))

			e.previousIndexes = indexes
			e.previousGets = gets
			e.previousSearches = searches

			break
		}
	}

	e.gauges["nodes"].Update(float64(cjson.GetPath("nodes", "count", "total").MustInt()))
	e.gauges["docs"].Update(float64(cjson.GetPath("indices", "docs", "count").MustInt()))

	status := cjson.Get("status").MustString()

	if status == "green" {
		e.gauges["status"].Update(float64(GREEN))
	} else if status == "yellow" {
		e.gauges["status"].Update(float64(YELLOW))
	} else {
		e.gauges["status"].Update(float64(RED))
	}
}
