package itswizard_m_bwStructs

import (
	"fmt"
	"github.com/itslearninggermany/itswizard_m_basic"
	"github.com/itslearninggermany/itswizard_m_imses"
	"github.com/jinzhu/gorm"
	"strconv"
)

func GetAllSchoolsToImportAndToUpdateAndToImport(db *gorm.DB) (schoolsToImport []BWSchool, schoolsToUpdate []BWSchool, schoolsToDelete []BWSchool, err error) {
	err = db.Where("to_update = ? and error = ?", true, false).Find(&schoolsToUpdate).Error
	if err != nil {
		return
	}

	err = db.Where("to_import = ? and error = ?", true, false).Find(&schoolsToImport).Error
	if err != nil {
		return
	}

	err = db.Where("to_delete = ? and success = ? and error = ?", true, false, false).Find(&schoolsToDelete).Error
	if err != nil {
		return
	}
	return
}

func (school *BWSchool) Save(db *gorm.DB) {
	err := db.Save(&school)
	if err != nil {
		fmt.Print(err)
	}
}

func (school *BWSchool) Import(itslspezial *itswizard_m_imses.Request, itsl *itswizard_m_imses.Request, db *gorm.DB) {
	resp, err := itslspezial.CreateGroup(itswizard_m_basic.DbGroup15{
		SyncID:        strconv.Itoa(int(school.OrganisationID)),
		Name:          school.SchoolNumber + " - " + school.SchoolName,
		ParentGroupID: "staff",
		Level:         1,
		IsCourse:      false,
	}, false)
	if err != nil {
		school.Error = true
		school.ErrorString = resp
		school.Save(db)
		SaveBwLog(LogStruct{
			Log:            "Fehler: Schule " + school.SchoolName + " wurde im LFB-Cluster nicht angelegt werden.",
			OrganisationID: school.OrganisationID,
			ClassID:        "",
			PersonSyncKey:  "",
			SchoolNumber:   school.SchoolNumber,
			Db:             db,
		})
		return
	}

	resp, err = itsl.CreateGroup(itswizard_m_basic.DbGroup15{
		SyncID:        strconv.Itoa(int(school.OrganisationID)),
		Name:          school.SchoolNumber + " - " + school.SchoolName,
		ParentGroupID: "0",
		Level:         1,
		IsCourse:      false,
	}, true)
	if err != nil {
		school.Error = true
		school.ErrorString = resp
		school.Save(db)
		SaveBwLog(LogStruct{
			Log:            "Fehler: Schule " + school.SchoolName + " wurde im Cluster nicht angelegt werden.",
			OrganisationID: school.OrganisationID,
			ClassID:        "",
			PersonSyncKey:  "",
			SchoolNumber:   school.SchoolNumber,
			Db:             db,
		})
		return
	}

	resp, err = itsl.CreateGroup(itswizard_m_basic.DbGroup15{
		SyncID:        strconv.Itoa(int(school.OrganisationID)) + "CLASS",
		Name:          "Lerngruppen",
		ParentGroupID: strconv.Itoa(int(school.OrganisationID)),
		Level:         2,
		IsCourse:      false,
	}, false)
	if err != nil {
		school.Error = true
		school.ErrorString = resp
		school.Save(db)
		SaveBwLog(LogStruct{
			Log:            "Fehler: Beim Erstellen der Gruppe >>Lerngruppe<< für Schule " + school.SchoolName + ".",
			OrganisationID: school.OrganisationID,
			ClassID:        "",
			PersonSyncKey:  "",
			SchoolNumber:   school.SchoolNumber,
			Db:             db,
		})
		return
	}

	resp, err = itsl.CreateGroup(itswizard_m_basic.DbGroup15{
		SyncID:        strconv.Itoa(int(school.OrganisationID)) + "Staff",
		Name:          "Lehrkräfte und Mitarbeiter:innen",
		ParentGroupID: strconv.Itoa(int(school.OrganisationID)),
		Level:         2,
		IsCourse:      false,
	}, false)
	if err != nil {
		school.Error = true
		school.ErrorString = resp
		school.Save(db)
		SaveBwLog(LogStruct{
			Log:            "Fehler: Beim Erstellen der Gruppe >>Lehrkräfte und Mitarbeiter:innen<< für Schule " + school.SchoolName,
			OrganisationID: school.OrganisationID,
			ClassID:        "",
			PersonSyncKey:  "",
			SchoolNumber:   school.SchoolNumber,
			Db:             db,
		})
		return
	}
	school.ToImport = false
	school.Success = true
	school.Error = false
	school.ErrorString = ""
	SaveBwLog(LogStruct{
		Log:            "Schule " + school.SchoolName + " wurde erfolgreich angelegt",
		OrganisationID: school.OrganisationID,
		ClassID:        "",
		PersonSyncKey:  "",
		SchoolNumber:   school.SchoolNumber,
		Db:             db,
	})
	school.Save(db)
}

