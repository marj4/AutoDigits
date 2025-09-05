package usecase

import (
	dto "auth-service/internal/application/DTO"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

type ResponseFromUS struct {
	Uuid     string `json:"UUID"`
	Password string `json:"Password"`
	Role     string `json:Role`
}

type ClientUseCase struct {
	client *http.Client
	url    string
}

func NewClientUseCase(url string, client *http.Client) *ClientUseCase {
	return &ClientUseCase{
		client: client,
		url:    url,
	}
}

func (u *ClientUseCase) SignUp(req *dto.Request) *dto.Response {
	username := req.Username
	log.Println("Username from middleware", username)

	password := req.Password
	role := req.Role

	// запрос в User-Service на проверку дубликатов
	resp, err := http.Get(u.url + "/user/" + username)
	if err != nil {
		log.Fatal("(SignUp) Error send GET-request: ", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusFound {
		return &dto.Response{
			Status:  strconv.Itoa(resp.StatusCode),
			Message: "This user is exist",
		}
	}

	log.Println("Error is here!!!")

	// хеширование пароля
	hashPass := HashPassword(password)

	request := &dto.Request{
		Username: username,
		Password: hashPass,
		Role:     role,
	}

	// запрос в User-Service на добавление пользователя
	data, err := json.Marshal(request)
	if err != nil {
		log.Fatal("(Sign) Error Marshaling: ", err)
	}

	dataReader := bytes.NewReader(data)

	if _, err = http.Post(u.url+"/user", "application/json", dataReader); err != nil {
		log.Fatal("(SignUp) Error send POST-request: ", err)
	}

	return &dto.Response{
		Status:  dto.Success,
		Message: "User is added!",
	}
}

func (u *ClientUseCase) SignIn(req *dto.Request) *dto.Response {
	username := req.Username

	// запрос в User-Service на проверку дубликатов и в случае если такой пользователь есть, то принимаем его UUDI and password, rolefrom US
	resp, err := http.Get(u.url + "/user/" + username)
	if err != nil {
		log.Fatal("(SignUp) Error send GET-request: ", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return &dto.Response{
			Status:  strconv.Itoa(resp.StatusCode),
			Message: "This user is not exist",
		}
	}

	// Распарсить данные в структуру и передать наверх
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// 4. Parse JSON in struct
	var data *ResponseFromUS
	if err := json.Unmarshal(body, &data); err != nil {
		log.Fatal("(SignIn) Error unmarshal: ", err)
	}

	// Передаем наверх
	return &dto.Response{
		Uuid:     data.Uuid,
		Password: data.Password,
		Role:     data.Role,
	}

}

func HashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}
