package sl

import (
	"fmt"
	"log/slog"
)

func Err(err error) slog.Attr {
	return slog.String("error", fmt.Sprintf("%v", err))
}
