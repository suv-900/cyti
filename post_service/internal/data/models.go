package data

import "go.mongodb.org/mongo-driver/mongo"

type Models struct {
	Posts interface {
		AddPost(post *Post) error
		GetPost(postid uint64) Post
		UpdatePost(post *Post) (Post, error)
		DeletePost(postid uint64) error
	}
}

func GetModels(client *mongo.Client) Models {
	return Models{
		Posts: PostModels{
			client:     client,
			database:   client.Database(""),
			collection: client.Database("").Collection(""),
		},
	}
}
