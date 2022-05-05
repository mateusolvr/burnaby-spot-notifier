package readpage

import (
	"context"
	"fmt"
	"log"
	"math"
	"regexp"
	"strconv"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
	"github.com/mateusolvr/web-scraper-go/domain"
)

type service struct {
	validationService domain.ValidationService
	activities        []domain.Activity
}

func NewService(validationService domain.ValidationService) *service {
	return &service{
		validationService: validationService,
	}
}

func (s *service) InitializeCrawler() {

	// create Chrome instance
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 300*time.Second)
	defer cancel()
	start := time.Now()

	// Go to main page and select Adults category
	err := chromedp.Run(ctx,
		emulation.SetUserAgentOverride("WebScraper 1.0"),
		chromedp.Navigate(`https://webreg.city.burnaby.bc.ca/webreg/Activities/ActivitiesAdvSearch.asp`),
		chromedp.Click(`//*[@id="browser"]/li[1]/span/a`, chromedp.BySearch),
	)
	if err != nil {
		log.Fatal(err)
	}

	totalActivities := s.getTotalActivities(ctx)
	totalPages := int(math.Ceil(float64(totalActivities) / 10))
	fmt.Println("Total Pages: ", totalPages)

	for i := 0; i <= totalPages; i++ {
		fmt.Println("Page: ", i)
		s.getActivitiesPage(ctx)
		if i == totalPages {
			break
		}
		s.clickNextPage(ctx)
		time.Sleep(2 * time.Second)
	}

	// Populate day field with time type from string
	for i, v := range s.activities {
		s.activities[i].Days = s.validationService.ParseDate(v.DaysStr[:11])
	}

	fmt.Printf("\nTook: %f secs\n", time.Since(start).Seconds())

}

func (s *service) clickNextPage(ctx context.Context) {
	err := chromedp.Run(ctx,
		chromedp.Click(`//*[@id="course-tab"]/div[2]/div/div[2]/div[2]/a[3]`, chromedp.BySearch),
	)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *service) getTotalActivities(ctx context.Context) (activitiesQty int) {
	var activitiesQtyStr string
	err := chromedp.Run(ctx,
		chromedp.Text("/html/body/div[1]/div[2]/div/div[2]/div/div/div[1]/div[2]/div/div[2]/div[1]", &activitiesQtyStr),
	)
	if err != nil {
		log.Fatal(err)
	}

	r := regexp.MustCompile(`Showing [\d]+ to +[\d]+ of +([\d]+) activities`)
	str := r.FindStringSubmatch(activitiesQtyStr)

	if len(str) >= 2 {
		activitiesQty, err = strconv.Atoi(str[1])
		if err != nil {
			log.Fatal(err)
		}
		return activitiesQty
	}

	return 0
}

func (s *service) getActivitiesPage(ctx context.Context) {
	var rows []*cdp.Node
	err := chromedp.Run(ctx,
		chromedp.Nodes(`//*[@id="course-tab"]/div[2]/div/table/tbody/tr`, &rows, chromedp.BySearch),
	)
	if err != nil {
		log.Fatal(err)
	}

	var sel = `/html/body/div[1]/div[2]/div/div[2]/div/div/div[1]/div[2]/div/table/tbody/tr[%d]/td/div[1]`
	var firstCol string
	for i := 1; i <= len(rows); i++ {
		if err := chromedp.Run(ctx,
			chromedp.Text(fmt.Sprintf(sel, i), &firstCol),
		); err != nil {
			log.Fatal(err)
		}
		if s.validationService.ValidateActivity(ctx, firstCol) {
			s.checkActivityAvailability(ctx, i)
		}

	}

}

func (s *service) checkActivityAvailability(ctx context.Context, index int) {
	sel := fmt.Sprintf("/html/body/div[1]/div[2]/div/div[2]/div/div/div[1]/div[2]/div/table/tbody/tr[%d]/td/div[3]/span[1]", index)

	var rows []*cdp.Node
	err := chromedp.Run(ctx,
		chromedp.Click(sel, chromedp.BySearch),
	)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(2 * time.Second)

	err = chromedp.Run(ctx,
		chromedp.Nodes(fmt.Sprintf(`/html/body/div[1]/div[2]/div/div[2]/div/div/div[1]/div[2]/div/table/tbody/tr[%d]/td/div[3]/div/div/table/tbody/tr`, index), &rows, chromedp.BySearch),
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Number rows: ", len(rows))

	if len(rows) > 0 {
		s.getDataFromActivityTable(ctx, rows, index)
	}
}

func (s *service) getDataFromActivityTable(ctx context.Context, rows []*cdp.Node, index int) {
	var courseName, weekDay, times, days, complexName, availableSpaces string

	rootPath := fmt.Sprintf("/html/body/div[1]/div[2]/div/div[2]/div/div/div[1]/div[2]/div/table/tbody/tr[%d]", index)
	tableSel := rootPath + "/td/div[3]/div/div/table/tbody/tr%s/td[%d]"

	for i := 1; i <= len(rows); i++ {
		rowIndex := fmt.Sprintf("[%d]", i)
		if len(rows) == 1 {
			rowIndex = ""
		}
		fmt.Println(fmt.Sprintf(tableSel, rowIndex, domain.CourseNameColumn))
		err := chromedp.Run(ctx,
			chromedp.Text(fmt.Sprintf(tableSel, rowIndex, domain.CourseNameColumn), &courseName),
			chromedp.Text(fmt.Sprintf(tableSel, rowIndex, domain.WeekDayColumn), &weekDay),
			chromedp.Text(fmt.Sprintf(tableSel, rowIndex, domain.TimesColumn), &times),
			chromedp.Text(fmt.Sprintf(tableSel, rowIndex, domain.DaysColumn), &days),
			chromedp.Text(fmt.Sprintf(tableSel, rowIndex, domain.ComplexNameColumn), &complexName),
			chromedp.Text(fmt.Sprintf(tableSel, rowIndex, domain.AvailableSpacesColumn), &availableSpaces),
		)
		if err != nil {
			log.Fatal(err)
		}

		courseName, weekDay, times, days, complexName, availableSpaces = s.validationService.CleanFields(courseName, weekDay, times, days, complexName, availableSpaces)
		availableSpacesInt, err := strconv.Atoi(availableSpaces)
		if err != nil {
			log.Fatal(err)
		}

		var newActivity = domain.Activity{
			CourseName:      courseName,
			WeekDay:         weekDay,
			Times:           times,
			DaysStr:         days,
			ComplexName:     complexName,
			AvailableSpaces: availableSpacesInt}
		s.activities = append(s.activities, newActivity)

		fmt.Printf("Activity: %s\nWeekDay: %s\nTimes: %s\nDays: %s\nComplex Name: %s\nAvailable Spaces: %s\n\n",
			courseName, weekDay, times, days, complexName, availableSpaces)
	}
}
