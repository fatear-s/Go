package main

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func unMarshFile2(filename string, data interface{}) (err error) {
	f, err := os.Open(filename)
	if err != nil {
		return err
	} else {
		buf := make([]byte, 1024)
		var n int
		n, err = f.Read(buf)
		fmt.Println(buf)
		return unMarsl(buf[:n], data)
	}
}

// 反序列化
// []byte  ---- >  结构体
func unMarsl(input []byte, result interface{}) (err error) {
	// 先判断是否是指针
	typeInfo := reflect.TypeOf(result)
	if typeInfo.Kind() != reflect.Ptr {
		return
	}
	// 判断下一层是否是结构体
	if typeInfo.Elem().Kind() != reflect.Struct {
		return
	}
	// 转类型，按行切割
	lineArr := strings.Split(string(input), "\n")
	// 定义全局标签名   也就是server和mysql
	var myFiledName string

	for _, line := range lineArr {
		// 各种严谨判断
		line = strings.TrimSpace(line)
		// 处理文档中有注释的情况
		if len(line) == 0 || line[0] == '#' || line[0] == ';' {
			continue
		}
		// 按照括号去判断
		if line[0] == '[' {
			myFiledName, err = myLabel(line, typeInfo.Elem())
			if err != nil {
				return
			}
			continue
		}
		// 按照大标签去处理
		err = myField0(myFiledName, line, result)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	return

}
func myLabel(line string, typeInfo reflect.Type) (myFiledName string, err error) {
	//有没有server
	name := line[1 : len(line)-1]
	for i := 0; i < typeInfo.NumField(); i++ {
		field := typeInfo.Field(i)
		if field.Tag.Get("ini") == name {
			myFiledName = field.Name
			break
		}
	}
	return
	//有赋予
}
func myField(fieldName string, line string, result interface{}) (err error) {
	fmt.Println(line)
	key := strings.TrimSpace(line[0:strings.Index(line, "=")])
	val := strings.TrimSpace(line[strings.Index(line, "=")+1:])
	// 解析到结构体
	//resultType := reflect.TypeOf(result)
	resultValue := reflect.ValueOf(result)
	// 拿到字段值，这里直接设置不知道类型
	labelValue := resultValue.Elem().FieldByName(fieldName)
	// 拿到该字段类型
	fmt.Println(labelValue)
	labelType := labelValue.Type()
	// 第一次进来应该是server
	// 存放取到的字段名
	var keyName string
	// 遍历server结构体的所有字段
	for i := 0; i < labelType.NumField(); i++ {
		// 获取结构体字段
		field := labelType.Field(i)
		tagVal := field.Tag.Get("ini")
		if tagVal == key {
			keyName = field.Name
			break
		}
	}

	// 给字段赋值
	// 取字段值
	filedValue := labelValue.FieldByName(keyName)
	// 修改值
	switch filedValue.Type().Kind() {
	case reflect.String:
		filedValue.SetString(val)
	case reflect.Int:
		i, err2 := strconv.ParseInt(val, 10, 64)
		if err2 != nil {
			fmt.Println(err2)
			return
		}
		filedValue.SetInt(i)
	case reflect.Uint:
		i, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			fmt.Println(err)
			return err
		}
		filedValue.SetUint(i)
	case reflect.Float32:
		f, _ := strconv.ParseFloat(val, 64)
		//if err != nil {
		//	fmt.Println(err)
		//	return
		//}
		filedValue.SetFloat(f)
	}
	return
}
func myField0(myFiledName string, line string, result interface{}) (err error) {
	print(myFiledName)
	print(line)
	//取待匹配字符
	name := strings.Replace(strings.Split(line, "=")[0], " ", "", 10)
	value := strings.Split(line, "=")[1]
	//取数据结构
	valueInfo := reflect.ValueOf(result)
	ptr1 := valueInfo.Elem().FieldByName(myFiledName)
	fmt.Println(ptr1)
	typeInfo := ptr1.Type()
	var keyName string
	for i := 0; i < typeInfo.NumField(); i++ {
		filed := typeInfo.Field(i)
		if filed.Tag.Get("ini") == name {
			keyName = filed.Name
			break
		}
	}
	// 取字段值
	filedValue := ptr1.FieldByName(keyName)
	//赋
	switch filedValue.Type().Kind() {
	case reflect.String:
		filedValue.SetString(value)
	}
	return
}
func MarshalFile(filename string, data interface{}) (err error) {

	result, err := Marshal(data)
	if err != nil {
		fmt.Println(err)
		return err
	} else {
		fmt.Println("write in ")
		fmt.Println(result)
		return os.WriteFile(filename, result, 0666)
	}
}

func Marshal(data interface{}) (result []byte, err error) {
	// 获取下类型
	typeInfo := reflect.TypeOf(data)
	valueInfo := reflect.ValueOf(data)
	// 判断类型
	if typeInfo.Kind() != reflect.Struct {
		return
	}
	var conf []string

	// 获取所有字段去处理
	for i := 0; i < typeInfo.NumField(); i++ {
		// 取字段
		data_key := typeInfo.Field(i)
		fieldType := data_key.Type
		// 取值
		data_value := valueInfo.Field(i)
		// 判断字段类型
		if data_value.Kind() != reflect.Struct {
			return
		}
		// 拼的是[server]和[mysql]
		// 获取个tag
		tagValue := data_key.Tag.Get("ini")
		label := fmt.Sprintf("\n[%s]\n", tagValue)
		conf = append(conf, label)
		// 拼 k-v
		for j := 0; j < fieldType.NumField(); j++ {
			// 这里取到的是大写
			data_key := fieldType.Field(j)
			data_value := data_value.Field(j)

			// 取tag
			tagValue := data_key.Tag.Get("ini")
			label := fmt.Sprintf("\n%s=%s\n", tagValue, data_value.Interface())
			conf = append(conf, label)
			// 取值

			// Interface()取真正对应的值

		}
	}
	// 遍历切片转类型
	fmt.Println(conf)

	for i := 0; i < len(conf); i++ {
		out := []byte(conf[i])
		result = append(result, out...)
	}

	/*
		for _, val := range conf {
			out := []byte(val)
			result = append(result, out...)
		}
	*/

	fmt.Println(result)
	return
}
