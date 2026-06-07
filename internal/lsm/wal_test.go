package lsm

import (
	"os"
	"testing"
)

func TestWALCreateAndClose(t *testing.T) {
	path := "/tmp/test_wal.log"

	os.Remove(path)

	wal, err := NewWAL(path)
	if err != nil {
		t.Fatalf("NewWAL failed: %v", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatalf("File was not created")
	}

	if err := wal.Close(); err != nil {
		t.Fatalf("Close failed: %v", err)
	}

	os.Remove(path)
}

func TestWALAppendAndRecover(t *testing.T) {
	path := "/tmp/test_append_wal.log"
	os.Remove(path)

	wal, err := NewWAL(path)
	if err != nil {
		t.Fatalf("error creating the file: %v", err)
	}

	records := []LogRecord{
		{Key: []byte("cpu"), Value: []byte("42")},
		{Key: []byte("mem"), Value: []byte("2048")},
		{Key: []byte("proc:neovim"), Value: []byte("128")},
	}

	for _, r := range records {
		if err := wal.Append(r.Key, r.Value); err != nil {
			t.Fatalf("Append failed: %v", err)
		}
	}

	wal.Close()

	wal2, err := NewWAL(path)
	if err != nil {
		t.Fatalf("error creating the file: %v", err)
	}

	got, err := wal2.Recover()
	if err != nil {
		t.Fatalf("Recover failed: %v", err)
	}

	wal2.Close()

	if len(got) != len(records) {
		t.Fatalf("got %d records, want %d", len(got), len(records))
	}

	for i, r := range got {
		if string(r.Key) != string(records[i].Key) {
			t.Errorf("record %d: key = %q, want %q", i, r.Key, records[i].Key)
		}
		if string(r.Value) != string(records[i].Value) {
			t.Errorf("record %d: value = %q, want %q", i, r.Value, records[i].Value)
		}
	}

	os.Remove(path)
}
