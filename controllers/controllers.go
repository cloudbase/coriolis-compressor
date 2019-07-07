// Copyright 2019 Cloudbase Solutions Srl
// All Rights Reserved.

package controllers

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"

	gzip "github.com/klauspost/pgzip"
)

// CompressorHandler echoes back a compressed version of whatever binary data
// is posted
func CompressorHandler(writer http.ResponseWriter, req *http.Request) {
	var buf bytes.Buffer
	tmpWr := bufio.NewWriter(&buf)
	var wr io.WriteCloser
	switch comp := req.Header.Get("X-Compression-Format"); comp {
	case "gzip":
		gz := gzip.NewWriter(tmpWr)
		defer gz.Close()
		wr = gz
	case "zlib", "":
		zl := gzip.NewWriterZlib(tmpWr)
		defer zl.Close()
		wr = zl
	default:
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, fmt.Sprintf("X-Compression-Format invalid: %s", comp))
		return
	}
	if _, err := io.Copy(wr, req.Body); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, fmt.Sprintf("Got error: %q", err))
		return
	}
	// Close the gzip/zlib writer
	if err := wr.Close(); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, fmt.Sprintf("Got error: %q", err))
		return
	}
	// flush the bufio writer
	if err := tmpWr.Flush(); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, fmt.Sprintf("Got error: %q", err))
		return
	}
	reader := bufio.NewReader(&buf)
	if _, err := io.Copy(writer, reader); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, fmt.Sprintf("Got error: %q", err))
		return
	}
	return
}
