package repository

import "category/domain/model"
import "github.com/jinzhu/gorm"

type ICategoryRepository interface{
	//初始化数据表
	InitTable() error
	//根据模块目录名称查找模块目录信息
	FindCategoryByID(int64)(*model.Category,error)
	//创建目录
	CreateCategory(*model.Category)(int64,error)
	//根据目录ID删除目录
	DeleteCategoryByID(int64)error
	//更新模块目录信息
	UpdateCategory(*model.Category)error
	//查找所有目录信息
	FindAll()([]model.Category,error)
	FindCategoryByName(string)(*model.Category,error)
	FindCategoryByLevel(uint32)([]model.Category,error)
	FindCategoryByParent(int64)([]model.Category,error)
}

//创建categoryRepository
func NewCategoryRepository(db *gorm.DB) ICategoryRepository{
	return &CategoryRepository{mysqlDB:db}
}
type CategoryRepository struct{
	mysqlDB *gorm.DB
}
//初始化表
func(u *CategoryRepository)InitTable() error{
	return u.mysqlDB.CreateTable(&model.Category{}).Error
}
//根据ID查询Category信息
func(u *CategoryRepository)FindCategoryByID(categoryID int64)(*model.Category,error){
	category := &model.Category{}
	return category,u.mysqlDB.First(category,categoryID).Error
}
//创建Category信息
func(u *CategoryRepository)CreateCategory(category *model.Category)(int64,error){
	return category.ID,u.mysqlDB.Create(category).Error
}
//根据ID删除Category信息
func(u *CategoryRepository)DeleteCategoryByID(categoryID int64)error{
	return u.mysqlDB.Where("id=?",categoryID).Delete(&model.Category{}).Error
}

//更新Category信息
func(u *CategoryRepository)UpdateCategory(category *model.Category)error{
	return u.mysqlDB.Model(category).Update(category).Error
}

//获取结果集
func(u *CategoryRepository)FindAll()(categoryAll []model.Category,err error){
	return categoryAll,u.mysqlDB.Find(&categoryAll).Error
}
//根据分类名称进行查找
func(u *CategoryRepository)FindCategoryByName(categoryName string)(category *model.Category,err error){
	category = &model.Category{}
	return category,u.mysqlDB.Where("category_name=?",categoryName).Find(category).Error
}

func (u *CategoryRepository)FindCategoryByLevel(level uint32)(categorySlice []model.Category,err error){
	return categorySlice,u.mysqlDB.Where("category_level=?",level).Find(categorySlice).Error
}
func (u *CategoryRepository) FindCategoryByParent(parent int64)(categorySlice []model.Category,err error){
	return categorySlice,u.mysqlDB.Where("category_parent=?",parent).Find(categorySlice).Error
}