package CodeWars

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
)

const baseURL = "https://www.codewars.com/api/v1/"

func GetUser(username string) (RespUser, error) {
	var response RespUser

	body, err := get(fmt.Sprintf("%susers/%s", baseURL, username))
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, err
	}

	if response.Reason != "" {
		return response, errors.New(response.Reason)
	}

	return response, nil
}

func GetCompletedChallenges(username string) (RespCompletedChallenges, error) {
	resChanel := make(chan respGetListChallenges, 10)

	// Getting the first page of solved tasks
	getListChallenges(username, 0, resChanel, nil)
	firstList := <-resChanel
	if firstList.err != nil {
		return RespCompletedChallenges{}, firstList.err
	}
	challenges := firstList.challenges

	// Retrieving the remaining pages
	go func() {
		defer close(resChanel)
		wg := sync.WaitGroup{}
		// Todo рассмотреть необходимость порционных запросов
		for page := 1; page < challenges.TotalPages; page++ {
			wg.Add(1)
			go getListChallenges(username, page, resChanel, &wg)
		}
		wg.Wait()
	}()

	// Merging all pages
	for res := range resChanel {
		if res.err != nil {
			return challenges, res.err
		}
		challenges.ListChallenges = append(challenges.ListChallenges, res.challenges.ListChallenges...)
	}

	return challenges, nil
}

func GetCodeChallenge(id string) (RespCodeChallenge, error) {
	var challenge RespCodeChallenge

	body, err := get(fmt.Sprintf("%scode-challenges/%s", baseURL, id))
	if err != nil {
		if err.Error() == "not found" {
			challenge.ID = id
			challenge.Rank.Name = "no rank"
			return challenge, nil
		}
		return challenge, err
	}

	err = json.Unmarshal(body, &challenge)
	if err != nil {
		return challenge, err
	}

	if challenge.Reason != "" {
		return challenge, errors.New(challenge.Reason)
	}

	if challenge.Rank.Name == "" {
		challenge.Rank.Name = "no rank"
	}

	return challenge, nil

}

func getListChallenges(username string, page int, ch chan respGetListChallenges, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}

	var challenges RespCompletedChallenges

	body, err := get(fmt.Sprintf("%susers/%s/code-challenges/completed?page=%d", baseURL, username, page))
	if err != nil {
		ch <- respGetListChallenges{challenges, err}
	}

	err = json.Unmarshal(body, &challenges)
	if err != nil {
		ch <- respGetListChallenges{challenges, err}
	}

	// Проверка на наличие ошибок
	if challenges.Reason != "" {
		ch <- respGetListChallenges{challenges, errors.New(challenges.Reason)}
	}

	ch <- respGetListChallenges{challenges, nil}
}

func get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		textErr := fmt.Sprintf("unknown status code: %d (%s)", resp.StatusCode, resp.Status)

		switch resp.StatusCode {
		case 429:
			textErr = "many requests"
		case 404:
			textErr = "not found"
		}

		return nil, errors.New(textErr)
	}

	return body, nil
}
