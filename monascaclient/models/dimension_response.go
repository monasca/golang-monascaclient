package models

type DimensionValueResponse struct {
	Links    []Link           `json:"links"`
	Elements []DimensionValue `json:"elements"`
}

type DimensionValue struct {
	Value string `json:"dimension_value"`
}