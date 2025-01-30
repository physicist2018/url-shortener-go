package compressor

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// Обертка для ResponseWriter, чтобы перехватывать вывод и сжимать его.
type GzipResponseWriter struct {
	Writer io.Writer
	http.ResponseWriter
}

func (w *GzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func RequestDecompressionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем, сжат ли запрос
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			// Распаковываем тело запроса
			reader, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "Failed to decompress request body", http.StatusBadRequest)
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
					http.Error(w, "Failed to zip data", http.StatusBadRequest)
					return
				}
				defer gzipWriter.Close()
				gzipResponseWriter := &GzipResponseWriter{Writer: gzipWriter, ResponseWriter: w}
				next.ServeHTTP(gzipResponseWriter, r)
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}
