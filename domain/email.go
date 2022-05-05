package domain

type EmailService interface {
	SendMail(cfg Config, htmlBody string)
	BuildHtmlBody(activities []Activity) string
}
