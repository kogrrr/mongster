package objects

type User struct {
	ID        string `json:"id" bson:"_id"`
	Name      string `json:"name" bson:"name"`
	OIDCtoken string `json:"oidc_token" bson:"oidc_token"`
}
