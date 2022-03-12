package allyes

import (
	"ddl-excel/library/mlog"
	"github.com/gogf/gf/os/gcmd"
	"github.com/gogf/gf/os/genv"
)

const EnvName = "GF_CLI_ALL_YES"

// Init initializes the package manually.
func Init() {
	if gcmd.ContainsOpt("y") {
		if err := genv.Set(EnvName, "1"); err != nil {
			mlog.Print(err)
		}
	}
}

// Check checks whether option allow all yes for command.
func Check() bool {
	return genv.Get(EnvName) == "1"
}
