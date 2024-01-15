package internal

import "errors"

var (
	// ErrFileNotFound is the error for when the file is not found
	ErrFileNotFound = errors.New("file not found")
	// ErrFileNotWritable is the error for when the file is not writable
	ErrFileNotWritable = errors.New("file not writable")
	// ErrFileNotReadable is the error for when the file is not readable
	ErrFileNotReadable = errors.New("file not readable")
	// ErrOpenFile is the error for when the file cannot be opened
	ErrOpenFile = errors.New("error opening the file")
	// ErrReadFile is the error for when the file cannot be read
	ErrReadFile = errors.New("error reading the file")
	// ErrWriteFile is the error for when the file cannot be written
	ErrWriteFile = errors.New("error writing the file")
	// ErrCreateFile is the error for when the file cannot be created
	ErrCreateFile = errors.New("error creating the file")
)

// ProductStorage is the interface for the product storage
type ProductStorage interface {
	// Store stores the map in the storage
	Write(data []byte) (err error)
	// Get gets the map from the storage
	Read() (data []byte, err error)
}