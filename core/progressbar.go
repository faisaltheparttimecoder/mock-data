package core

import (
	"time"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
)

// Progress bar for the app.
var (
	bar *mpb.Bar
	p *mpb.Progress
)
func ProgressBar(steps int, progressMsg string) {

	// Start a new bar
	p = mpb.New(
		mpb.WithWidth(100),
		mpb.WithRefreshRate(120*time.Millisecond),
	)

	// Total steps to take and the message of this bar
	total := steps
	name := "  " + progressMsg

	// Add a bar
	bar = p.AddBar(int64(total),

		// Prepending decorators
		mpb.PrependDecorators(
			decor.Elapsed(4, decor.DSyncSpace),
		),

		// Appending decorators
		mpb.AppendDecorators(
			decor.Percentage(5, 0),
			decor.StaticName(name, len(name), 0),
		),
	)
}

// Increment Progress bar
func IncrementBar() {
	bar.Incr(1)
}


// Close progress bar
func CloseProgressBar() {
	p.Stop()
}
