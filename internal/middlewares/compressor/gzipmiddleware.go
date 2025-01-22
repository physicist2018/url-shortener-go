package compressor

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// Middleware для поддержки gzip для запросов и ответов
func GzipMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Оригинальный http.ResponseWriter
		ow := w

		// Проверяем поддержку gzip у клиента для ответа
		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			// Оборачиваем ResponseWriter в поддержку gzip
			cw := newCompressWriter(w)
			ow = cw
			defer cw.Close() // Завершаем сжатие после выполнения middleware
		}

		// Проверяем, сжато ли тело запроса в gzip
		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			// Оборачиваем тело запроса в gzip.Reader для декомпрессии
			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer cr.Close() // Закрываем после завершения
		}

		// Передаём управление обработчику
		h.ServeHTTP(ow, r)
	}
}

// Создает новый gzip.Writer для сжатия ответа
func newCompressWriter(w http.ResponseWriter) *gzipResponseWriter {
	// Устанавливаем заголовок Content-Encoding для gzip
	w.Header().Set("Content-Encoding", "gzip")
	gz := gzip.NewWriter(w)
	return &gzipResponseWriter{ResponseWriter: w, Writer: gz}
}

// gzipResponseWriter оборачивает ResponseWriter для сжатия
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

// Создает новый gzip.Reader для декомпрессии тела запроса
func newCompressReader(r io.Reader) (io.ReadCloser, error) {
	return gzip.NewReader(r)
}
