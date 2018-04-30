package Util

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
)

func ReadRow(reader io.Reader) (res *Row, err error) {
	encodedBytes, err := ReadMessage(reader)
	if err == io.EOF {
		return nil, io.EOF
	}
	if err != nil {
		return nil, err
	}
	if res, err = DecodeRow(encodedBytes); err != nil {
		return nil, err
	}
	return res, nil
}

func DecodeRow(encodedBytes []byte) (res *Row, err error) {
	buf := bytes.NewBuffer(encodedBytes)
	decoder := gob.NewDecoder(buf)
	res = &Row{}
	err = decoder.Decode(res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func WriteRow(writer io.Writer, row *Row) (err error) {
	if row == nil {
		return fmt.Errorf("Nil Row")
	}
	encodedBytes, err := EncodeRow(row)
	if err != nil {
		return nil
	}
	return WriteMessage(writer, encodedBytes)
}

func EncodeRow(row *Row) (res []byte, err error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err = encoder.Encode(*row); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}