package static

type Options struct {
	// The directory where files will be written when building, default "build".
	OutputDir string
	// The number of files that will be built concurrently, default 50.
	Concurrency int
	// The filename to save director paths to, default "index.html".
	DirFilename string
}

func NewOptions() Options {
	return Options{
		OutputDir:   "build",
		Concurrency: 50,
		DirFilename: "index.html",
	}
}
