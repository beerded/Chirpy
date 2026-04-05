package main

import (
	"fmt"
	"strings"
)

func validateChirp(body string) (string, error) {
	const limit int = 140

	if len(body) > limit {
		return "", fmt.Errorf("Chirp is too long")
	}
	if body == "" {
		return "", fmt.Errorf("Missing Body from chirp. Be sure to add 'body:' to the request body")
	}
	badwords := map[string]struct{} {
		"kerfuffle": 	{},
		"sharbert": 	{},
		"fornax":		{},
	}
	return cleanLanguage(body, badwords), nil
}

func cleanLanguage(original string, badwords map[string]struct{}) string {
	newWords := []string{}
	for _, word := range strings.Split(original, " ") {
		tinyWord := strings.ToLower(word)
		if _, ok := badwords[tinyWord]; ok {
			newWords = append(newWords, "****")
		} else {
			newWords = append(newWords, word)
		}
	}
	return strings.Join(newWords, " ")
}
