package main

import (
	"errors"
	"flag"
	"net/http"
	"regexp"
	"time"

	"github.com/eternal-flame-AD/ifttt"
	"github.com/eternal-flame-AD/ifttt-netease/internal/api"
)

var idRegexp = regexp.MustCompile(`^\d+$`)

type NewSongTrigger struct {
	key string
}

func (c NewSongTrigger) Poll(req *ifttt.TriggerPollRequest, r *ifttt.Request) (ifttt.TriggerEventCollection, error) {
	if c.key != r.UserAccessToken {
		return nil, ifttt.ErrorInvalidToken
	}
	if id, ok := req.TriggerFields["id"]; !ok {
		return nil, errors.New("Missing ID")
	} else if !idRegexp.MatchString(id) {
		return nil, errors.New("ID should only contain numbers.")
	} else {
		list, err := api.GetSongList(id)
		if err != nil {
			return nil, err
		}

		res := make(ifttt.TriggerEventCollection, 0)
		now := time.Now().Unix()
		for index, song := range list {
			res = append(res, ifttt.TriggerEvent{
				Ingredients: map[string]string{
					"url":  song.URL(),
					"name": song.Name,
				},
				Meta: ifttt.TriggerEventMeta{
					ID:   song.ID,
					Time: time.Unix(now-int64(index), 0),
				},
			})
		}
		return res, nil
	}
}

func (c NewSongTrigger) Options(req *ifttt.Request) (*ifttt.DynamicOption, error) {
	return nil, nil
}

func (c NewSongTrigger) ValidateField(fieldslug string, value string, req *ifttt.Request) error {
	if fieldslug == "id" {
		if idRegexp.MatchString(value) {
			return nil
		} else {
			return errors.New("ID should only contain numbers.")
		}
	} else {
		return errors.New("Unknown field")
	}
}

func (c NewSongTrigger) ValidateContext(values map[string]string, req *ifttt.Request) (map[string]error, error) {
	return nil, nil
}

func (c NewSongTrigger) RemoveIdentity(triggerid string) error {
	return nil
}

func (c NewSongTrigger) RealTime() bool {
	return false
}

func main() {

	key := flag.String("key", "", "IFTTT Service key")
	debug := flag.Bool("debug", false, "debug")
	flag.Parse()

	service := ifttt.Service{
		ServiceKey: *key,
		Healthy:    func() bool { return true },
		UserInfo: func(req *ifttt.Request) (*ifttt.UserInfo, error) {
			return nil, errors.New("Not supported")
		},
	}
	if *debug {
		service.EnableDebug()
	}

	new_song := NewSongTrigger{*key}

	service.RegisterTrigger("new_song", new_song)

	http.ListenAndServe(":3999", service)
}
