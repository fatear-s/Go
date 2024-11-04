package main

import "fmt"

// 序列化
func paseFile1(filename string) {
	var conf Config
	conf.ServerConf.Ip = "127.0.0.1"
	conf.ServerConf.Port = 80
	conf.MysqlConf.Host = "127.0.0.1"
	conf.MysqlConf.Port = 80
	conf.MysqlConf.Database = "MySql"
	conf.MysqlConf.Passwd = "123456"
	conf.MysqlConf.Timeout = 10
	conf.MysqlConf.Username = "tau"
	err := MarshalFile(filename, conf) // make code
	if err != nil {
		fmt.Println(err)
		return
	}

}
func paseFile2(filename string) {
	var conf Config
	err := unMarshFile2(filename, &conf)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("反序列化成功  conf: %#v\n  port: %#v\n", conf, conf.ServerConf.Port)
}

//反序列化

func main() {
	//paseFile1("./my.ini")
	paseFile2("./config.ini")
}
