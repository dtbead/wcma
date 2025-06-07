package mock_file

import (
	"github.com/dtbead/wc-maps-archive/internal/entities"
	"github.com/dtbead/wc-maps-archive/internal/helper"
	"github.com/dtbead/wc-maps-archive/internal/helper/file"
)

// RandomFile returns a randomly generated file entity.
func RandomFile() entities.File {
	hashes := helper.RandomEntitiesHash()
	ext := helper.RandomFileExtension()
	path := file.BuildPath(hashes.SHA256, ext)

	return entities.File{
		PathRelative: path,
		Extension:    ext,
		Size:         int64(helper.RandomInt(16, 65565)),
		Hashes:       hashes,
	}
}
