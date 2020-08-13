package receipt

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"demo/inventoryservice/cors"
)

const receiptsBasePath = "receipts"

// SetupRoutes sets up the routes in Receipt Service
func SetupRoutes(apiBasePath string) {
	receiptsHandler := http.HandlerFunc(handleReceipts)
	downloadHandler := http.HandlerFunc(handleDownload)
	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, receiptsBasePath), cors.Middleware(receiptsHandler))
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, receiptsBasePath), cors.Middleware(downloadHandler))
}

func handleReceipts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// get all receipts
		receiptsList, err := GetReceipts()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		j, err := json.Marshal(receiptsList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
		w.Write(j)

	case http.MethodPost:
		// upload a new receipt
		r.ParseMultipartForm(5 << 20) // 5MB
		file, handler, err := r.FormFile("receipt")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}
		defer file.Close()

		f, err := os.OpenFile(filepath.Join(ReceiptDirectory, handler.Filename), os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}
		defer f.Close()

		io.Copy(f, file)
		w.WriteHeader(http.StatusCreated)

	case http.MethodOptions:
		// Check CORS headers - do nothing here; let Middleware handler handle this

	default:
		// Unhandled methods
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleDownload(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, fmt.Sprintf("%s/", receiptsBasePath))
	if len(urlPathSegments[1:]) > 1 {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("bad URL parameters:", urlPathSegments)
		return
	}

	fileName := urlPathSegments[1]
	file, err := os.Open(filepath.Join(ReceiptDirectory, fileName))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println(err)
		return
	}
	defer file.Close()

	fileHeader := make([]byte, 512)
	file.Read(fileHeader)
	contentType := http.DetectContentType(fileHeader)
	stat, err := file.Stat()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	contentLength := strconv.FormatInt(stat.Size(), 10)

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", contentLength)
	file.Seek(0, 0)
	io.Copy(w, file)
}
