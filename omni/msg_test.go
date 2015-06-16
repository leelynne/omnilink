package omni

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func testDeser(t *testing.T) {
	buf := &bytes.Buffer{}
	var seqnum uint16
	var mt uint8
	var res byte
	seqnum = 1
	mt = ControllerAckNewSession
	binary.Write(buf, binary.LittleEndian, seqnum)
	binary.Write(buf, binary.LittleEndian, mt)
	binary.Write(buf, binary.LittleEndian, res)
	//buf.Write(p []byte)
}
