package nft

import (
	"github.com/xuri/excelize/v2"
)

func ReadXLSXFile(file string) (*excelize.File, error) {
	return excelize.OpenFile(file)
}

func WriteXLSXFile(fileName string, file *excelize.File) error {
	return file.SaveAs(fileName)
}
