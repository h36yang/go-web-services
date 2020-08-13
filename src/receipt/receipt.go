package receipt

import (
	"io/ioutil"
	"path/filepath"
	"time"
)

// ReceiptDirectory : Receipt Upload Directory
var ReceiptDirectory string = filepath.Join("uploads")

// Receipt class
type Receipt struct {
	ReceiptName string    `json:"name"`
	UploadDate  time.Time `json:"uploadDate"`
}

// GetReceipts gets a slice of all receipts in the upload directory
func GetReceipts() ([]Receipt, error) {
	receipts := make([]Receipt, 0)
	files, err := ioutil.ReadDir(ReceiptDirectory)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		receipts = append(receipts, Receipt{ReceiptName: f.Name(), UploadDate: f.ModTime()})
	}
	return receipts, nil
}
