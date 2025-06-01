CREATE DOMAIN YoutubeVideoID AS TEXT NOT NULL CHECK (VALUE ~ '^[0-9A-Za-z_-]{10}[048AEIMQUYcgkosw]$');
CREATE DOMAIN YoutubeChannelID AS TEXT CHECK (VALUE ~ '^UC[0-9A-Za-z_-]{21}[AQgw]$');

CREATE TYPE ProjectType AS ENUM (
	'unknown',
	'other',
	'multi-animation',
	'multi-animation part',	
	'multi-edit', 
	'multi-edit part', 
	'animated music video',
	'picture music video',
	'animation meme'
);


CREATE TABLE "file" (
	"id" BIGINT NOT NULL UNIQUE GENERATED ALWAYS AS IDENTITY,
	"path" TEXT NOT NULL UNIQUE,
	"extension" TEXT NOT NULL CHECK (length(extension) >= 3 AND length(extension) <= 6),
	"md5" BYTEA NOT NULL UNIQUE CHECK (length(md5) = 16),
	"sha1" BYTEA NOT NULL UNIQUE CHECK (length(sha1) = 20),
	"sha256" BYTEA NOT NULL UNIQUE CHECK (length(sha256) = 32),
	"filesize" BIGINT NOT NULL CHECK (filesize >= 16),
	PRIMARY KEY("id")
);

