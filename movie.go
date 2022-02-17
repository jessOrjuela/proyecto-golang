package main

type Movie struct{
	Name string `json:"nombre"`
	Year int `json:"year"`
	Director string `json:"director"`
}
type Movies []Movie

