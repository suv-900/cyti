package data

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const context_timeout = 10 * time.Second

var (
	ErrNotFound = errors.New("not found.")
)

type Post struct {
	ID        primitive.ObjectID
	CreatedAt time.Time
	UpdatedAt time.Time

	AuthorID   uint64
	AuthorName string

	Title   string
	Content string
	Likes   uint16
}

type PostModels struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
}

func (p PostModels) AddPost(post *Post) (uint64, error) {
	var postid uint64

	ctx, cancel := context.WithTimeout(context.Background(), context_timeout)
	defer cancel()
	_, err := p.collection.InsertOne(ctx, bson.D{})
	if err != nil {
		return postid, err
	}
	return postid, nil
}

// FindOneAndDelete finds one deletes it returns the one before save
func (p PostModels) GetPost(postid primitive.ObjectID) (Post, error) {
	var post Post

	ctx, cancel := context.WithTimeout(context.Background(), context_timeout)
	defer cancel()

	var doc bson.D

	opt := options.FindOne()
	err := p.collection.FindOne(ctx, bson.D{{Key: "_id", Value: postid}}, opt).Decode(&doc)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return post, ErrNotFound
		}
		return post, err
	}

	doc_bytes, err := bson.Marshal(doc)
	if err != nil {
		return post, err
	}
	err = bson.Unmarshal(doc_bytes, &post)
	if err != nil {
		return post, err
	}

	return post, nil
}

func (p *Post) Validate() error {
	if len(p.Title) == 0 {
		return errors.New("post title cannot be empty")
	}
	if len(p.Content) == 0 {
		return errors.New("post content cannot be empty")
	}
	return nil
}

// creating a child context or context.TODO()

// make sure title and content are not empty
func (p PostModels) UpdatePost(post *Post) error {

	ctx, cancel := context.WithTimeout(context.Background(), context_timeout)
	defer cancel()

	options := options.Update()
	filter := bson.D{{Key: "_id", Value: post.ID}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "post_title", Value: post.Title}}},
		{Key: "$set", Value: bson.D{
			{Key: "content", Value: post.Content}}}}

	result, err := p.collection.UpdateOne(ctx, filter, update, options)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return ErrNotFound
	}

	return nil
}

func (p PostModels) DeletePost(postid primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), context_timeout)
	defer cancel()

	opts := options.Delete()
	filter := bson.D{{Key: "_id", Value: postid}}
	res, err := p.collection.DeleteOne(ctx, filter, opts)

	if res.DeletedCount == 0 {
		return ErrNotFound
	}

	return err
}

func GetFeaturedPosts(offset uint64) []Posts {
	var postMetaData []Posts
	limit := 3
	sql := "SELECT post_id,post_title FROM posts ORDER BY post_likes DESC OFFSET ? LIMIT ?"
	db.Raw(sql, offset, limit).Scan(&postMetaData)
	return postMetaData
}
func GetAllPostsMetaData() []Posts {
	var posts []Posts
	db.Raw(`SELECT post_id,post_title,author_name,post_likes FROM posts ORDER BY post_likes DESC`).Scan(&posts)
	return posts
}
func CheckPostTitleExists(posttitle string) (bool, error) {
	var exists bool
	r := db.Raw("SELECT EXISTS (SELECT 1 FROM posts WHERE post_title = ?)", posttitle).Scan(&exists)
	if r.Error != nil {
		return exists, r.Error
	}
	return exists, nil
}

func GetPostsByAuthorID(authorid uint64, limit uint64, offset uint64) []Posts {
	var posts []Posts
	db.Raw(`SELECT 
	post_id,post_title,author_name,post_likes,createdat 
	FROM posts WHERE author_id=? 
	ORDER BY post_likes DESC LIMIT ? OFFSET ?`, authorid, limit, offset*limit).Scan(&posts)
	return posts
}

func GetPostsMetaData(offset uint64, limit uint64) []Posts {
	var posts []Posts
	db.Raw(`SELECT post_id,post_title,
	author_name,post_likes 
	FROM posts ORDER BY post_likes DESC LIMIT ? OFFSET ?`, limit, offset*limit).Scan(&posts)
	return posts
}

