package zookeeper

import (
	"github.com/samuel/go-zookeeper/zk"
	"testing"
)

type logWriter struct {
	t *testing.T
	p string
}

func (lw logWriter) Write(b []byte) (int, error) {
	lw.t.Logf("%s%s", lw.p, string(b))
	return len(b), nil
}

func TestCollect(t *testing.T) {

	ts, err := zk.StartTestCluster(3, nil, logWriter{t: t, p: "[ZKERR] "})
	if err != nil {
		t.Fatal(err)
	}
	defer ts.Stop()

	z, err := ts.ConnectAll()
	if err != nil {
		t.Fatalf("Connect returned error: %+v", err)
	}
	defer z.Close()

	z.Create("/test", []byte("base"), 0, zk.WorldACL(zk.PermAll))
	z.Create("/test/path", []byte("sub"), 0, zk.WorldACL(zk.PermAll))

	testKeys := []string{"a", "b", "c"}

	for _, v := range testKeys {
		_, err = z.Create("/test/path/"+v, []byte("child"), 0, zk.WorldACL(zk.PermAll))
	}

	zkTest := Zookeeper{conn: z, paths: []string{"/test/path"}, stats: make(map[string]int)}

	zkTest.collect()

	if zkTest.stats["/test/path"] != len(testKeys) {
		t.Errorf("Expected %v children, got %v", len(testKeys), zkTest.stats["/test/path"])
	}
}
