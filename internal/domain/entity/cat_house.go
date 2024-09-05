package entity

type CatHouse struct {
	Nekojarashis []string // objectID
	Nikukyus     []string // userID
}

func (c *CatHouse) DeepCopy() *CatHouse {
	return &CatHouse{
		Nekojarashis: append([]string{}, c.Nekojarashis...),
		Nikukyus:     append([]string{}, c.Nikukyus...),
	}
}
