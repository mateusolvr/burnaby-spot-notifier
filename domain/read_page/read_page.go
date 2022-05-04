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
			s.checkAvailability(ctx, i)
		}

	}

}

func (s *service) checkAvailability(ctx context.Context, index int) {
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

	// var res string
	// if len(rows) > 0 {
	// 	err := chromedp.Run(ctx,
	// 		chromedp.OuterHTML(fmt.Sprintf(`/html/body/div[1]/div[2]/div/div[2]/div/div/div[1]/div[2]/div/table/tbody/tr[%d]/td/div[3]/div/div/table/tbody/tr`, index), &res),
	// 	)
	// 	if err != nil {
	// 		log.Fatal(err.Error())
	// 	}

	// 	fmt.Println(res)
	// }

	s.getDataFromActivityTable(ctx, rows, index)
}

func (s *service) getDataFromActivityTable(ctx context.Context, rows []*cdp.Node, index int) (courseName, weekDay, times, days, complexName, availableSpaces string) {
	const courseNameColumn int = 1
	const weekDayColumn int = 3
	const timesColumn int = 4
	const daysColumn int = 5
	const complexNameColumn int = 6
	const availableSpacesColumn int = 8

	tableSingleRowSel := "/html/body/div[1]/div[2]/div/div[2]/div/div/div[1]/div[2]/div/table/tbody/tr[%d]/td/div[3]/div/div/table/tbody/tr/td[%d]"
	tableMultRowSel := "/html/body/div[1]/div[2]/div/div[2]/div/div/div[1]/div[2]/div/table/tbody/tr[%d]/td/div[3]/div/div/table/tbody/tr[%d]/td[%d]"

	if len(rows) == 1 {
		err := chromedp.Run(ctx,
			chromedp.Text(fmt.Sprintf(tableSingleRowSel, index, courseNameColumn), &courseName),
			chromedp.Text(fmt.Sprintf(tableSingleRowSel, index, weekDayColumn), &weekDay),
			chromedp.Text(fmt.Sprintf(tableSingleRowSel, index, timesColumn), &times),
			chromedp.Text(fmt.Sprintf(tableSingleRowSel, index, daysColumn), &days),
			chromedp.Text(fmt.Sprintf(tableSingleRowSel, index, complexNameColumn), &complexName),
			chromedp.Text(fmt.Sprintf(tableSingleRowSel, index, availableSpacesColumn), &availableSpaces),
		)
		if err != nil {
			log.Fatal(err)
		}
		courseName = s.validationService.CleanString(courseName)
		weekDay = s.validationService.CleanString(weekDay)
		times = s.validationService.CleanString(times)
		days = s.validationService.CleanString(days)
		complexName = s.validationService.CleanString(complexName)
		availableSpaces = s.validationService.CleanString(availableSpaces)

		fmt.Printf("Activity: %s\nWeekDay: %s\nTimes: %s\nDays: %s\nComplex Name: %s\nAvailable Spaces: %s\n\n",
			courseName, weekDay, times, days, complexName, availableSpaces)
	}

	if len(rows) > 1 {
		for i := 1; i <= len(rows); i++ {
			err := chromedp.Run(ctx,
				chromedp.Text(fmt.Sprintf(tableMultRowSel, index, i, courseNameColumn), &courseName),
				chromedp.Text(fmt.Sprintf(tableMultRowSel, index, i, weekDayColumn), &weekDay),
				chromedp.Text(fmt.Sprintf(tableMultRowSel, index, i, timesColumn), &times),
				chromedp.Text(fmt.Sprintf(tableMultRowSel, index, i, daysColumn), &days),
				chromedp.Text(fmt.Sprintf(tableMultRowSel, index, i, complexNameColumn), &complexName),
				chromedp.Text(fmt.Sprintf(tableMultRowSel, index, i, availableSpacesColumn), &availableSpaces),
			)
			if err != nil {
				log.Fatal(err)
			}

			courseName = s.validationService.CleanString(courseName)
			weekDay = s.validationService.CleanString(weekDay)
			times = s.validationService.CleanString(times)
			days = s.validationService.CleanString(days)
			complexName = s.validationService.CleanString(complexName)
			availableSpaces = s.validationService.CleanString(availableSpaces)

			fmt.Printf("Activity: %s\nWeekDay: %s\nTimes: %s\nDays: %s\nComplex Name: %s\nAvailable Spaces: %s\n\n",
				courseName, weekDay, times, days, complexName, availableSpaces)
		}
	}
	return
}
