package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	viper.AddConfigPath("../config/")
	viper.SetConfigName("app")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Read config file fail: %s", err)
	}

	defaultConfig := viper.AllSettings()

	envMap := map[string]string{
		"debug":   "dev",
		"release": "prod",
	}

	viper.SetConfigName(fmt.Sprintf("app.%s", envMap[os.Getenv("GIN_MODE")]))
	if err := viper.ReadInConfig(); err == nil {
		newConfig := viper.AllSettings()

		// Traverse up to 3 levels
		for key, value := range newConfig {
			if strings.HasPrefix(fmt.Sprintf("%T", value), "map") {
				for subKey, subValue := range value.(map[string]interface{}) {
					if defaultConfig[key] == nil {
						defaultConfig[key] = map[string]interface{}{}
					}

					if strings.HasPrefix(fmt.Sprintf("%T", subValue), "map") {
						for childSubKey, childSubValue := range subValue.(map[string]interface{}) {
							if defaultConfig[key].(map[string]interface{})[subKey] == nil {
								defaultConfig[key].(map[string]interface{})[subKey] = map[string]interface{}{}
							}
							defaultConfig[key].(map[string]interface{})[subKey].(map[string]interface{})[childSubKey] = childSubValue
						}
					} else {
						defaultConfig[key].(map[string]interface{})[subKey] = subValue
					}
				}
			} else {
				defaultConfig[key] = value
			}
		}

		// Reset viper configuration
		if configData, err := json.Marshal(defaultConfig); err == nil {
			viper.ReadConfig(bytes.NewBuffer(configData))
		}
	}

	argsLen := len(os.Args)

	if argsLen > 0 {
		switch os.Args[1] {
		case "gm":
			if argsLen > 1 {
				generateModel(os.Args[2])
			}
		}
	}
}

// generate model: domain_modelname_type
func generateModel(table string) {
	strArrs := strings.Split(table, "_")
	arrLen := len(strArrs)
	if arrLen < 2 {
		fmt.Println("table name " + table + " not valid for domain_modelname_type")
		return
	}
	domain := strArrs[0]
	typeName := strArrs[arrLen-1]
	modelName := strings.Join(strArrs[1:len(strArrs)-1], "_")

	var tableStruct []struct {
		Field   string
		Type    string
		Null    string
		Key     string
		Default string
		Extra   string
	}
	db().Raw("show columns from " + table).Find(&tableStruct)

	if len(tableStruct) > 0 {
		tableNameCaptal := ""

		for _, str := range strArrs {
			tableNameCaptal += strings.ToUpper(str[0:1]) + str[1:]
		}

		modelText := []string{
			fmt.Sprintf("// %s\ntype %s struct {\n  gorm.Model", tableNameCaptal, tableNameCaptal),
		}

		// comment
		var commentList []struct {
			COLUMN_NAME    string
			COLUMN_COMMENT string
		}

		db().Raw("select COLUMN_NAME, COLUMN_COMMENT from information_schema.columns where table_schema = ? and table_name = ?", viper.GetString("db.database"), table).Find(&commentList)

		// count chars
		chars := []int{0, 1, 2}
		fieldLines := [][]string{}
		includeDataTypes := false

		for _, row := range tableStruct {
			fieldNameStrArr := strings.Split(row.Field, "_")
			fieldNameCaptal := ""
			for _, fieldStr := range fieldNameStrArr {
				fieldNameCaptal += strings.ToUpper(fieldStr[0:1]) + fieldStr[1:]
			}

			fieldChars := len(fieldNameCaptal)
			if chars[0] < fieldChars {
				chars[0] = fieldChars
			}

			fieldType := "string"
			switch row.Type {
			case "int unsigned":
				fieldType = "uint"
			case "int":
				fieldType = "int"
			case "tinyint unsigned":
				fieldType = "uint8"
			case "tinyint":
				fieldType = "int8"
			case "bigint unsigned":
				fieldType = "uint64"
			case "bigint":
				fieldType = "int64"
			case "datetime":
				fieldType = "time.Time"
			case "json":
				fieldType = "datatypes.JSON"
				includeDataTypes = true
			default:
				if strings.HasPrefix(fieldType, "decimal") {
					fieldType = "float64"
				}
			}

			typeChars := len(fieldType)
			if chars[1] < typeChars {
				chars[1] = typeChars
			}

			allowNull := ""
			if row.Null == "NO" {
				allowNull = ";not null"
			}

			defaultValue := ""
			if row.Default != "" {
				if fieldType == "string" {
					defaultValue = fmt.Sprintf(`;default:'%s'`, row.Default)
				} else {
					defaultValue = fmt.Sprintf(`;default:%s`, row.Default)
				}
			}

			comment := ""
			for _, cmt := range commentList {
				if cmt.COLUMN_NAME == row.Field && cmt.COLUMN_COMMENT != "" {
					comment = " // " + cmt.COLUMN_COMMENT
				}
			}

			column := fmt.Sprintf(`gorm:"column:%s;type:%s%s%s" json:"%s"`, row.Field, row.Type, allowNull, defaultValue, strings.ToLower(fieldNameCaptal[0:1])+fieldNameCaptal[1:])

			columnChars := len(column)
			if chars[2] < columnChars {
				chars[2] = columnChars
			}

			fieldLines = append(fieldLines, []string{fieldNameCaptal, fieldType, column, comment})
		}

		for _, fieldLine := range fieldLines {
			commentText := fieldLine[3]
			if commentText != "" {
				commentText = strings.Repeat(" ", chars[2]-len(fieldLine[2])) + commentText
			}
			modelText = append(modelText, fmt.Sprintf("  %s %s `%s`%s", fieldLine[0]+strings.Repeat(" ", chars[0]-len(fieldLine[0])), fieldLine[1]+strings.Repeat(" ", chars[1]-len(fieldLine[1])), fieldLine[2], commentText))
		}

		modelText = append(modelText, "}\n")
		modelText = append(modelText, fmt.Sprintf("// Table name\nfunc (m *%s) TableName() string {", tableNameCaptal))
		modelText = append(modelText, fmt.Sprintf(`  return "%s"`, table))
		modelText = append(modelText, "}\n")

		includeText := ""
		if includeDataTypes {
			includeText = `"gorm.io/datatypes"\n  `
		}

		text := fmt.Sprintf("package %sDomain\n\nimport (\n  \"time\"\n\n  %s\"gorm.io/gorm\"\n)\n\n", strings.ToUpper(domain[0:1])+domain[1:], includeText) + strings.Join(modelText, "\n")
		path := fmt.Sprintf("../app/domain/%s/%s", domain, typeName)
		err := os.MkdirAll(path, 0644)
		if err != nil {
			fmt.Print(err)
			return
		}
		err = os.WriteFile(path+"/"+modelName+".go", []byte(text), 0644)
		if err != nil {
			fmt.Print(err)
		}
	}
}

func db() *gorm.DB {
	// 从配置文件中获取参数
	host := viper.GetString("db.host")
	port := viper.GetString("db.port")
	database := viper.GetString("db.database")
	username := viper.GetString("db.username")
	password := viper.GetString("db.password")
	charset := viper.GetString("db.charset")
	loc := viper.GetString("db.loc")
	// fmt.Printf("prefix: %v\n", prefix)
	// 字符串拼接
	sqlStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=%s",
		username,
		password,
		host,
		port,
		database,
		charset,
		url.QueryEscape(loc),
	)
	handle, err := gorm.Open(mysql.Open(sqlStr))
	if err != nil {
		fmt.Println("打开数据库失败", err)
		panic("打开数据库失败" + err.Error())
	}

	return handle
}
