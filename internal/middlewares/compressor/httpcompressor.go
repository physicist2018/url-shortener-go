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

func CompressionMiddlewareForHandler(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Распаковка входящих запросов
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			reader, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "Failed to decompress request body", http.StatusBadRequest)
				return
			}
			defer reader.Close()
			r.Body = io.NopCloser(reader) // Оборачиваем сжатое тело в ReadCloser
		}

		// Перехватываем ответ для сжатия, если клиент поддерживает gzip
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			// Создаем gzip-обертку для ответа
			w.Header().Set("Content-Encoding", "gzip")
			gzipWriter := gzip.NewWriter(w)
			defer gzipWriter.Close()

			// Оборачиваем ResponseWriter в gzip и передаем управление следующему обработчику
			gzipResponseWriter := &GzipResponseWriter{
				Writer:         gzipWriter,
				ResponseWriter: w,
			}
			next.ServeHTTP(gzipResponseWriter, r)
		} else {
			// Если клиент не поддерживает gzip, просто передаем запрос дальше
			next.ServeHTTP(w, r)
		}
	})
}
