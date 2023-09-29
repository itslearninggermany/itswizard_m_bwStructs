package itswizard_m_bwStructs

import (
	"fmt"
	"github.com/itslearninggermany/itswizard_m_basic"
	"github.com/itslearninggermany/itswizard_m_imses"
	"github.com/jinzhu/gorm"
	"strconv"
	"strings"
)

func GetAllPerosnsToUpdateAndImportAndDelete(db *gorm.DB) (personsToImport []BWPerson, personsToUpdate []BWPerson, personsToDelete []BWPerson, err error) {
	err = db.Where("to_update = ? and error = ?", true, false).Find(&personsToUpdate).Error
	if err != nil {
		return
	}

	err = db.Where("to_import = ? and error = ?", true, false).Find(&personsToImport).Error
	if err != nil {
		return
	}

	err = db.Where("to_delete = ? and success = ? and error = ?", true, false, false).Find(&personsToDelete).Error
	if err != nil {
		return
	}
	return
}

func (person *BWPerson) Save(db *gorm.DB) {
	err := db.Save(&person)
	if err != nil {
		fmt.Print(err)
	}
}

func (person *BWPerson) Import(itslspecial *itswizard_m_imses.Request, itsl *itswizard_m_imses.Request, db *gorm.DB) {
	if person.Profile == "Staff" {
		resp, err := itslspecial.CreatePerson(itswizard_m_basic.DbPerson15{
			SyncPersonKey: person.PersonSyncKey,
			FirstName:     person.FirstName,
			LastName:      person.LastName,
			Username:      person.Username,
			Profile:       "Student",
			Email:         person.Email,
			Phone:         person.Phone,
			Mobile:        person.Mobile,
			Street1:       person.Street1,
			Street2:       person.Street2,
			Postcode:      person.Postcode,
			City:          person.City,
		})
		if err != nil {
			person.ErrorString = resp
			person.Error = true
			person.Save(db)
			SaveBwLog(LogStruct{
				Log:            "Fehler: Person " + person.FirstName + " " + person.LastName + " (" + person.PersonSyncKey + ") konnte nicht auf LFB importiert werden.",
				OrganisationID: person.OrganisationID,
				ClassID:        "",
				PersonSyncKey:  person.PersonSyncKey,
				SchoolNumber:   person.SchoolNumber,
				Db:             db,
			})
			return
		}
	}

	resp, err := itsl.CreatePerson(itswizard_m_basic.DbPerson15{
		SyncPersonKey: person.PersonSyncKey,
		FirstName:     person.FirstName,
		LastName:      person.LastName,
		Username:      person.Username,
		Profile:       person.Profile,
		Email:         person.Email,
		Phone:         person.Phone,
		Mobile:        person.Mobile,
		Street1:       person.Street1,
		Street2:       person.Street2,
		Postcode:      person.Postcode,
		City:          person.City,
	})
	if err != nil {
		person.ErrorString = resp
		person.Error = true
		person.Save(db)
		SaveBwLog(LogStruct{
			Log:            "Fehler: Person " + person.FirstName + " " + person.LastName + " (" + person.PersonSyncKey + ") konnte nicht importiert werden.",
			OrganisationID: person.OrganisationID,
			ClassID:        "",
			PersonSyncKey:  person.PersonSyncKey,
			SchoolNumber:   person.SchoolNumber,
			Db:             db,
		})
		return
	}
	person.ErrorString = ""
	person.Error = false
	person.ToImport = false
	person.Success = true
	SaveBwLog(LogStruct{
		Log:            "Person " + person.FirstName + " " + person.LastName + " (" + person.PersonSyncKey + ") wurde importiert.",
		OrganisationID: person.OrganisationID,
		ClassID:        "",
		PersonSyncKey:  person.PersonSyncKey,
		SchoolNumber:   person.SchoolNumber,
		Db:             db,
	})
	person.Save(db)
}

