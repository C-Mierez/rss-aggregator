package store

import (
	"github.com/c-mierez/rss-aggregator/internal/lib/queries"
	"github.com/google/uuid"
)

/* ---------------------------------- User ---------------------------------- */
type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	APIKey    string    `json:"api_key"`
}

func DBToStoreUser(dbUser queries.User) User {
	return User{
		ID:        dbUser.ID,
		Name:      dbUser.Name,
		CreatedAt: dbUser.CreatedAt.String(),
		UpdatedAt: dbUser.UpdatedAt.String(),
		APIKey:    dbUser.ApiKey,
	}
}

/* ---------------------------------- Feed ---------------------------------- */
type Feed struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
}

func DBToStoreFeed(dbFeed queries.Feed) Feed {
	return Feed{
		ID:        dbFeed.ID,
		Name:      dbFeed.Name,
		Url:       dbFeed.Url,
		CreatedAt: dbFeed.CreatedAt.String(),
		UpdatedAt: dbFeed.UpdatedAt.String(),
		UserID:    dbFeed.UserID,
	}
}

func DBToStoreFeeds(dbFeeds []queries.Feed) []Feed {
	feeds := []Feed{}
	for _, feed := range dbFeeds {
		feeds = append(feeds, DBToStoreFeed(feed))
	}
	return feeds
}

/* --------------------------------- Follow --------------------------------- */
type Follow struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	FeedID    uuid.UUID `json:"feed_id"`
}

func DBToStoreFollow(dbFollow queries.Follow) Follow {
	return Follow{
		ID:        dbFollow.ID,
		CreatedAt: dbFollow.CreatedAt.String(),
		UpdatedAt: dbFollow.UpdatedAt.String(),
		UserID:    dbFollow.UserID,
		FeedID:    dbFollow.FeedID,
	}
}

func DBToStoreFollows(dbFollows []queries.Follow) []Follow {
	follows := []Follow{}
	for _, follow := range dbFollows {
		follows = append(follows, DBToStoreFollow(follow))
	}
	return follows
}

/* ---------------------------------- Posts --------------------------------- */
type Post struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Url         string    `json:"url"`
	PublishedAt string    `json:"published_at"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
	FeedID      uuid.UUID `json:"feed_id"`
}

func DBToStorePost(dbPost queries.Post) Post {
	return Post{
		ID:          dbPost.ID,
		Title:       dbPost.Title,
		Description: dbPost.Description,
		Url:         dbPost.Url,
		PublishedAt: dbPost.PublishedAt.String(),
		CreatedAt:   dbPost.CreatedAt.String(),
		UpdatedAt:   dbPost.UpdatedAt.String(),
		FeedID:      dbPost.FeedID,
	}
}

func DBToStorePosts(dbPosts []queries.Post) []Post {
	posts := []Post{}
	for _, post := range dbPosts {
		posts = append(posts, DBToStorePost(post))
	}
	return posts
}
