package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your-secret-key-change-in-production")

type Claims struct {
	UserID         uint   `json:"user_id"`
	Username       string `json:"username"`
	UserMobile     string `json:"user_mobile"`
	DepartmentID   uint   `json:"department_id"`
	DepartmentName string `json:"department_name"`
	IsLeader       bool   `json:"is_leader"`        // 是否为部门负责人
	ManagedDeptIDs []uint `json:"managed_dept_ids"` // 负责的部门ID列表
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT Token
func GenerateToken(userID uint, username string, mobile string, deptID uint, deptName string, isLeader bool, managedDeptIDs []uint, expireHours int) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(time.Duration(expireHours) * time.Hour)

	claims := Claims{
		UserID:         userID,
		Username:       username,
		UserMobile:     mobile,
		DepartmentID:   deptID,
		DepartmentName: deptName,
		IsLeader:       isLeader,
		ManagedDeptIDs: managedDeptIDs,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(nowTime),
			Issuer:    "RHPRo-Task",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseToken 解析JWT Token
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// TaskClaims 任务专用 Token Claims
type TaskClaims struct {
	TaskID uint `json:"task_id"`
	jwt.RegisteredClaims
}

// GenerateTaskToken 生成任务专用 Token
func GenerateTaskToken(taskID uint, expireHours int) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(time.Duration(expireHours) * time.Hour)

	claims := TaskClaims{
		TaskID: taskID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(nowTime),
			Issuer:    "RHPRo-Task-Context",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseTaskToken 解析任务专用 Token
func ParseTaskToken(tokenString string) (*TaskClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TaskClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*TaskClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid task token")
}

// SetJWTSecret 设置JWT密钥
func SetJWTSecret(secret string) {
	jwtSecret = []byte(secret)
}
