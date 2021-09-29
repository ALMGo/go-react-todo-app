
-- +migrate Up

-- Password is secret for each user
INSERT INTO user (username, password)
VALUES
    ('user1', '$2a$14$v3pBPjTnFnkOwCvxDoGVQ.EwKc9A9NG9cUc4O0j3lwhi4WEg2wmVW'),
    ('user2', '$2a$14$0kChyxsiqxp8O5jW4Tztvex/Pf3FehS9A0Y7KfWKExik4mvNfyBsm'),
    ('user3', '$2a$14$xRmDmHKZ.JHfXAvG1m8Xq.4UjxW/zG7H6uEOrilyQQxaT0VK/lXCG');

INSERT INTO todo_item (user_id, text, due, category)
VALUES
    (1, 'Clean Pool', '2021-11-11 10:00:00', 'Home'),
    (1, 'Mow Lawn', '2021-11-12 10:00:00', 'Home'),
    (1, 'Update Computer', '2021-11-12 10:00:00', 'Work'),
    (1, 'Buy XRP', '2021-11-12 10:00:00', 'Investing'),

    (2, 'Wash Windows', '2021-11-11 10:00:00', 'Home'),
    (2, 'Bump dependencies versions', '2021-11-12 10:00:00', 'Work'),
    (2, 'Go Kayaking', '2021-11-12 10:00:00', 'Recreation'),

    (3, 'Wash Windows', '2021-11-11 10:00:00', 'Home'),
    (3, 'Bump dependencies versions', '2021-11-12 10:00:00', 'Work'),
    (3, 'Go Kayaking', '2021-11-12 10:00:00', 'Recreation');

-- +migrate Down
DELETE FROM user WHERE id in (0, 1, 2);
