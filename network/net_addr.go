package network

import (
	"fmt"
	"stapler/utils"
	"strconv"
	"strings"
)

type TunnelType byte

const (
	TUN_TCP TunnelType = 'T'
	TUN_UDP TunnelType = 'U'
)

type NetAddr struct {
	val []byte
}

func NewNetAddr() *NetAddr {
	return &NetAddr{val: make([]byte, 15)}
}

func WrapAddr(src []byte) *NetAddr {
	slice := src[:15]
	return &NetAddr{val: slice}
}

func MakeAddrByString(src *string) *NetAddr {
	slice := ChannelKeyToBytes(src)
	return &NetAddr{val: slice}
}

func (c *NetAddr) SetTunnel(tunnelType TunnelType) {
	c.val[0] = byte(tunnelType)
}

func (c *NetAddr) GetTunnel() TunnelType {
	return TunnelType(c.val[0])
}

func (c *NetAddr) GetValue() []byte {
	return c.val
}

func (c *NetAddr) String() string {
	return BytesToChannelKey(c.val)
}

func (c *NetAddr) GetLocalPort() int {
	return int(c.val[7]>>8 | c.val[8])
}

func (c *NetAddr) Clear() {
	for i := 0; i < 15; i++ {
		c.val[i] = 0
	}
}

func (c *NetAddr) GetRemoteAddr() []byte {
	return c.val[9:15]
}

func (c *NetAddr) CopyTo(dst []byte, offset int) {
	src := c.val[0:15]
	dst2 := dst[offset : offset+15]
	copy(dst2, src)
}

func (c *NetAddr) CopyFrom(src *NetAddr) {
	src.CopyTo(c.val, 0)
}

/**
 * 地址化为二进制表示        0 1    2   3  4 5 6     7   8  9  10  11
 * @param _session 会话地址 T:999-192.168.5.7:6000/202.98.232.22:4097
 * @return 字节数组 [ 1 + 2 + 6 + 6 ]
 */
func ChannelKeyToBytes(_key *string) []byte {
	buf := make([]byte, 15)
	tmp := channelKeyToAry(_key)

	if strings.Compare("T", tmp[0]) == 0 {
		buf[0] = byte(TUN_TCP)
	} else {
		buf[0] = byte(TUN_UDP)
	}
	seq := utils.HexToInt(tmp[1])
	buf[1] = byte(seq >> 8)
	buf[2] = byte(seq)

	buf[3] = utils.StrToByte(tmp[2])
	buf[4] = utils.StrToByte(tmp[3])
	buf[5] = utils.StrToByte(tmp[4])
	buf[6] = utils.StrToByte(tmp[5])

	port, _ := strconv.Atoi(tmp[6])
	buf[7] = byte(port >> 8)
	buf[8] = byte(port)

	buf[9] = utils.StrToByte(tmp[7])
	buf[10] = utils.StrToByte(tmp[8])
	buf[11] = utils.StrToByte(tmp[9])
	buf[12] = utils.StrToByte(tmp[10])

	port, _ = strconv.Atoi(tmp[11])
	buf[13] = byte(port >> 8)
	buf[14] = byte(port)

	return buf
}

/**
 * 地址化为字符数组	   0 1    2  3   4 5 6     7  8  9   10  11
 * @param _session 地址 T0000-192.168.5.7:6000-202.98.232.22:4097
 * @return 字符数据
 */
func channelKeyToAry(_key *string) []string {
	ary := make([]string, 12)

	var sx, ex int

	buf := []byte(*_key)
	lens := len(buf)

	ary[0] = string(buf[0:1])
	ary[1] = string(buf[1:5])

	ex = utils.IndexOf(&buf, 6, lens, '.')
	ary[2] = string(buf[6:ex])
	sx = ex + 1

	ex = utils.IndexOf(&buf, sx, lens, '.')
	ary[3] = string(buf[sx:ex])
	sx = ex + 1

	ex = utils.IndexOf(&buf, sx, lens, '.')
	ary[4] = string(buf[sx:ex])
	sx = ex + 1

	ex = utils.IndexOf(&buf, sx, lens, ':')
	ary[5] = string(buf[sx:ex])
	sx = ex + 1

	ex = utils.IndexOf(&buf, sx, lens, '/')
	ary[6] = string(buf[sx:ex])
	sx = ex + 1

	ex = utils.IndexOf(&buf, sx, lens, '.')
	ary[7] = string(buf[sx:ex])
	sx = ex + 1

	ex = utils.IndexOf(&buf, sx, lens, '.')
	ary[8] = string(buf[sx:ex])
	sx = ex + 1

	ex = utils.IndexOf(&buf, sx, lens, '.')
	ary[9] = string(buf[sx:ex])
	sx = ex + 1

	ex = utils.IndexOf(&buf, sx, lens, ':')
	ary[10] = string(buf[sx:ex])

	ary[11] = string(buf[ex+1:])

	return ary
}

func makeWord(_b1, _b2 byte) int {
	return int(int(_b1)<<8 | int(_b2))
}

func IpPortToStr(p []byte) string {
	return fmt.Sprintf("%d.%d.%d.%d:%d", p[0], p[1], p[2], p[3], makeWord(p[4], p[5]))
}

/**
 *                     0        12       3456       78           9012       34
 * 会话地址 bytes[15] = Tun[1] + SEQ[2] + SvrIP[4] + SvrPort[2] + CliIP[4] + CliPort[2]
 */
func BytesToChannelKey(p []byte) string {
	svr := IpPortToStr(p[3:9])
	cli := IpPortToStr(p[9:15])
	return fmt.Sprintf("%c%.4X-%s/%s", p[0], makeWord(p[1], p[2]), svr, cli)
}