// func GetPostAndUserPreferences(postid uint64, userid uint64) (Posts, error) {
// 	var post Posts
// 	r := db.Raw("SELECT * FROM posts WHERE post_id=?", postid).Scan(&post)
// 	if r.Error != nil {
// 		return post,r.Error
// 	}
// 	return post, nil
// }

func CheckUserReaction(userid uint64, postid uint64) (bool, bool, error) {
	var userLikedPost bool
	var userDislikedPost bool
	var count1 int
	var count2 int
	var err error
	r := db.Raw("SELECT COUNT(*) from posts_liked_by_users WHERE user_id=? AND post_id=?", userid, postid).Scan(&count1)
	if r.Error != nil {
		return userLikedPost, userDislikedPost, r.Error
	}
	r = db.Raw("SELECT COUNT(*) FROM posts_disliked_by_users WHERE user_id=? AND post_id=?", userid, postid).Scan(&count2)
	if r.Error != nil {
		return userLikedPost, userDislikedPost, r.Error
	}

	if count1 == 1 && count2 == 2 {
		return userLikedPost, userDislikedPost, errors.New("illegalState:User reaction like and dislike both exists.check CheckUserReaction() function")
	}

	if count1 == 1 {
		userLikedPost = true
		userDislikedPost = false
	} else if count2 == 1 {
		userLikedPost = false
		userDislikedPost = true
	}
	return userLikedPost, userDislikedPost, err
}

func LikePost(postid uint64, userid uint64) error {
	tx := db.Begin()
	r := tx.Exec("INSERT INTO posts_liked_by_users(user_id,post_id) VALUES(?,?) ", userid, postid)
	if r.Error != nil {
		tx.Rollback()
		return r.Error
	} else {
		tx.Commit()
	}

	tx = db.Begin()
	r = tx.Exec("UPDATE posts SET post_likes= post_likes + 1 WHERE post_id=?", postid)
	if r.Error != nil {
		tx.Rollback()
		return r.Error
	} else {
		tx.Commit()
	}
	return nil
}

func DislikePost(postid uint64, userid uint64) error {
	tx := db.Begin()
	r := tx.Exec("INSERT INTO posts_disliked_by_users(user_id,post_id) VALUES(?,?) ", userid, postid)
	if r.Error != nil {
		tx.Rollback()
		return r.Error
	} else {
		tx.Commit()
	}

	tx = db.Begin()
	r = tx.Exec("UPDATE posts SET post_likes=post_likes - 1 WHERE post_id=?", postid)
	if r.Error != nil {
		tx.Rollback()
		return r.Error
	} else {
		tx.Commit()
	}
	return nil
}
func RemoveLikeFromPost(postid uint64, userid uint64) error {
	tx := db.Begin()
	r := tx.Exec("DELETE FROM posts_liked_by_users WHERE user_id=? AND post_id=?", userid, postid)
	if r.Error != nil {
		tx.Rollback()
		return r.Error
	} else {
		tx.Commit()
	}

	tx = db.Begin()
	r = tx.Exec("UPDATE posts SET post_likes=post_likes - 1 WHERE post_id=?", postid)
	if r.Error != nil {
		tx.Rollback()
		return r.Error
	} else {
		tx.Commit()
	}
	return nil
}

func RemoveDislikeFromPost(postid uint64, userid uint64) error {
	tx := db.Begin()
	r := tx.Exec("DELETE FROM posts_disliked_by_users WHERE user_id=? AND post_id=?", userid, postid)
	if r.Error != nil {
		tx.Rollback()
		return r.Error
	} else {
		tx.Commit()
	}

	tx = db.Begin()
	r = tx.Exec("UPDATE posts SET post_likes=post_likes + 1 WHERE post_id=?", postid)
	if r.Error != nil {
		tx.Rollback()
		return r.Error
	} else {
		tx.Commit()
	}
	return nil
}
