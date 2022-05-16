package domain

type EmailService interface {
	SendEmail(cfg Config, htmlBody string)
	BuildHtmlBody(activities []Activity) string
}
