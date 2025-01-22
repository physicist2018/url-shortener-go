package compressor

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// Функция для проверки, нужно ли сжимать данные в ответе
func shouldCompress(contentType string) bool {
	return contentType == "application/json" || contentType == "text/html"
}

// Функция для создания gzip.Reader для декомпрессии
func newGzipReader(r io.Reader) (io.ReadCloser, error) {
	return gzip.NewReader(r)
}

// Функция для создания gzip.Writer для сжатия
func newGzipWriter(w http.ResponseWriter) *gzipResponseWriter {
	gz := gzip.NewWriter(w)
	return &gzipResponseWriter{
		ResponseWriter: w,
		Writer:         gz,
	}
}

// gzipResponseWriter оборачивает ResponseWriter и добавляет поддержку gzip-сжатия
type gzipResponseWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (w *gzipResponseWriter) Close() error {
	if gz, ok := w.Writer.(*gzip.Writer); ok {
		return gz.Close()
	}
	return nil
}

// GzipMiddleware для сжатия и декомпрессии запросов и ответов
func GzipMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Декодируем тело запроса, если оно сжато
		contentEncoding := r.Header.Get("Content-Encoding")
		if contentEncoding == "gzip" {
			cr, err := newGzipReader(r.Body)
			if err != nil {
				http.Error(w, "Failed to decompress request body", http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer cr.Close() // Закрываем декомпрессор после выполнения
		}

		// Проверяем, поддерживает ли клиент сжатие
		acceptEncoding := r.Header.Get("Accept-Encoding")
		if strings.Contains(acceptEncoding, "gzip") && shouldCompress(r.Header.Get("Content-Type")) {
			// Устанавливаем заголовок Content-Encoding для сжатого ответа
			w.Header().Set("Content-Encoding", "gzip")
			// Оборачиваем ResponseWriter в gzip-сжатие
			cw := newGzipWriter(w)
			defer cw.Close() // Закрываем сжатие после выполнения

			// Передаем управление обработчику с сжатым ResponseWriter
			h.ServeHTTP(cw, r)
		} else {
			// Если сжатие не требуется, передаем управление обычному обработчику
			h.ServeHTTP(w, r)
		}
	}
}
