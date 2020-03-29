package main

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
)

// MySQL 结构体
type MySQL struct {
	Addr string `ini:"addr"`
	Port int64  `ini:"port"`
	Name string `ini:"name"`
	Pwd  string `ini:"pwd"`
}

// Redis 结构体
type Redis struct {
	Username string `ini:"username"`
	Password string `ini:"password"`
}

// Config 配置文件结构体
type Config struct {
	MySQL `ini:"mysql"`
	Redis `ini:"redis"`
}

func loadConfig(fileName string, data interface{}) (err error) {
	/*
		1、判断传入的data是否是指针类型，只有指针类型才进行下面步骤
		2、判断传入的指针类型参数是否为结构体指针
		3、读取文件得到字节类型数据
		4、逐行读取数据
			4.1、如果有注释，则跳过处理
			4.2、如果是[]包含的则表示节点
			4.3、如果不是上述的则对其进行按'='切割
			4.4、对切割的数据进行存储到结构体
	*/
	t := reflect.TypeOf(data)
	if t.Kind() != reflect.Ptr {
		err = fmt.Errorf("data param must be a pointer")
		return
	}
	if t.Elem().Kind() != reflect.Struct {
		err = fmt.Errorf("data param must be a struct pointer")
		return
	}
	// 读取文本文件
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return
	}
	// 对读取出来的文本文件进行切片
	fileSlice := strings.Split(string(file), "\r\n")

	// 记录反射获取的结构体名
	var structName string
	// var valueType
	// 遍历切片
	for idx, line := range fileSlice {
		// 对line进行格式化
		line = strings.TrimSpace(line)
		// 如果是空行则跳过
		if len(line) == 0 {
			continue
		}
		// 如果为注释则跳过
		if strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
			continue
		}
		// 如果以[开头再进行进一步处理
		if strings.HasPrefix(line, "[") {
			// 如果结尾没有 ] ，则抛出错误
			if !strings.HasSuffix(line, "]") {
				err = fmt.Errorf("synx error in line %d", idx+1)
				return
			}
			// 如果以[开头，以]结尾，但是没有内容，则抛出错误
			selectName := strings.TrimSpace(line[1 : len(line)-1])
			if len(selectName) == 0 {
				err = fmt.Errorf("synx error in line %d,the content con't be empty", idx+1)
				return
			}
			// 如果上述条件都满足，则根据节点名称找到对应结构体
			// 读取data的value
			for i := 0; i < t.Elem().NumField(); i++ {
				field := t.Elem().Field(i)
				// 获取tag，反射获取其内容
				// fmt.Println(field.Tag.Get("ini"))
				// fmt.Println(selectName)
				if field.Tag.Get("ini") == selectName {
					// 说明找到了对应的结构体，然后讲结构体记录下来
					structName = field.Name
					fmt.Printf("find %s struct from %s in data\n", structName, selectName)
					break
				} else {
					structName = ""
				}
			}
		} else {
			// 上面通过反射找到structName，判断其是不是struct类型，如果不是抛出错误
			if structName != "" {
				v := reflect.ValueOf(data)
				sValue := v.Elem().FieldByName(structName)
				sType := sValue.Type()
				if sType.Kind() != reflect.Struct {
					err = fmt.Errorf("%s in data must be a struct type", structName)
					return
				}
				// 判断文本当前行是否有等号，并且等号左边是key，右边是value
				if strings.Contains(line, "=") {
					// 获取=所在的索引
					index := strings.Index(line, "=")
					key := strings.TrimSpace(line[:index])
					value := strings.TrimSpace(line[index+1:])
					if len(key) == 0 {
						err = fmt.Errorf("synx error. line: %d", idx)
						return
					}
					// fmt.Printf("key: %s value: %s\n", key, value)
					// 通过反射找到key对应的结构体中的变量名
					var fieldName string
					var filedType reflect.StructField
					for i := 0; i < sValue.NumField(); i++ {
						// fmt.Println(sValue.Field(i))
						field := sType.Field(i)
						filedType = field
						if field.Tag.Get("ini") == key {
							fieldName = field.Name
							break
						}
					}
					fileObj := sValue.FieldByName(fieldName)
					// fmt.Print(fieldName, filedType.Type.Kind())
					// 根据不同的类型进行写入操作
					switch filedType.Type.Kind() {
					case reflect.String:
						fileObj.SetString(value)
					case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8:
						// 对value进行转换，由string-->int
						var valueInt int64
						valueInt, err = strconv.ParseInt(value, 10, 64)
						if err != nil {
							err = fmt.Errorf("line: %d synx error", idx)
							return
						}
						fileObj.SetInt(valueInt)
					case reflect.Float32, reflect.Float64:
						// 对value进行转换，由string-->float
						var valueFloat float64
						valueFloat, err = strconv.ParseFloat(value, 64)
						if err != nil {
							err = fmt.Errorf("line: %d synx error", idx)
							return
						}
						fileObj.SetFloat(valueFloat)
					case reflect.Bool:
						var valueBool bool
						valueBool, err = strconv.ParseBool(value)
						if err != nil {
							err = fmt.Errorf("line: %d synx error", idx)
							return
						}
						fileObj.SetBool(valueBool)
					default:
						err = fmt.Errorf("line: %d value type error", idx)
						return
					}

				} else {
					err = fmt.Errorf("synx error. line: %d", idx)
					return
				}
			}
		}
	}
	return
}

func main() {
	// 声明结构体
	var cfg Config
	// fmt.Printf("%T\n", &cfg)
	err := loadConfig("./my.ini", &cfg)
	if err != nil {
		fmt.Println("conf load failed,err:", err)
	}
	fmt.Println(cfg)
}
