package SiegeStatus

import "testing"

func TestGet(t *testing.T) {
	resp, err := Get(t.Context(),
		APP_ID_SIEGE_PC,
		// APP_ID_SIEGE_ORBIS,
		// APP_ID_SIEGE_PS5,
		// APP_ID_SIEGE_SCARLETT,
		// APP_ID_SIEGE_DURANGO,
	)
	if err != nil {
		t.Fatalf("%T: %v", err, err)
	}
	t.Log(resp.LastModifiedAt.Local())
	for _, status := range resp.GameStatuses {
		t.Logf("%+v", status)
	}
}
