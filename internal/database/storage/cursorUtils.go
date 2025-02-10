package storage

import "ozon-GraphQL/graph/model"

func (r *InMemoryRepository) findIndex(cursor string, posts []*model.Post) int {
	for i, post := range posts {
		if post.ID == cursor {
			return i
		}
	}
	return -1
}

func (r *InMemoryRepository) findCommentIndex(cursor string, comments []*model.Comment) int {
	for i, comment := range comments {
		if comment.ID == cursor {
			return i
		}
	}
	return -1
}

func (r *InMemoryRepository) findReplyIndex(cursor string, edges []*model.CommentEdge) int {
	for i, edge := range edges {
		if edge.Cursor == cursor {
			return i
		}
	}
	return -1
}

func (r *InMemoryRepository) getEndCursor(posts []*model.Post) *string {
	if len(posts) == 0 {
		return nil
	}
	lastPost := posts[len(posts)-1]
	return &lastPost.ID
}

func (r *InMemoryRepository) getEndCursorComments(comments []*model.Comment) *string {
	if len(comments) == 0 {
		return nil
	}
	lastComment := comments[len(comments)-1]
	return &lastComment.ID
}

func (r *InMemoryRepository) getEndCursorEdges(edges []*model.CommentEdge) *string {
	if len(edges) == 0 {
		return nil
	}
	lastEdge := edges[len(edges)-1]
	return &lastEdge.Cursor
}

func (r *InMemoryRepository) postsToSlice() []*model.Post {
	var posts []*model.Post
	for _, post := range r.posts {
		posts = append(posts, post)
	}
	return posts
}
