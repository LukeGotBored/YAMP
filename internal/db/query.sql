-- name: CreateArtist :one
INSERT INTO artists (name) VALUES (?)
ON CONFLICT(name) DO UPDATE SET name=excluded.name
RETURNING *;

-- name: CreateAlbum :one
INSERT INTO albums (title, artist_id, year) VALUES (?, ?, ?)
ON CONFLICT(title, artist_id) DO UPDATE SET title=excluded.title
RETURNING *;

-- name: CreateTrack :one
INSERT INTO tracks (title, album_id, artist_id, track_number, duration, file_path) 
VALUES (?, ?, ?, ?, ?, ?)
ON CONFLICT(file_path) DO UPDATE SET 
    title=excluded.title, 
    album_id=excluded.album_id, 
    artist_id=excluded.artist_id, 
    track_number=excluded.track_number, 
    duration=excluded.duration
RETURNING *;

-- name: ListTracks :many
SELECT 
    tracks.*,
    COALESCE(artists.name, 'Unknown Artist') as artist_name,
    COALESCE(albums.title, 'Unknown Album') as album_title
FROM tracks 
LEFT JOIN artists ON tracks.artist_id = artists.id
LEFT JOIN albums ON tracks.album_id = albums.id
ORDER BY tracks.title;

-- name: ListArtists :many
SELECT * FROM artists ORDER BY name;

-- name: ListAlbums :many
SELECT * FROM albums ORDER BY year DESC, title ASC;
