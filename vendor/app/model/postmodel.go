package model

import (
	"app/shared/database"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Post contains the information for each post
type Post struct {
	PostID      bson.ObjectId `json:"id" bson:"_id"`
	FeedID      string        `json:"feed_id,omitempty" bson:"feed_id,omitempty"` // Feed Owner
	CreatorID   string        `json:"creator_id,omitempty" bson:"creator_id,omitempty"`
	PostContent string        `json:"post_con,omitempty" bson:"post_con,omitempty"`
	CreatedAt   time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at" bson:"updated_at"`
	Deleted     uint8         `json:"deleted" bson:"deleted"`
}

// GetPostByID : gets post by postid
func GetPostByID(postid string) (Post, error) {
	var err error

	result := Post{}

	if database.CheckConnection() {
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().Database).C("post")
		if bson.IsObjectIdHex(postid) {
			err = c.FindId(bson.ObjectIdHex(postid)).One(&result)
		} else {
			err = ErrNoResult
		}
	} else {
		err = ErrUnavailable
	}

	return result, standardizeError(err)
}

// GetPostByFeed gets all post from a user/feed
func GetPostByFeed(userorfeed string) ([]Post, error) {
	var err error

	var result []Post

	if database.CheckConnection() {
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().Database).C("post")
		err = c.Find(bson.M{"feed_id": userorfeed}).All(&result)
	} else {
		err = ErrUnavailable
	}

	return result, standardizeError(err)
}

// CreatePostFeed : creating post
func CreatePostFeed(feedid, visitorid, content string) error {
	var err error

	now := time.Now()

	if database.CheckConnection() {
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().Database).C("post")
		post := &Post{
			PostID:      bson.NewObjectId(),
			FeedID:      feedid,
			CreatorID:   visitorid,
			PostContent: content,
			CreatedAt:   now,
			UpdatedAt:   now,
			Deleted:     0,
		}
		err = c.Insert(post)
	} else {
		err = ErrUnavailable
	}

	return standardizeError(err)
}

// UpdatePost : updates a post
func UpdatePost(content, visitorid, postid string) error {
	var err error

	now := time.Now()

	if database.CheckConnection() {
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().Database).C("post")
		var post Post
		post, err = GetPostByID(postid)
		if err == nil {
			if post.CreatorID == visitorid {
				post.UpdatedAt = now
				post.PostContent = content
				err = c.UpdateId(bson.ObjectIdHex(postid), &post)
			} else {
				err = ErrUnauthorized
			}
		}
	} else {
		err = ErrUnavailable
	}

	return standardizeError(err)
}

// DeletePost : deletes a post
func DeletePost(visitorid, postid string) error {
	var err error

	if database.CheckConnection() {
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().Database).C("post")

		var post Post
		post, err = GetPostByID(postid)
		if err == nil {
			if post.FeedID == visitorid || post.CreatorID == visitorid {
				err = c.RemoveId(bson.ObjectIdHex(postid))
			} else {
				err = ErrUnauthorized
			}
		}
	} else {
		err = ErrUnavailable
	}

	return standardizeError(err)
}
