package models

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/labstack/echo"
	"gitlab.com/teamliquid-dev/decks-of-runeterra/doruneterraapi-go/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type SocialLinks struct {
	Instagram string `json:"instagram,omitempty" bson:"instagram,omitempty"`
	Facebook  string `json:"facebook,omitempty" bson:"facebook,omitempty"`
	Twitter   string `json:"twitter,omitempty" bson:"twitter,omitempty"`
	Discord   string `json:"discord,omitempty" bson:"discord,omitempty"`
	Twitch    string `json:"twitch,omitempty" bson:"twitch,omitempty"`
}

type User struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Access      int                `json:"access" bson:"access"`
	Username    string             `json:"username" bson:"username"`
	Email       string             `json:"email" bson:"email"`
	Password    string             `json:"password" bson:"password"`
	DateCreated time.Time          `json:"date_created" bson:"date_created"`
	DateUpdated time.Time          `json:"date_updated,omitempty" bson:"date_updated,omitempty"`
	Socials     SocialLinks        `json:"socials,omitempty" bson:"socials,omitempty"`
}

func (u User) UserID() string {
	return u.ID.Hex()
}

type UserModel struct {
	collection *mongo.Collection
}

func InitUserModel(d *db.Database) *UserModel {
	collection := d.Collection("users")
	indices := make([]mongo.IndexModel, 2)
	indices[0] = mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	indices[1] = mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := collection.Indexes().CreateMany(
		context.Background(),
		indices,
	)
	if err != nil {
		panic(err)
	}

	model := NewUserModel(collection)

	return model
}

func NewUserModel(c *mongo.Collection) *UserModel {
	return &UserModel{
		collection: c,
	}
}

func (u *UserModel) Login(email string, password string) (*User, error) {
	user, err := u.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if !CheckPasswordHash(password, user.Password) {
		return nil, echo.ErrUnauthorized
	}

	return user, nil
}

func (u *UserModel) Register(username, email, password string) (*User, error) {
	emailUser, _ := u.GetUserByEmail(email)
	if emailUser != nil {
		return nil, errors.New("An account with that email address already exists")
	}

	usernameUser, _ := u.GetUserByUsername(username)
	if usernameUser != nil {
		return nil, errors.New("An account with that username already exists")
	}

	hash, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	newUser := User{
		Username:    username,
		Email:       email,
		Password:    hash,
		Access:      0,
		DateCreated: time.Now(),
		DateUpdated: time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cur, err := u.collection.InsertOne(ctx, newUser)
	if err != nil {
		return nil, err
	}

	newUser.ID = cur.InsertedID.(primitive.ObjectID)

	return &newUser, nil
}

func (u *UserModel) GetUserByEmail(email string) (*User, error) {
	var user User
	pattern := fmt.Sprintf(`^%s$`, email)
	regex := bson.D{{Key: "$regex", Value: primitive.Regex{Pattern: pattern, Options: "i"}}}
	result := u.collection.FindOne(context.Background(), bson.D{primitive.E{Key: "email", Value: regex}})
	err := result.Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserModel) GetUserById(id string) (*User, error) {
	var user User
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	result := u.collection.FindOne(context.Background(), bson.D{primitive.E{Key: "_id", Value: userID}})
	err = result.Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserModel) GetUserByUsername(username string) (*User, error) {
	var user User
	pattern := fmt.Sprintf(`^%s$`, username)
	regex := bson.D{{Key: "$regex", Value: primitive.Regex{Pattern: pattern, Options: "i"}}}
	result := u.collection.FindOne(context.Background(), bson.D{primitive.E{Key: "username", Value: regex}})
	err := result.Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserModel) SearchUsers(username, email string) ([]*User, error) {
	if len(username) == 0 && len(email) == 0 {
		return make([]*User, 0), nil
	}

	var users []*User

	regexUsername := bson.D{{Key: "$regex", Value: primitive.Regex{Pattern: username, Options: "i"}}}
	regexEmail := bson.D{{Key: "$regex", Value: primitive.Regex{Pattern: email, Options: "i"}}}
	regexOr := *new([]bson.D)
	if len(username) > 0 {
		regexOr = append(regexOr, bson.D{{Key: "username", Value: regexUsername}})
	}
	if len(email) > 0 {
		regexOr = append(regexOr, bson.D{{Key: "email", Value: regexEmail}})
	}

	filter := bson.D{{Key: "$or", Value: regexOr}}

	cur, err := u.collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	defer cur.Close(context.Background())

	if err := cur.All(context.Background(), &users); err != nil {
		return nil, err
	}

	if users == nil {
		users = make([]*User, 0)
	}

	return users, nil
}

func (u *UserModel) UpdateUser(user *User) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var updatedUser User
	after := options.After
	options := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}
	userID := user.ID

	user.ID = primitive.NilObjectID
	user.DateUpdated = time.Now()

	update := bson.M{"$set": user}
	filter := bson.M{"_id": userID}

	res := u.collection.FindOneAndUpdate(ctx, filter, update, &options)
	err := res.Decode(&updatedUser)
	if err != nil {
		return nil, err
	}

	return &updatedUser, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
