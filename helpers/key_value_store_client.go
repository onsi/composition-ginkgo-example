package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func KeyValueStorePinger(host string) func() error {
	return func() error {
		response, err := http.Get(host + "/ping")
		if err != nil {
			return err
		}
		defer response.Body.Close()
		if response.StatusCode != http.StatusOK {
			return fmt.Errorf("unexpected status code: %d", response.StatusCode)
		}
		return nil
	}
}

type KeyValueStoreClient struct {
	URL string
}

func (s *KeyValueStoreClient) Get(key string) (string, error) {
	response, err := http.Get(fmt.Sprintf("%s/get?key=%s", s.URL, url.QueryEscape(key)))
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	value, err := ioutil.ReadAll(response.Body)
	return string(value), err
}

func (s *KeyValueStoreClient) Set(key string, value string) error {
	response, err := http.Get(fmt.Sprintf("%s/set?key=%s&value=%s", s.URL, url.QueryEscape(key), url.QueryEscape(value)))
	if err != nil {
		return err
	}
	response.Body.Close()
	return nil
}

func (s *KeyValueStoreClient) Delete(key string) error {
	response, err := http.Get(fmt.Sprintf("%s/delete?key=%s", s.URL, url.QueryEscape(key)))
	if err != nil {
		return err
	}
	response.Body.Close()
	return nil
}

func (s *KeyValueStoreClient) GetPrefix(prefix string) ([]string, error) {
	response, err := http.Get(fmt.Sprintf("%s/get-prefix?prefix=%s", s.URL, url.QueryEscape(prefix)))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var data []string
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *KeyValueStoreClient) DeletePrefix(prefix string) error {
	response, err := http.Get(fmt.Sprintf("%s/delete-prefix?prefix=%s", s.URL, url.QueryEscape(prefix)))
	if err != nil {
		return err
	}
	response.Body.Close()
	return nil
}
