package config

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// DumpConfig try to dump config in proper formatted text
func DumpConfig(w io.Writer) {
	lock.Lock()
	defer lock.Unlock()

	tab := tabwriter.NewWriter(w, 0, 8, 0, '\t', 0)
	fmt.Fprint(w, "Key\tDescription\tField\tValue")
	for key := range configs {
		d, ok := o.Get(key)
		fmt.Fprintf(tab, "%s\t%s\t%b\t%v\n", key, configs[key], ok, d)
	}
	tab.Flush()
}
