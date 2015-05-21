package zookeeper

func (z *Zookeeper) collect() {

	for _, path := range z.paths {
		if children, _, err := z.conn.Children(path); err == nil {
			z.gauges[path].Update(float64(len(children)))
		} else {
			z.gauges[path].Update(0)
		}
	}
}