func (person *BWPerson) Update(itslspecial *itswizard_m_imses.Request, itsl *itswizard_m_imses.Request, db *gorm.DB) {
	if person.Profile == "Staff" {
		if person.UpdateBirthday {
			g := strings.Split(person.Birthday, "-")
			if len(g) < 4 {
				year, err := strconv.Atoi(g[0])
				month, err := strconv.Atoi(g[1])
				day, err := strconv.Atoi(g[2])
				if err != nil {
					person.ErrorString = "Wrong Birthdayformat it must be YYYY-MM-DD"
					person.Error = true
					person.Save(db)
					SaveBwLog(LogStruct{
						Log:            "Fehler: Das Geburtstag von " + person.FirstName + " " + person.LastName + " (" + person.PersonSyncKey + ") konnte nicht in LFB aktualisiert werden.",
						OrganisationID: person.OrganisationID,
						ClassID:        "",
						PersonSyncKey:  person.PersonSyncKey,
						SchoolNumber:   person.SchoolNumber,
						Db:             db,
					})
					return
				}
				resp, err := itslspecial.UpdateBirthday(person.PersonSyncKey, uint(day), uint(month), uint(year))
				if err != nil {
					person.ErrorString = resp
					person.Error = true
					person.Save(db)
					SaveBwLog(LogStruct{
						Log:            "Fehler: Das Geburtstag von " + person.FirstName + " " + person.LastName + " (" + person.PersonSyncKey + ") konnte nicht in LFB aktualisiert werden.",
						OrganisationID: person.OrganisationID,
						ClassID:        "",
						PersonSyncKey:  person.PersonSyncKey,
						SchoolNumber:   person.SchoolNumber,
						Db:             db,
					})
					return
				}
			} else {
				person.ErrorString = "Wrong Birthdayformat it must be YYYY-MM-DD"
				person.Error = true
				person.Save(db)
				SaveBwLog(LogStruct{
					Log:            "Fehler: Das Geburtstag von " + person.FirstName + " " + person.LastName + " (" + person.PersonSyncKey + ") konnte nicht in LFB aktualisiert werden.",
					OrganisationID: person.OrganisationID,
					ClassID:        "",
					PersonSyncKey:  person.PersonSyncKey,
					SchoolNumber:   person.SchoolNumber,
					Db:             db,
				})
				return
			}
		}
		if person.UpdateFirstname {
			resp, err := itslspecial.UpdateFirstName(person.PersonSyncKey, person.FirstName)
			if err != nil {
				person.ErrorString = resp
				person.Error = true
				person.Save(db)
				SaveBwLog(LogStruct{
					Log:            "Fehler: Der Vorname von " + person.FirstName + " " + person.LastName + " (" + person.PersonSyncKey + ") konnte nicht in LFB aktualisiert werden.",
					OrganisationID: person.OrganisationID,
					ClassID:        "",
					PersonSyncKey:  person.PersonSyncKey,
					SchoolNumber:   person.SchoolNumber,
					Db:             db,
				})
				return
			}
		}
		if person.UpdateLastname {
			resp, err := itslspecial.UpdateLastName(person.PersonSyncKey, person.LastName)
			if err != nil {
				person.ErrorString = resp
				person.Error = true
				person.Save(db)
				SaveBwLog(LogStruct{
					Log:            "Fehler: Der Nachname von " + person.FirstName + " " + person.LastName + " (" + person.PersonSyncKey + ") konnte nicht in LFB aktualisiert werden.",
					OrganisationID: person.OrganisationID,
					ClassID:        "",
					PersonSyncKey:  person.PersonSyncKey,
					SchoolNumber:   person.SchoolNumber,
					Db:             db,
				})
				return
			}
		}
		if person.UpdateUsername {
			resp, err := itslspecial.UpdateUsername(person.PersonSyncKey, person.Username)
			if err != nil {
				person.ErrorString = resp
				person.Error = true
				person.Save(db)
				SaveBwLog(LogStruct{
					Log:            "Fehler: Der Username von " + person.FirstName + " " + person.LastName + " (" + person.PersonSyncKey + ") konnte nicht in LFB aktualisiert werden.",
					OrganisationID: person.OrganisationID,
					ClassID:        "",
					PersonSyncKey:  person.PersonSyncKey,
					SchoolNumber:   person.SchoolNumber,
					Db:             db,
				})
				return
			}
		}
		if person.UpdateProfile {
			resp, err := itslspecial.UpdateProfile(person.PersonSyncKey, "Student")
			if err != nil {
				person.ErrorString = resp
				person.Error = true
				person.Save(db)
				SaveBwLog(LogStruct{
					Log:            "Fehler: Das Profil von " + person.FirstName + " " + person.LastName + " (" + person.PersonSyncKey + ") konnte nicht in LFB aktualisiert werden.",
					OrganisationID: person.OrganisationID,
					ClassID:        "",
					PersonSyncKey:  person.PersonSyncKey,
					SchoolNumber:   person.SchoolNumber,
					Db:             db,
				})
				return
			}
		}
		if person.UpdateEmail {
			resp, err := itslspecial.UpdateEmail(person.PersonSyncKey, person.Email)
			if err != nil {
				person.ErrorString = resp
				person.Error = true
				person.Save(db)
				SaveBwLog(LogStruct{
					Log:            "Fehler: Die Email von " + person.FirstName + " " + person.LastName + " (" + person.PersonSyncKey + ") konnte nicht in LFB aktualisiert werden.",
					OrganisationID: person.OrganisationID,
					ClassID:        "",
					PersonSyncKey:  person.PersonSyncKey,
					SchoolNumber:   person.SchoolNumber,
					Db:             db,
				})
				return
			}
		}
		if person.UpdatePhone {
			resp, err := itslspecial.UpdatePhone(person.PersonSyncKey, person.Phone)
			if err != nil {
				person.ErrorString = resp
				person.Error = true
				person.Save(db)
				SaveBwLog(LogStruct{
					Log:            "Fehler: Die Festnetznummr von " + person.FirstName + " " + person.LastName + " (" + person.PersonSyncKey + ") konnte nicht in LFB aktualisiert werden.",
					OrganisationID: person.OrganisationID,
					ClassID:        "",
					PersonSyncKey:  person.PersonSyncKey,
					SchoolNumber:   person.SchoolNumber,
					Db:             db,
				})
				return
			}
		}
		if person.UpdateMobile {
			resp, err := itslspecial.UpdateMobilePhone(person.PersonSyncKey, person.Mobile)
			if err != nil {
				person.ErrorString = resp
				person.Error = true
				person.Save(db)
				SaveBwLog(LogStruct{
					Log:            "Fehler: Die Mobilnummer von " + person.FirstName + " " + person.LastName + " (" + person.PersonSyncKey + ") konnte nicht in LFB aktualisiert werden.",
					OrganisationID: person.OrganisationID,
					ClassID:        "",
					PersonSyncKey:  person.PersonSyncKey,
					SchoolNumber:   person.SchoolNumber,
					Db:             db,
				})
				return
			}
		}
		if person.UpdateStreet1 {
			resp, err := itslspecial.UpdateAdresse(person.PersonSyncKey, person.Postcode, person.City, person.Street1, person.Street2)
			if err != nil {
				person.ErrorString = resp
				person.Error = true
				person.Save(db)
				SaveBwLog(LogStruct{
					Log:            "Fehler: Die Straße1 von " + person.FirstName + " " + person.LastName + " (" + person.PersonSyncKey + ") konnte nicht in LFB aktualisiert werden.",
					OrganisationID: person.OrganisationID,
					ClassID:        "",
					PersonSyncKey:  person.PersonSyncKey,
					SchoolNumber:   person.SchoolNumber,
					Db:             db,
				})
				SaveBwLog(LogStruct{
					Log:            "Fehler: Die Straße1 von " + person.FirstName + " " + person.LastName + " (" + person.PersonSyncKey + ") konnte nicht in LFB aktualisiert werden.",
					OrganisationID: person.OrganisationID,
					ClassID:        "",
					PersonSyncKey:  person.PersonSyncKey,
					SchoolNumber:   person.SchoolNumber,
					Db:             db,
				})
				return
			}
		}
		if person.UpdateStreet2 {
			resp, err := itslspecial.UpdateAdresse(person.PersonSyncKey, person.Postcode, person.City, person.Street1, person.Street2)
			if err != nil {
				person.ErrorString = resp
				person.Error = true
				person.Save(db)
				SaveBwLog(LogStruct{
					Log:            "Fehler: Die Straße2 von " + person.FirstName + " " + person.LastName + " (" + person.PersonSyncKey + ") konnte nicht in LFB aktualisiert werden.",
					OrganisationID: person.OrganisationID,
					ClassID:        "",
					PersonSyncKey:  person.PersonSyncKey,
					SchoolNumber:   person.SchoolNumber,
					Db:             db,
				})
				return
			}
		}
		if person.UpdatePostcode {
			resp, err := itslspecial.UpdateAdresse(person.PersonSyncKey, person.Postcode, person.City, person.Street1, person.Street2)
			if err != nil {
				person.ErrorString = resp
				person.Error = true
				person.Save(db)
				SaveBwLog(LogStruct{
					Log:            "Fehler: Die Postleitzahl von " + person.FirstName + " " + person.LastName + " (" + person.PersonSyncKey + ") konnte nicht in LFB aktualisiert werden.",
					OrganisationID: person.OrganisationID,
					ClassID:        "",
					PersonSyncKey:  person.PersonSyncKey,
					SchoolNumber:   person.SchoolNumber,
					Db:             db,
				})
				return
			}
		}
		if person.UpdateCity {
			resp, err := itslspecial.UpdateAdresse(person.PersonSyncKey, person.Postcode, person.City, person.Street1, person.Street2)
			if err != nil {
				person.ErrorString = resp
				person.Error = true
				person.Save(db)
				SaveBwLog(LogStruct{
					Log:            "Fehler: Die Stadt von " + person.FirstName + " " + person.LastName + " (" + person.PersonSyncKey + ") konnte nicht in LFB aktualisiert werden.",
					OrganisationID: person.OrganisationID,
					ClassID:        "",
					PersonSyncKey:  person.PersonSyncKey,
					SchoolNumber:   person.SchoolNumber,
					Db:             db,
				})
				return
			}
		}
	}
	if person.UpdateBirthday {
		g := strings.Split(person.Birthday, "-")
		if len(g) < 4 {
			year, err := strconv.Atoi(g[0])
			month, err := strconv.Atoi(g[1])
			day, err := strconv.Atoi(g[2])
			if err != nil {
				person.ErrorString = "Wrong Birthdayformat it must be YYYY-MM-DD"
				person.Error = true
				person.Save(db)
				SaveBwLog(LogStruct{
					Log:            "Fehler: Das Geburtstag von " + person.FirstName + " " + person.LastName + " (" + person.PersonSyncKey + ") konnte nicht in LFB aktualisiert werden.",
					OrganisationID: person.OrganisationID,
					ClassID:        "",
					PersonSyncKey:  person.PersonSyncKey,
					SchoolNumber:   person.SchoolNumber,
					Db:             db,
				})
				return
			}
			resp, err := itslspecial.UpdateBirthday(person.PersonSyncKey, uint(day), uint(month), uint(year))
			if err != nil {
				person.ErrorString = resp
				person.Error = true
				person.Save(db)
				SaveBwLog(LogStruct{
					Log:            "Fehler: Das Geburtstag von " + person.FirstName + " " + person.LastName + " (" + person.PersonSyncKey + ") konnte nicht aktualisiert werden.",
					OrganisationID: person.OrganisationID,
					ClassID:        "",
					PersonSyncKey:  person.PersonSyncKey,
					SchoolNumber:   person.SchoolNumber,
					Db:             db,
				})
				return
			}
		} else {
			person.ErrorString = "Wrong Birthdayformat it must be YYYY-MM-DD"
			person.Error = true
			person.Save(db)
			SaveBwLog(LogStruct{
				Log:            "Fehler: Das Geburtstag von " + person.FirstName + " " + person.LastName + " (" + person.PersonSyncKey + ") konnte nicht aktualisiert werden.",
				OrganisationID: person.OrganisationID,
				ClassID:        "",
				PersonSyncKey:  person.PersonSyncKey,
				SchoolNumber:   person.SchoolNumber,
				Db:             db,
			})
			return
		}
	}
	if person.UpdateFirstname {
		resp, err := itsl.UpdateFirstName(person.PersonSyncKey, person.FirstName)
		if err != nil {
			person.ErrorString = resp
			person.Error = true
			person.Save(db)
			SaveBwLog(LogStruct{
				Log:            "Fehler: Der Vorname von " + person.FirstName + " " + person.LastName + " (" + person.PersonSyncKey + ") konnte nicht aktualisiert werden.",
				OrganisationID: person.OrganisationID,
				ClassID:        "",
				PersonSyncKey:  person.PersonSyncKey,
				SchoolNumber:   person.SchoolNumber,
				Db:             db,
			})
			return
		}
	}
	if person.UpdateLastname {
		resp, err := itsl.UpdateLastName(person.PersonSyncKey, person.LastName)
		if err != nil {
			person.ErrorString = resp
			person.Error = true
			person.Save(db)
			SaveBwLog(LogStruct{
				Log:            "Fehler: Der Nachname von " + person.FirstName + " " + person.LastName + " (" + person.PersonSyncKey + ") konnte nicht aktualisiert werden.",
				OrganisationID: person.OrganisationID,
				ClassID:        "",
				PersonSyncKey:  person.PersonSyncKey,
				SchoolNumber:   person.SchoolNumber,
				Db:             db,
			})
			return
		}
	}
	if person.UpdateUsername {
		resp, err := itsl.UpdateUsername(person.PersonSyncKey, person.Username)
		if err != nil {
			person.ErrorString = resp
			person.Error = true
			person.Save(db)
			SaveBwLog(LogStruct{
				Log:            "Fehler: Der Username von " + person.FirstName + " " + person.LastName + " (" + person.PersonSyncKey + ") konnte nicht aktualisiert werden.",
				OrganisationID: person.OrganisationID,
				ClassID:        "",
				PersonSyncKey:  person.PersonSyncKey,
				SchoolNumber:   person.SchoolNumber,
				Db:             db,
			})
			return
		}
	}
	if person.UpdateProfile {
		resp, err := itsl.UpdateProfile(person.PersonSyncKey, person.Profile)
		if err != nil {
			person.ErrorString = resp
			person.Error = true
			person.Save(db)
			SaveBwLog(LogStruct{
				Log:            "Fehler: Das Profile von " + person.FirstName + " " + person.LastName + " (" + person.PersonSyncKey + ") konnte nicht aktualisiert werden.",
				OrganisationID: person.OrganisationID,
				ClassID:        "",
				PersonSyncKey:  person.PersonSyncKey,
				SchoolNumber:   person.SchoolNumber,
				Db:             db,
			})
			return
		}
	}
	if person.UpdateEmail {
		resp, err := itsl.UpdateEmail(person.PersonSyncKey, person.Email)
		if err != nil {
			person.ErrorString = resp
			person.Error = true
			person.Save(db)
			SaveBwLog(LogStruct{
				Log:            "Fehler: Die Email von " + person.FirstName + " " + person.LastName + " (" + person.PersonSyncKey + ") konnte nicht aktualisiert werden.",
				OrganisationID: person.OrganisationID,
				ClassID:        "",
				PersonSyncKey:  person.PersonSyncKey,
				SchoolNumber:   person.SchoolNumber,
				Db:             db,
			})
			return
		}
	}
	if person.UpdatePhone {
		resp, err := itsl.UpdatePhone(person.PersonSyncKey, person.Phone)
		if err != nil {
			person.ErrorString = resp
			person.Error = true
			person.Save(db)
			SaveBwLog(LogStruct{
				Log:            "Fehler: Die Festnetznummer von " + person.FirstName + " " + person.LastName + " (" + person.PersonSyncKey + ") konnte nicht aktualisiert werden.",
				OrganisationID: person.OrganisationID,
				ClassID:        "",
				PersonSyncKey:  person.PersonSyncKey,
				SchoolNumber:   person.SchoolNumber,
				Db:             db,
			})
			return
		}
	}
	if person.UpdateMobile {
		resp, err := itsl.UpdateMobilePhone(person.PersonSyncKey, person.Mobile)
		if err != nil {
			person.ErrorString = resp
			person.Error = true
			person.Save(db)
			SaveBwLog(LogStruct{
				Log:            "Fehler: Die Mobilnummer von " + person.FirstName + " " + person.LastName + " (" + person.PersonSyncKey + ") konnte nicht aktualisiert werden.",
				OrganisationID: person.OrganisationID,
				ClassID:        "",
				PersonSyncKey:  person.PersonSyncKey,
				SchoolNumber:   person.SchoolNumber,
				Db:             db,
			})
			return
		}
	}
	if person.UpdateStreet1 {
		resp, err := itsl.UpdateAdresse(person.PersonSyncKey, person.Postcode, person.City, person.Street1, person.Street2)
		if err != nil {
			person.ErrorString = resp
			person.Error = true
			person.Save(db)
			SaveBwLog(LogStruct{
				Log:            "Fehler: Die Straße1 von " + person.FirstName + " " + person.LastName + " (" + person.PersonSyncKey + ") konnte nicht aktualisiert werden.",
				OrganisationID: person.OrganisationID,
				ClassID:        "",
				PersonSyncKey:  person.PersonSyncKey,
				SchoolNumber:   person.SchoolNumber,
				Db:             db,
			})
			return
		}
	}
	if person.UpdateStreet2 {
		resp, err := itsl.UpdateAdresse(person.PersonSyncKey, person.Postcode, person.City, person.Street1, person.Street2)
		if err != nil {
			person.ErrorString = resp
			person.Error = true
			person.Save(db)
			SaveBwLog(LogStruct{
				Log:            "Fehler: Die Straße2 von " + person.FirstName + " " + person.LastName + " (" + person.PersonSyncKey + ") konnte nicht aktualisiert werden.",
				OrganisationID: person.OrganisationID,
				ClassID:        "",
				PersonSyncKey:  person.PersonSyncKey,
				SchoolNumber:   person.SchoolNumber,
				Db:             db,
			})
			return
		}
	}
	if person.UpdatePostcode {
		resp, err := itsl.UpdateAdresse(person.PersonSyncKey, person.Postcode, person.City, person.Street1, person.Street2)
		if err != nil {
			person.ErrorString = resp
			person.Error = true
			person.Save(db)
			SaveBwLog(LogStruct{
				Log:            "Fehler: Die Postleitzahl von " + person.FirstName + " " + person.LastName + " (" + person.PersonSyncKey + ") konnte nicht aktualisiert werden.",
				OrganisationID: person.OrganisationID,
				ClassID:        "",
				PersonSyncKey:  person.PersonSyncKey,
				SchoolNumber:   person.SchoolNumber,
				Db:             db,
			})
			return
		}
	}
	if person.UpdateCity {
		resp, err := itsl.UpdateAdresse(person.PersonSyncKey, person.Postcode, person.City, person.Street1, person.Street2)
		if err != nil {
			person.ErrorString = resp
			person.Error = true
			person.Save(db)
			SaveBwLog(LogStruct{
				Log:            "Fehler: Die Stadt von " + person.FirstName + " " + person.LastName + " (" + person.PersonSyncKey + ") konnte nicht aktualisiert werden.",
				OrganisationID: person.OrganisationID,
				ClassID:        "",
				PersonSyncKey:  person.PersonSyncKey,
				SchoolNumber:   person.SchoolNumber,
				Db:             db,
			})
			return
		}
	}
	person.ToUpdate = false
	person.Error = false
	person.ErrorString = ""
	person.UpdateBirthday = false
	person.UpdateFirstname = false
	person.UpdateLastname = false
	person.UpdateUsername = false
	person.UpdateProfile = false
	person.UpdateEmail = false
	person.UpdatePhone = false
	person.UpdateMobile = false
	person.UpdateStreet1 = false
	person.UpdateStreet2 = false
	person.UpdatePostcode = false
	person.UpdateCity = false
	person.Success = true
	SaveBwLog(LogStruct{
		Log:            "Person " + person.FirstName + " " + person.LastName + " (" + person.PersonSyncKey + ") wurde aktualisiert.",
		OrganisationID: person.OrganisationID,
		ClassID:        "",
		PersonSyncKey:  person.PersonSyncKey,
		SchoolNumber:   person.SchoolNumber,
		Db:             db,
	})
	person.Save(db)
}

