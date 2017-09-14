package model

// Form structurized the form model
type Form struct {
	UserName        string `json:"username,omitempty" bson:"username,omitempty"`
	FirstName       string `json:"first_name,omitempty" bson:"first_name,omitempty"`
	LastName        string `json:"last_name,omitempty" bson:"last_name,omitempty"`
	Email           string `json:"email,omitempty" bson:"email,omitempty"`
	Password        string `json:"password,omitempty" bson:"password,omitempty"`
	ConfirmPassword string `json:"con_pass,omitempty" bson:"con_pass,omitempty"`
}

// UpdateForm structurized the form model
type UpdateForm struct {
	UserName        string `json:"username,omitempty" bson:"username,omitempty"`
	NewFirstName    string `json:"new_first,omitempty" bson:"new_first,omitempty"`
	NewLastName     string `json:"new_last,omitempty" bson:"new_last,omitempty"`
	NewEmail        string `json:"new_email,omitempty" bson:"new_email,omitempty"`
	OldPassword     string `json:"old_pass,omitempty" bson:"old_pass,omitempty"`
	NewPassword     string `json:"new_pass,omitempty" bson:"new_pass,omitempty"`
	ConfirmPassword string `json:"con_pass,omitempty" bson:"con_pass,omitempty"`
}
