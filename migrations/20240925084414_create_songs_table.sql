-- +goose Up
CREATE TABLE IF NOT EXISTS songs (
    id SERIAL PRIMARY KEY,
    group_name VARCHAR(255) NOT NULL,
    song_name VARCHAR(255) NOT NULL,
    release_date DATE,
    text TEXT,
    link VARCHAR(255)
);

-- Insert example data
INSERT INTO songs (group_name, song_name, release_date, text, link) VALUES
('Muse', 'Supermassive Black Hole', '2006-07-16', 'Ooh baby, dont you know I suffer?', 'https://www.youtube.com/watch?v=Xsp3_a-PMTw'),
('Radiohead', 'Creep', '1992-09-21', 'When you were here before, couldnt look you in the eye', 'https://www.youtube.com/watch?v=XFkzRNyygfk'),
('Rammstein', 'Du Hast', '1997-07-18', 'Du, du hast, du hast mich', 'https://www.youtube.com/watch?v=W3q8Od5qJio'),
('Rammstein', 'Sonne', '2001-07-13', 'Hier kommt die Sonne', 'https://www.youtube.com/watch?v=StZcUAPRRac'),
('Pink Floyd', 'Comfortably Numb', '1979-11-30', 'Hello, is there anybody in there?', 'https://www.youtube.com/watch?v=_FrOQC-zEog'),
('Nirvana', 'Smells Like Teen Spirit', '1991-09-10', 'Load up on guns, bring your friends', 'https://www.youtube.com/watch?v=hTWKbfoikeg');

-- +goose Down
DROP TABLE IF EXISTS songs;
