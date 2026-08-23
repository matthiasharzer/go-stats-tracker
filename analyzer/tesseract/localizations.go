package tesseract

type localization struct {
	You     string
	Friends string
	Groups  string
}

var localizations = []localization{
	{
		You:     "You",
		Friends: "Friends",
		Groups:  "Groups",
	},
	{
		You:     "Du",
		Friends: "Freunde",
		Groups:  "Gruppen",
	},
}
