package example

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type Payload struct {
	Number int `json:"number"`
}

func PanicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func MustReadCompressedBody[T any](r io.Reader) *T {
	gr, err := gzip.NewReader(r)
	PanicIfErr(err)
	defer gr.Close()

	var t T
	PanicIfErr(json.NewDecoder(gr).Decode(&t))
	return &t
}

func MustWriteCompressedResponse(w http.ResponseWriter, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Encoding", "gzip")

	gw := gzip.NewWriter(w)
	defer gw.Close()
	PanicIfErr(json.NewEncoder(gw).Encode(body))
}

func MustCreateCompressedRequest(ctx context.Context, method, url string, body any) *http.Request {
	pr, pw := io.Pipe()

	go func() {
		gw := gzip.NewWriter(pw)
		err := json.NewEncoder(gw).Encode(body)
		defer PanicIfErr(gw.Close())
		defer pw.CloseWithError(err)
	}()

	r, err := http.NewRequestWithContext(ctx, method, url, pr)
	PanicIfErr(err)

	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Content-Encoding", "gzip")

	// if this header is set
	// the response won't be decompressed automatically (resp.Uncompressed == false)
	// hence you need to use example.MustReadCompressedBody[example.Payload](resp.Body)
	// ref: https://github.com/golang/go/blob/master/src/net/http/response.go#L89-L96
	// ref: https://github.com/golang/go/blob/master/src/net/http/transport.go#L182-L190
	//r.Header.Set("Accept-Encoding", "gzip")

	return r
}

//func MustCreateCompressedRequest(ctx context.Context, method, url string, body any) *http.Request {
//	var b bytes.Buffer
//
//	gw := gzip.NewWriter(&b)
//	err := json.NewEncoder(gw).Encode(body)
//	defer PanicIfErr(gw.Close())
//
//	r, err := http.NewRequestWithContext(ctx, method, url, &b)
//	PanicIfErr(err)
//
//	r.Header.Set("Content-Type", "application/json")
//	r.Header.Set("Content-Encoding", "gzip")
//
//	return r
//}
