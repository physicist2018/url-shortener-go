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

// CompressionMiddleware is a middleware function that compresses the response body using gzip.
// It takes a compression level as an argument and returns a function that takes an http.Handler and returns an http.Handler.
// The returned function checks if the request body is compressed and decompresses it if necessary.
// It also checks if the client accepts gzip compression and compresses the response body if necessary.
// If there is an error during the process, it returns a 400 Bad Request error.
func CompressionMiddleware(compressionLevel int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Проверяем, сжат ли запрос
			if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
				// Распаковываем тело запроса
				reader, err := gzip.NewReader(r.Body)
				if err != nil {
					http.Error(w, "Failed to decompress request body", http.StatusBadRequest)
					return
				}
				defer reader.Close()
				r.Body = io.NopCloser(reader)
			}

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
				// Если клиент не поддерживает сжатие, отправляем обычный ответ
				next.ServeHTTP(w, r)
			}
		})
	}
}
