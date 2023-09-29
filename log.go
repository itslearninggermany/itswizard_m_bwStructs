package itswizard_m_bwStructs

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

func SaveBwLog(input LogStruct) {
	l := BWLog{
		Log:            input.Log,
		OrganisationID: input.OrganisationID,
		ClassID:        input.ClassID,
		PersonSyncKey:  input.PersonSyncKey,
		SchoolNumber:   input.SchoolNumber,
	}
	err := input.Db.Save(&l).Error
	if err != nil {
		fmt.Println(err)
	}
}

type LogStruct struct {
	Log            string
	OrganisationID uint
	ClassID        string
	PersonSyncKey  string
	SchoolNumber   string
	Db             *gorm.DB
}
