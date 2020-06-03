package tidio

import (
	"os"

	"github.com/gregoryv/fox"
)

var warn = fox.NewSyncLog(os.Stdout).FilterEmpty().Log
