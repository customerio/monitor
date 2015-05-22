package zookeeper

import "fmt"

func (z *Zookeeper) collect() {

	for i, path := range z.paths {
		if children, _, err := z.conn.Children(path); err == nil {
			z.gauges[i].Update(float64(len(children)))
		} else {
			fmt.Printf("panic: Zookeeper: %v\n", err)
			z.gauges[i].Update(0)
		}
	}
}
