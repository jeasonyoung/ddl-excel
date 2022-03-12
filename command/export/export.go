package export

import (
	"ddl-excel/library/mlog"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcmd"
	"github.com/gogf/gf/text/gstr"
)

func Help() {
	mlog.Print(gstr.TrimLeft(`
	USAGE
		ddl-excel export [OPTION]

	OPTION
		-l,--link		数据库连接字符串
		-h,--host		数据库服务器IP地址
		-P,--port		数据库服务器端口,默认为3306
		-u,--user		数据库用户名
		-p,--pass		数据库密码
		-n,--name		数据库名称

	DESCRIPTION
		ddl-excel export [OPTION]	将数据库表DDL导出生成Excel

	EXAMPLES
		ddl-excel export -l mysql:root:123456@tcp(127.0.0.1:3306)/test
		ddl-excel export -link="mysql:root:123456@tcp(127.0.0.1:3306)/test"
		ddl-excel export -h 127.0.0.1 -P 3306 -u root -p 123456 -n test
		ddl-excel export --host=127.0.0.1 --port=3306 --user=root --pass=123456 --name=test
`))
}

func Run() {
	parser, err := gcmd.Parse(g.MapStrBool{
		"l,link": true,
		"h,host": true,
		"P,port": true,
		"u,user": true,
		"p,pass": true,
		"n,name": true,
	})
	if err != nil {
		mlog.Fatal(err)
	}
	//数据库连接字符串
	DaoConn.SetLink(parser, "link")
	//数据库服务器地址
	DaoConn.setHost(parser, "host")
	//数据库服务器端口
	DaoConn.setPort(parser, "port")
	//数据库用户名
	DaoConn.setUser(parser, "user")
	//数据库密码
	DaoConn.setPass(parser, "pass")
	//数据库名称
	DaoConn.setName(parser, "name")
	//检查数据库配置
	if !DaoConn.check() {
		mlog.Fatal("数据库相关配置不完整!")
	}
	//开始生成excel
	mlog.Print("开始生成导出Excel...")
	if err := DaoConn.createExcel(); err != nil {
		mlog.Fatal(err)
	}
	mlog.Print("生成导出Excel成功.")
}
