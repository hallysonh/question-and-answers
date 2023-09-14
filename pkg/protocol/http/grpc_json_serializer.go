package http

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type GrpcJSONSerializer struct{}

// Serialize converts an interface into a json and writes it to the response.
// You can optionally use the indent parameter to produce pretty JSONs.
func (d GrpcJSONSerializer) Serialize(c echo.Context, i interface{}, indent string) error {
	if value, ok := i.(protoreflect.ProtoMessage); ok {
		enc := &protojson.MarshalOptions{Multiline: indent != ""}
		data, err := enc.Marshal(value)
		if err != nil {
			return err
		}
		_, err = c.Response().Write(data)
		return err
	}

	enc := json.NewEncoder(c.Response())
	if indent != "" {
		enc.SetIndent("", indent)
	}
	return enc.Encode(i)
}

// Deserialize reads a JSON from a request body and converts it into an interface.
func (d GrpcJSONSerializer) Deserialize(c echo.Context, i interface{}) error {
	if value, ok := i.(protoreflect.ProtoMessage); ok {
		data, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return err
		}
		if err := protojson.Unmarshal(data, value); err != nil {
			slog.ErrorContext(c.Request().Context(), err.Error())
			return echo.NewHTTPError(http.StatusBadRequest, "Unmarshal GRPC message error")
		}
		return err
	}

	err := json.NewDecoder(c.Request().Body).Decode(i)
	if ute, ok := err.(*json.UnmarshalTypeError); ok {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Unmarshal type error: expected=%v, got=%v, field=%v, offset=%v", ute.Type, ute.Value, ute.Field, ute.Offset)).SetInternal(err)
	} else if se, ok := err.(*json.SyntaxError); ok {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Syntax error: offset=%v, error=%v", se.Offset, se.Error())).SetInternal(err)
	}
	return err
}
