package api

import (
	"a21hc3NpZ25tZW50/entity"
	"a21hc3NpZ25tZW50/service"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type CategoryAPI interface {
	GetCategory(w http.ResponseWriter, r *http.Request)
	CreateNewCategory(w http.ResponseWriter, r *http.Request)
	DeleteCategory(w http.ResponseWriter, r *http.Request)
	GetCategoryWithTasks(w http.ResponseWriter, r *http.Request)
}

type categoryAPI struct {
	categoryService service.CategoryService
}

func NewCategoryAPI(categoryService service.CategoryService) *categoryAPI {
	return &categoryAPI{categoryService}
}

func (c *categoryAPI) GetCategory(w http.ResponseWriter, r *http.Request) {
	userId := fmt.Sprintf("%s", r.Context().Value("id"))
	if len(userId) == 0 {
		w.WriteHeader(400) // WriteHeader ditulis di awal!
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error": "invalid user id"}`))
		return
	}
	convUserId, err := strconv.Atoi(userId)
	if err != nil {
		panic(err)
	}
	// fmt.Println("ini get category(bedasarkan id) :", convUserId)
	categori, err := c.categoryService.GetCategories(r.Context(), convUserId)
	if err != nil {
		w.WriteHeader(500) // WriteHeader ditulis di awal!
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error": "error internal server"}`))
		return
	}
	// fmt.Println("ini GetCategories : ", categori)
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(categori)
	// TODO: answer here
}

func (c *categoryAPI) CreateNewCategory(w http.ResponseWriter, r *http.Request) {
	var category entity.CategoryRequest

	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid category request"))
		return
	}

	addCategory := entity.Category{}
	addCategory.Type = category.Type
	if len(addCategory.Type) == 0 {
		w.WriteHeader(400) // WriteHeader ditulis di awal!
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error": "invalid category request"}`))
		return
	}

	userId := fmt.Sprintf("%s", r.Context().Value("id"))
	if len(userId) == 0 {
		w.WriteHeader(400) // WriteHeader ditulis di awal!
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error": "invalid user id"}`))
		return
	}
	convUserId, err := strconv.Atoi(userId)
	if err != nil {
		panic(err)
	}
	addCategory.UserID = convUserId
	cat, err := c.categoryService.StoreCategory(r.Context(), &addCategory)
	if err != nil {
		w.WriteHeader(500) // WriteHeader ditulis di awal!
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error": "error internal server"}`))
		return
	}
	// fmt.Println("ini API CreateNewCategory : ", cat)
	res := make(map[string]interface{})
	res["user_id"] = cat.UserID
	res["category_id"] = cat.ID
	res["message"] = "success create new category"
	jsonInBytes, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(201) // WriteHeader ditulis di awal!
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonInBytes)
	// TODO: answer here
}

func (c *categoryAPI) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	userId := fmt.Sprintf("%s", r.Context().Value("id"))
	convUserId, err := strconv.Atoi(userId)
	if err != nil {
		panic(err)
	}
	// fmt.Println("ini DeleteCategory (userID) :", userId)
	categoryID := r.URL.Query().Get("category_id")
	// fmt.Println("ini DeleteCategory (categoryID) :", categoryID)
	convCategoryID, err := strconv.Atoi(categoryID)
	if err != nil {
		panic(err)
	}
	err = c.categoryService.DeleteCategory(r.Context(), convCategoryID)
	if err != nil {
		w.WriteHeader(500) // WriteHeader ditulis di awal!
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error": "error internal server"}`))
	}
	res := make(map[string]interface{})
	res["user_id"] = convUserId
	res["category_id"] = convCategoryID
	res["message"] = "success delete category"
	jsonInBytes, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(200) // WriteHeader ditulis di awal!
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonInBytes)
	// TODO: answer here
}

func (c *categoryAPI) GetCategoryWithTasks(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("id")

	idLogin, err := strconv.Atoi(userId.(string))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("get category task", err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	categories, err := c.categoryService.GetCategoriesWithTasks(r.Context(), int(idLogin))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("internal server error"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(categories)

}
