package httplogger

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

type (
	// берём структуру для хранения сведений об ответе
	responseData struct {
		status int
		size   int
	}

	// добавляем реализацию http.ResponseWriter
	loggingResponseWriter struct {
		http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
		responseData        *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	// записываем ответ, используя оригинальный http.ResponseWriter
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size // захватываем размер
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	// записываем код статуса, используя оригинальный http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // захватываем код статуса
}

func LoggerMiddleware(logger *zerolog.Logger) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			respData := &responseData{
				status: 0,
				size:   0,
			}

			lw := loggingResponseWriter{
				ResponseWriter: w,
				responseData:   respData,
			}

			uri := r.RequestURI
			method := r.Method

			// если сервер выдаст ошибку, то в этом случае все равно выведет статус
			defer func() {
				duration := time.Since(start)
				logger.Debug().
					Str("uri", uri).
					Str("method", method).
					Int("status", respData.status).
					Dur("duration", duration).
					Int("size", respData.size).
					Send()
			}()

			next.ServeHTTP(&lw, r)
		})
	}
}
