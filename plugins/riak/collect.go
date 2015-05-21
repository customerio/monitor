package riak

import (
	"io/ioutil"
	"net/http"

	"github.com/bitly/go-simplejson"
)

func grabInt(json *simplejson.Json, metric string) float64 {
	m, err := json.Get(metric).Int()
	if err != nil {
		panic(err)
	}
	return float64(m)
}

func (r *Riak) collect() {
	defer func() {
		if rr := recover(); rr != nil {
			r.clear()
		}
	}()

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

	r.memory.Update(grabInt(json, "memory_total"))
	r.gets.Update(grabInt(json, "vnode_gets"))
	r.puts.Update(grabInt(json, "vnode_puts"))
	r.index_gets.Update(grabInt(json, "vnode_index_reads"))
}
