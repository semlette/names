package names

import (
	"reflect"
	"testing"
)

var allowedCharacters = "ABCDEFGHIJKLMNOPQRSTUVWXYZWabcdefghijklmnopqrstuvwzyx1234567890_$"

func run(tok *Tokeniser) []Token {
	var slice []Token
	for {
		tkn := tok.Token()
		if tkn.Kind == Invalid {
			slice = append(slice, tkn)
			return slice
		}
		if tkn.Kind == EOF {
			return slice
		}
		slice = append(slice, tkn)
	}
}

func TestToken(t *testing.T) {
	tests := []struct {
		input  string
		tokens []Token
	}{
		{
			"top_posts",
			[]Token{
				{Word, "top", 0},
				{Symbol, "_", 0},
				{Word, "posts", 0},
			},
		},
		{
			"topPosts",
			[]Token{
				{Word, "top", 0},
				{Word, "Posts", 0},
			},
		},
		{
			"TopPosts",
			[]Token{
				{Word, "Top", 0},
				{Word, "Posts", 0},
			},
		},
		{
			"TOPPOSTS",
			[]Token{
				{Word, "TOPPOSTS", flags(true)},
			},
		},
		{
			"TOP_POSTS",
			[]Token{
				{Word, "TOP", flags(true)},
				{Symbol, "_", 0},
				{Word, "POSTS", flags(true)},
			},
		},
		{
			"Property$Wrapper",
			[]Token{
				{Word, "Property", 0},
				{Symbol, "$", 0},
				{Word, "Wrapper", 0},
			},
		},
		{
			"One2Three4Five6",
			[]Token{
				{Word, "One", 0},
				{Symbol, "2", 0},
				{Word, "Three", 0},
				{Symbol, "4", 0},
				{Word, "Five", 0},
				{Symbol, "6", 0},
			},
		},
		{
			"_$123",
			[]Token{
				{Symbol, "_", 0},
				{Symbol, "$", 0},
				{Symbol, "1", 0},
				{Symbol, "2", 0},
				{Symbol, "3", 0},
			},
		},
		{
			"W420",
			[]Token{
				{Word, "W", flags(true)},
				{Symbol, "4", 0},
				{Symbol, "2", 0},
				{Symbol, "0", 0},
			},
		},
		{
			"isValid",
			[]Token{
				{Word, "is", 0},
				{Word, "Valid", 0},
			},
		},
		{
			"ofType",
			[]Token{
				{Word, "of", 0},
				{Word, "Type", 0},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			var tok Tokeniser
			tok.Input = test.input
			tok.Characters = allowedCharacters
			result := run(&tok)
			if !reflect.DeepEqual(test.tokens, result) {
				t.Log("expected:")
				for idx, tkn := range test.tokens {
					t.Logf("\t%d: %v", idx, tkn)
				}
				t.Log("result:")
				for idx, tkn := range result {
					t.Logf("\t%d: %v", idx, tkn)
				}
				t.Fail()
			}
		})
	}
}
