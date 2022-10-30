// Package conf TODO
package conf

import (
	"errors"
	"fmt"

	"github.com/jessevdk/go-flags"
)

// ErrHelp is returned when --help flag is
// used and application should not launch.
var ErrHelp = errors.New("help")

// New reads flags and envs and returns AppConfig
// that corresponds to the values read.
func New(config interface{}) error {
	if _, err := flags.Parse(config); err != nil {
		var flagsErr *flags.Error
		if errors.As(err, &flagsErr) && flagsErr.Type == flags.ErrHelp {
			return ErrHelp
		}
		return fmt.Errorf("failed to parse config: %w", err)
	}
	return nil
}
