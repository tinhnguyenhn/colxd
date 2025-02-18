// Copyright (c) 2013-2015 The btcsuite developers
// Copyright (c) 2016 The Dash developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package wire_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/tinhnguyenhn/colxd/wire"
)

// TestGetAddr tests the MsgGetAddr API.
func TestGetAddr(t *testing.T) {
	pver := wire.ProtocolVersion

	// Ensure the command is expected value.
	wantCmd := "getaddr"
	msg := wire.NewMsgGetAddr()
	if cmd := msg.Command(); cmd != wantCmd {
		t.Errorf("NewMsgGetAddr: wrong command - got %v want %v",
			cmd, wantCmd)
	}

	// Ensure max payload is expected value for latest protocol version.
	// Num addresses (varInt) + max allowed addresses.
	wantPayload := uint32(0)
	maxPayload := msg.MaxPayloadLength(pver)
	if maxPayload != wantPayload {
		t.Errorf("MaxPayloadLength: wrong max payload length for "+
			"protocol version %d - got %v, want %v", pver,
			maxPayload, wantPayload)
	}

	return
}

// TestGetAddrWire tests the MsgGetAddr wire encode and decode for various
// protocol versions.
func TestGetAddrWire(t *testing.T) {
	msgGetAddr := wire.NewMsgGetAddr()
	msgGetAddrEncoded := []byte{}

	tests := []struct {
		in   *wire.MsgGetAddr // Message to encode
		out  *wire.MsgGetAddr // Expected decoded message
		buf  []byte           // Wire encoding
		pver uint32           // Protocol version for wire encoding
	}{
		// Latest protocol version.
		{
			msgGetAddr,
			msgGetAddr,
			msgGetAddrEncoded,
			wire.ProtocolVersion,
		},

		// Protocol version BIP0035Version.
		{
			msgGetAddr,
			msgGetAddr,
			msgGetAddrEncoded,
			wire.BIP0035Version,
		},

		// Protocol version BIP0031Version.
		{
			msgGetAddr,
			msgGetAddr,
			msgGetAddrEncoded,
			wire.BIP0031Version,
		},

		// Protocol version NetAddressTimeVersion.
		{
			msgGetAddr,
			msgGetAddr,
			msgGetAddrEncoded,
			wire.NetAddressTimeVersion,
		},

		// Protocol version MultipleAddressVersion.
		{
			msgGetAddr,
			msgGetAddr,
			msgGetAddrEncoded,
			wire.MultipleAddressVersion,
		},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Encode the message to wire format.
		var buf bytes.Buffer
		err := test.in.BtcEncode(&buf, test.pver)
		if err != nil {
			t.Errorf("BtcEncode #%d error %v", i, err)
			continue
		}
		if !bytes.Equal(buf.Bytes(), test.buf) {
			t.Errorf("BtcEncode #%d\n got: %s want: %s", i,
				spew.Sdump(buf.Bytes()), spew.Sdump(test.buf))
			continue
		}

		// Decode the message from wire format.
		var msg wire.MsgGetAddr
		rbuf := bytes.NewReader(test.buf)
		err = msg.BtcDecode(rbuf, test.pver)
		if err != nil {
			t.Errorf("BtcDecode #%d error %v", i, err)
			continue
		}
		if !reflect.DeepEqual(&msg, test.out) {
			t.Errorf("BtcDecode #%d\n got: %s want: %s", i,
				spew.Sdump(msg), spew.Sdump(test.out))
			continue
		}
	}
}
