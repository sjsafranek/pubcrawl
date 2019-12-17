package api

const (
	CategoryBrewery string = "brewery"
	CategoryBar     string = "bar"
)

var categoryName = map[string]string{
	CategoryBrewery: "50327c8591d4c4b30a586d5d",
	CategoryBar:     "4bf58dd8d48988d116941735",
}

func CategoryCode(name string) string {
	return categoryName[name]
}

var searchCategeories = []string{"50327c8591d4c4b30a586d5d", "4bf58dd8d48988d116941735"}
