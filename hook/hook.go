package hook

import (
	"errors"
	"net/url"
	"fmt"
	"crypto/md5"
	"encoding/json"
	"strings"
)

type objHook struct {
	Url          string `json:"url"`
	Topic        string `json:"topic"`
	Failures     int `json:"failures"`     //not clear what to do with this param
	Successes    int `json:"successes"`    //not clear what to do with this param
	Max_failures int `json:"max_failures"` //not clear what to do with this param
}

func (h *objHook) UpdateMaxFailures(newVal int) {
	h.Max_failures = newVal
}

var all_hooks map[string]*objHook

func Init() {
	all_hooks = map[string]*objHook{}
}

func CreateHook(i_url string, i_topic string, i_max_failures int) (success bool, err error) {
	success = false

	if len(i_url) == 0 {
		err = errors.New("Url length is 0!")
		return
	}

	if len(i_topic) == 0 {
		err = errors.New("Topic should length is 0!")
		return
	}

	//parse url
	u, u_err := url.Parse(i_url)
	if u_err != nil {
		err = u_err
		return
	}
	if u.Scheme != "http" {
		err = errors.New("Unsupported url scheme! http is supported only.")
		return
	}

	h := objHook{i_url, i_topic, 0, 0, i_max_failures}

	//calculate url + topic md5 to check if hook exists already
	key := fmt.Sprintf("%x", md5.Sum([]byte(i_url + i_topic)))

	if val, ok := all_hooks[string(key)]; ok {
		val.UpdateMaxFailures(i_max_failures)
	} else {
		all_hooks[key] = &h
	}
	success = true
	return
}

func GetHooks(url string, topic string) (response []byte) {
	if len(url) == 0 && len(topic) == 0 {
		response, _ = json.Marshal(all_hooks)
		return
	}

	tempHooks := map[string]objHook{}
	for k, v := range all_hooks {
		tempHooks[k] = *v
	}

	if len(url) > 0 {
		for key, value := range tempHooks {
			if strings.Compare(value.Url, url) != 0 {
				delete(tempHooks, key)
			}
		}
	}
	if len(topic) > 0 {
		for key, value := range tempHooks {
			if strings.Compare(value.Topic, topic) != 0 {
				delete(tempHooks, key)
			}
		}
	}
	response, _ = json.Marshal(tempHooks)
	return
}

func DeleteHooks(url string, topic string) {
	if len(url) == 0 && len(topic) == 0 {
		for k := range all_hooks {
			delete(all_hooks, k)
		}
		return
	}
	//we have md5 hash here :)
	if len(url) > 0 && len(topic) > 0 {
		keyHash := fmt.Sprintf("%x", md5.Sum([]byte(url + topic)))
		if _, ok := all_hooks[keyHash]; ok {
			delete(all_hooks, keyHash)
		}
	}

	if len(url) > 0 && len(topic) == 0 {
		for key, value := range all_hooks {
			if strings.Compare(value.Url, url) == 0 {
				delete(all_hooks, key)
			}
		}
	}
	if len(topic) > 0 && len(url) == 0 {
		for key, value := range all_hooks {
			if strings.Compare(value.Topic, topic) == 0 {
				delete(all_hooks, key)
			}
		}
	}
}

func PutTopics(url string, text string) (response []byte) {
	// not clear what to do with text - send somewhere?
	topic := strings.TrimPrefix(url, "/topics") ;
	urlsSlice := []string{}
	for _, v := range all_hooks {
		if strings.HasPrefix(v.Topic, topic) {
			urlsSlice = append(urlsSlice, v.Url)
		}
	}
	response, _ = json.Marshal(urlsSlice)
	return
}