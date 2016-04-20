package static

// The options for configuring the behavior of the Build functions. Get the default options by calling DefaultOptions().
type Options struct {
	// The directory where files will be written when building, default "build".
	OutputDir string
	// The number of files that will be built concurrently, default 50.
	Concurrency int
	// The filename to save director paths to, default "index.html".
	DirFilename string
}

// The default Options.
func DefaultOptions() Options {
	return Options{
		OutputDir:   "build",
		Concurrency: 50,
		DirFilename: "index.html",
	}
}
