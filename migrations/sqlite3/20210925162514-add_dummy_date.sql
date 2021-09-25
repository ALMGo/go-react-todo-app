
-- +migrate Up
INSERT INTO user (id, username, password)
VALUES
    ('ad2a2d', 'user1', '$2a$14$JJhy4G7Tlm0GtVXCoFElEeqmDkXgTQfsP7hxwcGzVAkZv67wrIKGe'),
    ('bsr5s5', 'user2', '$2a$14$O.zjBprfO.NxgCJovBS8JuZqekgyh.VwrZhkPF73b47QXQ.H0Lf9a'),
    ('z4t434', 'user3', '$2a$14$O.zjBprfO.NxgCJovBS8JuZqekgyh.VwrZhkPF73b47QXQ.H0Lf9a');

INSERT INTO todo_item (user_id, completed, text, due, category)
VALUES
    ('ad2a2d', 0, 'Clean Pool', '2021-11-11 10:00:00', 'Home'),
    ('ad2a2d', 0, 'Mow Lawn', '2021-11-12 10:00:00', 'Home'),
    ('ad2a2d', 0, 'Update Computer', '2021-11-12 10:00:00', 'Work'),
    ('ad2a2d', 0, 'Buy XRP', '2021-11-12 10:00:00', 'Investing'),

    ('bsr5s5', 0, 'Wash Windows', '2021-11-11 10:00:00', 'Home'),
    ('bsr5s5', 0, 'Bump dependencies versions', '2021-11-12 10:00:00', 'Work'),
    ('bsr5s5', 0, 'Go Kayaking', '2021-11-12 10:00:00', 'Recreation'),

    ('bsr5s5', 0, 'Wash Windows', '2021-11-11 10:00:00', 'Home'),
    ('bsr5s5', 0, 'Bump dependencies versions', '2021-11-12 10:00:00', 'Work'),
    ('bsr5s5', 0, 'Go Kayaking', '2021-11-12 10:00:00', 'Recreation');

-- +migrate Down
DELETE FROM user WHERE id in ('ad2a2d', 'bsr5s5', 'z4t434');
