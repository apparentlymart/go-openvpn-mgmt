package openvpn

import (
	"bytes"
	"fmt"
)

var eventSep = []byte(":")
var bytecountEventKW = []byte("BYTECOUNT")
var bytecountCliEventKW = []byte("BYTECOUNT_CLI")
var clientEventKW = []byte("CLIENT")
var echoEventKW = []byte("ECHO")
var fatalEventKW = []byte("FATAL")
var holdEventKW = []byte("HOLD")
var infoEventKW = []byte("INFO")
var logEventKW = []byte("LOG")
var needOkEventKW = []byte("NEED-OK")
var needStrEventKW = []byte("NEED-STR")
var passwordEventKW = []byte("PASSWORD")
var stateEventKW = []byte("STATE")

type Event interface {
	String() string
}

type UnknownEvent struct {
	keyword []byte
	body    []byte
}

func (e *UnknownEvent) Type() string {
	return string(e.keyword)
}

func (e *UnknownEvent) Body() string {
	return string(e.body)
}

func (e *UnknownEvent) String() string {
	return fmt.Sprintf("%s: %s", e.keyword, e.body)
}

type MalformedEvent struct {
	raw []byte
}

func (e *MalformedEvent) String() string {
	return fmt.Sprintf("Malformed Event %q", e.raw)
}

func upgradeEvent(raw []byte) Event {
	splitIdx := bytes.Index(raw, eventSep)
	if splitIdx == -1 {
		// Should never happen, but we'll handle it robustly if it does.
		return &MalformedEvent{raw}
	}

	keyword := raw[:splitIdx]
	body := raw[splitIdx+1:]

	switch {
	default:
		return &UnknownEvent{keyword, body}
	}
}
