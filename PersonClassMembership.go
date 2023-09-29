package itswizard_m_bwStructs

import (
	"fmt"
	"github.com/itslearninggermany/itswizard_m_imses"
	"github.com/jinzhu/gorm"
	"strconv"
)

func GetAllPcmToImportAndPcmToUpdateAndPcmToDelete(db *gorm.DB) (pcmToImport []BWPersonClassMembership, pcmToUpdate []BWPersonClassMembership, pcmToDelete []BWPersonClassMembership, err error) {
	err = db.Where("to_import = ? and error = ?", true, false).Find(&pcmToImport).Error
	if err != nil {
		return
	}
	err = db.Where("to_update = ? and error = ?", true, false).Find(&pcmToUpdate).Error
	if err != nil {
		return
	}
	err = db.Where("to_delete = ? and success = ? and error = ?", true, false, false).Find(&pcmToDelete).Error
	if err != nil {
		return
	}
	return
}

func (pcm *BWPersonClassMembership) Save(db *gorm.DB) {
	err := db.Save(&pcm)
	if err != nil {
		fmt.Print(err)
	}
}

func (pcm *BWPersonClassMembership) Import(itsl *itswizard_m_imses.Request, db *gorm.DB) {
	resp, err := itsl.CreateMembership(pcm.GroupSyncKey, pcm.PersonSyncKey, pcm.ClassRole)
	if err != nil {
		pcm.ErrorString = resp
		pcm.Error = true
		pcm.Save(db)
		SaveBwLog(LogStruct{
			Log:            "Fehler: Mitgliedschaft zur Klasse " + strconv.Itoa(int(pcm.OrganisationID)) + " von " + pcm.PersonSyncKey + " konnte nicht importiert werden.",
			OrganisationID: pcm.OrganisationID,
			ClassID:        pcm.GroupSyncKey,
			PersonSyncKey:  pcm.PersonSyncKey,
			SchoolNumber:   pcm.SchoolNumber,
			Db:             db,
		})
		return
	}

	pcm.ToImport = false
	pcm.Success = true
	SaveBwLog(LogStruct{
		Log:            "Mitgliedschaft zur Klasse " + pcm.GroupSyncKey + " von " + pcm.PersonSyncKey + " wurde importiert.",
		OrganisationID: pcm.OrganisationID,
		ClassID:        pcm.GroupSyncKey,
		PersonSyncKey:  pcm.PersonSyncKey,
		SchoolNumber:   pcm.SchoolNumber,
		Db:             db,
	})
	pcm.Save(db)
}

func (pcm *BWPersonClassMembership) Update(itsl *itswizard_m_imses.Request, db *gorm.DB) {
	resp, err := itsl.CreateMembership(pcm.GroupSyncKey, pcm.PersonSyncKey, pcm.ClassRole)
	if err != nil {
		pcm.ErrorString = resp
		pcm.Error = true
		pcm.Save(db)
		SaveBwLog(LogStruct{
			Log:            "Fehler: Mitgliedschaft zur Klasse " + pcm.GroupSyncKey + " von " + pcm.PersonSyncKey + " wurde nicht aktualisiert.",
			OrganisationID: pcm.OrganisationID,
			ClassID:        pcm.GroupSyncKey,
			PersonSyncKey:  pcm.PersonSyncKey,
			SchoolNumber:   pcm.SchoolNumber,
			Db:             db,
		})
		return
	}

	pcm.ToImport = false
	pcm.Success = true
	SaveBwLog(LogStruct{
		Log:            "Mitgliedschaft zur Klasse " + pcm.GroupSyncKey + " von " + pcm.PersonSyncKey + " wurde aktualisiert.",
		OrganisationID: pcm.OrganisationID,
		ClassID:        pcm.GroupSyncKey,
		PersonSyncKey:  pcm.PersonSyncKey,
		SchoolNumber:   pcm.SchoolNumber,
		Db:             db,
	})
	pcm.Save(db)
}

func (pcm *BWPersonClassMembership) Delete(itsl *itswizard_m_imses.Request, db *gorm.DB) {
	//personID + groupID
	resp, err := itsl.DeleteMembership(pcm.PersonSyncKey + pcm.GroupSyncKey)
	if err != nil {
		pcm.ErrorString = resp
		pcm.Error = true
		pcm.Save(db)
		SaveBwLog(LogStruct{
			Log:            "Fehler: Mitgliedschaft zur Klasse " + pcm.GroupSyncKey + " von " + pcm.PersonSyncKey + " konnte nicht gelöscht werden.",
			OrganisationID: pcm.OrganisationID,
			ClassID:        pcm.GroupSyncKey,
			PersonSyncKey:  pcm.PersonSyncKey,
			SchoolNumber:   pcm.SchoolNumber,
			Db:             db,
		})
		return
	}

	pcm.ToImport = false
	pcm.Success = true
	SaveBwLog(LogStruct{
		Log:            "Mitgliedschaft zu " + pcm.GroupSyncKey + " von " + pcm.PersonSyncKey + " wurde gelöscht.",
		OrganisationID: pcm.OrganisationID,
		ClassID:        pcm.GroupSyncKey,
		PersonSyncKey:  pcm.PersonSyncKey,
		SchoolNumber:   pcm.SchoolNumber,
		Db:             db,
	})
	pcm.Save(db)
}
