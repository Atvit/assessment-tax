package tax

import (
	"github.com/gocarina/gocsv"
	"io/ioutil"
	"mime/multipart"
)

func ReadCSV[T any](file *multipart.FileHeader, data T) (T, error) {
	f, err := file.Open()
	if err != nil {
		return *new(T), err
	}

	fileBytes, err := ioutil.ReadAll(f)
	if err != nil {
		return *new(T), err
	}

	if err := gocsv.UnmarshalBytes(fileBytes, &data); err != nil {
		return *new(T), err
	}

	return data, nil
}
