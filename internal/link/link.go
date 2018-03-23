package link

type Link struct {
	ID        int    `json:"id"`
	URL       string `json:"url"`
	UserID    uint   `json:"iduser"`
	CreatedAt int    `json:"created_at"`
}

func Save(l Link) error {

	return nil
}
