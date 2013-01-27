package gozdcf

import (
	"testing"
)

func TestunmarshalZpl(t *testing.T) {
	raw := `
#   Notice that indentation is always 4 spaces, there are no tabs.
#
context
    iothreads = 1
    verbose = 1      #   Ask for a trace

main
    type = zmq_queue
    frontend
        option
            hwm = 1000
            swap = 25000000
            subscribe = "#2"
        bind = tcp://eth0:5555
    backend
        bind = tcp://eth0:5556`
	conf := make(map[string]interface{})
	err := Unmarshal([]byte(raw), data)
	if err != nil {
		t.Fatalf("failed to unmarshal: %s", err)
	}
	if conf == nil {
		t.Fatalf("unmarshal returned two nils.")
	}
	tmp, err := conf["context"]
	if err != nil {
		t.Fatalf("context not found.")
	}
	context := tmp.(map[string]interface{})
	tmp, err := context["iothreads"]
	if err != nil {
		t.Fatalf("iothreads not found.")
	}
}
