package domain

type EmailService interface {
	SendEmail(htmlBody string)
	BuildHtmlBody(activities []Activity) string
	SendErrorEmail(err error)
}
