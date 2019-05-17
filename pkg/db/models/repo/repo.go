package repo

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/ryantking/marina/pkg/db"
)

// Org represents an entry in the database
type Org struct {
	gorm.Model
	Name    string `gorm:"column:name"`
	OrgID   uint   ``
	OrgName string `gorm:"column:org_name"`
}

// TableName returns the organization table name
func (Org) TableName() string {
	return "repositories"
}

// OrgExists checks whether or not a given organization exists
func OrgExists(name, orgName string) (bool, error) {
	ms := []Org{}
	res := db.Get().Model(&Org{}).Association("organizations").Find(&ms)
	if res.Error != nil {
		fmt.Println(res.Error.Error())
	}
	fmt.Println(ms)

	return false, nil
}

// // GetNames returns the names of all available repos
// func GetNames() ([]string, error) {
// 	col, err := getCollection(CollectionName)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	repos := make([]Model, 0)
// 	err = col.Find().OrderBy("org_name").OrderBy("repo_name").All(&repos)
// 	if err != nil {
// 		return []string{}, errors.Wrap(err, "error retrieving repos")
// 	}
//
// 	return repoNames(repos), nil
// }
//
// // GetNamesPaginated
// func GetNamesPaginated(pageSize uint, last string) ([]string, string, error) {
// 	col, err := getCollection(CollectionName)
// 	if err != nil {
// 		return nil, "", err
// 	}
//
// 	repos := make([]Model, 0, pageSize)
// 	res := col.Find().OrderBy("org_name").OrderBy("repo_name").Paginate(pageSize).Cursor("name")
// 	err = res.NextPage(last).All(repos)
// 	if err != nil {
// 		return []string{}, "", errors.Wrap(err, "error retrieving repos")
// 	}
// 	nextLast := repos[len(repos)-1].Name
// 	exists, err := res.NextPage(nextLast).Exists()
// 	if err != nil {
// 		return []string{}, "", errors.Wrap(err, "error checking if next page exists")
// 	}
// 	if !exists {
// 		nextLast = ""
// 	}
//
// 	return repoNames(repos), nextLast, nil
//
// }
//
// func repoNames(repos []Model) []string {
// 	names := make([]string, len(repos))
// 	for i, repo := range repos {
// 		names[i] = fmt.Sprintf("%s/%s", repo.OrgName, repo.Name)
// 	}
//
// 	return names
// }
