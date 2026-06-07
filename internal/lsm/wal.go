package lsm

import (
	"encoding/binary"
	"io"
	"os"
)

type LogRecord struct {
	Key   []byte
	Value []byte
}

type WAL struct {
	file *os.File
}

func NewWAL(path string) (*WAL, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	return &WAL{file: f}, nil
}

func (w *WAL) Append(key, value []byte) error {

	total := 4 + len(key) + 4 + len(value)
	buf := make([]byte, total)

	binary.LittleEndian.PutUint32(buf[0:4], uint32(len(key)))
	copy(buf[4:], key)

	binary.LittleEndian.PutUint32(buf[4+len(key):8+len(key)], uint32(len(value)))
	copy(buf[8+len(key):], value)

	_, err := w.file.Write(buf)
	if err != nil {
		return err
	}

	return w.file.Sync()
}

func (w *WAL) Recover() ([]LogRecord, error) {

	_, err := w.file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	var records []LogRecord
	var lastOffset int64

	for {

		var keyLen uint32
		if err := binary.Read(w.file, binary.LittleEndian, &keyLen); err != nil {
			if err == io.EOF {
				break
			}
			break
		}

		key := make([]byte, keyLen)
		if _, err := io.ReadFull(w.file, key); err != nil {
			break
		}

		var valLen uint32
		if err := binary.Read(w.file, binary.LittleEndian, &valLen); err != nil {
			break
		}

		value := make([]byte, valLen)
		if _, err := io.ReadFull(w.file, value); err != nil {
			break
		}

		records = append(records, LogRecord{Key: key, Value: value})

		pos, err := w.file.Seek(0, io.SeekCurrent)
		if err != nil {
			return nil, err
		}

		lastOffset = pos
	}

	if err := w.file.Truncate(lastOffset); err != nil {
		return nil, err
	}

	return records, nil
}

func (w *WAL) Close() error {
	return w.file.Close()
}
