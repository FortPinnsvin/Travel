package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/fort-pinnsvin/travel/models"
	gooauth2 "github.com/golang/oauth2"
	"github.com/martini-contrib/oauth2"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"io/ioutil"
	"log"
	"net/http"
)

func GetData(session sessions.Session, tokens oauth2.Tokens, rnd render.Render) {
	tr := gooauth2.NewTransport(http.DefaultTransport, nil, &gooauth2.Token{AccessToken: tokens.Access()})
	client := http.Client{Transport: tr}
	res, err := client.Get("https://www.googleapis.com/oauth2/v1/userinfo?alt=json")
	if err != nil {
		log.Fatal(err)
	}
	robots, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	redirectUrl := "/"
	var dat map[string]string
	user := &models.User{}
	if err := json.Unmarshal(robots, &dat); err == nil {
		user.Id = string(dat["id"])
		user.FirstName = string(dat["given_name"])
		user.LastName = string(dat["family_name"])
		user.Email = string("")
		user.Avatar = string(dat["picture"])
		user.Longitude = 10000
		user.Latitude = 10000

		findUser := &models.User{}
		models.UserCollection.FindId(user.Id).One(&findUser)
		fmt.Println(findUser)

		if findUser.Id == "" {
			models.UserCollection.Insert(&user)
			redirectUrl = "/?newbie"
		} else {
			user = findUser
		}
	}
	session.Set("auth_id", user.Id)
	session.Set("first_name", user.FirstName)
	session.Set("last_name", user.LastName)
	session.Set("avatar", user.Avatar)
	session.Set("lang", user.Language)

	rnd.Redirect(redirectUrl)
}
