package common

import (
	"fmt"
	"io"

	stdErrors "errors"

	"github.com/google/uuid"
)

type ID uuid.UUID

var ErrNonStringUUID = stdErrors.New("uuid is not a string")

func (id ID) MarshalGQL(w io.Writer) {
	str := uuid.UUID(id).String()
	if _, err := w.Write([]byte(fmt.Sprintf("\"%s\"", str))); err != nil {
		fmt.Print(err)
	}
}

func (id *ID) UnmarshalGQL(v interface{}) error {
	switch v := v.(type) {
	case string:
		u, err := uuid.Parse(v)
		if err != nil {
			return err
		}

		*id = ID(u)

		return nil
	default:
		return fmt.Errorf("%w: %T", ErrNonStringUUID, v)
	}
}
