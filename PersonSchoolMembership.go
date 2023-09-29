package itswizard_m_bwStructs

import (
	"errors"
	"fmt"
	"github.com/itslearninggermany/itswizard_m_imses"
	"github.com/jinzhu/gorm"
	"strconv"
)

func GetAllPsmToImportAndPsmToUpdateAndPsmToDelete(db *gorm.DB) (psmToImport []BWPersonSchoolMembership, psmToUpdate []BWPersonSchoolMembership, psmToDelete []BWPersonSchoolMembership, err error) {
	err = db.Where("to_import = ? and error = ?", true, false).Find(&psmToImport).Error
	if err != nil {
		return
	}
	err = db.Where("to_update = ? and error = ?", true, false).Find(&psmToUpdate).Error
	if err != nil {
		return
	}
	err = db.Where("to_delete = ? and success = ? and error = ?", true, false, false).Find(&psmToDelete).Error
	if err != nil {
		return
	}
	return
}

func (psm *BWPersonSchoolMembership) Save(db *gorm.DB) {
	err := db.Save(&psm)
	if err != nil {
		fmt.Print(err)
	}
}

func (psm *BWPersonSchoolMembership) Import(itslspecial *itswizard_m_imses.Request, itsl *itswizard_m_imses.Request, db *gorm.DB) {
	if psm.Role == "Instructor" {
		resp, err := itslspecial.CreateMembership(strconv.Itoa(int(psm.OrganisationID)), psm.PersonSyncKey, "Learner")
		if err != nil {
			psm.ErrorString = resp
			psm.Error = true
			psm.Save(db)
			SaveBwLog(LogStruct{
				Log:            "Fehler: Mitgliedschaft zu" + strconv.Itoa(int(psm.OrganisationID)) + " von " + psm.PersonSyncKey + " konnte nicht im LFB importiert werden.",
				OrganisationID: psm.OrganisationID,
				ClassID:        "",
				PersonSyncKey:  psm.PersonSyncKey,
				SchoolNumber:   psm.SchoolNumber,
				Db:             db,
			})
			return
		}
	}
	var resp string
	var err error
	if psm.Role == "Instructor" {
		resp, err = itsl.CreateMembership(strconv.Itoa(int(psm.OrganisationID))+"Staff", psm.PersonSyncKey, psm.Role)
	} else {
		resp, err = itsl.CreateMembership(strconv.Itoa(int(psm.OrganisationID)), psm.PersonSyncKey, psm.Role)
	}
	if err != nil {
		psm.ErrorString = resp
		psm.Error = true
		psm.Save(db)
		SaveBwLog(LogStruct{
			Log:            "Fehler: Mitgliedschaft zu " + strconv.Itoa(int(psm.OrganisationID)) + " von " + psm.PersonSyncKey + " konnte nicht importiert werden.",
			OrganisationID: psm.OrganisationID,
			ClassID:        "",
			PersonSyncKey:  psm.PersonSyncKey,
			SchoolNumber:   psm.SchoolNumber,
			Db:             db,
		})
		return
	}
	psm.ToImport = false
	psm.Success = true
	SaveBwLog(LogStruct{
		Log:            "Mitgliedschaft zu " + strconv.Itoa(int(psm.OrganisationID)) + " von " + psm.PersonSyncKey + " wurde importiert.",
		OrganisationID: psm.OrganisationID,
		ClassID:        "",
		PersonSyncKey:  psm.PersonSyncKey,
		SchoolNumber:   psm.SchoolNumber,
		Db:             db,
	})
	psm.Save(db)
}

