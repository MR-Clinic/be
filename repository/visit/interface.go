package visit

import "be/entities"

type Visit interface {
	CreateVal(doctor_uid, patient_uid string, req entities.Visit) (entities.Visit, error)
	Update(visit_uid string, req entities.Visit) (entities.Visit, error)
	Delete(visit_uid string) (entities.Visit, error)
	GetVisitsVer1(kind, uid, status, date, grouped string) (Visits, error)
	GetVisitList(visit_uid string) (VisitCalendar, error)
}
