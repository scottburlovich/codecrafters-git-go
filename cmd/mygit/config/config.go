package config

// Git directory structure
const (
	GitDir       = ".git"
	ObjectsDir   = GitDir + "/objects"
	RefsDir      = GitDir + "/refs"
	HeadFilePath = GitDir + "/HEAD"
)

// Git object types
const (
	Blob = "blob"
	Tree = "tree"

	ModeBlob     = "100644"
	ModeTree     = "40000"
	ModeBlobExec = "100755"
	ModeSymLink  = "120000"
)
