package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/GerardoHP/inventory-service/data"
)

type Receipts struct {
}

func NewReceiptsHandler() *Receipts {
	return &Receipts{}
}

func (receipt *Receipts) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		receiptList, err := data.GetReceipts()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		j, err := json.Marshal(receiptList)
		if err != nil {
			log.Fatal(err)
		}

		_, err = w.Write(j)
		if err != nil {
			log.Fatal(err)
		}
	case http.MethodPost:
		r.ParseMultipartForm(5 << 20) // 5 Mb
		file, handler, err := r.FormFile("receipt")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		defer file.Close()
		if _, err := os.Stat(data.ReceiptDirectory); errors.Is(err, os.ErrNotExist) {
			err := os.Mkdir(data.ReceiptDirectory, os.ModePerm)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		f, _ := os.OpenFile(filepath.Join(data.ReceiptDirectory, handler.Filename), os.O_WRONLY|os.O_CREATE, 0666)
		defer f.Close()
		io.Copy(f, file)
		w.WriteHeader(http.StatusCreated)
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