CREATE TABLE "file_video" (
	"file_id" BIGINT NOT NULL UNIQUE,
	"duration" INTEGER NOT NULL CHECK (duration >= 0),
	"width" SMALLINT NOT NULL CHECK (width >= 0),
	"height" SMALLINT NOT NULL CHECK (height >= 0),
	"fps" SMALLINT CHECK (fps >= 0),
	"video_codec" TEXT,
	"audio_codec" TEXT,
	FOREIGN KEY ("file_id") REFERENCES file("id")
	ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE "youtube_channel" (
	"id" YoutubeChannelID NOT NULL UNIQUE DEFAULT ('UC000000000000000000000A'),
	PRIMARY KEY("id")
);


CREATE RULE youtube_channel_protect_unknown_id as
ON DELETE TO youtube_channel
WHERE OLD.id = 'UC000000000000000000000A'
DO INSTEAD nothing;


CREATE TABLE "youtube_channel_uploader_id" (
	"channel_id" YoutubeChannelID NOT NULL DEFAULT ('UC000000000000000000000A'),
	"uploader_id" TEXT NOT NULL,
	UNIQUE ("channel_id", "uploader_id"),
	PRIMARY KEY("channel_id", "uploader_id"),
	FOREIGN KEY ("channel_id") REFERENCES youtube_channel("id")
	ON UPDATE CASCADE ON DELETE CASCADE 
);

CREATE RULE youtube_channel_uploader_id_protect_unknown_id as
ON DELETE TO youtube_channel_uploader_id
WHERE OLD.channel_id = 'UC000000000000000000000A'
DO INSTEAD nothing;


CREATE TABLE "youtube_channel_uploader_name" (
	"channel_id" YoutubeChannelID NOT NULL DEFAULT ('UC000000000000000000000A'),
	"uploader" TEXT NOT NULL,
	UNIQUE ("channel_id", "uploader"),
	PRIMARY KEY("channel_id", "uploader"),
	FOREIGN KEY ("channel_id") REFERENCES youtube_channel("id")
	ON UPDATE CASCADE ON DELETE CASCADE 
);

CREATE RULE youtube_channel_uploader_name_protect_unknown_id as
ON DELETE TO youtube_channel_uploader_name
WHERE OLD.channel_id = 'UC000000000000000000000A'
DO INSTEAD nothing;


CREATE TABLE "youtube_video" (
	"id" YoutubeVideoID NOT NULL UNIQUE DEFAULT ('00000000000'),
	"upload_date" TIMESTAMP NOT NULL,
	"duration" INTEGER NOT NULL CHECK (duration >= 0),
	"view_count" INTEGER CHECK (view_count >= 0),
	"like_count" INTEGER CHECK (like_count >= 0),
	"dislike_count" INTEGER CHECK (dislike_count >= 0),
	"is_live" BOOLEAN,
	"is_restricted" BOOLEAN,
	PRIMARY KEY ("id")
);

CREATE TABLE "youtube_title" (
	"youtube_id" YoutubeVideoID NOT NULL,
	"title" TEXT NOT NULL,
	"title_md5" BYTEA NOT NULL UNIQUE CHECK (length(title_md5) = 16),
	"date_added" TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'), 
	UNIQUE("youtube_id", "title_md5"),
	PRIMARY KEY ("youtube_id", "title_md5"),
	FOREIGN KEY ("youtube_id") REFERENCES "youtube_video"("id")
	ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE "youtube_description" (
	"youtube_id" YoutubeVideoID NOT NULL,
	"description" TEXT NOT NULL,
	"description_md5" BYTEA NOT NULL UNIQUE CHECK (length(description_md5) = 16),
	"date_added" TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'), 
	UNIQUE("youtube_id", "description_md5"),
	PRIMARY KEY ("youtube_id", "description_md5"),
	FOREIGN KEY ("youtube_id") REFERENCES "youtube_video"("id")
	ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE "youtube_channel_youtube_video" (
	"channel_id" YoutubeChannelID NOT NULL DEFAULT ('UC000000000000000000000A'),
	"youtube_id" YoutubeVideoID NOT NULL UNIQUE DEFAULT ('00000000000'),
	UNIQUE ("channel_id", "youtube_id"),
	PRIMARY KEY("channel_id", "youtube_id"),
	FOREIGN KEY ("channel_id") REFERENCES "youtube_channel"("id")
	ON UPDATE CASCADE ON DELETE CASCADE,
	FOREIGN KEY ("youtube_id") REFERENCES "youtube_video"("id")
	ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE "youtube_file" (
	"youtube_id" YoutubeVideoID NOT NULL,
	"file_id" BIGINT NOT NULL UNIQUE,
	UNIQUE("youtube_id", "file_id"),
	PRIMARY KEY ("youtube_id", "file_id"),
	FOREIGN KEY ("youtube_id") REFERENCES youtube_video("id")
	ON UPDATE CASCADE ON DELETE CASCADE,
	FOREIGN KEY ("file_id") REFERENCES file("id")
	ON UPDATE CASCADE ON DELETE CASCADE
);


CREATE TABLE "youtube_video_format" (
	"youtube_id" YoutubeVideoID NOT NULL,
	"file_id" BIGINT NOT NULL UNIQUE,
	"format_id" TEXT NOT NULL,
	"format" TEXT NOT NULL,
	UNIQUE("youtube_id", "format_id"),
	PRIMARY KEY("youtube_id", "format_id"),
	FOREIGN KEY ("youtube_id") REFERENCES "youtube_video"("id")
	ON UPDATE CASCADE ON DELETE CASCADE,
	FOREIGN KEY ("file_id") REFERENCES "file"("id")
	ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE "youtube_video_ytdlp_version" (
	"file_id" BIGINT NOT NULL UNIQUE,
	"youtube_id" YoutubeVideoID NOT NULL,
	"repository" TEXT NOT NULL,
	"release_git_head" TEXT NOT NULL,
	"version" TEXT NOT NULL,
	UNIQUE("file_id", "youtube_id"),
	PRIMARY KEY("file_id", "youtube_id"),
	FOREIGN KEY("file_id") REFERENCES "file"("id")
	ON UPDATE CASCADE ON DELETE CASCADE,
	FOREIGN KEY ("youtube_id") REFERENCES "youtube_video"("id")
	ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE "project" (
	"id" BIGINT NOT NULL UNIQUE GENERATED ALWAYS AS IDENTITY,
	"uuid" TEXT NOT NULL UNIQUE CHECK (uuid != ''),
	"type" ProjectType NOT NULL DEFAULT 'unknown',
	"date_announced" TIMESTAMP,
	"date_completed" TIMESTAMP,
	"date_archived" TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
	PRIMARY KEY("id")
);

CREATE TABLE "project_title" (
	"project_id" BIGINT NOT NULL,
	"title" TEXT NOT NULL,
	"title_md5" BYTEA NOT NULL UNIQUE CHECK (length(title_md5) = 16),
	"date_added" TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'), 
	UNIQUE("project_id", "title_md5"),
	PRIMARY KEY ("project_id", "title_md5"),
	FOREIGN KEY ("project_id") REFERENCES "project"("id")
	ON UPDATE CASCADE ON DELETE CASCADE
);


CREATE TABLE "project_description" (
	"project_id" BIGINT NOT NULL UNIQUE GENERATED ALWAYS AS IDENTITY,
	"description" TEXT NOT NULL,
	"description_md5" BYTEA NOT NULL UNIQUE CHECK (length(description_md5) = 16),
	"date_added" TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'), 
	UNIQUE("project_id", "description_md5"),
	PRIMARY KEY ("project_id", "description_md5"),
	FOREIGN KEY ("project_id") REFERENCES "project"("id")
	ON UPDATE CASCADE ON DELETE CASCADE
);


CREATE TABLE "project_file" (
	"project_id" BIGINT NOT NULL,
	"file_id"  BIGINT NOT NULL,
	UNIQUE("project_id", "file_id"),
	PRIMARY KEY("project_id", "file_id"),
	FOREIGN KEY ("project_id") REFERENCES "project"("id")
	ON UPDATE CASCADE ON DELETE CASCADE,
	FOREIGN KEY ("file_id") REFERENCES "file"("id")
	ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE "project_participant" (
	"project_id" INTEGER NOT NULL,
	"file_id" BIGINT NOT NULL,
	UNIQUE("project_id", "file_id"),
	PRIMARY KEY("project_id", "file_id"),
	FOREIGN KEY ("project_id") REFERENCES "project"("id")
	ON UPDATE CASCADE ON DELETE CASCADE,
	FOREIGN KEY ("file_id") REFERENCES "file"("id")
	ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE "music" (
	"id" INTEGER NOT NULL UNIQUE GENERATED ALWAYS AS IDENTITY,
	"artist" citext NOT NULL,
	"title" citext NOT NULL,
	UNIQUE("artist", "title"),
	PRIMARY KEY("id")
);

CREATE TABLE "project_music" (
	"project_id" INTEGER NOT NULL,
	"music_id" INTEGER NOT NULL,
	UNIQUE ("project_id", "music_id"),
	PRIMARY KEY("project_id", "music_id"),
	FOREIGN KEY ("project_id") REFERENCES "project"("id")
	ON UPDATE CASCADE ON DELETE CASCADE,
	FOREIGN KEY ("music_id") REFERENCES "music"("id")
	ON UPDATE CASCADE ON DELETE CASCADE
);


CREATE TABLE "character" (
	"id" INTEGER NOT NULL UNIQUE GENERATED ALWAYS AS IDENTITY,
	"name" TEXT NOT NULL,
	"series" TEXT NOT NULL,
	"is_original" BOOLEAN NOT NULL DEFAULT (false),
	UNIQUE("name", "series", "is_original"),
	PRIMARY KEY("id")
);


CREATE TABLE "project_character" (
	"project_id" INTEGER NOT NULL,
	"character_id" INTEGER NOT NULL,
	UNIQUE("project_id", "character_id"),
	PRIMARY KEY("project_id", "character_id"),
	FOREIGN KEY ("project_id") REFERENCES "project"("id")
	ON UPDATE CASCADE ON DELETE CASCADE,
	FOREIGN KEY ("character_id") REFERENCES "character"("id")
	ON UPDATE CASCADE ON DELETE CASCADE
);