CREATE VIEW threads_with_user AS
SELECT
    thread_id,
    title,
    created_at,
    users.user_id,
    users.name AS user_name
FROM threads
NATURAL JOIN post_threads
NATURAL JOIN users
ORDER BY created_at DESC
;

CREATE VIEW threads_with_user_category AS
SELECT
    thread_id,
    title,
    created_at,
    user_id,
    user_name,
    categories.category_id,
    categories.name AS category_name
FROM threads_with_user
NATURAL JOIN link_categories
NATURAL JOIN categories
;

CREATE VIEW tag_with_thread_id AS
SELECT
    tag_id,
    name,
    thread_id
FROM tags
NATURAL JOIN add_tags
;

CREATE VIEW num_comments AS
SELECT
    thread_id,
    COUNT(comment_id) AS num_comment
FROM comments
GROUP BY thread_id
;

CREATE VIEW comments_with_user_thread AS
SELECT
    comments.content,
    comments.comment_id,
    comments.thread_id,
    threads.title AS thread_title,
    post_comments.created_at,
    post_comments.user_id,
    users.name AS user_name
FROM comments
JOIN post_comments
    ON comments.thread_id = post_comments.thread_id
    AND comments.comment_id = post_comments.comment_id
JOIN threads ON comments.thread_id = threads.thread_id
JOIN users ON post_comments.user_id = users.user_id
ORDER BY created_at DESC
;
