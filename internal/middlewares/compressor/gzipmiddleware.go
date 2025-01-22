package compressor

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// func shouldCompress(contentType string) bool {
// 	switch contentType {
// 	case "application/javascript", "application/json", "text/css", "text/html", "text/plain", "text/xml":
// 		return true
// 	default:
// 		return false
// 	}
// }

// // Middleware для поддержки gzip для запросов и ответов
// func GzipMiddleware(h http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		// Оригинальный http.ResponseWriter
// 		ow := w

// 		// Проверяем поддержку gzip у клиента для ответа
// 		acceptEncoding := r.Header.Get("Accept-Encoding")
// 		supportsGzip := strings.Contains(acceptEncoding, "gzip")
// 		if supportsGzip {
// 			// Оборачиваем ResponseWriter в поддержку gzip
// 			cw := newCompressWriter(w)
// 			ow = cw
// 			defer cw.Close() // Завершаем сжатие после выполнения middleware
// 		}

// 		// Проверяем, сжато ли тело запроса в gzip
// 		contentEncoding := r.Header.Get("Content-Encoding")
// 		sendsGzip := strings.Contains(contentEncoding, "gzip")
// 		if sendsGzip {
// 			// Оборачиваем тело запроса в gzip.Reader для декомпрессии
// 			cr, err := newCompressReader(r.Body)
// 			if err != nil {
// 				w.WriteHeader(http.StatusInternalServerError)
// 				return
// 			}
// 			r.Body = cr
// 			defer cr.Close() // Закрываем после завершения
// 		}

// 		// Передаём управление обработчику
// 		h.ServeHTTP(ow, r)
// 	}
// }

// // Создает новый gzip.Writer для сжатия ответа
// func newCompressWriter(w http.ResponseWriter) *gzipResponseWriter {
// 	// Устанавливаем заголовок Content-Encoding для gzip
// 	w.Header().Set("Content-Encoding", "gzip")
// 	gz := gzip.NewWriter(w)
// 	return &gzipResponseWriter{ResponseWriter: w, Writer: gz}
// }

// // gzipResponseWriter оборачивает ResponseWriter для сжатия
// type gzipResponseWriter struct {
// 	http.ResponseWriter
// 	Writer io.Writer
// }

// func (w *gzipResponseWriter) Write(b []byte) (int, error) {
// 	return w.Writer.Write(b)
// }

// func (w *gzipResponseWriter) Close() error {
// 	if gz, ok := w.Writer.(*gzip.Writer); ok {
// 		return gz.Close()
// 	}
// 	return nil
// }

// // Создает новый gzip.Reader для декомпрессии тела запроса
// func newCompressReader(r io.Reader) (io.ReadCloser, error) {
// 	return gzip.NewReader(r)
// }

// Функция для проверки Content-Type, которые должны быть сжаты
func shouldCompress(contentType string) bool {
	switch contentType {
	case "application/javascript", "application/json", "text/css", "text/html", "text/plain", "text/xml":
		return true
	default:
		return false
	}
}

// Middleware для gzip-сжатия запросов и ответов
func GzipMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Оригинальный http.ResponseWriter
		ow := w

		// Проверяем, поддерживает ли клиент сжатие gzip
		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")

		// Проверяем Content-Type ответа
		contentType := r.Header.Get("Content-Type")
		shouldCompressResponse := shouldCompress(contentType)

		// Если клиент поддерживает gzip и Content-Type допустим для сжатия
		if supportsGzip && shouldCompressResponse {
			// Оборачиваем ResponseWriter для сжатия
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
