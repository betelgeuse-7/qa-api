CREATE TABLE users (
    user_id serial primary key,
    username varchar(255) unique not null, 
    email varchar(255) unique not null,
    password text not null,
    handle varchar(255) unique not null,
    last_online timestamp with time zone, 
    created_at timestamp with time zone default CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);

CREATE TABLE questions (
    question_id serial primary key,
    title varchar(500) not null,
    text text not null,
    question_by int references users(user_id),
    created_at timestamp with time zone default CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);

CREATE TABLE tags (
    tag_id serial primary key,
    tag varchar(99) not null
);

CREATE TABLE question_tags (
    question_id int references questions(question_id),
    tag_id int references tags(tag_id),

    PRIMARY KEY(question_id, tag_id)
);

CREATE TABLE answers (
    answer_id serial primary key,
    answer_by int references users(user_id),
    text text not null,
    to_question int references questions(question_id),
    created_at timestamp with time zone default CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);

CREATE TABLE question_upvotes (
    question_id int references questions(question_id),
    upvote_by int references users(user_id),

    PRIMARY KEY(question_id, upvote_by)
);

CREATE TABLE answer_upvotes (
    answer_id int references answers(answer_id),
    upvote_by int references users(user_id),

    PRIMARY KEY(answer_id, upvote_by)
);

CREATE TABLE question_downvotes (
    question_id int references questions(question_id),
    downvote_by int references users(user_id),

    PRIMARY KEY(question_id, downvote_by)
);

CREATE TABLE answer_downvotes (
    answer_id int references answers(answer_id),
    downvote_by int references users(user_id),

    PRIMARY KEY(answer_id, downvote_by)
);

CREATE TABLE comments_to_question (
    comment_id serial primary key,
    text text not null,
    to_question int references questions(question_id),
    comment_by int references users(user_id),
    created_at timestamp with time zone default CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);

CREATE TABLE comments_to_answer (
    comment_id serial primary key,
    text text not null,
    to_answer int references answers(answer_id),
    comment_by int references users(user_id),
    created_at timestamp with time zone default CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);