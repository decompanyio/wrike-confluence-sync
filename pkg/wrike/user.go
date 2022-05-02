package wrike

type Users struct {
	Kind string `json:"kind"`
	Data []User `json:"data"`
}

type User struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Type      string `json:"type"`
	Profiles  []struct {
		AccountID string `json:"accountId"`
		Email     string `json:"email"`
		Role      string `json:"role"`
		External  bool   `json:"external"`
		Admin     bool   `json:"admin"`
		Owner     bool   `json:"owner"`
	} `json:"profiles"`
	AvatarURL string `json:"avatarUrl"`
	Timezone  string `json:"timezone"`
	Locale    string `json:"locale"`
	Deleted   bool   `json:"deleted"`
	Me        bool   `json:"me"`
	Phone     string `json:"phone"`
}

func (w *WrikeClient) User(userId string) User {
	users := Users{}
	w.newAPI("/users/"+userId, nil, &users)

	return users.Data[0]
}
