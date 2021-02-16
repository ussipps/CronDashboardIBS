package main

import (
	_ "CronDashboardIBS/functions"
	"CronDashboardIBS/scheduler"
	_ "github.com/jasonlvhit/gocron"
	_ "github.com/joho/godotenv"
	_ "log"
	"os"
	_ "strconv"
)

func main() {
	//err := godotenv.Load()
	//if err != nil {
	//	log.Fatal("Error loading .env file")
	//}
	cmdArgs := os.Args[1]
	//IntervalTimeDashboardReal, err := strconv.ParseUint(os.Getenv("INTERVAL_TIME_DASHBOARD_REAL"), 10, 64)
	//TimeDashboardChart := functions.ParseTimeScheduler(os.Getenv("TIME_GET_DASHBOARD_CHART"))
	//TimeDashboardReport := functions.ParseTimeScheduler(os.Getenv("TIME_GET_DASHBOARD_REPORT"))
	//functions.Logger().Info("Starting Scheduler Cron Dashboard IBS")
	//if uint64(IntervalTimeDashboardReal) == 0 {
	//}else{
	//		gocron.Every(uint64(IntervalTimeDashboardReal)).Minutes().Do(scheduler.GetDataDashboardReal)
	//}
	//gocron.Every(1).Day().At(string(TimeDashboardChart)).Do(scheduler.GetDataDashboardChart)
	//gocron.Every(1).Day().At(string(TimeDashboardReport)).Do(scheduler.GetDataDashboardReport)
	if cmdArgs == "real" {
		scheduler.GetDataDashboardReal()
	} else if cmdArgs == "report" {
		scheduler.GetDataDashboardReport()
	} else if cmdArgs == "chart" {
		scheduler.GetDataDashboardChart()
	} else if cmdArgs == "charttahun" {
		scheduler.GetDataDashboardChartTahun()
	} else if cmdArgs == "chartbulan" {
		scheduler.GetDataDashboardChartBulan()
	} else {
		scheduler.GetDataDashboardReal()
	}

	//scheduler.GetDataDashboardReal()
	//scheduler.GetDataDashboardChart()
	//scheduler.GetDataDashboardReport()
	//<-gocron.Start()
}
