package recipeprovider

type Dish struct {
	Name   string `json:"name"`
	Recipe string `json:"recipe"`
}

type DishResponse struct {
	Dishes []Dish `json:"dishes"`
}
