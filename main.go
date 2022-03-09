package main

import (
	"be/configs"
	"be/delivery/controllers/auth"
	"be/delivery/controllers/doctor"
	"be/delivery/controllers/patient"
	"be/delivery/controllers/visit"
	"be/delivery/routes"
	authRepo "be/repository/auth"
	doctorRepo "be/repository/doctor"
	patientRepo "be/repository/patient"
	visitRepo "be/repository/visit"
	"be/utils"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func main() {
	var config = configs.GetConfig()
	var db = utils.InitDB(config)
	// log.Info(db)

	var authRepo = authRepo.New(db)
	var authCont = auth.New(authRepo)

	var doctorRepo = doctorRepo.New(db)
	var doctorCont = doctor.New(doctorRepo)

	var patientRepo = patientRepo.New(db)
	var patientCont = patient.New(patientRepo)

	var visitRepo = visitRepo.New(db)
	var visitCont = visit.New(visitRepo)

	var e = echo.New()

	routes.RoutesPath(e, authCont, doctorCont, patientCont, visitCont)

	log.Fatal(e.Start(fmt.Sprintf(":%d", config.PORT)))

}
