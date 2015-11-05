package riak

import (
	"io/ioutil"
	"net/http"

	"github.com/bitly/go-simplejson"
	"github.com/customerio/monitor/metrics"
	"github.com/customerio/monitor/plugins"
)

func grabInt(json *simplejson.Json, metric string) float64 {
	m, err := json.Get(metric).Int()
	if err != nil {
		panic(err)
	}
	return float64(m)
}

func (f *Riak) Collect(b *metrics.Batch) {
	f.collect()

	for _, u := range f.updaters {
		u.Fill(b)
	}
}

func (r *Riak) collect() {
	defer func() {
		if rr := recover(); rr != nil {
			plugins.Logger.Printf("panic: Riak: %v\n", rr)
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

	r.updaters[memoryGauge].Update(grabInt(json, "memory_total"))
	r.updaters[getsGauge].Update(grabInt(json, "vnode_gets"))
	r.updaters[putsGauge].Update(grabInt(json, "vnode_puts"))
	r.updaters[indexGetsGauge].Update(grabInt(json, "vnode_index_reads"))
}
