package request

import (
	"bytes"
	"fmt"
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/encoding/form"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/transport/http"
	"io"
)

type Request struct {
	Debug bool
}

func (c *Request) JSONDecoder(r *http.Request, v interface{}) error {
	codec, ok := http.CodecForRequest(r, "Content-Type")
	if !ok {
		return errors.BadRequest("CODEC", fmt.Sprintf("unregister Content-Type: %s", r.Header.Get("Content-Type")))
	}
	data, err := io.ReadAll(r.Body)

	// reset body.
	r.Body = io.NopCloser(bytes.NewBuffer(data))

	if c.Debug {
		println(r.URL.String(), " -> ", string(data))
	}

	if err != nil {
		return errors.BadRequest("CODEC", err.Error())
	}
	if len(data) == 0 {
		return nil
	}
	switch v.(type) {
	//	Notice Add type
	//case *pb.EquipmentMaintenanceOptRequest:
	//	r.Header.Set("raw_body", string(data))
	//	return nil
	}

	if err = codec.Unmarshal(data, v); err != nil {
		return errors.BadRequest("CODEC", fmt.Sprintf("body unmarshal %s", err.Error()))
	}
	return nil
}

func (c *Request) QueryDecoder(r *http.Request, target interface{}) error {
	if c.Debug {
		println(r.URL.String(), " -> ", string(r.URL.Query().Encode()))
	}

	if err := encoding.GetCodec(form.Name).Unmarshal([]byte(r.URL.Query().Encode()), target); err != nil {
		return errors.BadRequest("CODEC", err.Error())
	}
	return nil
}
