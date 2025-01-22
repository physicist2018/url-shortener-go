package compressor

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"strings"
)

// Функция для проверки, сжаты ли данные или имеют тип application/x-gzip
func isGzipCompressed(contentType string, contentEncoding string) bool {
	return contentEncoding == "gzip" || contentType == "application/x-gzip"
}

// Функция для проверки, поддерживает ли клиент сжатие
func clientSupportsGzip(acceptEncoding string) bool {
	return strings.Contains(acceptEncoding, "gzip")
}

func GZipMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w

		log.Println("accessContentTypes")
		acceptEncoding := r.Header.Get("Accept-Encoding")
		isSupportGZip := strings.Contains(acceptEncoding, "gzip")
		if isSupportGZip {
			log.Println("Accept-Encoding run")

			cw := newCompressWriter(w)
			ow = cw
			defer cw.Close()
		}

		contentEncoding := r.Header.Get("Content-Encoding")
		isSendGZip := strings.Contains(contentEncoding, "gzip")
		if isSendGZip {
			log.Println("Content-Encoding run")

			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			r.Body = cr
			defer cr.Close()
		}

		h.ServeHTTP(ow, r)
	})
}

// Middleware для обработки сжатых запросов и сжатия ответов
func GzipMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получаем заголовки запроса
		contentEncoding := r.Header.Get("Content-Encoding")
		contentType := r.Header.Get("Content-Type")
		acceptEncoding := r.Header.Get("Accept-Encoding")

		// Если запрос сжат или имеет тип application/x-gzip, декомпрессируем его
		if isGzipCompressed(contentType, contentEncoding) {
			// Декомпрессируем тело запроса
			cr, err := newCompressReader(r.Body)
			if err != nil {
				http.Error(w, "Failed to decompress request body", http.StatusInternalServerError)
				return
			}
			// Заменяем тело запроса на декомпрессированное
			r.Body = cr
			defer cr.Close() // Закрываем декомпрессор после выполнения
		}

		// Если клиент поддерживает gzip и контент сжимаемый, сжимаем ответ
		if clientSupportsGzip(acceptEncoding) && shouldCompress(contentType) {
			// Устанавливаем Content-Encoding для сжатого ответа
			w.Header().Set("Content-Encoding", "gzip")
			// Оборачиваем ResponseWriter для сжатия
			cw := newCompressWriter(w)
			defer cw.Close() // Закрываем сжатие после выполнения

			// Передаем управление обработчику с сжатым ResponseWriter
			h.ServeHTTP(cw, r)
		} else {
			// Если не нужно сжимать, просто передаем управление обработчику
			h.ServeHTTP(w, r)
		}
	}
}

// Функция, которая проверяет, нужно ли сжимать данный контент (например, текстовые данные или JSON)
func shouldCompress(contentType string) bool {
	switch contentType {
	case "application/javascript", "application/json", "text/css", "text/html", "text/plain", "text/xml", "application/x-gzip":
		return true
	default:
		return false
	}
}

// Создает новый gzip.Writer для сжатия ответа
func newCompressWriter(w http.ResponseWriter) *gzipResponseWriter {
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
