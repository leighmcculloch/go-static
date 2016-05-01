package static

// Options for configuring the behavior of the Build functions. Get the default options with DefaultOptions.
type Options struct {
	// The directory where files will be written when building.
	OutputDir string
	// The number of files that will be built concurrently.
	Concurrency int
	// The filename to use when saving directory paths. e.g. index.html
	DirFilename string
}

// DefaultOptions contain the default recommended Options.
var DefaultOptions = Options{
	OutputDir:   "build",
	Concurrency: 50,
	DirFilename: "index.html",
}
