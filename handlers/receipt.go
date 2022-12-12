package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/GerardoHP/inventory-service/data"
)

const ReceiptPath = "receipts"

type Receipt struct {
}

func NewReceiptHandler() *Receipt {
	return &Receipt{}
}

func (receipt Receipt) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, fmt.Sprintf("%s/", ReceiptPath))
	if len(urlPathSegments[:1]) > 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fileName := urlPathSegments[1:][0]
	file, err := os.Open(filepath.Join(data.ReceiptDirectory, fileName))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	defer file.Close()
	fHeader := make([]byte, 512)
	file.Read(fHeader)
	fContentType := http.DetectContentType(fHeader)
	stat, err := file.Stat()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fSize := strconv.FormatInt(stat.Size(), 10)
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", fContentType)
	w.Header().Set("Content-Length", fSize)
	file.Seek(0, 0)
	io.Copy(w, file)
}
