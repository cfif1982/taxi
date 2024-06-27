package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {

	return w.Writer.Write(b)
}

// func GzipCompressMiddleware(h http.Handler) http.Handler {
func GzipCompressMiddleware(h http.Handler) http.Handler {
	gzipFn := func(rw http.ResponseWriter, req *http.Request) {

		// проверяем, что клиент поддерживает gzip-сжатие
		if !strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") {
			// если gzip не поддерживается, передаём управление дальше без изменений
			h.ServeHTTP(rw, req)
			return
		}

		// создаём gzip.Writer, который будет писать данные в rw
		gz, err := gzip.NewWriterLevel(rw, gzip.BestSpeed)
		if err != nil {
			io.WriteString(rw, err.Error())
			return
		}
		defer gz.Close()

		// устанавливаем заголовок
		rw.Header().Set("Content-Encoding", "gzip")

		// передаём обработчику страницы переменную типа gzipWriter для вывода данных
		h.ServeHTTP(gzipWriter{
			ResponseWriter: rw,
			Writer:         gz,
		}, req)

	}

	return http.HandlerFunc(gzipFn)
}
