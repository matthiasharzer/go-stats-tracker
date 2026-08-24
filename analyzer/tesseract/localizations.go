package tesseract

type localization struct {
	Me      string
	Friends string
	Social  string
}

var localizations = []localization{
	{
		Me:      "Me",
		Friends: "Friends",
		Social:  "Social",
	},
	{
		Me:      "Du",
		Friends: "Freunde",
		Social:  "Gruppen",
	},
}
