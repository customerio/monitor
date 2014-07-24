package riak

import (
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"net/http"
)

func grabInt(json *simplejson.Json, metric string) int {
	m, err := json.Get(metric).Int()
	if err != nil {
		panic(err)
	}
	return m
}

func (r *Riak) collect() {

	resp, _ := http.Get("http://" + r.server + "/stats")
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	json, err := simplejson.NewJson(body)
	if err != nil {
		panic(err)
	}

	r.memory = grabInt(json, "memory_total")
	r.gets = grabInt(json, "node_gets")
	r.puts = grabInt(json, "node_puts")
	r.index_gets = grabInt(json, "vnode_index_reads")
}
