package main

import (
	"ddl-excel/command/export"
	"ddl-excel/library/allyes"
	"ddl-excel/library/mlog"
	"fmt"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/os/gbuild"
	"github.com/gogf/gf/os/gcmd"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/text/gregex"
	"github.com/gogf/gf/text/gstr"
)

const VERSION = "v1.0.0"

var helpContent = gstr.TrimLeft(`
	USAGE:
		ddl-excel COMMAND [ARGUMENT] [OPTION]

	COMMAND
		export 导出数据库DDL到EXCEL
		help 显示命令帮助信息
		version 显示当前版本信息
	
	OPTION
		-y		all yes for all command without prompt ask
		-?,-h	show this help or detail for specified command
		-v,-i	show version information
		-debug	show internal detailed debugging information

	ADDITIONAL
		use 'ddl-excel help COMMAND' or 'ddl-excel COMMAND -h' for detail about a command,which has '...'
		in the tail of their comments.
`)

func main() {
	defer func() {
		if exp := recover(); exp != nil {
			if err, ok := exp.(error); ok {
				mlog.Print(gerror.Current(err).Error())
			} else {
				panic(exp)
			}
		}
	}()
	//
	allyes.Init()
	command := gcmd.GetArg(1)
	// help information
	if command == "" {
		help(command)
		return
	}
	switch command {
	case "help":
		help(gcmd.GetArg(2))
	case "export":
		export.Run()
	default:
		for k := range gcmd.GetOptAll() {
			switch k {
			case "?", "h":
				mlog.Print(helpContent)
				return

			case "i", "v":
				version()
				return
			}
		}
		mlog.Printf(helpContent)
	}
}

// help shows more information for specified command.
func help(command string) {
	switch command {
	case "export":
		export.Help()
	default:
		mlog.Print(helpContent)
	}
}

func version() {
	info := gbuild.Info()
	if info["git"] == "" {
		info["git"] = "none"
	}
	mlog.Printf(`ddl-excel CLI Tool %s,https://github.com/jeasonyoung/ddl-excel`, VERSION)
	gfVersion, err := getGFVersionOfCurrentProject()
	if err != nil {
		gfVersion = err.Error()
	} else {
		gfVersion = gfVersion + " in current go.mod"
	}
	mlog.Printf(`GoFrame Version: %s`, gfVersion)
	mlog.Printf(`CLI Installed At: %s`, gfile.SelfPath())
	if info["gf"] == "" {
		mlog.Print(`Current is a custom installed version, no installation information.`)
		return
	}

	mlog.Print(gstr.Trim(fmt.Sprintf(`
CLI Built Detail:
  Go Version:  %s
  GF Version:  %s
  Git Commit:  %s
  Build Time:  %s
`, info["go"], info["gf"], info["git"], info["time"])))
}

// getGFVersionOfCurrentProject checks and returns the GoFrame version current project using.
func getGFVersionOfCurrentProject() (string, error) {
	goModPath := gfile.Join(gfile.Pwd(), "go.mod")
	if gfile.Exists(goModPath) {
		match, err := gregex.MatchString(`github.com/gogf/gf\s+([\w\d\.]+)`, gfile.GetContents(goModPath))
		if err != nil {
			return "", err
		}
		if len(match) > 1 {
			return match[1], nil
		}
		return "", gerror.New("cannot find goframe requirement in go.mod")
	} else {
		return "", gerror.New("cannot find go.mod")
	}
}
