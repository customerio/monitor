package zookeeper

func (z *Zookeeper) collect() {

	for _, path := range z.paths {
		if children, _, err := z.conn.Children(path); err == nil {
			z.stats[path] = len(children)
		}
	}
}
