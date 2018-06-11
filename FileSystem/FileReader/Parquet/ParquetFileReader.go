package Parquet

import (
	"io"

	"github.com/xitongsys/guery/FileSystem"
	"github.com/xitongsys/guery/Util"
	. "github.com/xitongsys/parquet-go/ParquetFile"
	. "github.com/xitongsys/parquet-go/ParquetReader"
	. "github.com/xitongsys/parquet-go/ParquetType"
	"github.com/xitongsys/parquet-go/parquet"
)

type PqFile struct {
	FileName string
	VF       FileSystem.VirtualFile
}

func (self *PqFile) Create(name string) (ParquetFile, error) {
	return nil, nil
}

func (self *PqFile) Open(name string) (ParquetFile, error) {
	if name == "" {
		name = self.FileName
	}
	vf, err := FileSystem.Open(name)
	if err != nil {
		return nil, err
	}
	res := &PqFile{
		VF:       vf,
		FileName: name,
	}
	return res, nil
}
func (self *PqFile) Seek(offset int64, pos int) (int64, error) {
	return self.VF.Seek(offset, pos)
}
func (self *PqFile) Read(b []byte) (n int, err error) {
	return self.VF.Read(b)
}
func (self *PqFile) Write(b []byte) (n int, err error) {
	return 0, nil
}
func (self *PqFile) Close() error { return nil }

type ParquetFileReader struct {
	pqReader *ParquetReader
	NumRows  int
	Cursor   int

	ReadColumnIndexes        []int
	ReadColumnTypes          []*parquet.Type
	ReadColumnConvertedTypes []*parquet.ConvertedType
}

func New(fileName string) *ParquetFileReader {
	parquetFileReader := new(ParquetFileReader)
	var pqFile ParquetFile = &PqFile{}
	pqFile, _ = pqFile.Open(fileName)
	parquetFileReader.pqReader, _ = NewParquetColumnReader(pqFile, 1)
	parquetFileReader.NumRows = int(parquetFileReader.pqReader.GetNumRows())
	return parquetFileReader
}

func (self *ParquetFileReader) ReadHeader() (fieldNames []string, err error) {
	return self.pqReader.SchemaHandler.ValueColumns, nil
}

func (self *ParquetFileReader) Read() (row *Util.Row, err error) {
	if self.Cursor >= self.NumRows {
		return nil, io.EOF
	}
	if len(self.ReadColumnIndexes) <= 0 {
		indexes := make([]int, len(self.pqReader.SchemaHandler.ValueColumns))
		for i := 0; i < len(indexes); i++ {
			indexes[i] = i
		}
		self.SetReadColumns(indexes)
	}
	objects := make([]interface{}, 0)
	for i, index := range self.ReadColumnIndexes {
		values, _, _ := self.pqReader.ReadColumnByIndex(index, 1)
		objects = append(objects, ParquetTypeToGoType(values[0],
			self.ReadColumnTypes[i],
			self.ReadColumnConvertedTypes[i],
		))
	}
	self.Cursor++
	row = &Util.Row{}
	row.AppendVals(objects...)
	return row, nil
}

func (self *ParquetFileReader) SetReadColumns(indexes []int) {
	for _, index := range indexes {
		fieldName := self.pqReader.SchemaHandler.ValueColumns[index]
		schemaIndex := self.pqReader.SchemaHandler.MapIndex[fieldName]
		t := self.pqReader.SchemaHandler.SchemaElements[schemaIndex].Type
		ct := self.pqReader.SchemaHandler.SchemaElements[schemaIndex].ConvertedType

		self.ReadColumnIndexes = append(self.ReadColumnIndexes, index)
		self.ReadColumnTypes = append(self.ReadColumnTypes, t)
		self.ReadColumnConvertedTypes = append(self.ReadColumnConvertedTypes, ct)
	}
}

func (self *ParquetFileReader) ReadByColumns(indexes []int) (row *Util.Row, err error) {
	if self.Cursor >= self.NumRows {
		return nil, io.EOF
	}
	if len(self.ReadColumnIndexes) <= 0 {
		self.SetReadColumns(indexes)
	}
	objects := make([]interface{}, 0)
	for i, index := range self.ReadColumnIndexes {
		values, _, _ := self.pqReader.ReadColumnByIndex(index, 1)
		objects = append(objects, ParquetTypeToGueryType(values[0],
			self.ReadColumnTypes[i],
			self.ReadColumnConvertedTypes[i],
		))
	}
	self.Cursor++
	row = &Util.Row{}
	row.AppendVals(objects...)
	return row, nil
}
