package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestWave3ListAndShowClients(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name, listPath, showPath string
		list                     func(*Client) (int, int, error)
		show                     func(*Client) (string, error)
	}{
		{"vehicles", "/accounts/acc/vehicles", "/vehicles/v1",
			func(c *Client) (int, int, error) {
				v, n, e := c.ListVehicles(context.Background(), "acc", 2, 25)
				return len(v), n, e
			},
			func(c *Client) (string, error) { v, e := c.GetVehicle(context.Background(), "v1"); return v.ID, e }},
		{"real-estates", "/accounts/acc/real-estates", "/real-estates/r1",
			func(c *Client) (int, int, error) {
				v, n, e := c.ListRealEstates(context.Background(), "acc", 2, 25)
				return len(v), n, e
			},
			func(c *Client) (string, error) { v, e := c.GetRealEstate(context.Background(), "r1"); return v.ID, e }},
		{"online-accounts", "/accounts/acc/online-accounts", "/online-accounts/o1",
			func(c *Client) (int, int, error) {
				v, n, e := c.ListOnlineAccounts(context.Background(), "acc", 2, 25)
				return len(v), n, e
			},
			func(c *Client) (string, error) {
				v, e := c.GetOnlineAccount(context.Background(), "o1")
				return v.ID, e
			}},
		{"insurances", "/accounts/acc/insurances", "/insurances/i1",
			func(c *Client) (int, int, error) {
				v, n, e := c.ListInsurances(context.Background(), "acc", 2, 25)
				return len(v), n, e
			},
			func(c *Client) (string, error) { v, e := c.GetInsurance(context.Background(), "i1"); return v.ID, e }},
		{"wills", "/accounts/acc/wills", "/wills/w1",
			func(c *Client) (int, int, error) {
				v, n, e := c.ListWills(context.Background(), "acc", 2, 25)
				return len(v), n, e
			},
			func(c *Client) (string, error) { v, e := c.GetWill(context.Background(), "w1"); return v.ID, e }},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			requests := 0
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests++
				switch r.URL.Path {
				case tt.listPath:
					if r.URL.Query().Get("page") != "2" || r.URL.Query().Get("size") != "25" {
						t.Errorf("query = %s", r.URL.RawQuery)
					}
					_ = json.NewEncoder(w).Encode(map[string]any{"items": []map[string]any{{"id": tt.showPath[len(tt.showPath)-2:], "name": "Item", "password": "secret"}}, "total": 7})
				case tt.showPath:
					_ = json.NewEncoder(w).Encode(map[string]any{"id": tt.showPath[len(tt.showPath)-2:], "name": "Item", "password": "secret"})
				default:
					t.Errorf("unexpected path %s", r.URL.Path)
				}
			}))
			defer srv.Close()
			c := New(srv.URL, staticToken("t"))
			count, total, err := tt.list(c)
			if err != nil || count != 1 || total != 7 {
				t.Fatalf("list count=%d total=%d err=%v", count, total, err)
			}
			id, err := tt.show(c)
			if err != nil || id == "" || requests != 2 {
				t.Fatalf("show id=%q requests=%d err=%v", id, requests, err)
			}
		})
	}
}

func TestOnlineAccountStructHasNoPassword(t *testing.T) {
	if _, ok := reflect.TypeOf(OnlineAccount{}).FieldByName("Password"); ok {
		t.Fatal("OnlineAccount must not expose a password field")
	}
}

func TestInsuranceProvidersAndExpiring(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/accounts/acc/insurance-providers":
			if r.URL.Query().Get("size") != "20" {
				t.Errorf("providers query = %s", r.URL.RawQuery)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"items": []map[string]any{{"id": "p1", "name": "Acme"}}, "total": 3})
		case "/accounts/acc/insurances/expiring":
			if r.URL.Query().Get("page") != "1" || r.URL.Query().Get("size") != "20" {
				t.Errorf("expiring query = %s", r.URL.RawQuery)
			}
			_ = json.NewEncoder(w).Encode([]map[string]any{{"id": "i1", "name": "Policy"}})
		default:
			t.Errorf("unexpected path %s", r.URL.Path)
		}
	}))
	defer srv.Close()
	c := New(srv.URL, staticToken("t"))
	providers, total, err := c.ListInsuranceProviders(context.Background(), "acc", 1, 20)
	if err != nil || len(providers) != 1 || total != 3 {
		t.Fatalf("providers=%v total=%d err=%v", providers, total, err)
	}
	expiring, err := c.ListExpiringInsurances(context.Background(), "acc", 1, 20)
	if err != nil || len(expiring) != 1 {
		t.Fatalf("expiring=%v err=%v", expiring, err)
	}
}

func TestShareLinksListAndFind(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/accounts/acc/shares" {
			t.Errorf("path = %s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"items": []map[string]any{{"id": "s1", "mode": "one_time", "status": "active"}},
			"total": 1, "page": 1, "size": 100, "pages": 1,
		})
	}))
	defer srv.Close()
	c := New(srv.URL, staticToken("t"))
	page, err := c.ListShareLinks(context.Background(), "accounts", "acc", 1, 25)
	if err != nil || len(page.Items) != 1 || page.Total != 1 {
		t.Fatalf("page=%+v err=%v", page, err)
	}
	item, err := c.FindShareLink(context.Background(), "accounts", "acc", "s1")
	if err != nil || item.ID != "s1" {
		t.Fatalf("item=%+v err=%v", item, err)
	}
}
