package entities

import (
	"context"
	"errors"
	"io"
	"regexp"
	"time"
)

type FileID int64
type ProjectType int
type ProjectUUID string
type YoutubeVideoID string
type YoutubeChannelID string

const InvalidProjectUUID ProjectUUID = ""
const InvalidFileID FileID = -1
const UnknownYoutubeChannelID YoutubeChannelID = "UC000000000000000000000A"
const UnknownYoutubeID YoutubeVideoID = "00000000000"

const (
	ProjectTypeUnknown ProjectType = iota
	ProjectTypeOther
	ProjectMultiAnimation
	ProjectMultiAnimationPart
	ProjectMultiEdit
	ProjectMultiEditPart
	ProjectAnimatedMusicVideo
	ProjectPictureMusicVideo
	ProjectAnimationMeme
)

func (p ProjectType) ToString() string {
	switch p {
	case ProjectTypeOther:
		return "other"
	case ProjectMultiAnimation:
		return "multi-animation"
	case ProjectMultiAnimationPart:
		return "multi-animation-part"
	case ProjectMultiEdit:
		return "multi-edit"
	case ProjectMultiEditPart:
		return "multi-edit-part"
	case ProjectAnimatedMusicVideo:
		return "animated music video"
	case ProjectPictureMusicVideo:
		return "picture music video"
	case ProjectAnimationMeme:
		return "animation meme"
	default:
		return "unknown"
	}
}

func NewProjectType(s string) (ProjectType, error) {
	switch s {
	case "other":
		return ProjectTypeOther, nil
	case "multi-animation":
		return ProjectMultiAnimation, nil
	case "multi-animation-part":
		return ProjectMultiAnimationPart, nil
	case "multi-edit":
		return ProjectMultiEdit, nil
	case "multi-edit-part":
		return ProjectMultiEditPart, nil
	case "animated music video":
		return ProjectAnimatedMusicVideo, nil
	case "picture music video":
		return ProjectPictureMusicVideo, nil
	case "animation meme":
		return ProjectAnimationMeme, nil
	default:
		return ProjectTypeUnknown, errors.New("unknown project type")
	}
}

func (f FileID) IsValid() bool {
	return f > 0
}

func (y YoutubeVideoID) IsValid() bool {
	m, _ := regexp.MatchString(`^[0-9A-Za-z_-]{10}[048AEIMQUYcgkosw]$`, string(y))
	return m
}

func (y YoutubeChannelID) IsValid() bool {
	m, _ := regexp.MatchString(`^UC[0-9A-Za-z_-]{21}[AQgw]$`, string(y))
	return m
}

type Video struct {
	VideoCodec, AudioCodec string
	Duration               int
	Width, Height, Fps     int16
}

type Youtube struct {
	YouTube            YoutubeVideo
	Channel            *VideoYoutubeChannel
	Format             *VideoYoutubeFormat
	DlpVersion         *VideoYoutubeDlpVersion
	Title, Description string
}

type YoutubeVideo struct {
	YoutubeID               YoutubeVideoID
	Video                   Video
	UploadDate              time.Time
	Duration                int
	ViewCount               int
	LikeCount, DislikeCount int
	IsLive, IsRestricted    bool
}

type VideoYoutubeChannel struct {
	ChannelID  YoutubeChannelID
	UploaderID string
	Uploader   string
}

type VideoYoutubeFormat struct {
	YoutubeID        YoutubeVideoID
	FileID           FileID
	Format, FormatID string
}

type VideoYoutubeDlpVersion struct {
	YoutubeID                           YoutubeVideoID
	FileID                              FileID
	Version, ReleaseGitHead, Repository string
}

type ProjectTitle struct {
	ProjectUUID ProjectUUID
	Title       string
	project_id  int32
	title_md5   []byte
}

type ProjectDescription struct {
	ProjectUUID     ProjectUUID
	Description     string
	description_md5 []byte
	project_id      int32
}

type Project struct {
	UUID          string
	project_id    int32
	FileIDs       []FileID
	ProjectType   ProjectType
	DateAnnounced time.Time
	DateCompleted time.Time
	DateArchived  time.Time
}

type ProjectImport struct {
	ProjectType                                ProjectType
	DateAnnounced, DateCompleted, DateArchived time.Time
	FileIDs                                    []FileID
}

type ProjectYoutube struct {
	Project *Project
	Youtube *Youtube
}

type ProjectMusic struct {
	ProjectUUID ProjectUUID
	id          int32
	MusicID     int32
}

type Music struct {
	ID            int32
	Artist, Title string
}

type Character struct {
	ID           int
	Name, Series string
	IsOriginal   bool
}

type Hashes struct {
	SHA256, SHA1, MD5 []byte
}

type File struct {
	PathAbsolute, PathRelative string
	Extension                  string
	Size                       int64
	Hashes                     Hashes
}

type VideoImport struct {
	Video     io.Reader
	Extension string
	Duration  int
}

var (
	ErrorInvalidFilePtr          = errors.New("nil file pointer")
	ErrorInvalidDuration         = errors.New("invalid duration")
	ErrorInvalidFileSize         = errors.New("invalid filesize")
	ErrorInvalidExtensionName    = errors.New("invalid extension")
	ErrorInvalidYoutubeDuration  = errors.New("invalid youtube duration")
	ErrorInvalidYoutubeLikes     = errors.New("invalid youtube like_count")
	ErrorInvalidYoutubeDislikes  = errors.New("invalid youtube dislike_count")
	ErrorInvalidYoutubeID        = errors.New("invalid youtube video id")
	ErrorInvalidYoutubeChannelID = errors.New("invalid youtube channel id")
	ErrorVideoNotFound           = errors.New("video not found")
	ErrorInvalidVideoID          = errors.New("invalid video id")
	ErrorInvalidVideoPtr         = errors.New("nil video pointer")
	ErrorInvalidYoutubeVideoPtr  = errors.New("nil youtube video pointer")
	ErrorNotFound                = errors.New("not found")
)

type YoutubeDownloader interface {
	Download(ctx context.Context, url string, output io.Writer) (youtube *Youtube, extension string, err error)
}

type FileRelationship struct {
	FileID  FileID
	Youtube YoutubeVideoID
	Project ProjectUUID
}
