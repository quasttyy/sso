package jwt

import (
	"fmt"
	"sso/internal/domain/models"
	"testing"
	"time"
	"github.com/golang-jwt/jwt/v5"
)

func TestNewToken(t *testing.T) {
	// Создаём тестовые данные
	user1 := models.User{
		ID: 1,
		Email: "1@mail.ru",
		PassHash: []byte("passHash1"),
	}
	app1 := models.App{
		App_id: 1,
		Name: "phone",
		Secret: "",
	}
	dur1 := time.Second * 30

	// Прописываем ожидаемый результат
	expectedRes := jwt.MapClaims{
		"userID": float64(user1.ID),
		"email": user1.Email,
		"appID": float64(app1.App_id),
	}

	// Вызываем метод с тестовыми данными
	token1, err := NewToken(user1, app1, dur1)
	if err != nil {
		t.Errorf("failed to generate token in test")
	}

	// Парсим полученный токен
	parsed1, err := jwt.Parse(token1, func(token *jwt.Token) (interface{}, error) {
    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
        return nil, fmt.Errorf("unexpected signing method: %T", token.Method)
    }
    	return []byte(app1.Secret), nil
	})

	// Получаем claims с спаршенного токена
	claims1, ok := parsed1.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatalf("unexpected claims type %T", parsed1.Claims)
	}

	// Сверяем claims
	for k, want := range expectedRes {
		got, ok := claims1[k]
		if !ok || got != want {
			t.Fatalf("claim %s: got %v want %v", k, got, want)
		}
	}
}
