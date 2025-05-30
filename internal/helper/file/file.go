package file

import (
	"bufio"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/dtbead/wc-maps-archive/internal/entities"
)

var hashPool sync.Pool

type hasher struct {
	buf    bufio.Reader
	md5    hash.Hash
	sha1   hash.Hash
	sha256 hash.Hash
}

func (h *hasher) Reset() {
	h.md5.Reset()
	h.sha1.Reset()
	h.sha256.Reset()
}

func init() {
	hashPool = sync.Pool{
		New: func() any {
			h := hasher{}
			h.md5 = md5.New()
			h.sha1 = sha1.New()
			h.sha256 = sha256.New()
			return &h
		},
	}
}

func getHashpool() *hasher {
	h := hashPool.Get()
	if h != nil {
		return h.(*hasher)
	}
	return nil
}

func putHashPool(h hasher) {
	h.Reset()
	hashPool.Put(&h)
}

type Hashes struct {
	MD5    []byte
	SHA1   []byte
	SHA256 []byte
}

func Copy(destination string, r io.Reader) error {
	rootDirectory := filepath.Dir(destination)
	if !DoesPathExist(rootDirectory) {
		err := os.MkdirAll(rootDirectory, 0664)
		if err != nil {
			return err
		}
	}

	file, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer file.Close()

	buf := bufio.NewReader(r)

	w, err := buf.WriteTo(file)
	if err != nil {
		return err
	}

	if w <= 0 {
		return errors.New("copied 0 bytes")
	}

	return nil
}

// GetHash returns the MD5, SHA1, SHA256 hash and total bytes read of a given io.Reader.
func GetHash(r io.Reader) (hashes entities.Hashes, read int64, err error) {
	hashPool := getHashpool()
	hashPool.buf.Reset(r)

	mw := io.MultiWriter(hashPool.md5, hashPool.sha1, hashPool.sha256)
	bytesRead, err := io.Copy(mw, &hashPool.buf)
	if err != nil {
		return entities.Hashes{}, bytesRead, err
	}

	h := entities.Hashes{
		MD5:    hashPool.md5.Sum(nil),
		SHA1:   hashPool.sha1.Sum(nil),
		SHA256: hashPool.sha256.Sum(nil),
	}

	putHashPool(*hashPool)
	return h, bytesRead, nil
}

// BuildPath builds a path to store media. md5 gets encoded to a hexidecimal string
// to create a storage path such as "f15/f15f38b5cfdbfd56aeb6da48b65d3d6f.png".
// BuildPath expects an extension to have a period prefix already added by caller
func BuildPath(hash []byte, extension string) string {
	return fmt.Sprintf("%s/%s.%s", string(ByteToHexString(hash[:1])), string(ByteToHexString(hash[:])), extension)
}

func SanitizePath(s string) string {
	return strings.TrimSuffix(strings.ReplaceAll(path.Clean(s), "\\", "/"), "/")
}

func DoesPathExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}

	return false
}

func GetSize(r io.Reader) (bytes int64, err error) {
	f, ok := r.(*os.File)
	if ok {
		st, err := f.Stat()
		if err != nil {
			return -1, err
		}

		return st.Size(), nil
	}

	n, err := io.Copy(io.Discard, r)
	if err != nil {
		return -1, err
	}

	return n, nil
}

// resetFileSeek checks whether a given io.Reader is of *os.File
// and resets the file pointer for future read/write ops.
func ResetFileSeek(r io.Reader) {
	f, ok := r.(*os.File)
	if ok {
		f.Seek(0, io.SeekStart)
	}
}

func ByteToHexString(h []byte) string {
	return hex.EncodeToString(h)
}

func HexStringToByte(s string) []byte {
	h, err := hex.DecodeString(s)
	if err != nil {
		return []byte{}
	}

	return h
}
