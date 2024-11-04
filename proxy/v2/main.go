package main

import (
	"bufio"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

const socketv5 = 0x05
const cmdBind = 0x01
const atypeIPV4 = 0x01
const atypeHOST = 0x03
const atypeIPV6 = 0x04

func main() {
	b, err := net.Listen("tcp", "127.0.0.1:1080")
	if err != nil {
		log.Printf("监听失败:%w", err)
		fmt.Print("err1")
		panic(err)
	}
	//单协程只能够代理一个,我的理解
	for {
		Conn, err := b.Accept()
		if err != nil {
			fmt.Print("err2")
			log.Printf("连接失败：%w", err)
			panic(err)
		}
		go process(Conn)
	}
}
func process(Conn net.Conn) {
	defer Conn.Close()
	//解析收到的socket5
	render := bufio.NewReader(Conn)
	err := auto(render, Conn)
	if err != nil {
		fmt.Print("err3")
		log.Printf("握手失败:%w", err)
		return
	}
	//连接服务器,包含转发给端口
	err = connect(render, Conn)
	if err != nil {
		fmt.Print("err4")
		log.Printf("连接服务器失败:%w", err)
		return
	}
}
func auto(render *bufio.Reader, Conn net.Conn) (err error) {
	// |VER | NMETHODS | METHODS  |
	// +----+----------+----------+
	// | 1  |    1     | 1 to 255 |
	// +----+----------+----------+
	// VER: 协议版本，socks5为0x05
	// NMETHODS: 支持认证的方法数量
	// METHODS: 对应NMETHODS，NMETHODS的值为多少，METHODS就有多少个字节。RFC预定义了一些值的含义，内容如下:
	// X’00’ NO AUTHENTICATION REQUIRED
	// X’02’ USERNAME/PASSWORD
	buf, err := render.ReadByte()
	if err != nil {
		log.Printf("读取失败:%w", err)
		return err
	}
	if buf != byte(socketv5) {
		log.Printf("不是socketv5%w", err)
		return err
	}
	_, err = render.ReadByte()
	if err != nil {
		log.Printf("读取失败:%w", err)
		return err
	}

	// +----+--------+
	// |VER | METHOD |
	// +----+--------+
	// | 1  |   0    |
	// +----+--------+
	_, err = Conn.Write([]byte{socketv5, 0x00})
	if err != nil {
		log.Printf("返回失败:%w", err)
		return err
	}
	return nil

}
func connect(render *bufio.Reader, Conn net.Conn) (err error) {
	// +----+-----+-------+------+----------+----------+
	// |VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
	// +----+-----+-------+------+----------+----------+
	// | 1  |  1  | X'00' |  1   | Variable |    2     |
	// +----+-----+-------+------+----------+----------+
	// VER 版本号，socks5的值为0x05
	// CMD 0x01表示CONNECT请求
	// RSV 保留字段，值为0x00
	// ATYP 目标地址类型，DST.ADDR的数据对应这个字段的类型。
	//   0x01表示IPv4地址，DST.ADDR为4个字节
	//   0x03表示域名，DST.ADDR是一个可变长度的域名
	// DST.ADDR 一个可变长度的值
	// DST.PORT 目标端口，固定2个字节
	buf := make([]byte, 4)
	_, err = io.ReadFull(render, buf)
	if err != nil {
		return err
	}

	Ver, Cmd, _, Atyp := buf[0], buf[1], buf[2], buf[3] //不考虑保留字段
	if Ver != byte(socketv5) {
		return err
	}
	if Cmd != byte(cmdBind) {
		return err
	}
	addr := ""

	if err != nil {
		return err
	}
	switch Atyp {

	case atypeIPV4:
		_, err = io.ReadFull(render, buf)
		if err != nil {
			return err
		}
		addr = fmt.Sprintf("%d.%d.%d.%d", buf[0], buf[1], buf[2], buf[3])
	case atypeHOST:
		hostSize, err := render.ReadByte()
		if err != nil {
			return err
		}
		host := make([]byte, hostSize)
		_, err = io.ReadFull(render, host)
		if err != nil {
			return err
		}
		addr = fmt.Sprintf("%d.%d.%d.%d", host[0], host[1], host[2], host[3])
	case atypeIPV6:
		return fmt.Errorf("IPV6不支持")
	default:
		return fmt.Errorf("不支持版本")
	}
	_, err = io.ReadFull(render, buf[:2])
	if err != nil {
		return err
	}
	PORT := binary.BigEndian.Uint16(buf[:2])

	//连接
	server, err := net.DialTimeout("tcp", fmt.Sprint(addr, PORT), time.Microsecond*200)
	if err != nil {
		return err
	}

	log.Println("dial", addr, PORT)

	// +----+-----+-------+------+----------+----------+
	// |VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
	// +----+-----+-------+------+----------+----------+
	// | 1  |  0  | X'00' |  1   | Variable |    2     |
	// +----+-----+-------+------+----------+----------+
	// VER socks版本，这里为0x05
	// REP Relay field,内容取值如下 X’00’ succeeded
	// RSV 保留字段
	// ATYPE 地址类型
	// BND.ADDR 服务绑定的地址
	// BND.PORT 服务绑定的端口DST.PORT
	_, err = Conn.Write([]byte{socketv5, 0x00, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
	if err != nil {
		return err
	}
	//通道缓冲区
	//不太懂
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		_, err = io.Copy(server, render)
		cancel()
	}()

	go func() {
		_, err = io.Copy(Conn, server)
		cancel()
	}()
	//不太懂
	<-ctx.Done()
	return nil

}
