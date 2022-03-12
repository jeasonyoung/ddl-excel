package export

import (
	"ddl-excel/library/mlog"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/xuri/excelize/v2"
	"strings"
)

const (
	sqlTables  = "select table_name,table_comment from information_schema.tables where table_schema=(select database());"
	sqlColumns = `select column_name,column_comment,column_type,
					   if(is_nullable = 'no' and column_key != 'PRI', 1, 0) as is_required,
					   if(column_key = 'PRI', 1, 0) as is_pk,
					   if(extra = 'auto_increment', 1, 0) as is_increment
				from information_schema.columns
				where table_schema = (select database()) and table_name = ?
				order by ordinal_position`
)

//获取表集合
func (dc *daoConn) getTables() (map[string]string, error) {
	if rows, err := dc.db.GetAll(sqlTables); err != nil {
		return nil, err
	} else {
		tables := map[string]string{}
		for _, row := range rows {
			tableName := row["table_name"].String()
			tableComment := row["table_comment"].String()
			if tableName != "" {
				tables[tableName] = tableComment
			}
		}
		return tables, nil
	}
}

func (dc *daoConn) createExcel() error {
	mlog.Print("开始获取表信息...")
	//获取表集合
	tables, err := dc.getTables()
	if err != nil {
		mlog.Fatal(err)
	}
	//创建Excel对象
	mlog.Print("开始创建Excel")
	f := excelize.NewFile()
	//sheet
	for name, comment := range tables {
		mlog.Printf("开始生成表数据:[%s]%s...", name, comment)
		if comment != "" {
			f.NewSheet(comment)
			dc.createSheet(f, comment, name)
			continue
		}
		f.NewSheet(name)
		dc.createSheet(f, name, name)
	}
	//删除Sheet1
	f.DeleteSheet("Sheet1")
	mlog.Print("开始保存Excel.")
	//保存Excel
	return f.SaveAs(fmt.Sprintf("%s.xlsx", dc.name))
}

func (dc *daoConn) createSheet(f *excelize.File, sheetName, tableName string) {
	setCellValue(f, sheetName, "A1", "字段名")
	setCellValue(f, sheetName, "B1", "中文名称")
	setCellValue(f, sheetName, "C1", "类型")
	setCellValue(f, sheetName, "D1", "说明")
	//查询字段集合
	items, err := dc.db.GetAll(sqlColumns, tableName)
	if err != nil {
		mlog.Fatal(err)
	}
	if len(items) > 0 {
		//写入数据
		for idx, item := range items {
			name := item["column_name"].String()
			setCellValue(f, sheetName, fmt.Sprintf("A%d", idx+2), name)
			comment := item["column_comment"].String()
			setCellValue(f, sheetName, fmt.Sprintf("B%d", idx+2), comment)
			colType := item["column_type"].String()
			setCellValue(f, sheetName, fmt.Sprintf("C%d", idx+2), colType)
			//
			remarks := g.ArrayStr{}
			isPk := item["is_pk"].Uint()
			if isPk > 0 {
				remarks = append(remarks, "主键")
			}
			isIncrement := item["is_increment"].Uint()
			if isIncrement > 0 {
				remarks = append(remarks, "自增")
			}
			isRequired := item["is_required"].Uint()
			if isRequired > 0 {
				remarks = append(remarks, "唯一")
			}
			setCellValue(f, sheetName, fmt.Sprintf("D%d", idx+2), strings.Join(remarks, ","))
		}
	}
}

func setCellValue(f *excelize.File, sheetName, axis, cellVal string) {
	if err := f.SetCellStr(sheetName, axis, cellVal); err != nil {
		mlog.Printf("sheet: %s[%s = %s]=> %s", sheetName, axis, cellVal, err.Error())
	}
}
