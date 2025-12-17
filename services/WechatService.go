package services

import (
	"RHPRo-Task/config"
	"RHPRo-Task/database"
	"RHPRo-Task/dto"
	"RHPRo-Task/models"
	"RHPRo-Task/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type WechatService struct{}

// WechatAccessTokenResponse 微信access_token响应
type WechatAccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"openid"`
	Scope        string `json:"scope"`
	UnionID      string `json:"unionid"`
	ErrCode      int    `json:"errcode"`
	ErrMsg       string `json:"errmsg"`
}

// WechatUserInfoResponse 微信用户信息响应
type WechatUserInfoResponse struct {
	OpenID     string `json:"openid"`
	Nickname   string `json:"nickname"`
	HeadImgURL string `json:"headimgurl"`
	UnionID    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

// WechatPhoneResponse 微信手机号响应（小程序）
type WechatPhoneResponse struct {
	ErrCode   int    `json:"errcode"`
	ErrMsg    string `json:"errmsg"`
	PhoneInfo struct {
		PhoneNumber     string `json:"phoneNumber"`
		PurePhoneNumber string `json:"purePhoneNumber"`
		CountryCode     string `json:"countryCode"`
	} `json:"phone_info"`
}

// Login 微信登录
func (s *WechatService) Login(req *dto.WechatLoginRequest) (*dto.WechatLoginResponse, error) {
	// 1. 根据登录类型获取对应的AppID和Secret
	appID, appSecret := s.getWechatConfig(req.LoginType)
	if appID == "" || appSecret == "" {
		return nil, errors.New("微信配置未设置")
	}

	// 2. 用code换取access_token和openid
	tokenResp, err := s.getAccessToken(appID, appSecret, req.Code)
	if err != nil {
		return nil, err
	}

	// 3. 获取微信用户信息
	userInfo, err := s.getUserInfo(tokenResp.AccessToken, tokenResp.OpenID)
	if err != nil {
		return nil, err
	}

	// 4. 查找是否已通过微信绑定用户（优先用unionid，其次openid）
	var user models.User
	var foundByWechat bool
	if tokenResp.UnionID != "" {
		if err := database.DB.Where("wechat_unionid = ?", tokenResp.UnionID).First(&user).Error; err == nil {
			foundByWechat = true
		}
	}
	if !foundByWechat && tokenResp.OpenID != "" {
		if err := database.DB.Where("wechat_openid = ?", tokenResp.OpenID).First(&user).Error; err == nil {
			foundByWechat = true
		}
	}

	// 5. 已绑定微信的用户，直接登录
	if foundByWechat {
		return s.loginExistingUser(&user)
	}

	// 6. 未绑定微信，尝试获取微信绑定的手机号
	wechatMobile := ""
	if req.PhoneCode != "" {
		// 小程序场景：用 phone_code 获取手机号
		phoneResp, err := s.getPhoneNumber(appID, appSecret, req.PhoneCode)
		if err == nil && phoneResp.PhoneInfo.PurePhoneNumber != "" {
			wechatMobile = phoneResp.PhoneInfo.PurePhoneNumber
		}
	}

	// 7. 如果获取到微信手机号，检查该手机号是否已注册
	if wechatMobile != "" {
		var existingUser models.User
		if err := database.DB.Where("mobile = ?", wechatMobile).First(&existingUser).Error; err == nil {
			// 手机号已注册，绑定微信信息并直接登录
			existingUser.WechatUnionID = tokenResp.UnionID
			existingUser.WechatOpenID = tokenResp.OpenID
			if existingUser.Nickname == "" {
				existingUser.Nickname = userInfo.Nickname
			}
			if existingUser.Avatar == "" {
				existingUser.Avatar = userInfo.HeadImgURL
			}
			if err := database.DB.Save(&existingUser).Error; err != nil {
				return nil, err
			}
			return s.loginExistingUser(&existingUser)
		}
	}

	// 8. 未绑定用户，返回微信信息和手机号，要求补充密码
	tempToken, err := s.generateTempToken(tokenResp.UnionID, tokenResp.OpenID, userInfo.Nickname, userInfo.HeadImgURL, wechatMobile)
	if err != nil {
		return nil, err
	}

	return &dto.WechatLoginResponse{
		NeedBind:  true,
		TempToken: tempToken,
		WechatInfo: &dto.WechatUserInfo{
			UnionID:  tokenResp.UnionID,
			OpenID:   tokenResp.OpenID,
			Nickname: userInfo.Nickname,
			Avatar:   userInfo.HeadImgURL,
			Mobile:   wechatMobile,
		},
	}, nil
}

// BindMobile 微信绑定手机号（新用户注册）
func (s *WechatService) BindMobile(req *dto.WechatBindRequest) (*dto.LoginResponse, error) {
	// 1. 解析临时token
	claims, err := s.parseTempToken(req.TempToken)
	if err != nil {
		return nil, errors.New("临时凭证无效或已过期")
	}

	// 2. 如果微信已绑定手机号，必须使用该手机号注册
	if claims.Mobile != "" && req.Mobile != claims.Mobile {
		return nil, errors.New("请使用微信绑定的手机号进行注册")
	}

	// 3. 检查手机号是否已注册
	var existingUser models.User
	if err := database.DB.Where("mobile = ?", req.Mobile).First(&existingUser).Error; err == nil {
		// 手机号已存在，绑定微信到现有账号
		existingUser.WechatUnionID = claims.UnionID
		existingUser.WechatOpenID = claims.OpenID
		if existingUser.Nickname == "" {
			existingUser.Nickname = claims.Nickname
		}
		if existingUser.Avatar == "" {
			existingUser.Avatar = claims.Avatar
		}
		if err := database.DB.Save(&existingUser).Error; err != nil {
			return nil, err
		}
		return s.generateLoginResponse(&existingUser)
	}

	// 4. 创建新用户
	user := &models.User{
		Mobile:        req.Mobile,
		Username:      req.UserName,
		Nickname:      claims.Nickname,
		Avatar:        claims.Avatar,
		WechatUnionID: claims.UnionID,
		WechatOpenID:  claims.OpenID,
		Status:        models.UserStatusPending,
	}
	if err := user.SetPassword(req.Password); err != nil {
		return nil, err
	}
	if err := database.DB.Create(user).Error; err != nil {
		return nil, err
	}

	// 5. 分配默认角色
	var userRole models.Role
	if err := database.DB.Where("name = ?", "user").First(&userRole).Error; err == nil {
		database.DB.Model(user).Association("Roles").Append(&userRole)
	}

	return s.generateLoginResponse(user)
}

// getWechatConfig 获取微信配置
func (s *WechatService) getWechatConfig(loginType string) (appID, appSecret string) {
	cfg := config.GetConfig()
	switch loginType {
	case "scan": // 扫码登录（开放平台）
		return cfg.Wechat.OpenAppID, cfg.Wechat.OpenAppSecret
	case "mp": // 小程序
		return cfg.Wechat.MpAppID, cfg.Wechat.MpAppSecret
	case "h5": // 公众号H5
		return cfg.Wechat.H5AppID, cfg.Wechat.H5AppSecret
	}
	return "", ""
}

// getAccessToken 用code换取access_token
func (s *WechatService) getAccessToken(appID, appSecret, code string) (*WechatAccessTokenResponse, error) {
	url := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
		appID, appSecret, code,
	)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result WechatAccessTokenResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	if result.ErrCode != 0 {
		return nil, fmt.Errorf("微信授权失败: %s", result.ErrMsg)
	}
	return &result, nil
}

// getUserInfo 获取微信用户信息
func (s *WechatService) getUserInfo(accessToken, openID string) (*WechatUserInfoResponse, error) {
	url := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN",
		accessToken, openID,
	)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result WechatUserInfoResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	if result.ErrCode != 0 {
		return nil, fmt.Errorf("获取微信用户信息失败: %s", result.ErrMsg)
	}
	return &result, nil
}

// getPhoneNumber 获取微信绑定的手机号（小程序专用）
func (s *WechatService) getPhoneNumber(appID, appSecret, phoneCode string) (*WechatPhoneResponse, error) {
	// 先获取 access_token（小程序的 client_credential 模式）
	tokenURL := fmt.Sprintf(
		"https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s",
		appID, appSecret,
	)
	tokenResp, err := http.Get(tokenURL)
	if err != nil {
		return nil, err
	}
	defer tokenResp.Body.Close()

	tokenBody, _ := io.ReadAll(tokenResp.Body)
	var tokenResult struct {
		AccessToken string `json:"access_token"`
		ErrCode     int    `json:"errcode"`
		ErrMsg      string `json:"errmsg"`
	}
	if err := json.Unmarshal(tokenBody, &tokenResult); err != nil {
		return nil, err
	}
	if tokenResult.ErrCode != 0 {
		return nil, fmt.Errorf("获取access_token失败: %s", tokenResult.ErrMsg)
	}

	// 用 phone_code 获取手机号
	phoneURL := fmt.Sprintf(
		"https://api.weixin.qq.com/wxa/business/getuserphonenumber?access_token=%s",
		tokenResult.AccessToken,
	)
	reqBody := fmt.Sprintf(`{"code":"%s"}`, phoneCode)
	phoneResp, err := http.Post(phoneURL, "application/json", strings.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	defer phoneResp.Body.Close()

	phoneBody, _ := io.ReadAll(phoneResp.Body)
	var result WechatPhoneResponse
	if err := json.Unmarshal(phoneBody, &result); err != nil {
		return nil, err
	}
	if result.ErrCode != 0 {
		return nil, fmt.Errorf("获取手机号失败: %s", result.ErrMsg)
	}
	return &result, nil
}

// loginExistingUser 已绑定用户登录
func (s *WechatService) loginExistingUser(user *models.User) (*dto.WechatLoginResponse, error) {
	// 检查用户状态
	if user.Status == models.UserStatusDisabled {
		return nil, errors.New("用户已被禁用")
	}
	if user.Status == models.UserStatusPending {
		return nil, errors.New("用户待审核，请联系管理员")
	}

	loginResp, err := s.generateLoginResponse(user)
	if err != nil {
		return nil, err
	}

	return &dto.WechatLoginResponse{
		NeedBind: false,
		Token:    loginResp.Token,
		UserInfo: loginResp.UserInfo,
	}, nil
}

// generateLoginResponse 生成登录响应
func (s *WechatService) generateLoginResponse(user *models.User) (*dto.LoginResponse, error) {
	// 重新加载用户关联数据
	database.DB.Preload("Roles.Permissions").Preload("Department").Preload("ManagedDepartments").First(user, user.ID)

	// var deptID uint
	// var deptName string
	// if user.Department != nil {
	// 	deptID = user.Department.ID
	// 	deptName = user.Department.Name
	// }

	// var managedDeptIDs []uint
	// isLeader := len(user.ManagedDepartments) > 0
	// for _, dept := range user.ManagedDepartments {
	// 	managedDeptIDs = append(managedDeptIDs, dept.ID)
	// }

	token, err := utils.GenerateToken(user.ID, user.Username, user.Mobile, 24)
	if err != nil {
		return nil, err
	}

	userInfo := dto.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		Email:    user.Email,
		Mobile:   user.Mobile,
		Status:   user.Status,
	}

	return &dto.LoginResponse{
		Token:    token,
		UserInfo: userInfo,
	}, nil
}

// TempTokenClaims 临时token的claims
type TempTokenClaims struct {
	UnionID  string `json:"unionid"`
	OpenID   string `json:"openid"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Mobile   string `json:"mobile"` // 微信绑定的手机号
	jwt.RegisteredClaims
}

// generateTempToken 生成临时token（用于绑定手机号）
func (s *WechatService) generateTempToken(unionID, openID, nickname, avatar, mobile string) (string, error) {
	claims := TempTokenClaims{
		UnionID:  unionID,
		OpenID:   openID,
		Nickname: nickname,
		Avatar:   avatar,
		Mobile:   mobile,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)), // 10分钟有效
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.GetConfig().JWT.Secret))
}

// parseTempToken 解析临时token
func (s *WechatService) parseTempToken(tokenString string) (*TempTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TempTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetConfig().JWT.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*TempTokenClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
