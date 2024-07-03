CREATE TABLE posts (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    user_id UUID NOT NULL,
    allow_comments BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL
);
CREATE TABLE comments (
    id UUID PRIMARY KEY,
    post_id UUID NOT NULL REFERENCES posts(id),
    parent_id UUID,
    content TEXT NOT NULL,
    user_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL
);
CREATE TABLE structure_tree (
    ancestor_id UUID NOT NULL,
    descendant_id UUID NOT NULL,
    nearest_ancestor_id UUID NOT NULL,
    level INT NOT NULL,
    subject_id UUID NOT NULL,
    PRIMARY KEY (ancestor_id, descendant_id)
);
