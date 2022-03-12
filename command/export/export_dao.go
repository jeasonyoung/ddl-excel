package export

import (
	"ddl-excel/library/mlog"
	"fmt"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/os/gcmd"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/text/gregex"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
)

//DaoConn 数据库连接对象
var DaoConn = &daoConn{
	db:   nil,
	port: 3306,
}

// 数据库连接
type daoConn struct {
	db gdb.DB

	link string //数据库连接字符串
	host string //数据库服务器地址
	port uint   //数据库服务器
	user string //数据库用户名
	pass string //数据库密码
	name string //数据库名称
}

func setOptionValueHandler(parser *gcmd.Parser, name string, valFn func(val string)) {
	if parser != nil && name != "" && valFn != nil {
		value := parser.GetOpt(name)
		if value != "" {
			valFn(value)
		}
	}
}

func (dc *daoConn) SetLink(parser *gcmd.Parser, name string) {
	setOptionValueHandler(parser, name, func(val string) {
		dc.link = val
	})
}

func (dc *daoConn) setHost(parser *gcmd.Parser, name string) {
	setOptionValueHandler(parser, name, func(val string) {
		dc.host = val
	})
}

func (dc *daoConn) setPort(parser *gcmd.Parser, name string) {
	setOptionValueHandler(parser, name, func(val string) {
		dc.port = gconv.Uint(val)
	})
}

func (dc *daoConn) setUser(parser *gcmd.Parser, name string) {
	setOptionValueHandler(parser, name, func(val string) {
		dc.user = val
	})
}

func (dc *daoConn) setPass(parser *gcmd.Parser, name string) {
	setOptionValueHandler(parser, name, func(val string) {
		dc.pass = val
	})
}

func (dc *daoConn) setName(parser *gcmd.Parser, name string) {
	setOptionValueHandler(parser, name, func(val string) {
		dc.name = val
	})
}

func (dc *daoConn) buildLink() {
	if dc.host != "" && dc.user != "" && dc.pass != "" && dc.name != "" {
		if dc.port <= 0 {
			dc.port = 3306
		}
		dc.link = fmt.Sprintf("mysql:%s:%s@tcp(%s:%d)/%s", dc.user, dc.pass, dc.host, dc.port, dc.name)
	}
}

func (dc *daoConn) check() bool {
	dc.buildLink()
	link := dc.link
	if link != "" {
		tempGroup := gtime.TimestampNanoStr()
		match, _ := gregex.MatchString(`([a-z]+):(.+)`, link)
		if len(match) == 3 {
			gdb.AddConfigNode(tempGroup, gdb.ConfigNode{
				Type: gstr.Trim(match[1]),
				Link: gstr.Trim(match[2]),
			})
			if db, err := gdb.Instance(tempGroup); err != nil {
				mlog.Fatal(err)
			} else {
				dc.db = db
			}
		}
		if dc.db == nil {
			mlog.Fatal("数据库初始化失败")
		}
		if dc.name == "" {
			if val, err := dc.db.GetValue("select database()"); err != nil {
				mlog.Fatal(err)
			} else {
				dc.name = val.String()
			}
		}
		return true
	}
	return false
}
