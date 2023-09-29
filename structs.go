package itswizard_m_bwStructs

import "github.com/jinzhu/gorm"

type BWLog struct {
	gorm.Model
	Log            string
	OrganisationID uint
	ClassID        string
	PersonSyncKey  string
	SchoolNumber   string
}

type DeletedBWSchool struct {
	gorm.Model
	OrganisationID uint `gorm:"unique"`
	SchoolNumber   string
	SchoolName     string
	SchoolType     string
	Cluster        int
	ClusterName    string
	EMail          string
	GdeSchl        string
}

type BWSchool struct {
	gorm.Model
	OrganisationID uint   `gorm:"unique"`
	ADKey          string `gorm:"unique"`
	BWKey          string `gorm:"unique"`
	SchoolNumber   string
	SchoolName     string
	SchoolType     string
	Cluster        int
	ClusterName    string
	EMail          string
	GdeSchl        string
	ToImport       bool
	ToUpdate       bool
	ToDelete       bool
	Success        bool
	Error          bool
	ErrorString    string `gorm:"type:MEDIUMTEXT"`
	ToSync         bool
}

type BWPerson struct {
	gorm.Model
	PersonSyncKey   string `gorm:"unique; type:VARCHAR(100)"`
	ADKey           string `gorm:"unique"`
	BWKey           string `gorm:"unique"`
	ToImport        bool
	ToUpdate        bool
	ToDelete        bool
	Success         bool
	Error           bool
	ErrorString     string `gorm:"type:MEDIUMTEXT"`
	Birthday        string
	UpdateBirthday  bool
	FirstName       string `gorm:"not null; type:VARCHAR(45)" json:"first_name"`
	UpdateFirstname bool
	LastName        string `gorm:"not null; type:VARCHAR(45)" json:"last_name"`
	UpdateLastname  bool
	Username        string `gorm:"unique; not null; type:VARCHAR(45)" json:"username"`
	UpdateUsername  bool
	Profile         string `gorm:"not null; type:VARCHAR(45)" json:"profile"`
	UpdateProfile   bool
	Email           string `gorm:"type:VARCHAR(45); primary key" json:"email"`
	UpdateEmail     bool
	Phone           string `gorm:"type:VARCHAR(45)" json:"phone"`
	UpdatePhone     bool
	Mobile          string `gorm:"type:VARCHAR(45)" json:"mobile"`
	UpdateMobile    bool
	Street1         string `gorm:"type:VARCHAR(45)" json:"street_1"`
	UpdateStreet1   bool
	Street2         string `gorm:"type:VARCHAR(45)" json:"street_2"`
	UpdateStreet2   bool
	Postcode        string `gorm:"type:VARCHAR(45)" json:"postcode"`
	UpdatePostcode  bool
	City            string `gorm:"type:VARCHAR(45)" json:"city"`
	UpdateCity      bool
	OriginalData    string `gorm:"type:LONGTEXT"`
	SchoolID        string
	OrganisationID  uint
	SchoolNumber    string
}

// Update:
func (p *BWPerson) ToSlice(password string) (out []string) {
	out = append(out, p.FirstName)
	out = append(out, p.LastName)
	out = append(out, p.Username)
	out = append(out, p.Email)
	out = append(out, password)
	return
}

/*
func (p *BWPerson) ToSlice(password string) (out []string) {
	out = append(out,p.PersonSyncKey)
	out = append(out,p.FirstName)
	out = append(out,p.LastName)
	out = append(out,p.Username)
	out = append(out,p.Email)
	out = append(out,password)
	return
}

*/

type BWClass struct {
	gorm.Model
	GroupSyncKey   string `gorm:"unique"`
	ADKey          string `gorm:"unique"`
	BWKey          string `gorm:"unique"`
	ToImport       bool
	ToUpdate       bool
	ToDelete       bool
	Success        bool
	Name           string
	Error          bool
	ErrorString    string `gorm:"type:MEDIUMTEXT"`
	OrganisationID uint
	SchoolNumber   string
}

type BWPersonSchoolMembership struct {
	gorm.Model
	Role                string
	IDMRoles            string `gorm:"type:MEDIUMTEXT"`
	PersonSyncKey       string
	OrganisationsyncKey string
	ToImport            bool
	ToUpdate            bool
	ToDelete            bool
	Success             bool
	Error               bool
	ErrorString         string `gorm:"type:MEDIUMTEXT"`
	OrganisationID      uint
	SchoolNumber        string
}

type BWPersonClassMembership struct {
	gorm.Model
	ClassRole      string
	PersonSyncKey  string
	GroupSyncKey   string
	ToImport       bool
	ToUpdate       bool
	ToDelete       bool
	Success        bool
	Error          bool
	ErrorString    string `gorm:"type:MEDIUMTEXT"`
	OrganisationID uint
	SchoolNumber   string
}

type BWClusterInformation struct {
	gorm.Model
	ClusterNumber uint
	ClusterSite   string
	SAMLMetadata  string `gorm:"type:MEDIUMTEXT"`
}
