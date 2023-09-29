package itswizard_m_bwStructs

import (
	"fmt"
	"github.com/itslearninggermany/itswizard_m_basic"
	"github.com/itslearninggermany/itswizard_m_imses"
	"github.com/jinzhu/gorm"
	"strconv"
)

func GetAllClassesToImortAndToUpdateAndToDelete(db *gorm.DB) (classsesToImport []BWClass, classesToUpdate []BWClass, classesToDelete []BWClass, err error) {
	err = db.Where("to_update = ? and error = ?", true, false).Find(&classesToUpdate).Error
	if err != nil {
		return
	}

	err = db.Where("to_import = ? and error = ?", true, false).Find(&classsesToImport).Error
	if err != nil {
		return
	}

	err = db.Where("to_delete = ? and success = ? and error = ?", true, false, false).Find(&classesToDelete).Error
	if err != nil {
		return
	}
	return
}

func (class *BWClass) Save(db *gorm.DB) {
	err := db.Save(&class)
	if err != nil {
		fmt.Print(err)
	}
}

func (class *BWClass) Import(itsl *itswizard_m_imses.Request, db *gorm.DB) {
	fmt.Println(class)
	fmt.Println(itsl)
	resp, err := itsl.CreateGroup(itswizard_m_basic.DbGroup15{
		SyncID:        class.GroupSyncKey,
		Name:          class.Name,
		ParentGroupID: strconv.Itoa(int(class.OrganisationID)) + "CLASS",
		Level:         2,
		IsCourse:      false,
	}, false)
	if err != nil {
		class.ErrorString = resp
		class.Error = true
		class.Save(db)
		SaveBwLog(LogStruct{
			Log:            "Fehler: Klasse " + class.Name + " konnte nicht angelegt werden.",
			OrganisationID: class.OrganisationID,
			ClassID:        class.GroupSyncKey,
			PersonSyncKey:  "",
			SchoolNumber:   class.SchoolNumber,
			Db:             db,
		})
		return
	}
	class.ToImport = false
	class.Success = true
	SaveBwLog(LogStruct{
		Log:            "Klasse " + class.Name + " wurde erfolgreich angelegt",
		OrganisationID: class.OrganisationID,
		ClassID:        class.GroupSyncKey,
		PersonSyncKey:  "",
		SchoolNumber:   class.SchoolNumber,
		Db:             db,
	})
	class.Save(db)

}

func (class *BWClass) Update(itsl *itswizard_m_imses.Request, db *gorm.DB) {
	resp, err := itsl.UpdateGroup(itswizard_m_basic.Group{
		GroupSyncKey:  class.GroupSyncKey,
		Name:          class.Name,
		ParentGroupID: "",
		IsCourse:      false,
		AzureGroupID:  strconv.Itoa(int(class.OrganisationID)) + "CLASS",
	}, false)
	if err != nil {
		class.ErrorString = resp
		class.Error = true
		class.Save(db)
		SaveBwLog(LogStruct{
			Log:            "Fehler: Klasse " + class.Name + " konnte nicht aktualisiert werden.",
			OrganisationID: class.OrganisationID,
			ClassID:        class.GroupSyncKey,
			PersonSyncKey:  "",
			SchoolNumber:   class.SchoolNumber,
			Db:             db,
		})
		return
	}
	class.ToUpdate = false
	class.Success = true
	SaveBwLog(LogStruct{
		Log:            "Klasse " + class.Name + " wurde erfolgreich aktualisiert.",
		OrganisationID: class.OrganisationID,
		ClassID:        class.GroupSyncKey,
		PersonSyncKey:  "",
		SchoolNumber:   class.SchoolNumber,
		Db:             db,
	})
	class.Save(db)
}

func (class *BWClass) Delete(itsl *itswizard_m_imses.Request, db *gorm.DB) {
	resp, err := itsl.DeleteGroup(class.GroupSyncKey)
	if err != nil {
		class.ErrorString = resp
		class.Save(db)
		SaveBwLog(LogStruct{
			Log:            "Fehler: Klasse " + class.Name + " konnte nicht gelöscht werden.",
			OrganisationID: class.OrganisationID,
			ClassID:        class.GroupSyncKey,
			PersonSyncKey:  "",
			SchoolNumber:   class.SchoolNumber,
			Db:             db,
		})
		return
	}
	class.ToDelete = true
	class.Error = false
	class.ErrorString = ""
	class.Success = true
	SaveBwLog(LogStruct{
		Log:            "Klasse " + class.Name + " konnte gelöscht werden.",
		OrganisationID: class.OrganisationID,
		ClassID:        class.GroupSyncKey,
		PersonSyncKey:  "",
		SchoolNumber:   class.SchoolNumber,
		Db:             db,
	})
	class.Save(db)
}
