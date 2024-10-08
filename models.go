package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/haedarrfd/simple-rss-aggregator/internal/database"
)

// Costum key name in JSON
type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdateAt  time.Time `json:"updated_at"`
	APIKey    string    `json:"api_key"`
}

type Feed struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	UpdateAt  time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
}

type FeedFollow struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdateAt  time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	FeedID    uuid.UUID `json:"feed_id"`
}

type Post struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	PublishedAt time.Time `json:"published_at"`
	Url         string    `json:"url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	FeedID      uuid.UUID `json:"feed_id"`
}

// convert objects from the database format generated by sqlc
func databaseUserToUser(dbUser database.User) User {
	return User{
		ID:        dbUser.ID,
		Name:      dbUser.Name,
		CreatedAt: dbUser.CreatedAt,
		UpdateAt:  dbUser.UpdatedAt,
		APIKey:    dbUser.ApiKey,
	}
}

func databaseFeedToFeed(dbFeed database.Feed) Feed {
	return Feed{
		ID:        dbFeed.ID,
		Name:      dbFeed.Name,
		URL:       dbFeed.Url,
		CreatedAt: dbFeed.CreatedAt,
		UpdateAt:  dbFeed.UpdatedAt,
		UserID:    dbFeed.UserID,
	}
}

func databaseFeedFolToFeedFol(dbFeedFol database.FeedFollow) FeedFollow {
	return FeedFollow{
		ID:        dbFeedFol.ID,
		CreatedAt: dbFeedFol.CreatedAt,
		UpdateAt:  dbFeedFol.UpdatedAt,
		UserID:    dbFeedFol.UserID,
		FeedID:    dbFeedFol.FeedID,
	}
}

func databasePostToPost(dbPost database.Post) Post {
	// if the decription exists, assisgn its value to description pointer
	var description *string
	if dbPost.Description.Valid {
		description = &dbPost.Description.String
	}

	return Post{
		ID:          dbPost.ID,
		Title:       dbPost.Title,
		Description: description,
		PublishedAt: dbPost.PublishedAt,
		Url:         dbPost.Url,
		CreatedAt:   dbPost.CreatedAt,
		UpdatedAt:   dbPost.UpdatedAt,
		FeedID:      dbPost.FeedID,
	}
}

// Converts into a slice of data
func databaseFeedsToFeeds(dbFeeds []database.Feed) []Feed {
	feeds := []Feed{}

	for _, dbFeed := range dbFeeds {
		feeds = append(feeds, databaseFeedToFeed(dbFeed))
	}

	return feeds
}

func databaseFeedFolsToFeedFols(dbFeedFols []database.FeedFollow) []FeedFollow {
	feedFollows := []FeedFollow{}

	for _, dbFeedFol := range dbFeedFols {
		feedFollows = append(feedFollows, databaseFeedFolToFeedFol(dbFeedFol))
	}

	return feedFollows
}

func databasePostsToPosts(dbPosts []database.Post) []Post {
	posts := []Post{}

	for _, dbpost := range dbPosts {
		posts = append(posts, databasePostToPost(dbpost))
	}

	return posts
}