func (psm *BWPersonSchoolMembership) Update(itslspecial *itswizard_m_imses.Request, itsl *itswizard_m_imses.Request, db *gorm.DB) {
	if psm.Role == "Instructor" {
		resp, err := itslspecial.CreateMembership(strconv.Itoa(int(psm.OrganisationID)), psm.PersonSyncKey, "Learner")
		if err != nil {
			psm.ErrorString = resp
			psm.Error = true
			psm.Save(db)
			SaveBwLog(LogStruct{
				Log:            "Fehler: Mitgliedschaft zu " + strconv.Itoa(int(psm.OrganisationID)) + " von " + psm.PersonSyncKey + " konnte nicht im LFB importiert werden.",
				OrganisationID: psm.OrganisationID,
				ClassID:        "",
				PersonSyncKey:  psm.PersonSyncKey,
				SchoolNumber:   psm.SchoolNumber,
				Db:             db,
			})
			return
		}
	}
	resp, err := itsl.CreateMembership(strconv.Itoa(int(psm.OrganisationID)), psm.PersonSyncKey, psm.Role)
	if err != nil {
		psm.ErrorString = resp
		psm.Error = true
		psm.Save(db)
		SaveBwLog(LogStruct{
			Log:            "Fehler: Mitgliedschaft zu " + strconv.Itoa(int(psm.OrganisationID)) + " von " + psm.PersonSyncKey + " konnte nicht importiert werden.",
			OrganisationID: psm.OrganisationID,
			ClassID:        "",
			PersonSyncKey:  psm.PersonSyncKey,
			SchoolNumber:   psm.SchoolNumber,
			Db:             db,
		})
		return
	}

	psm.ToImport = false
	psm.Success = true
	SaveBwLog(LogStruct{
		Log:            "Mitgliedschaft zu " + strconv.Itoa(int(psm.OrganisationID)) + " von " + psm.PersonSyncKey + " wurde importiert.",
		OrganisationID: psm.OrganisationID,
		ClassID:        "",
		PersonSyncKey:  psm.PersonSyncKey,
		SchoolNumber:   psm.SchoolNumber,
		Db:             db,
	})
	psm.Save(db)
}

func getSchoolMembership(itsl *itswizard_m_imses.Request, schoolid string, personID string) (memid string, err error) {
	mem := itsl.ReadMembershipsForPerson(personID)

	for _, k := range mem {
		if k.GroupID == schoolid {
			memid = k.ID
		}
	}

	if memid == "" {
		err = errors.New("No Membership for " + personID + " whith " + schoolid)
	}
	return
}

func (psm *BWPersonSchoolMembership) Delete(itslspecial *itswizard_m_imses.Request, itsl *itswizard_m_imses.Request, db *gorm.DB) {
	//personID + groupID
	if psm.Role == "Instructor" {
		resp, err := itslspecial.DeleteMembership(psm.PersonSyncKey + strconv.Itoa(int(psm.OrganisationID)))
		if err != nil {
			psm.ErrorString = resp
			psm.Error = true
			psm.Save(db)
			SaveBwLog(LogStruct{
				Log:            "Fehler: Mitgliedschaft zu " + strconv.Itoa(int(psm.OrganisationID)) + " von " + psm.PersonSyncKey + " konnte nicht im LFB gelöscht werden.",
				OrganisationID: psm.OrganisationID,
				ClassID:        "",
				PersonSyncKey:  psm.PersonSyncKey,
				SchoolNumber:   psm.SchoolNumber,
				Db:             db,
			})
			return
		}
	}

	mem, err := getSchoolMembership(itsl, strconv.Itoa(int(psm.OrganisationID)), psm.PersonSyncKey)
	if err != nil {
		psm.ErrorString = err.Error()
		psm.Error = true
		psm.Save(db)
		SaveBwLog(LogStruct{
			Log:            "Fehler: Mitgliedschaft zu " + strconv.Itoa(int(psm.OrganisationID)) + " von " + psm.PersonSyncKey + " konnte nicht gelöscht werden.",
			OrganisationID: psm.OrganisationID,
			ClassID:        "",
			PersonSyncKey:  psm.PersonSyncKey,
			SchoolNumber:   psm.SchoolNumber,
			Db:             db,
		})
		return
	}

	resp, err := itsl.DeleteMembership(mem)
	if err != nil {
		psm.ErrorString = resp
		psm.Error = true
		psm.Save(db)
		SaveBwLog(LogStruct{
			Log:            "Fehler: Mitgliedschaft zu " + strconv.Itoa(int(psm.OrganisationID)) + " von " + psm.PersonSyncKey + " konnte nicht gelöscht werden.",
			OrganisationID: psm.OrganisationID,
			ClassID:        "",
			PersonSyncKey:  psm.PersonSyncKey,
			SchoolNumber:   psm.SchoolNumber,
			Db:             db,
		})
		return
	}

	psm.ToImport = false
	psm.Success = true
	SaveBwLog(LogStruct{
		Log:            "Mitgliedschaft zu " + strconv.Itoa(int(psm.OrganisationID)) + " von " + psm.PersonSyncKey + " wurde gelöscht.",
		OrganisationID: psm.OrganisationID,
		ClassID:        "",
		PersonSyncKey:  psm.PersonSyncKey,
		SchoolNumber:   psm.SchoolNumber,
		Db:             db,
	})
	psm.Save(db)
}
