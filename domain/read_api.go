package domain

type PostResponseBody struct {
	Body struct {
		CenterEvents []CentreEvent `json:"center_events"`
	} `json:"body"`
}

type CentreEvent struct {
	Events     []Event `json:"events"`
	CentreName string  `json:"center_name"`
}

type Event struct {
	Title             string `json:"title"`
	StartTime         string `json:"start_time"`
	EndTime           string `json:"end_time"`
	EventItemID       int    `json:"event_item_id"`
	ActivityDetailURL string `json:"activity_detail_url"`
}

type GetResponseBody struct {
	Body struct {
		ActivityDetail ActivityDetail `json:"activity_detail"`
	} `json:"body"`
}

type ActivityDetail struct {
	RegistrationDate struct {
		EnrollmentDate []struct {
			InternetTime string `json:"first_daytime_internet"`
		} `json:"enrollment_datetimes"`
	} `json:"meeting_and_registration_dates"`
	SpaceStatus string `json:"space_status"`
}
