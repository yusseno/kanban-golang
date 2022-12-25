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

type UserAPI interface {
	Login(w http.ResponseWriter, r *http.Request)
	Register(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)

	Delete(w http.ResponseWriter, r *http.Request)
}

type userAPI struct {
	userService service.UserService
}

func NewUserAPI(userService service.UserService) *userAPI {
	return &userAPI{userService}
}

func (u *userAPI) Login(w http.ResponseWriter, r *http.Request) {
	var user entity.UserLogin

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}
	userLogin := entity.User{}
	userLogin.Email = user.Email
	userLogin.Password = user.Password
	if len(userLogin.Email) == 0 || len(user.Password) == 0 {
		w.WriteHeader(400) // WriteHeader ditulis di awal!
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error": "email or password is empty"}`))
		return
	}
	id, err := u.userService.Login(r.Context(), &userLogin)
	if err != nil {
		w.WriteHeader(500)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error": "error internal server"}`))
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:   "user_id",
		Value:  strconv.Itoa(id),
		Path:   "/",
		MaxAge: 31536000,
		Domain: "",
	})
	res := make(map[string]interface{})
	res["message"] = "login success"
	res["user_id"] = id
	jsonInBytes, _ := json.Marshal(res)
	w.WriteHeader(200) // WriteHeader ditulis di awal!
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonInBytes)
	// TODO: answer here
}

func (u *userAPI) Register(w http.ResponseWriter, r *http.Request) {
	var user entity.UserRegister

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}
	userReg := entity.User{}
	userReg.Fullname = user.Fullname
	userReg.Email = user.Email
	userReg.Password = user.Password
	if len(userReg.Email) == 0 || len(userReg.Password) == 0 || userReg.Fullname == "" {
		w.WriteHeader(400) // WriteHeader ditulis di awal!
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error": "register data is empty"}`))
		return
	}
	fmt.Println("userReg: ", userReg)
	userAvailable, err := u.userService.Register(r.Context(), &userReg)
	// fmt.Println("ini userRegister : ", userAvailable)
	if err != nil {
		w.WriteHeader(500)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error": "error internal server"}`))
		return
	}

	res := make(map[string]interface{})
	res["message"] = "register success"
	res["user_id"] = userAvailable.ID
	jsonInBytes, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(201) // WriteHeader ditulis di awal!
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonInBytes)
	// TODO: answer here
}

func (u *userAPI) Logout(w http.ResponseWriter, r *http.Request) {
	// username := fmt.Sprintf("%s", r.Context().Value("username"))
	// fmt.Println("ini logout (username) : ", username)
	c, err := r.Cookie("user_id")
	if err != nil {
		panic(err)
	}
	// fmt.Println("ini logout (cookies) : ", c.Value)
	id, err := strconv.Atoi(c.Value)
	if err != nil {
		panic(err)
	}
	// fmt.Println("ini logout (cookies) : ", id)

	err = u.userService.Delete(r.Context(), id)
	if err != nil {
		panic(err)
	}
	http.SetCookie(w, &http.Cookie{
		Name:   "user_id",
		Value:  "",
		Path:   "/",
		MaxAge: 31536000,
		Domain: "",
	})
	w.WriteHeader(200) // WriteHeader ditulis di awal!
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "logout success"}`))
	// TODO: answer here
}

func (u *userAPI) Delete(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("user_id")

	if userId == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("user_id is empty"))
		return
	}

	deleteUserId, _ := strconv.Atoi(userId)

	err := u.userService.Delete(r.Context(), int(deleteUserId))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "delete success"})
}
