package patient

import (
	"be/configs"
	"be/entities"
	"be/repository/doctor"
	"be/utils"
	"testing"

	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	var config = configs.GetConfig()
	var db = utils.InitDB(config)
	var r = New(db)
	db.Migrator().DropTable(&entities.Patient{})
	db.Migrator().DropTable(&entities.Doctor{})
	db.Migrator().DropTable(&entities.Visit{})
	db.AutoMigrate(&entities.Doctor{})
	db.AutoMigrate(&entities.Patient{})

	t.Run("succress run Create", func(t *testing.T) {
		var mock1 = entities.Patient{UserName: "anonim", Email: "anonim@", Password: "anonim", Nik: "1"}

		var res, err = r.Create(mock1)
		assert.Nil(t, err)
		assert.NotNil(t, res)
		// log.Info(res)
	})

	t.Run("success handle username", func(t *testing.T) {
		var mock = entities.Doctor{UserName: "patient2", Email: "patient@", Password: "patient"}

		if _, err := doctor.New(db).Create(mock); err != nil {
			t.Log()
			t.Fatal()
		}

		var mock1 = entities.Patient{UserName: "patient2", Email: "clinic@", Password: "clinic", Nik: "1"}

		var _, err = r.Create(mock1)
		assert.NotNil(t, err)
	})
}

func TestUpdate(t *testing.T) {
	var config = configs.GetConfig()
	var db = utils.InitDB(config)
	var r = New(db)
	db.Migrator().DropTable(&entities.Patient{})
	db.Migrator().DropTable(&entities.Doctor{})
	db.Migrator().DropTable(&entities.Visit{})
	db.AutoMigrate(&entities.Doctor{})
	db.AutoMigrate(&entities.Patient{})

	t.Run("success update", func(t *testing.T) {
		var mock1 = entities.Patient{UserName: "clinic1", Email: "clinic@", Password: "clinic", Nik: "1"}

		var res, err = r.Create(mock1)
		if err != nil {
			log.Info(err)
			t.Fatal()
		}

		mock1 = entities.Patient{Name: "clinic"}

		res, err = r.Update(res.Patient_uid, mock1)
		assert.Nil(t, err)
		assert.NotNil(t, res)
		// log.Info(res.ClinicName)
	})

	t.Run("success handle username", func(t *testing.T) {
		var mock = entities.Doctor{UserName: "patient2", Email: "patient@", Password: "patient"}

		if _, err := doctor.New(db).Create(mock); err != nil {
			t.Log()
			t.Fatal()
		}

		var mock1 = entities.Patient{UserName: "clinic5", Email: "clinic@", Password: "clinic", Nik: "1"}

		var res, err = r.Create(mock1)
		if err != nil {
			log.Info(err)
			t.Fatal()
		}

		mock1 = entities.Patient{Name: "clinic", UserName: "patient2"}

		_, err = r.Update(res.Patient_uid, mock1)
		assert.NotNil(t, err)
		// log.Info(err)
	})

	t.Run("error input uid", func(t *testing.T) {
		var mock1 = entities.Patient{UserName: "clinic7", Email: "clinic@", Password: "clinic", Nik: "1"}

		var _, err = r.Create(mock1)
		if err != nil {
			log.Info(err)
			t.Fatal()
		}

		mock1 = entities.Patient{Name: "clinic"}

		_, err = r.Update("", mock1)
		assert.NotNil(t, err)
		// log.Info(err)
	})

}

func TestDelete(t *testing.T) {
	var config = configs.GetConfig()
	var db = utils.InitDB(config)
	var r = New(db)
	db.Migrator().DropTable(&entities.Patient{})
	db.Migrator().DropTable(&entities.Doctor{})
	db.Migrator().DropTable(&entities.Visit{})
	db.AutoMigrate(&entities.Doctor{})
	db.AutoMigrate(&entities.Patient{})

	t.Run("success delete", func(t *testing.T) {
		var mock1 = entities.Patient{UserName: "clinic1", Email: "clinic@", Password: "clinic", Nik: "1"}

		var res, err = r.Create(mock1)
		if err != nil {
			log.Info(err)
			t.Fatal()
		}

		res, err = r.Delete(res.Patient_uid)
		assert.Nil(t, err)
		assert.Equal(t, true, res.DeletedAt.Valid)
		// log.Info(res.ClinicName)
	})

	t.Run("error input uid", func(t *testing.T) {
		var mock1 = entities.Patient{UserName: "clinic3", Email: "clinic@", Password: "clinic", Nik: "1"}

		var _, err = r.Create(mock1)
		if err != nil {
			log.Info(err)
			t.Fatal()
		}

		_, err = r.Delete("")
		assert.NotNil(t, err)
		// log.Info(err)
	})
}

func TestGetProfile(t *testing.T) {
	var config = configs.GetConfig()
	var db = utils.InitDB(config)
	var r = New(db)
	db.Migrator().DropTable(&entities.Patient{})
	db.Migrator().DropTable(&entities.Doctor{})
	db.Migrator().DropTable(&entities.Visit{})
	db.AutoMigrate(&entities.Doctor{})
	db.AutoMigrate(&entities.Patient{})

	t.Run("success get profile", func(t *testing.T) {
		var mock1 = entities.Patient{UserName: "clinic1", Email: "clinic@", Password: "clinic", Nik: "1"}

		var res, err = r.Create(mock1)
		if err != nil {
			log.Info(err)
			t.Fatal()
		}

		var res1, err1 = r.GetProfile(res.Patient_uid)
		assert.Nil(t, err1)
		assert.NotNil(t, res1)
		// log.Info(res1)
	})

	t.Run("error input uit", func(t *testing.T) {
		var mock1 = entities.Patient{UserName: "clinic2", Email: "clinic@", Password: "clinic", Nik: "1"}

		var _, err = r.Create(mock1)
		if err != nil {
			log.Info(err)
			t.Fatal()
		}

		var _, err1 = r.GetProfile("")
		assert.NotNil(t, err1)
		// log.Info(res1)
	})
}
