package main

import (
	"bytes"
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"strings"
	"text/template"
)

var (
	Func = template.FuncMap{
		"lower":       strings.ToLower,
		"line2c":      lineToCamel,
		"type2":       GetType,
		"dealGormTag": dealGormTag,
	}
	headerTemplate = template.Must(template.New("header").Funcs(Func).Parse(`
	package model
	import (
	"{{.Package}}"
	_ "fmt"
	)

	type {{line2c .Name}} struct {
	{{range .Table}}{{line2c  .Field}} {{type2 .Type .Null}} {{ dealGormTag .Key .Field }}     //{{.Comment}}
	{{end}}
	}

	func ({{line2c .Name}}) TableName() string {
		return "{{lower .Name}}"
	}

	`))
	Name = "DefaultTable"
)

type TableDesc struct {
	Field   string
	Type    string
	Null    string
	Key     string
	Default sql.NullString
	Extra   string
	Comment string
}

type TableDescArray struct {
	Table   []TableDesc
	Name    string
	Package string
	Cmd     string
}

var dsn string
var TableName string

func init() {
	flag.StringVar(&dsn, "db", "root:123456@tcp(127.0.0.1:3306)/test", "db dsn")

	flag.StringVar(&TableName, "t", "trade", "table name")
}
func main() {

	flag.Parse()
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	TableInfo := TableDescArray{}

	TableInfo.Name = TableName
	if TableName == "" {
		fmt.Println("表名不能为空. -t ")
		return
	}
	TableInfo.Cmd = strings.Join(os.Args, " ")

	rows, err := db.Query("show full fields from " + TableName)
	if err != nil {
		fmt.Println("没有这张表.")
		return
	}

	for rows.Next() {
		var Field, Type, Null, Key, Extra, Privileges, Comment string
		var Collation, Default sql.NullString
		rows.Scan(&Field, &Type, &Collation, &Null, &Key, &Default, &Extra, &Privileges, &Comment)
		TableInfo.Table = append(TableInfo.Table, TableDesc{
			Field:   Field,
			Type:    Type,
			Null:    Null,
			Key:     Key,
			Default: Default,
			Extra:   Extra,
			Comment: Comment,
		})
		if Type == "datetime" {
			TableInfo.Package = `time`
		}

	}

	applyTemplate(TableInfo)
}

func applyTemplate(info TableDescArray) {

	w := bytes.NewBuffer(nil)

	err := headerTemplate.Execute(w, info)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(w.String())

}

var specilField = map[string]string{
	"id":   "",
	"no":   "",
	"url":  "",
	"uuid": "",
}

func lineToCamel(src string) string {
	result := ""
	for _, v := range strings.Split(src, "_") {
		tmp := strings.Title(v)
		if _, ok := specilField[strings.ToLower(v)]; ok {
			tmp = strings.ToUpper(v)
		}
		result += tmp
	}
	return result

}
func dealGormTag(key, field string) (dt string) {
	if key == "PRI" {
		dt = "primary_key;"
	}
	dt += "column:" + field
	dt = "`sql:\"" + dt + "\" json:\"" + field + ",omitempty\"`"
	return

}

func GetType(src, null string) string {
	null = "NO"
	if src == "decimal(10,2)" {
		return "money.Money2"
	}
	if index := strings.Index(src, "("); index > 0 {
		switch src[:index] {
		case "bigint", "tinyint", "int":
			if null == "NO" {
				return "int64"
			} else {
				return "*int64"
			}
		case "varchar":
			if null == "NO" {
				return "string"
			} else {
				return "*string"
			}
		case "float", "double":
			if null == "NO" {
				return "float64"
			} else {
				return "*float64"
			}
		default:
			if null == "NO" {
				return "string"
			} else {
				return "*string"
			}
		}

	} else {
		switch src {
		case "datetime", "timestamp":
			if null == "NO" {
				return "time.Time"
			} else {
				return "*time.Time"
			}
		case "text", "longtext":
			if null == "NO" {
				return "string"
			} else {
				return "*string"
			}
		default:
			if null == "NO" {
				return "string"
			} else {
				return "*string"
			}
		}
	}

}
