package SiegeStatus

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// += "?applicationIds=ID1,ID2,ID3"
const API_URL = `https://public-ubiservices.ubi.com/v1` +
	`/applications/gameStatuses`

var headers = map[string]string{
	// "access-control-request-headers": "ubi-appid",
	// "access-control-request-method":  "GET",
	// "origin":  "https://www.ubisoft.com",
	// "referer": "https://www.ubisoft.com/",
	// "priority": "u=1, i",
	// "sec-fetch-dest": "empty",
	// "sec-fetch-mode": "cors",
	// "sec-fetch-site": "cross-site",
	// "cache-control": "no-cache",
	// "pragma":        "no-cache",
	"ubi-appid": APP_ID_GAME_STATUSES,
}

type AppId = string

const (
	APP_ID_GAME_STATUSES AppId = "f612511e-58a2-4e9a-831f-61838b1950bb"

	// PC              // "spaceId": ""
	APP_ID_SIEGE_PC AppId = `e3d5ea9e-50bd-43b7-88bf-39794f4e3d40`
	// PS4             // "spaceId": "05bfb3f7-6c21-4c42-be1f-97a33fb5cf66"
	APP_ID_SIEGE_ORBIS AppId = `fb4cc4c9-2063-461d-a1e8-84a7d36525fc`
	// PS5             // "spaceId": "96c1d424-057e-4ff7-860b-6b9c9222bdbf"
	APP_ID_SIEGE_PS5 AppId = `6e3c99c9-6c3f-43f4-b4f6-f1a3143f2764`
	// XBOX SERIES X|S // "spaceId": "631d8095-c443-4e21-b301-4af1a0929c27"
	APP_ID_SIEGE_SCARLETT AppId = `76f580d5-7f50-47cc-bbc1-152d000bfe59`
	// XBOX ONE        // "spaceId": "98a601e5-ca91-4440-b1c5-753f601a2c90"
	APP_ID_SIEGE_DURANGO AppId = `4008612d-3baf-49e4-957a-33066726a7bc`

	// [TODO]
)

type Feature = string

const (
	// "Connectivity" | "连接"
	FEATURE_LEADERBOARD Feature = "Leaderboard"
	// "Authentication" | "认证"
	FEATURE_AUTHENTICATION Feature = "Authentication"
	// "In-Game Store"
	FEATURE_PURCHASE Feature = "Purchase"
	// "Matchmaking" | "匹配"
	FEATURE_MATCHMAKING Feature = "Matchmaking"
)

type GameStatus struct {
	ApplicationId AppId  `json:"applicationId"`
	SpaceId       string `json:"spaceId"`
	Name          string `json:"name"`
	PlatformType  string `json:"platformType"`

	// "interrupted" && false -> "Unplanned Outage" | "预期外停机"
	// [TODO] "" && true -> "" | ""
	Status        string `json:"status"`
	IsMaintenance bool   `json:"isMaintenance"`

	// unsorted but with fixed order // "Degraded" | "降级"
	ImpactedFeatures []Feature `json:"impactedFeatures"`
}

type Response struct {
	// "2006-01-02T15:04:05Z"
	LastModifiedAt time.Time `json:"lastModifiedAt"`

	GameStatuses []GameStatus `json:"gameStatuses"`

	RequestError
}

type RequestError struct {
	ErrorCode       int       `json:"errorCode"`
	Message         string    `json:"message"`
	HttpCode        int       `json:"httpCode"`
	ErrorContext    string    `json:"errorContext"`
	MoreInfo        string    `json:"moreInfo"`
	TransactionTime time.Time `json:"transactionTime"`
	TransactionId   string    `json:"transactionId"`
}

func (e *RequestError) Error() string {
	return fmt.Sprintf("%+v", *e)
}

func Get(ctx context.Context, ids ...AppId) (r Response, err error) {
	url := API_URL + "?applicationIds=" + strings.Join(ids, ",")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Response{}, err
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return Response{}, err
	}

	if r.Message != "" {
		e := r.RequestError
		return Response{}, &e
	}

	return r, nil
}
