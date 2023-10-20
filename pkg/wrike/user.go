package wrike

import "fmt"

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
	AvatarURL   string   `json:"avatarUrl"`
	Timezone    string   `json:"timezone"`
	Locale      string   `json:"locale"`
	Deleted     bool     `json:"deleted"`
	Title       string   `json:"title,omitempty"`
	CompanyName string   `json:"companyName,omitempty"`
	Phone       string   `json:"phone,omitempty"`
	Location    string   `json:"location,omitempty"`
	Me          bool     `json:"me,omitempty"`
	MemberIds   []string `json:"memberIds,omitempty"`
	MyTeam      bool     `json:"myTeam,omitempty"`
}

type AllUserMap map[string]User

func (aum *AllUserMap) findUser(userId string) User {
	return (*aum)[userId]
}

func (w *Client) GetAllUsers() (AllUserMap, error) {
	users := Users{}

	// Fetch data from Wrike
	urlQuery := map[string]string{
		"deleted": `false`,
	}

	resp, err := w.httpClient.R().SetQueryParams(urlQuery).
		SetResult(&users).
		Get("/contacts")

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode(), resp.String())
	}

	if err != nil {
		return nil, err
	}

	// userID를 키로 하는 map 생성
	userAll := AllUserMap{}
	for _, user := range users.Data {
		userAll[user.ID] = user
	}

	return userAll, err
}
