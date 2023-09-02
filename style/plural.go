package style

func SelectPlural(singular string, plural string, count int) string {
	if count == 1 {
		return singular
	} else {
		return plural
	}
}
