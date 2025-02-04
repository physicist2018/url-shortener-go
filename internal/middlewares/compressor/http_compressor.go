package compressor

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

const (
	BestCompression = gzip.BestCompression
)

// Обертка для ResponseWriter, чтобы перехватывать вывод и сжимать его.
type gzipResponseWriter struct {
	writer io.Writer
	http.ResponseWriter
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	return w.writer.Write(b)
}

func RequestDecompressionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем, сжат ли запрос
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			// Распаковываем тело запроса
			reader, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "Ошибка восстановления данных", http.StatusBadRequest)
				return
			}

			r.Body = reader
		}
		next.ServeHTTP(w, r)
	})
}

func ResponseCompressionMiddleware(compressionLevel int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Перехватываем ответ для сжатия
			if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				w.Header().Set("Content-Encoding", "gzip")
				gzipWriter, err := gzip.NewWriterLevel(w, compressionLevel)
				if err != nil {
					http.Error(w, "Ошибка сжатия данных", http.StatusBadRequest)
					return
				}
				defer gzipWriter.Close()
				gzipResponseWriter := &gzipResponseWriter{writer: gzipWriter, ResponseWriter: w}
				next.ServeHTTP(gzipResponseWriter, r)
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}
