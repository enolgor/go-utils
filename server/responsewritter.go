package server

import (
	"io"
	"net/http"
)

type ResponseWriter interface {
	http.ResponseWriter
	Status() int
	Written() bool
	Size() int
}

func NewResponseWriter(rw http.ResponseWriter) ResponseWriter {
	nrw := &responseWriter{
		ResponseWriter: rw,
	}

	return nrw
}

type responseWriter struct {
	http.ResponseWriter
	pendingStatus int
	status        int
	size          int
}

func (rw *responseWriter) WriteHeader(s int) {
	if rw.Written() {
		return
	}

	rw.pendingStatus = s

	if rw.Written() {
		return
	}

	rw.status = s
	rw.ResponseWriter.WriteHeader(s)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.Written() {
		rw.WriteHeader(http.StatusOK)
	}
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

func (rw *responseWriter) ReadFrom(r io.Reader) (n int64, err error) {
	if !rw.Written() {
		rw.WriteHeader(http.StatusOK)
	}
	n, err = io.Copy(rw.ResponseWriter, r)
	rw.size += int(n)
	return
}

func (rw *responseWriter) Unwrap() http.ResponseWriter {
	return rw.ResponseWriter
}

func (rw *responseWriter) Status() int {
	if rw.Written() {
		return rw.status
	}

	return rw.pendingStatus
}

func (rw *responseWriter) Size() int {
	return rw.size
}

func (rw *responseWriter) Written() bool {
	return rw.status != 0
}
