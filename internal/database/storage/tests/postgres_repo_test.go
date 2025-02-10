package tests

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"ozon-GraphQL/internal/database/storage"
	"ozon-GraphQL/internal/database/storage/mocks"
	"testing"
)

func TestPostgresCreatePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDatabase(ctrl)

	repo := storage.NewPostgresSQLRepository(mockDB)

	authorID := "author123"
	title := "Test Post"
	content := "Lorem ipsum"
	allowComments := true

	mockRow := mocks.NewMockRow(ctrl)
	mockRow.EXPECT().Scan(
		gomock.Any(), // id
		gomock.Any(), // author_id
		gomock.Any(), // title
		gomock.Any(), // content
		gomock.Any(), // allow_comments
	).Return(nil).Times(1)

	mockDB.EXPECT().
		QueryRow(gomock.Any(), gomock.Any(), authorID, title, content, allowComments).
		Return(mockRow).
		Times(1)

	_, err := repo.CreatePost(authorID, title, content, allowComments)

	assert.NoError(t, err, "Expected no error")
}

func TestPostgresGetPostByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDatabase(ctrl)

	repo := storage.NewPostgresSQLRepository(mockDB)

	postID := "post123"

	mockRow := mocks.NewMockRow(ctrl)
	mockRow.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

	mockDB.EXPECT().
		QueryRow(gomock.Any(), gomock.Any(), postID).
		Return(mockRow).
		Times(1)

	_, err := repo.GetPostByID(postID)

	assert.NoError(t, err, "Expected no error")
}

func TestPostgresCreateComment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDatabase(ctrl)

	repo := storage.NewPostgresSQLRepository(mockDB)

	authorID := "author123"
	postID := "post123"
	content := "This is a comment"

	mockRow := mocks.NewMockRow(ctrl)
	mockRow.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

	mockDB.EXPECT().
		QueryRow(gomock.Any(), gomock.Any(), authorID, postID, content, gomock.Any()).
		Return(mockRow).
		Times(1)

	_, err := repo.CreateComment(authorID, postID, content)

	assert.NoError(t, err, "Expected no error")
}
