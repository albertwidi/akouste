package archive

import "github.com/mholt/archiver"

// Unarchive unarchives the given archive file into the destination folder.
// The archive format is selected implicitly.
func Unarchive(source, destination string) error {
	return archiver.Unarchive(source, destination)
}

// Archive creates an archive of the source files to a new file at destination.
// The archive format is chosen implicitly by file extension.
func Archive(sources []string, destination string) error {
	return archiver.Archive(sources, destination)
}
