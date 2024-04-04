package compressor

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type GzipCompressor struct {
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
	return w.Writer.Write(b)
}

func isAllowedContentType(contentType string) bool {
	availableContentTypes := []string{"application/json", "text/html"}
	for _, ct := range availableContentTypes {
		if strings.Contains(contentType, ct) {
			return true
		}
	}
	return false
}

func (compressor GzipCompressor) RequestCompressor(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		// проверяем, что клиент поддерживает gzip-сжатие
		// это упрощённый пример. В реальном приложении следует проверять все
		// значения r.Header.Values("Accept-Encoding") и разбирать строку
		// на составные части, чтобы избежать неожиданных результатов

		isAcceptEncoding := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
		isAllowedContentType := isAllowedContentType(r.Header.Get("Content-Type"))
		isValidRequest := isAcceptEncoding && isAllowedContentType

		if !isValidRequest {
			h.ServeHTTP(w, r)
			return
		}

		// создаём gzip.Writer поверх текущего w
		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		// передаём обработчику страницы переменную типа gzipWriter для вывода данных
		h.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	}
	return http.HandlerFunc(logFn)
}
