package zookeeper

func (z *Zookeeper) collect() {

	for i, path := range z.paths {
		if children, _, err := z.conn.Children(path); err == nil {
			z.gauges[i].Update(float64(len(children)))
		} else {
			z.gauges[i].Update(0)
		}
	}
}
