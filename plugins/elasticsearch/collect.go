package elasticsearch

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/bitly/go-simplejson"
	"github.com/customerio/monitor/plugins"
)

func (e *Elasticsearch) Collect() {
	defer func() {
		if r := recover(); r != nil {
			plugins.Logger.Printf("panic: Elasticsearch: %v\n", r)
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

			e.updaters[cpuGauge].Update(float64(node.GetPath("process", "cpu", "percent").MustInt()))
			e.updaters[memoryGauge].Update(float64(node.GetPath("jvm", "mem", "heap_used_in_bytes").MustInt()))

			indexes := node.GetPath("indices", "indexing", "index_total").MustInt()
			gets := node.GetPath("indices", "get", "total").MustInt()
			searches := node.GetPath("indices", "search", "query_total").MustInt()

			e.updaters[indexesGauge].Update(float64(indexes - e.previousIndexes))
			e.updaters[getsGauge].Update(float64(gets - e.previousGets))
			e.updaters[searchesGauge].Update(float64(searches - e.previousSearches))

			e.previousIndexes = indexes
			e.previousGets = gets
			e.previousSearches = searches

			break
		}
	}

	e.updaters[nodesGauge].Update(float64(cjson.GetPath("nodes", "count", "total").MustInt()))
	e.updaters[docsGauge].Update(float64(cjson.GetPath("indices", "docs", "count").MustInt()))

	status := cjson.Get("status").MustString()

	if status == "green" {
		e.updaters[statusGauge].Update(float64(statusGreen))
	} else if status == "yellow" {
		e.updaters[statusGauge].Update(float64(statusYellow))
	} else {
		e.updaters[statusGauge].Update(float64(statusRed))
	}
}