func (school *BWSchool) Update(itslspezial *itswizard_m_imses.Request, itsl *itswizard_m_imses.Request, db *gorm.DB) {
	resp, err := itslspezial.UpdateGroup(itswizard_m_basic.Group{
		GroupSyncKey:  strconv.Itoa(int(school.OrganisationID)),
		Name:          school.SchoolNumber + " - " + school.SchoolName,
		ParentGroupID: "staff",
		IsCourse:      false,
	}, false)
	if err != nil {
		school.Error = true
		school.ErrorString = resp
		school.Save(db)
		SaveBwLog(LogStruct{
			Log:            "Fehler: Schule " + school.SchoolName + " wurde im LFB-Cluster nicht upgedated.",
			OrganisationID: school.OrganisationID,
			ClassID:        "",
			PersonSyncKey:  "",
			SchoolNumber:   school.SchoolNumber,
			Db:             db,
		})
		return
	}

	//Create Group and Subgroups in the target Cluster

	resp, err = itsl.UpdateGroup(itswizard_m_basic.Group{
		GroupSyncKey:  strconv.Itoa(int(school.OrganisationID)),
		Name:          school.SchoolNumber + " - " + school.SchoolName,
		ParentGroupID: "0",
		IsCourse:      false,
	}, true)
	if err != nil {
		school.Error = true
		school.ErrorString = resp
		school.Save(db)
		SaveBwLog(LogStruct{
			Log:            "Fehler: Schule " + school.SchoolName + " wurde in Cluster nicht upgedated.",
			OrganisationID: school.OrganisationID,
			ClassID:        "",
			PersonSyncKey:  "",
			SchoolNumber:   school.SchoolNumber,
			Db:             db,
		})
		return
	}

	resp, err = itsl.UpdateGroup(itswizard_m_basic.Group{
		GroupSyncKey:  strconv.Itoa(int(school.OrganisationID)) + "CLASS",
		Name:          "Lerngruppen",
		ParentGroupID: strconv.Itoa(int(school.OrganisationID)),
		IsCourse:      false,
	}, false)
	if err != nil {
		school.Error = true
		school.ErrorString = resp
		school.Save(db)
		SaveBwLog(LogStruct{
			Log:            "Fehler: Update der Gruppe >>Lerngruppe<< für Schule " + school.SchoolName + " war nicht erfolgreich.",
			OrganisationID: school.OrganisationID,
			ClassID:        "",
			PersonSyncKey:  "",
			SchoolNumber:   school.SchoolNumber,
			Db:             db,
		})
		return
	}

	resp, err = itsl.UpdateGroup(itswizard_m_basic.Group{
		GroupSyncKey:  strconv.Itoa(int(school.OrganisationID)) + "Staff",
		Name:          "Lehrkräfte und Mitarbeiter:innen",
		ParentGroupID: strconv.Itoa(int(school.OrganisationID)),
		IsCourse:      false,
	}, false)
	if err != nil {
		school.Save(db)
		school.Error = true
		school.ErrorString = resp
		SaveBwLog(LogStruct{
			Log:            "Fehler: Update der Gruppe >>Lehrkräfte und Mitarbeiter:innen<< für Schule " + school.SchoolName + " war nicht erfolgreich",
			OrganisationID: school.OrganisationID,
			ClassID:        "",
			PersonSyncKey:  "",
			SchoolNumber:   school.SchoolNumber,
			Db:             db,
		})
		return
	}

	school.ToUpdate = false
	school.Success = true
	school.Error = false
	school.ErrorString = ""
	SaveBwLog(LogStruct{
		Log:            "Schule " + school.SchoolName + " wurde erfolgreich upgedated",
		OrganisationID: school.OrganisationID,
		ClassID:        "",
		PersonSyncKey:  "",
		SchoolNumber:   school.SchoolNumber,
		Db:             db,
	})
	school.Save(db)
}

