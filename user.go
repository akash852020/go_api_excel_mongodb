package models

type Users struct {
	ID    string `bson:"_id,omitempty" json:"id,omitempty"`
	Name  string `bson:"name" json:"name"`
	Email string `bson:"email" json:"email"`
}
