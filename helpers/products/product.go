package products

type ProductImages struct {
	BaseImage        string
	SmallImage       string
	ThumbnailImage   string
	RolloverImage    string
	AdditionalImages []string
}

type Product struct {
	Sku    string
	Images ProductImages
}
