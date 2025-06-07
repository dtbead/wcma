package helper

import (
	"crypto/md5"
	"crypto/rand"
	"errors"
	"fmt"
	randv2 "math/rand/v2"
	"time"

	"github.com/dtbead/wc-maps-archive/internal/entities"
	"github.com/google/uuid"
	"github.com/lithammer/shortuuid"
)

func init() {
	uuid.EnableRandPool()
}

func RandomInt(min, max int) int {
	return randv2.IntN(max-min) + min
}

func RandomString(length int) string {
	b := make([]byte, length+2)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[2 : length+2]
}

func RandomHash(length int) []byte {
	b := make([]byte, length)
	rand.Read(b)
	return b
}

func RandomUUID() string {
	return shortuuid.New()
}

func RandomEntitiesHash() entities.Hashes {
	return entities.Hashes{
		SHA256: RandomHash(32),
		SHA1:   RandomHash(20),
		MD5:    RandomHash(16),
	}
}

func RandomRune(list []rune, amount int) []rune {
	r := make([]rune, 0, amount)

	for range amount {
		r = append(r, list[RandomInt(0, len(list))])
	}

	return r
}

func RandomYoutubeID() entities.YoutubeVideoID {
	// [0-9A-Za-z_-]{10}[048AEIMQUYcgkosw] regex
	var part1 = []rune{
		'-', '_', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L',
		'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X',
		'Y', 'Z', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j',
		'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v',
		'w', 'x', 'y', 'z',
	}

	var part2 = []rune{
		'0', '4', '8', 'A', 'E', 'I', 'M', 'Q', 'U', 'Y', 'c', 'g',
		'k', 'o', 's', 'w',
	}

	youtube_id := entities.YoutubeVideoID(string(RandomRune(part1, 10)) + string(RandomRune(part2, 1)))
	if !youtube_id.IsValid() {
		panic("unexpectedly generated a random invalid youtube_id")
	}

	return entities.YoutubeVideoID(youtube_id)
}

func RandomYoutubeChannelID() entities.YoutubeChannelID {
	// UC[0-9A-Za-z_-]{21}[AQgw] regex
	var part1 = []rune{
		'-', '_', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L',
		'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X',
		'Y', 'Z', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j',
		'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v',
		'w', 'x', 'y', 'z',
	}

	var part2 = []rune{
		'A', 'Q', 'g', 'w',
	}

	channel_id := entities.YoutubeChannelID("UC" + string(RandomRune(part1, 21)) + string(RandomRune(part2, 1)))
	if !channel_id.IsValid() {
		panic("unexpectedly generated a random invalid channel_id")
	}
	return entities.YoutubeChannelID(channel_id)
}

func RandomFileExtension() string {
	var ext = []string{
		"mp4",
		"webm",
		"mkv",
	}

	return ext[RandomInt(0, len(ext))]
}

func RandomBool() bool {
	return randv2.IntN(2) == 1
}

func UnixEpochToTime(epoch int) time.Time {
	return time.Unix(int64(epoch), 0)
}

func RandomTime(start, end time.Time) time.Time {
	return UnixEpochToTime(RandomInt(int(start.Unix()), int(end.Unix())))
}

func RandomSlice(s []any) (any, error) {
	if s == nil || len(s) < 1 {
		return nil, errors.New("no slice given")
	}

	r := RandomInt(0, len(s))
	return s[r], nil
}

func GetMD5HashFromString(text string) []byte {
	hash := md5.Sum([]byte(text))
	return hash[:]
}
