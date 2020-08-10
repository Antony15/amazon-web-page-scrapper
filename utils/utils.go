package utils

import (
	"regexp"
	"strings"
)

// FormatPrice Function for formating the price
func FormatPrice(price *string) {
	r := regexp.MustCompile(`\$(\d+(\.\d+)?).*$`)

	newPrices := r.FindStringSubmatch(*price)

	if len(newPrices) > 1 {
		*price = newPrices[1]
	} else {
		*price = "Unknown"
	}

}

// FormatStars Function for formating the stars
func FormatStars(stars *string) {
	if len(*stars) >= 3 {
		*stars = (*stars)[0:3]
	} else {
		*stars = "Unknown"
	}
}

// FormatReviews Function for formating the reviews
func FormatReviews(totalReviews *string) {
	s := strings.Split(*totalReviews, " ")
	*totalReviews = s[0]
}
