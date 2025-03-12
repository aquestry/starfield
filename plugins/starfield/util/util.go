package util

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/aquestry/starfield/plugins/starfield/orch/node"
	"io"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

func RandomString() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	s := ""
	for i := 0; i < 6; i++ {
		c := 'a' + rune(r.Intn(26))
		if r.Intn(2) == 0 {
			c -= 32
		}
		s += string(c)
	}
	return s
}

func GetPort(n node.Node, name string) (int, error) {
	p, e := n.Run("docker", "port", name, "25565")
	if e != nil {
		return 0, e
	}
	parts := strings.SplitN(p, ":", 2)
	if len(parts) < 2 {
		return 0, e
	}
	port, _ := strconv.Atoi(parts[1])
	return port, nil
}

func GetState(addr string) (string, bool) {
	if !strings.Contains(addr, ":") {
		addr = addr + ":25565"
	}
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return err.Error(), false
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(5 * time.Second))
	host, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		return err.Error(), false
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return err.Error(), false
	}
	var handshakeBody bytes.Buffer
	tmp := int32(754)
	for {
		b := byte(tmp & 0x7F)
		tmp >>= 7
		if tmp != 0 {
			b |= 0x80
		}
		handshakeBody.WriteByte(b)
		if tmp == 0 {
			break
		}
	}
	hostBytes := []byte(host)
	tmp = int32(len(hostBytes))
	for {
		b := byte(tmp & 0x7F)
		tmp >>= 7
		if tmp != 0 {
			b |= 0x80
		}
		handshakeBody.WriteByte(b)
		if tmp == 0 {
			break
		}
	}
	handshakeBody.Write(hostBytes)
	var portBytes [2]byte
	binary.BigEndian.PutUint16(portBytes[:], uint16(port))
	handshakeBody.Write(portBytes[:])
	tmp = int32(1)
	for {
		b := byte(tmp & 0x7F)
		tmp >>= 7
		if tmp != 0 {
			b |= 0x80
		}
		handshakeBody.WriteByte(b)
		if tmp == 0 {
			break
		}
	}
	var handshakePacket bytes.Buffer
	packetLen := handshakeBody.Len() + 1
	tmp = int32(packetLen)
	for {
		b := byte(tmp & 0x7F)
		tmp >>= 7
		if tmp != 0 {
			b |= 0x80
		}
		handshakePacket.WriteByte(b)
		if tmp == 0 {
			break
		}
	}
	handshakePacket.WriteByte(0x00)
	handshakePacket.Write(handshakeBody.Bytes())
	if _, err := conn.Write(handshakePacket.Bytes()); err != nil {
		return err.Error(), false
	}
	var requestPacket bytes.Buffer
	tmp = 1
	for {
		b := byte(tmp & 0x7F)
		tmp >>= 7
		if tmp != 0 {
			b |= 0x80
		}
		requestPacket.WriteByte(b)
		if tmp == 0 {
			break
		}
	}
	requestPacket.WriteByte(0x00)
	if _, err := conn.Write(requestPacket.Bytes()); err != nil {
		return err.Error(), false
	}
	readVarInt := func() (int32, error) {
		var numRead int32
		var result int32
		for {
			var buf [1]byte
			if _, err := conn.Read(buf[:]); err != nil {
				return 0, err
			}
			value := buf[0] & 0x7F
			result |= int32(value) << (7 * numRead)
			numRead++
			if numRead > 5 {
				return 0, fmt.Errorf("varint too big")
			}
			if (buf[0] & 0x80) == 0 {
				break
			}
		}
		return result, nil
	}
	_, err = readVarInt()
	if err != nil {
		return err.Error(), false
	}
	_, err = readVarInt()
	if err != nil {
		return err.Error(), false
	}
	length, err := readVarInt()
	if err != nil {
		return err.Error(), false
	}
	responseData := make([]byte, length)
	if _, err := io.ReadFull(conn, responseData); err != nil {
		return err.Error(), false
	}
	var result map[string]interface{}
	if err := json.Unmarshal(responseData, &result); err != nil {
		return err.Error(), false
	}
	description, ok := result["description"]
	if !ok {
		return "No MOTD found", false
	}
	switch d := description.(type) {
	case string:
		return d, true
	case map[string]interface{}:
		if text, ok := d["text"].(string); ok {
			return text, true
		}
		j, err := json.Marshal(d)
		if err != nil {
			return err.Error(), false
		}
		return string(j), false
	default:
		return fmt.Sprintf("%v", d), true
	}
}