func (person *BWPerson) Delete(itslspecial *itswizard_m_imses.Request, itsl *itswizard_m_imses.Request, db *gorm.DB) {
	if person.Profile == "Staff" {
		resp, err := itslspecial.DeletePerson(person.PersonSyncKey)
		if err != nil {
			person.ErrorString = resp
			person.Error = true
			person.Save(db)
			SaveBwLog(LogStruct{
				Log:            "Fehler: Person " + person.FirstName + " " + person.LastName + " (" + person.PersonSyncKey + ") konnte nicht in LFB gelöscht werden.",
				OrganisationID: person.OrganisationID,
				ClassID:        "",
				PersonSyncKey:  person.PersonSyncKey,
				SchoolNumber:   person.SchoolNumber,
				Db:             db,
			})
			return
		}
	}
	resp, err := itsl.DeletePerson(person.PersonSyncKey)
	if err != nil {
		person.ErrorString = resp
		person.Error = true
		person.Save(db)
		SaveBwLog(LogStruct{
			Log:            "Fehler: Person " + person.FirstName + " " + person.LastName + " (" + person.PersonSyncKey + ") konnte nicht gelöscht werden.",
			OrganisationID: person.OrganisationID,
			ClassID:        "",
			PersonSyncKey:  person.PersonSyncKey,
			SchoolNumber:   person.SchoolNumber,
			Db:             db,
		})
		return
	}
	person.ErrorString = ""
	person.Error = false
	person.ToDelete = true
	person.Success = true
	fmt.Println("Person "+person.FirstName+" "+person.LastName+" ("+person.PersonSyncKey+") wurde gelöscht.", person)

	SaveBwLog(LogStruct{
		Log:            "Person " + person.FirstName + " " + person.LastName + " (" + person.PersonSyncKey + ") wurde gelöscht.",
		OrganisationID: person.OrganisationID,
		ClassID:        "",
		PersonSyncKey:  person.PersonSyncKey,
		SchoolNumber:   person.SchoolNumber,
		Db:             db,
	})
	person.Save(db)
}
