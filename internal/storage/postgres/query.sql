-- name: NewFile :one
INSERT INTO file (path, extension, md5, sha1, sha256, filesize) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;

-- name: DeleteFileByID :exec
DELETE FROM file WHERE id = $1;

-- name: GetFileByID :one
SELECT * FROM file WHERE id = $1;

-- name: NewProject :one
INSERT INTO project (uuid, type, date_announced, date_completed) VALUES ($1, $2, $3, $4) RETURNING uuid;

-- name: DeleteProjectByUUID :exec
DELETE FROM project WHERE uuid = $1;

-- name: GetProjectByUUID :one
SELECT * FROM project WHERE uuid = $1;

-- name: AssignProjectFile :exec
INSERT INTO project_file (project_id, file_id) VALUES ((SELECT id FROM project WHERE uuid = $1), $2);

-- name: UnassignProjectFile :exec
DELETE FROM project_file WHERE file_id = $1;

-- name: GetProjectFile :many
SELECT file_id FROM project_file WHERE project_id = (SELECT id FROM project WHERE uuid = $1);

-- name: NewYoutube :exec
INSERT INTO youtube_video (id, upload_date, duration, view_count, like_count, dislike_count, is_live, is_restricted)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: NewYoutubeChannelVideo :exec
INSERT INTO youtube_channel_youtube_video (channel_id, youtube_id) VALUES ($1, $2);

-- name: NewYoutubeChannel :exec
INSERT INTO youtube_channel (id) VALUES ($1) ON CONFLICT DO NOTHING;

-- name: NewYoutubeChannelUploaderName :exec
INSERT INTO youtube_channel_uploader_name (channel_id, uploader) VALUES ($1, $2) ON CONFLICT DO NOTHING;

-- name: NewYoutubeChannelUploaderID :exec
INSERT INTO youtube_channel_uploader_id (channel_id, uploader_id) VALUES ($1, $2) ON CONFLICT DO NOTHING;

-- name: NewYoutubeFormat :exec
INSERT INTO youtube_video_format (youtube_id, file_id, format_id, format)
VALUES ($1, $2, $3, $4);

-- name: GetYoutubeVideo :one
SELECT * FROM youtube_video WHERE id = $1;


-- name: GetYoutubeVideoFormatByYoutubeID :many
SELECT * FROM youtube_video_format WHERE youtube_id = $1;

-- name: GetYoutubeTitle :many
SELECT title FROM youtube_title WHERE youtube_id = $1;

-- name: GetYoutubeDescription :many
SELECT description FROM youtube_description WHERE youtube_id = $1;

-- name: AssignYoutubeTitle :exec
INSERT INTO youtube_title (youtube_id, title, title_md5) VALUES ($1, $2, $3);

-- name: AssignYoutubeDescription :exec
INSERT INTO youtube_description (youtube_id, description, description_md5) VALUES ($1, $2, $3);

-- name: AssignYoutubeFileID :exec
INSERT INTO youtube_file (youtube_id, file_id) VALUES ($1, $2);

-- name: GetYoutubeFileID :many
SELECT file_id FROM youtube_file WHERE youtube_id = $1;

-- name: GetYoutubeChannelByID :one
SELECT 
    youtube_channel_youtube_video.channel_id AS channel_id, 
    youtube_channel_uploader_id.uploader_id AS uploader_id,
    youtube_channel_uploader_name.uploader AS uploader_name 
FROM youtube_channel_youtube_video
    INNER JOIN youtube_channel_uploader_id ON youtube_channel_uploader_id.channel_id = youtube_channel_youtube_video.channel_id
    INNER JOIN youtube_channel_uploader_name ON youtube_channel_uploader_name.channel_id = youtube_channel_youtube_video.channel_id
WHERE youtube_channel_youtube_video.youtube_id = $1 LIMIT 1;

-- name: GetYoutubeYtdlpVersion :one
SELECT * FROM youtube_video_ytdlp_version WHERE youtube_id = $1 AND file_id = $2;

-- name: NewYoutubeYtdlpVersion :exec
INSERT INTO youtube_video_ytdlp_version ("file_id", "youtube_id", "repository", "release_git_head", "version") VALUES ($1, $2, $3, $4, $5);

-- name: GetProjectByYoutubeID :one
SELECT project.* FROM project 
INNER JOIN project_file ON project.id = project_file.project_id
INNER JOIN youtube_file ON project_file.file_id = youtube_file.file_id
WHERE youtube_file.youtube_id = $1;

-- name: GetProjectTypeByYoutubeID :one
SELECT project.type FROM project 
INNER JOIN project_file ON project.id = project_file.project_id
INNER JOIN youtube_file ON project_file.file_id = youtube_file.file_id
WHERE youtube_file.youtube_id = $1;