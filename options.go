package static

// The options for configuring the behavior of the Build functions. Get the default options with DefaultOptions.
type Options struct {
	// The directory where files will be written when building.
	OutputDir string
	// The number of files that will be built concurrently.
	Concurrency int
	// The filename to save director paths to.
	DirFilename string
}

// The default recommended Options
var DefaultOptions = Options{
	OutputDir:   "build",
	Concurrency: 50,
	DirFilename: "index.html",
}
