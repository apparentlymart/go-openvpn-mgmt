package openvpn

import (
	"fmt"
	"testing"
)

func TestMalformedEvent(t *testing.T) {
	testCases := [][]byte{
		[]byte(""),
		[]byte("HTTP/1.1 200 OK"),
		[]byte("     "),
		[]byte("\x00"),
	}

	for i, testCase := range testCases {
		event := upgradeEvent(testCase)

		var malformed *MalformedEvent
		var ok bool
		if malformed, ok = event.(*MalformedEvent); !ok {
			t.Errorf("test %d got %T; want %T", i, event, malformed)
			continue
		}

		wantString := fmt.Sprintf("Malformed Event %q", testCase)
		if gotString := malformed.String(); gotString != wantString {
			t.Errorf("test %d String returned %q; want %q", i, gotString, wantString)
		}
	}
}

func TestUnknownEvent(t *testing.T) {
	type TestCase struct {
		Input    []byte
		WantType string
		WantBody string
	}
	testCases := []TestCase{
		{
			Input:    []byte("DUMMY:baz"),
			WantType: "DUMMY",
			WantBody: "baz",
		},
		{
			Input:    []byte("DUMMY:"),
			WantType: "DUMMY",
			WantBody: "",
		},
		{
			Input:    []byte("DUMMY:abc,123,456"),
			WantType: "DUMMY",
			WantBody: "abc,123,456",
		},
	}

	for i, testCase := range testCases {
		event := upgradeEvent(testCase.Input)

		var unk *UnknownEvent
		var ok bool
		if unk, ok = event.(*UnknownEvent); !ok {
			t.Errorf("test %d got %T; want %T", i, event, unk)
			continue
		}

		if got, want := unk.Type(), testCase.WantType; got != want {
			t.Errorf("test %d Type returned %q; want %q", i, got, want)
		}
		if got, want := unk.Body(), testCase.WantBody; got != want {
			t.Errorf("test %d Body returned %q; want %q", i, got, want)
		}
	}
}

func TestHoldEvent(t *testing.T) {
	testCases := [][]byte{
		[]byte("HOLD:"),
		[]byte("HOLD:waiting for hold release"),
	}

	for i, testCase := range testCases {
		event := upgradeEvent(testCase)

		var hold *HoldEvent
		var ok bool
		if hold, ok = event.(*HoldEvent); !ok {
			t.Errorf("test %d got %T; want %T", i, event, hold)
			continue
		}
	}
}

func TestEchoEvent(t *testing.T) {
	type TestCase struct {
		Input         []byte
		WantTimestamp string
		WantMessage   string
	}
	testCases := []TestCase{
		{
			Input:         []byte("ECHO:123,foo"),
			WantTimestamp: "123",
			WantMessage:   "foo",
		},
		{
			Input:         []byte("ECHO:123,"),
			WantTimestamp: "123",
			WantMessage:   "",
		},
		{
			Input:         []byte("ECHO:,foo"),
			WantTimestamp: "",
			WantMessage:   "foo",
		},
		{
			Input:         []byte("ECHO:,"),
			WantTimestamp: "",
			WantMessage:   "",
		},
		{
			Input:         []byte("ECHO:"),
			WantTimestamp: "",
			WantMessage:   "",
		},
	}

	for i, testCase := range testCases {
		event := upgradeEvent(testCase.Input)

		var echo *EchoEvent
		var ok bool
		if echo, ok = event.(*EchoEvent); !ok {
			t.Errorf("test %d got %T; want %T", i, event, echo)
			continue
		}

		if got, want := echo.RawTimestamp(), testCase.WantTimestamp; got != want {
			t.Errorf("test %d RawTimestamp returned %q; want %q", i, got, want)
		}
		if got, want := echo.Message(), testCase.WantMessage; got != want {
			t.Errorf("test %d Message returned %q; want %q", i, got, want)
		}
	}
}

func TestStateEvent(t *testing.T) {
	type TestCase struct {
		Input          []byte
		WantTimestamp  string
		WantState      string
		WantDesc       string
		WantLocalAddr  string
		WantRemoteAddr string
	}
	testCases := []TestCase{
		{
			Input:          []byte("STATE:"),
			WantTimestamp:  "",
			WantState:      "",
			WantDesc:       "",
			WantLocalAddr:  "",
			WantRemoteAddr: "",
		},
		{
			Input:          []byte("STATE:,"),
			WantTimestamp:  "",
			WantState:      "",
			WantDesc:       "",
			WantLocalAddr:  "",
			WantRemoteAddr: "",
		},
		{
			Input:          []byte("STATE:,,,,"),
			WantTimestamp:  "",
			WantState:      "",
			WantDesc:       "",
			WantLocalAddr:  "",
			WantRemoteAddr: "",
		},
		{
			Input:          []byte("STATE:123,CONNECTED,good,172.16.0.1,192.168.4.1"),
			WantTimestamp:  "123",
			WantState:      "CONNECTED",
			WantDesc:       "good",
			WantLocalAddr:  "172.16.0.1",
			WantRemoteAddr: "192.168.4.1",
		},
		{
			Input:          []byte("STATE:123,RECONNECTING,SIGHUP,,"),
			WantTimestamp:  "123",
			WantState:      "RECONNECTING",
			WantDesc:       "SIGHUP",
			WantLocalAddr:  "",
			WantRemoteAddr: "",
		},
		{
			Input:          []byte("STATE:123,RECONNECTING,SIGHUP,,,extra"),
			WantTimestamp:  "123",
			WantState:      "RECONNECTING",
			WantDesc:       "SIGHUP",
			WantLocalAddr:  "",
			WantRemoteAddr: "",
		},
	}

	for i, testCase := range testCases {
		event := upgradeEvent(testCase.Input)

		var st *StateEvent
		var ok bool
		if st, ok = event.(*StateEvent); !ok {
			t.Errorf("test %d got %T; want %T", i, event, st)
			continue
		}

		if got, want := st.RawTimestamp(), testCase.WantTimestamp; got != want {
			t.Errorf("test %d RawTimestamp returned %q; want %q", i, got, want)
		}
		if got, want := st.NewState(), testCase.WantState; got != want {
			t.Errorf("test %d NewState returned %q; want %q", i, got, want)
		}
		if got, want := st.Description(), testCase.WantDesc; got != want {
			t.Errorf("test %d Description returned %q; want %q", i, got, want)
		}
		if got, want := st.LocalTunnelAddr(), testCase.WantLocalAddr; got != want {
			t.Errorf("test %d LocalTunnelAddr returned %q; want %q", i, got, want)
		}
		if got, want := st.RemoteAddr(), testCase.WantRemoteAddr; got != want {
			t.Errorf("test %d RemoteAddr returned %q; want %q", i, got, want)
		}
	}
}
