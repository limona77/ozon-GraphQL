package storage

import (
	"context"
	"ozon-GraphQL/graph/model"
	"ozon-GraphQL/internal/database"
	"time"
)

type PostgresSQLRepository struct {
	db database.Database
}

func NewPostgresSQLRepository(db database.Database) *PostgresSQLRepository {
	return &PostgresSQLRepository{db: db}
}

func (r *PostgresSQLRepository) CreatePost(authorID, title, content string, allowComments bool) (*model.Post, error) {
	query := `INSERT INTO posts (author_id, title, content, allow_comments) 
			  VALUES ($1, $2, $3, $4) RETURNING id, author_id, title, content, allow_comments`

	var post model.Post
	err := r.db.QueryRow(context.Background(), query, authorID, title, content, allowComments).Scan(
		&post.ID, &post.AuthorID, &post.Title, &post.Content, &post.AllowComments,
	)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *PostgresSQLRepository) GetPosts(limit int, after *string) (*model.PostConnection, error) {
	query := `SELECT id, author_id ,title, content, allow_comments FROM posts ORDER BY id LIMIT $1`
	args := []interface{}{limit}
	if after != nil {
		query = `SELECT id, author_id, title, content, allow_comments FROM posts WHERE id > $2 ORDER BY id LIMIT $1`
		args = append(args, *after)
	}

	rows, err := r.db.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*model.Post
	for rows.Next() {
		var post model.Post
		if err := rows.Scan(&post.ID, &post.AuthorID, &post.Title, &post.Content, &post.AllowComments); err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}

	edges := make([]*model.PostEdge, len(posts))
	for i, post := range posts {
		edges[i] = &model.PostEdge{
			Cursor: post.ID,
			Node:   post,
		}
	}

	var endCursor *string

	if len(edges) > 0 {
		endCursor = &edges[len(edges)-1].Cursor
	}

	hasNextPage := len(posts) == limit

	return &model.PostConnection{
		Edges: edges,
		PageInfo: &model.PageInfo{
			EndCursor:   endCursor,
			HasNextPage: hasNextPage,
		},
	}, nil
}

func (r *PostgresSQLRepository) GetPostByID(id string) (*model.Post, error) {
	query := `SELECT id, author_id, title, content, allow_comments FROM posts WHERE id = $1`
	var post model.Post
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&post.ID, &post.AuthorID, &post.Title, &post.Content, &post.AllowComments,
	)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *PostgresSQLRepository) CreateComment(authorID, postID string, content string) (*model.Comment, error) {
	query := `
		INSERT INTO comments (author_id, post_id, content, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, author_id, post_id, content, created_at
	`

	var comment model.Comment
	var createdAt time.Time

	err := r.db.QueryRow(context.Background(), query, authorID, postID, content, time.Now()).Scan(
		&comment.ID, &comment.AuthorID, &comment.PostID, &comment.Content, &createdAt,
	)
	if err != nil {
		return nil, err
	}

	comment.CreatedAt = createdAt.Format(time.RFC3339)

	return &comment, nil
}

func (r *PostgresSQLRepository) GetComments(postID string, limit int, after *string) (*model.CommentConnection, error) {
	query := `SELECT id, author_id, post_id, content, created_at FROM comments WHERE post_id = $1 ORDER BY id LIMIT $2`
	args := []interface{}{postID, limit}
	if after != nil {
		query = `SELECT id, author_id, post_id, content, created_at FROM comments WHERE post_id = $1 AND id > $3 ORDER BY id LIMIT $2`
		args = append(args, *after)
	}

	rows, err := r.db.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*model.Comment
	var createdAt time.Time
	for rows.Next() {
		var comment model.Comment
		if err := rows.Scan(&comment.ID, &comment.AuthorID, &comment.PostID, &comment.Content, &createdAt); err != nil {
			return nil, err
		}
		comment.CreatedAt = createdAt.Format(time.RFC3339)
		comments = append(comments, &comment)
	}

	edges := make([]*model.CommentEdge, len(comments))
	for i, comment := range comments {
		edges[i] = &model.CommentEdge{
			Cursor: comment.ID,
			Node:   comment,
		}
	}

	var endCursor *string
	if len(edges) > 0 {
		endCursor = &edges[len(edges)-1].Cursor
	}

	hasNextPage := len(comments) == limit

	return &model.CommentConnection{
		Edges: edges,
		PageInfo: &model.PageInfo{
			EndCursor:   endCursor,
			HasNextPage: hasNextPage,
		},
	}, nil
}

func (r *PostgresSQLRepository) CreateReply(authorID, postID string, content string, parentID *string) (*model.Comment, error) {
	comment, err := r.CreateComment(authorID, postID, content)
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO replies_comments (parent_comment_id, reply_comment_id)
		VALUES ($1, $2)
		RETURNING parent_comment_id
		`

	var createdAt time.Time

	err = r.db.QueryRow(context.Background(), query, parentID, comment.ID).Scan(
		&comment.ParentID,
	)
	if err != nil {
		return nil, err
	}

	comment.CreatedAt = createdAt.Format(time.RFC3339)

	return comment, nil
}

func (r *PostgresSQLRepository) GetRepliesByCommentID(commentID string, limit int, after *string) (*model.CommentConnection, error) {
	query := `SELECT c.id, c.author_id, c.post_id, rc.parent_comment_id, c.content, c.created_at
			  FROM comments c JOIN replies_comments rc ON c.id = rc.reply_comment_id
			  WHERE rc.parent_comment_id = $1 ORDER BY c.id LIMIT $2`
	args := []interface{}{commentID, limit}
	if after != nil {
		query = `SELECT c.id, c.author_id, c.post_id, rc.parent_comment_id, c.content, c.created_at
				 FROM comments c JOIN replies_comments rc ON c.id = rc.reply_comment_id
				 WHERE rc.parent_comment_id = $1 AND c.id > $3 ORDER BY c.id LIMIT $2`
		args = append(args, *after)
	}

	rows, err := r.db.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var replies []*model.Comment
	var createdAt time.Time

	for rows.Next() {
		var reply model.Comment
		if err := rows.Scan(&reply.ID, &reply.AuthorID, &reply.PostID, &reply.ParentID, &reply.Content, &createdAt); err != nil {
			return nil, err
		}

		reply.CreatedAt = createdAt.Format(time.RFC3339)
		replies = append(replies, &reply)
	}

	edges := make([]*model.CommentEdge, len(replies))
	for i, reply := range replies {
		edges[i] = &model.CommentEdge{
			Cursor: reply.ID,
			Node:   reply,
		}
	}

	var endCursor *string
	if len(edges) > 0 {
		endCursor = &edges[len(edges)-1].Cursor
	}

	hasNextPage := len(replies) == limit

	return &model.CommentConnection{
		Edges: edges,
		PageInfo: &model.PageInfo{
			EndCursor:   endCursor,
			HasNextPage: hasNextPage,
		},
	}, nil
}
