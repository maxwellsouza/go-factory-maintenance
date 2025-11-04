package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/maxwellsouza/go-factory-maintenance/internal/http/handlers"
	"github.com/maxwellsouza/go-factory-maintenance/internal/repository/memory"
	"github.com/maxwellsouza/go-factory-maintenance/internal/service"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(gin.Recovery())

	assetRepo := memory.NewAssetMemoryRepo()
	workOrderRepo := memory.NewWorkOrderMemoryRepo()

	assetSvc := service.NewAssetService(assetRepo)
	workOrderSvc := service.NewWorkOrderService(workOrderRepo)

	assetH := handlers.NewAssetHandler(assetSvc)
	woH := handlers.NewWorkOrderHandler(workOrderSvc)

	// healthz p/ sanity
	r.GET("/healthz", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

	assetH.RegisterRoutes(r)
	woH.RegisterRoutes(r)

	return r
}

func TestHealthz(t *testing.T) {
	r := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if body := w.Body.String(); body != "ok" {
		t.Fatalf("expected body 'ok', got %q", body)
	}
}

func TestAssets_CreateAndList(t *testing.T) {
	r := setupRouter()

	// POST /assets
	payload := []byte(`{"name":"Cortadeira","location":"Galpao A","criticality":"A"}`)
	req := httptest.NewRequest(http.MethodPost, "/assets", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("POST /assets expected 201, got %d; body=%s", w.Code, w.Body.String())
	}

	// GET /assets
	reqList := httptest.NewRequest(http.MethodGet, "/assets", nil)
	wList := httptest.NewRecorder()
	r.ServeHTTP(wList, reqList)

	if wList.Code != http.StatusOK {
		t.Fatalf("GET /assets expected 200, got %d", wList.Code)
	}

	var assets []map[string]any
	if err := json.Unmarshal(wList.Body.Bytes(), &assets); err != nil {
		t.Fatalf("failed to unmarshal assets: %v; body=%s", err, wList.Body.String())
	}
	if len(assets) != 1 {
		t.Fatalf("expected 1 asset, got %d", len(assets))
	}
	if assets[0]["name"] != "Cortadeira" {
		t.Fatalf("expected asset name 'Cortadeira', got %v", assets[0]["name"])
	}
}

func TestWorkOrders_CreateAndFilter(t *testing.T) {
	r := setupRouter()

	// cria um asset primeiro (asset_id=1)
	reqAsset := httptest.NewRequest(http.MethodPost, "/assets",
		bytes.NewReader([]byte(`{"name":"Rebobinadeira","location":"Galpao B"}`)))
	reqAsset.Header.Set("Content-Type", "application/json")
	wAsset := httptest.NewRecorder()
	r.ServeHTTP(wAsset, reqAsset)
	if wAsset.Code != http.StatusCreated {
		t.Fatalf("POST /assets expected 201, got %d; body=%s", wAsset.Code, wAsset.Body.String())
	}

	// cria duas OS, uma default open/corrective e outra in_progress
	reqWO1 := httptest.NewRequest(http.MethodPost, "/work-orders",
		bytes.NewReader([]byte(`{"asset_id":1,"title":"Trocar rolete","description":"barulho"}`)))
	reqWO1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, reqWO1)
	if w1.Code != http.StatusCreated {
		t.Fatalf("POST /work-orders expected 201, got %d; body=%s", w1.Code, w1.Body.String())
	}

	reqWO2 := httptest.NewRequest(http.MethodPost, "/work-orders",
		bytes.NewReader([]byte(`{"asset_id":1,"type":"corrective","status":"in_progress","title":"Ajuste correia"}`)))
	reqWO2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, reqWO2)
	if w2.Code != http.StatusCreated {
		t.Fatalf("POST /work-orders expected 201, got %d; body=%s", w2.Code, w2.Body.String())
	}

	// GET /work-orders?status=open → deve trazer pelo menos 1
	reqOpen := httptest.NewRequest(http.MethodGet, "/work-orders?status=open", nil)
	wOpen := httptest.NewRecorder()
	r.ServeHTTP(wOpen, reqOpen)
	if wOpen.Code != http.StatusOK {
		t.Fatalf("GET /work-orders?status=open expected 200, got %d", wOpen.Code)
	}
	var open []map[string]any
	if err := json.Unmarshal(wOpen.Body.Bytes(), &open); err != nil {
		t.Fatalf("unmarshal open list: %v; body=%s", err, wOpen.Body.String())
	}
	if len(open) < 1 {
		t.Fatalf("expected at least 1 open work order, got %d", len(open))
	}

	// GET /work-orders sem filtro → deve trazer 2
	reqAll := httptest.NewRequest(http.MethodGet, "/work-orders", nil)
	wAll := httptest.NewRecorder()
	r.ServeHTTP(wAll, reqAll)
	if wAll.Code != http.StatusOK {
		t.Fatalf("GET /work-orders expected 200, got %d", wAll.Code)
	}
	var all []map[string]any
	if err := json.Unmarshal(wAll.Body.Bytes(), &all); err != nil {
		t.Fatalf("unmarshal all list: %v; body=%s", err, wAll.Body.String())
	}
	if len(all) != 2 {
		t.Fatalf("expected 2 work orders, got %d", len(all))
	}
}

func TestAssets_Create_Validation422(t *testing.T) {
	r := setupRouter()

	// JSON vazio -> falta "name" (required)
	payload := []byte(`{}`)
	req := httptest.NewRequest(http.MethodPost, "/assets", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("POST /assets expected 422, got %d; body=%s", w.Code, w.Body.String())
	}

	// Opcional: valida que vem um JSON com "error":"erro de validação"
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if body["error"] != "erro de validação" {
		t.Fatalf("expected error 'erro de validação', got %v", body["error"])
	}
}

func TestWorkOrders_Create_Validation422(t *testing.T) {
	r := setupRouter()

	// Primeiro, cria um asset válido (para depois testar outros cenários)
	reqAsset := httptest.NewRequest(http.MethodPost, "/assets",
		bytes.NewReader([]byte(`{"name":"Asset Test 422"}`)))
	reqAsset.Header.Set("Content-Type", "application/json")
	wAsset := httptest.NewRecorder()
	r.ServeHTTP(wAsset, reqAsset)
	if wAsset.Code != http.StatusCreated {
		t.Fatalf("POST /assets expected 201, got %d; body=%s", wAsset.Code, wAsset.Body.String())
	}

	tests := []struct {
		name    string
		payload string
	}{
		{
			name:    "missing asset_id and title",
			payload: `{}`,
		},
		{
			name:    "invalid asset_id (zero)",
			payload: `{"asset_id":0,"title":"x"}`,
		},
		{
			name:    "missing title",
			payload: `{"asset_id":1}`,
		},
		{
			name:    "short title (min=3)",
			payload: `{"asset_id":1,"title":"ab"}`,
		},
		{
			name:    "invalid status value",
			payload: `{"asset_id":1,"title":"Valid","status":"WRONG"}`,
		},
		{
			name:    "invalid type value",
			payload: `{"asset_id":1,"title":"Valid","type":"WRONG"}`,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/work-orders", bytes.NewReader([]byte(tc.payload)))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != http.StatusUnprocessableEntity {
				t.Fatalf("POST /work-orders (%s) expected 422, got %d; body=%s",
					tc.name, w.Code, w.Body.String())
			}

			var body map[string]any
			if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
				t.Fatalf("failed to unmarshal response: %v", err)
			}
			if body["error"] != "erro de validação" {
				t.Fatalf("expected error 'erro de validação', got %v", body["error"])
			}
		})
	}
}
