package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	dict_path         string = "/usr/share/dict/words"
	encrypted_keyword string = "JICAHUHN"
	cipher_text       string = "EFMZRNQMWOBEBUIXDDMTRDGAXGJUBKNEIVWATHLHJDZUOAOENVIBROCXCSLZUIFUDCSPHSGMDTQTQAUNVSWTVOAKXUPZVNGATAQJQLRVGSIYROSMWSJUBKGATAITHFNVIIZKEOSMWSJUBKUTHWVIYUQXSHPKBNVHCOBSLRRJJSAZFOLHJFMRVKRPDKBNVSOHDYKUZEFPXHPGAOABDBMBRNVYNCCJBNGIPFBOPUYTGZGRVKRHCWWTFIZLJFMEBUPTCOXVEEPBPHMZUEYHVWAZVCFHUGPOCPVGVOVEFOEMDTXXBDHVTRQYPRRXIZGOASVWTCNGAAYETUMJCRBZGOUSVNTFPBCGY"
)

func main() {
	unencrypted_keys := most_likely_keys_for_encrypted_key(encrypted_keyword)
	fmt.Println("Candidate keys:", unencrypted_keys)

	for _, key := range unencrypted_keys {
		next_key_rune := rotator_for(key)

		for _, char := range cipher_text {
			next_key_char := next_key_rune()
			fmt.Print(string(rotate_str_by(rune(char), offset_for(next_key_char, 'A'))))
		}
		fmt.Println("")
	}
}

func rotate_str_by(char rune, rotate_value int) rune {
	offset_position := int(char) + rotate_value
	if offset_position < 65 {
		offset_position = 26 + offset_position
	} else if offset_position > 65+26 {
		offset_position = offset_position + 26
	}
	return rune(offset_position)
}

func most_likely_keys_for_encrypted_key(e_key string) []string {
	words := make(chan string, 10000)
	go func() {
		f, err := os.Open(dict_path)
		if err != nil {
			panic(err)
		}

		f_reader := bufio.NewReader(f)
		for {
			word, err := f_reader.ReadString('\n')
			if err != nil && err != io.EOF {
				panic(err)
			}

			stripped_word := strings.TrimSpace(word)
			if len(stripped_word) > 0 {
				words <- strings.ToUpper(stripped_word)
			}

			if err == io.EOF {
				close(words)
				break
			}
		}
	}()

	var candidates []string
	key_len := len(e_key)
Words:
	for word := range words {
		// First check is simple, is the word the same length?
		if len(word) == key_len {
			/* 	Next check it more complex, do all letters in the candidate
			word and encrypted key have the same normalised offset */
			first_letter_offset := wrapped_offset_for(rune(e_key[0]), rune(word[0]))
			for i := 1; i < key_len; i++ {
				if wrapped_offset_for(rune(e_key[i]), rune(word[i])) != first_letter_offset {
					continue Words
				}
			}
			candidates = append(candidates, word)
		}
	}
	return candidates
}

func rotator_for(s string) func() rune {
	current_i := 0
	return func() rune {
		next_rune := rune(s[current_i])
		current_i = current_i + 1
		if current_i == len(s) {
			current_i = 0
		}
		return next_rune
	}
}

func wrapped_offset_for(r1, r2 rune) int {
	v := offset_for(r1, r2)
	if v < 0 {
		v = 26 + v
	}
	return v
}

func offset_for(r1, r2 rune) int {
	return int(r2 - r1)
}
