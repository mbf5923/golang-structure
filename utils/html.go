package util

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"text/template"
)

type BodyRequest struct {
	To    string
	Token string
}

func ParseHtml(fileName string, data map[string]string) string {
	html, errParse := template.ParseFiles("templates/" + fileName + ".html")

	if errParse != nil {
		defer fmt.Println("parser file html failed")
		logrus.Fatal(errParse.Error())
	}

	body := BodyRequest{To: data["to"], Token: data["token"]}

	buf := new(bytes.Buffer)
	errExecute := html.Execute(buf, body)

	if errExecute != nil {
		defer fmt.Println("execute html file failed")
		logrus.Fatal(errExecute.Error())
	}

	return buf.String()
}
