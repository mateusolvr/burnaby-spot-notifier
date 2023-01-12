package domain

type EmailService interface {
	SendEmailCache(activities []Activity)
	SendErrorEmailCache(err error)
	SendErrorEmail(err error)
}
