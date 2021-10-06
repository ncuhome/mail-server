package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
)

// Postman 邮差的工作内容
type Postman interface {
	SendMail(string)
}

// postman 邮差要带的东西
type postman struct {
	message    Message
	PostOffice PostOffice
}

// SendMail 邮差将邮件运送到 PostOffice
func (p *postman) SendMail(to string) {
	p.message.To = to
	p.PostOffice.ReciveMail(p.message)
}

// Message 消息模板
type Message struct {
	From    string
	To      string
	Subject string
	Body    string
}

// Byte 将 Message 转换为可直接发送的字节数组
func (m *Message) Byte() []byte {
	return []byte(fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"\r\n"+
			"%s", m.From, m.To, m.Subject, m.Body))
}

// HirePostman 招聘一名专职邮差负责向 PostOffice 投递 Message
func HirePostman(message Message, PostOffice PostOffice) Postman {
	return &postman{
		message:    message,
		PostOffice: PostOffice,
	}
}

// PostOffice 邮局的工作内容
type PostOffice interface {
	ReciveMail(Message)
}

// postOffice 邮局的地址以及电子锁的帐号和密码
type postOffice struct {
	server   string
	port     string
	username string
	password string
	mails    []byte
}

func (p *postOffice) servername() string {
	return fmt.Sprintf("%s:%s", p.server, p.port)
}

// ReciveMail 接受邮件
func (m *postOffice) ReciveMail(message Message) {

	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         m.server,
	}

	conn, err := tls.Dial("tcp", m.servername(), tlsconfig)
	if err != nil {
		log.Fatal(err)
	}

	c, err := smtp.NewClient(conn, m.server)
	if err != nil {
		log.Fatal(err)
	}

	auth := smtp.PlainAuth("", m.username, "ncuhome123456", m.server)

	if err = c.Auth(auth); err != nil {
		log.Fatal(err)
	}

	if err = c.Mail(m.username); err != nil {
		log.Fatal(err)
	}

	if err = c.Rcpt(message.To); err != nil {
		log.Fatal(err)
	}

	w, err := c.Data()
	if err != nil {
		log.Fatal(err)
	}

	_, err = w.Write(message.Byte())
	if err != nil {
		log.Fatal(err)
	}
	err = w.Close()
	if err != nil {
		log.Fatal(err)
	}

	c.Quit()
}

// NewPostOffice 建立一个邮局
func NewPostOffice(server, port, username, password string) PostOffice {
	return &postOffice{
		server:   server,
		port:     port,
		username: username,
		password: password,
	}
}

func main() {
	// PostOffice Config
	server := ""
	port := ""
	username := ""
	password := ""

	// Message info
	from := ""
	subject := ""
	body := ""
	to := ""
	PostOffice := NewPostOffice(server, port, username, password)
	message := Message{
		From:    from,
		Subject: subject,
		Body:    body,
	}
	usPostman := HirePostman(message, PostOffice)
	usPostman.SendMail(to)
}
