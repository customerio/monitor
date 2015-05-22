package riak

import (
	"fmt"
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
			fmt.Printf("panic: Riak: %v\n", rr)
			r.clear()
		}
	}()

	resp, err := http.Get("http://" + r.server + "/stats")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	json, err := simplejson.NewJson(body)
	if err != nil {
		panic(err)
	}

	r.gauges[memoryGauge].Update(grabInt(json, "memory_total"))
	r.gauges[getsGauge].Update(grabInt(json, "vnode_gets"))
	r.gauges[putsGauge].Update(grabInt(json, "vnode_puts"))
	r.gauges[indexGetsGauge].Update(grabInt(json, "vnode_index_reads"))
}