func (school *BWSchool) Delete(itslspezial *itswizard_m_imses.Request, itsl *itswizard_m_imses.Request, db *gorm.DB) {
	/*
		Im LSB Cluster löschen
	*/
	resp, err := itslspezial.DeleteGroup(strconv.Itoa(int(school.OrganisationID)))
	if err != nil {
		school.Error = true
		school.ErrorString = resp
		school.Save(db)
		SaveBwLog(LogStruct{
			Log:            "Fehler: Die Schule " + school.SchoolName + " war im LFB-Cluster nicht gelöscht",
			OrganisationID: school.OrganisationID,
			ClassID:        "",
			PersonSyncKey:  "",
			SchoolNumber:   school.SchoolNumber,
			Db:             db,
		})
		return
	}

	/*
		Im Hauptcluster löschen
	*/
	resp, err = itsl.DeleteGroup(strconv.Itoa(int(school.OrganisationID)) + "CLASS")
	if err != nil {
		school.Error = true
		school.ErrorString = resp
		school.Save(db)
		SaveBwLog(LogStruct{
			Log:            "Fehler: Löschen der Gruppe >>Lerngruppe<< für Schule " + school.SchoolName + " war nicht erfolgreich.",
			OrganisationID: school.OrganisationID,
			ClassID:        "",
			PersonSyncKey:  "",
			SchoolNumber:   school.SchoolNumber,
			Db:             db,
		})
		return
	}

	resp, err = itsl.DeleteGroup(strconv.Itoa(int(school.OrganisationID)) + "Staff")
	if err != nil {
		school.Error = true
		school.ErrorString = resp
		school.Save(db)
		SaveBwLog(LogStruct{
			Log:            "Fehler: Löschen der Gruppe >>Lehrkräfte und Mitarbeiter:innen<< für Schule " + school.SchoolName + " war nicht erfolgreich.",
			OrganisationID: school.OrganisationID,
			ClassID:        "",
			PersonSyncKey:  "",
			SchoolNumber:   school.SchoolNumber,
			Db:             db,
		})
		return
	}

	resp, err = itsl.DeleteGroup(strconv.Itoa(int(school.OrganisationID)))
	if err != nil {
		school.Error = true
		school.ErrorString = resp
		school.Save(db)
		SaveBwLog(LogStruct{
			Log:            "Fehler: Schule " + school.SchoolName + " wurde in Cluster nicht gelöscht.",
			OrganisationID: school.OrganisationID,
			ClassID:        "",
			PersonSyncKey:  "",
			SchoolNumber:   school.SchoolNumber,
			Db:             db,
		})
		return
	}

	school.ToDelete = true
	school.Success = true
	school.Error = true
	school.ErrorString = ""
	SaveBwLog(LogStruct{
		Log:            "Schule " + school.SchoolName + " wurde erfolgreich gelöscht",
		OrganisationID: school.OrganisationID,
		ClassID:        "",
		PersonSyncKey:  "",
		SchoolNumber:   school.SchoolNumber,
		Db:             db,
	})
	school.Save(db)

}
