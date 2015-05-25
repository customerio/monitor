### WIP

Example:

    package main
    
    import (
    	"github.com/customerio/monitor/plugins/cpu"
	"github.com/customerio/monitor/plugins"
    
    	"log"
    	"os"
    	"time"
    )
    
    func main() {
	c := cpu.New()
	plugins.AddCollector(c)

	plugins.Collect("localhost", "", "", time.Second)
    }

