package lsm

import (
	"testing"
)

func TestMemTablePutAndGet(t *testing.T) {

	mt := NewMemTable()

	mt.Put([]byte("cpu"), []byte("42"))
	mt.Put([]byte("mem"), []byte("2048"))
	mt.Put([]byte("proc:neovim"), []byte("128"))

	if mt.Len() != 3 {
		t.Fatalf("expected 3 entries, got %d", mt.Len())
	}

	tests := []struct {
		key    string
		want   string
		wantOK bool
	}{
		{"cpu", "42", true},
		{"mem", "2048", true},
		{"proc:neovim", "128", true},
	}

	for _, tt := range tests {
		val, ok := mt.Get([]byte(tt.key))
		if ok != tt.wantOK {
			t.Errorf("get(%q) ok = %v, want %v", tt.key, ok, tt.wantOK)
			continue
		}

		if string(val) != tt.want {
			t.Errorf("get(%q) = %q, want %q", tt.key, val, tt.want)
		}
	}

	if val, ok := mt.Get([]byte("aaaaaaa")); ok {
		t.Errorf("expected false. got value %q", val)
	}

	mt.Put([]byte("cpu"), []byte("99"))
	if val, ok := mt.Get([]byte("cpu")); !ok || string(val) != "99" {
		t.Errorf("after update, Get(cpu) = %q, %v, want %q, true", val, ok, "99")
	}

	if mt.Len() != 3 {
		t.Fatalf("after update, expected 3 entries, got %d", mt.Len())
	}

}
