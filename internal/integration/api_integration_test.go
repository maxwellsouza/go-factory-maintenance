//go:build integration

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/maxwellsouza/go-factory-maintenance/internal/http/handlers"
	"github.com/maxwellsouza/go-factory-maintenance/internal/repository/postgres"
	"github.com/maxwellsouza/go-factory-maintenance/internal/service"
)

func setupAPI(t *testing.T) *gin.Engine {
	t.Helper()

	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "dev")
	os.Setenv("DB_PASS", "dev")
	os.Setenv("DB_NAME", "maintenance")

	ctx := context.Background()

	db, err := postgres.New(ctx)
	if err != nil {
		t.Fatalf("failed to connect to DB: %v", err)
	}

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(gin.Recovery())

	assetRepo := postgres.NewAssetRepo(db)
	workOrderRepo := postgres.NewWorkOrderRepo(db)

	assetSvc := service.NewAssetService(assetRepo)
	workOrderSvc := service.NewWorkOrderService(workOrderRepo)

	assetHandler := handlers.NewAssetHandler(assetSvc)
	workOrderHandler := handlers.NewWorkOrderHandler(workOrderSvc)

	assetHandler.RegisterRoutes(r)
	workOrderHandler.RegisterRoutes(r)

	return r
}

func TestIntegration_AssetsAndWorkOrders(t *testing.T) {
	r := setupAPI(t)

	// Criar um ativo
	assetPayload := []byte(`{"name":"IntegrTest Cortadeira","location":"Galpão A","criticality":"A"}`)
	req := httptest.NewRequest(http.MethodPost, "/assets", bytes.NewReader(assetPayload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d; body=%s", w.Code, w.Body.String())
	}

	// Listar ativos
	reqList := httptest.NewRequest(http.MethodGet, "/assets", nil)
	wList := httptest.NewRecorder()
	r.ServeHTTP(wList, reqList)

	if wList.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", wList.Code)
	}

	var assets []map[string]any
	if err := json.Unmarshal(wList.Body.Bytes(), &assets); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if len(assets) == 0 {
		t.Fatalf("expected at least 1 asset, got 0")
	}

	id := int(assets[len(assets)-1]["id"].(float64))

	// Criar OS vinculada ao ativo
	payloadWO := []byte(fmt.Sprintf(`{"asset_id":%d,"type":"corrective","title":"Lubrificar rolamento","description":"rolamento ruidoso"}`, id))
	reqWO := httptest.NewRequest(http.MethodPost, "/work-orders", bytes.NewReader(payloadWO))
	reqWO.Header.Set("Content-Type", "application/json")
	wWO := httptest.NewRecorder()
	r.ServeHTTP(wWO, reqWO)
	if wWO.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d; body=%s", wWO.Code, wWO.Body.String())
	}

	// Listar OS abertas
	reqListWO := httptest.NewRequest(http.MethodGet, "/work-orders?status=open", nil)
	wListWO := httptest.NewRecorder()
	r.ServeHTTP(wListWO, reqListWO)

	if wListWO.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", wListWO.Code)
	}

	var workOrders []map[string]any
	if err := json.Unmarshal(wListWO.Body.Bytes(), &workOrders); err != nil {
		t.Fatalf("failed to unmarshal work orders: %v", err)
	}
	if len(workOrders) == 0 {
		t.Fatalf("expected at least 1 work order, got 0")
	}
}
