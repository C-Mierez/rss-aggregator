// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: follows.sql

package queries

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createFollow = `-- name: CreateFollow :one
INSERT INTO
  follows (id, created_at, updated_at, feed_id, user_id)
VALUES
  ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at, feed_id, user_id
`

type CreateFollowParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	FeedID    uuid.UUID
	UserID    uuid.UUID
}

func (q *Queries) CreateFollow(ctx context.Context, arg CreateFollowParams) (Follow, error) {
	row := q.db.QueryRowContext(ctx, createFollow,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.FeedID,
		arg.UserID,
	)
	var i Follow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.FeedID,
		&i.UserID,
	)
	return i, err
}

const deleteUserFollowByID = `-- name: DeleteUserFollowByID :exec
DELETE FROM
  follows
WHERE
  id = $1
  AND user_id = $2
`

type DeleteUserFollowByIDParams struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

func (q *Queries) DeleteUserFollowByID(ctx context.Context, arg DeleteUserFollowByIDParams) error {
	_, err := q.db.ExecContext(ctx, deleteUserFollowByID, arg.ID, arg.UserID)
	return err
}

const getFollowsByUserID = `-- name: GetFollowsByUserID :many
SELECT
  id, created_at, updated_at, feed_id, user_id
FROM
  follows
WHERE
  user_id = $1
`

func (q *Queries) GetFollowsByUserID(ctx context.Context, userID uuid.UUID) ([]Follow, error) {
	rows, err := q.db.QueryContext(ctx, getFollowsByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Follow
	for rows.Next() {
		var i Follow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.FeedID,
			&i.UserID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
