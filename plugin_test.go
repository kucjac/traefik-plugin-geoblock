package traefik_plugin_geoblock_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	geoblock "github.com/kucjac/traefik-plugin-geoblock"
)

const (
	pluginName = "geoblock"
	dbFilePath = "./IP2LOCATION-LITE-DB1.IPV6.BIN"
)

type noopHandler struct{}

func (n noopHandler) ServeHTTP(rw http.ResponseWriter, _ *http.Request) {
	rw.WriteHeader(http.StatusTeapot)
}

func TestNew(t *testing.T) {
	t.Run("Disabled", func(t *testing.T) {
		plugin, err := geoblock.New(context.TODO(), &noopHandler{}, &geoblock.Config{Enabled: false}, pluginName)
		if err != nil {
			t.Fatal(err)
		}

		req := httptest.NewRequest(http.MethodGet, "/foobar", nil)

		rr := httptest.NewRecorder()
		plugin.ServeHTTP(rr, req)

		if http.StatusTeapot != rr.Code {
			t.Fatalf("expected: %v is %v", http.StatusTeapot, rr.Code)
		}
	})

	t.Run("NoNextHandler", func(t *testing.T) {
		plugin, err := geoblock.New(context.TODO(), nil, &geoblock.Config{Enabled: true}, pluginName)
		if err == nil {
			t.Fatal("an error is expected but is nil")
		}
		if plugin != nil {
			t.Error("plugin is expected to be nil")
		}
	})

	t.Run("NoConfig", func(t *testing.T) {
		plugin, err := geoblock.New(context.TODO(), &noopHandler{}, nil, pluginName)
		if err == nil {
			t.Fatal("an error is expected but is nil")
		}
		if plugin != nil {
			t.Error("plugin is expected to be nil")
		}
	})

	t.Run("NoDatabaseFilePath", func(t *testing.T) {
		plugin, err := geoblock.New(context.TODO(), &noopHandler{}, &geoblock.Config{Enabled: true}, pluginName)
		if err == nil {
			t.Fatal(err)
		}
		if plugin != nil {
			t.Error("plugin is expected to be nil")
		}
	})
}

func TestPlugin_ServeHTTP(t *testing.T) {
	t.Run("Allowed", func(t *testing.T) {
		cfg := &geoblock.Config{Enabled: true, DatabaseFilePath: dbFilePath, AllowedCountries: []string{"US"}}
		plugin, err := geoblock.New(context.TODO(), &noopHandler{}, cfg, pluginName)
		if err != nil {
			t.Fatal(err)
		}

		req := httptest.NewRequest(http.MethodGet, "/foobar", nil)
		req.Header.Set("X-Real-IP", "1.1.1.1")

		rr := httptest.NewRecorder()
		plugin.ServeHTTP(rr, req)

		if http.StatusTeapot != rr.Code {
			t.Fatalf("expected: %v is %v", http.StatusTeapot, rr.Code)
		}
	})

	t.Run("AllowedPrivate", func(t *testing.T) {
		cfg := &geoblock.Config{Enabled: true, DatabaseFilePath: dbFilePath, AllowedCountries: []string{}, AllowPrivate: true}
		plugin, err := geoblock.New(context.TODO(), &noopHandler{}, cfg, pluginName)
		if err != nil {
			t.Fatal(err)
		}

		req := httptest.NewRequest(http.MethodGet, "/foobar", nil)
		req.Header.Set("X-Real-IP", "192.168.178.66")

		rr := httptest.NewRecorder()
		plugin.ServeHTTP(rr, req)

		if http.StatusTeapot != rr.Code {
			t.Fatalf("expected: %v is %v", http.StatusTeapot, rr.Code)
		}
	})

	t.Run("Disallowed", func(t *testing.T) {
		cfg := &geoblock.Config{Enabled: true, DatabaseFilePath: dbFilePath, AllowedCountries: []string{"DE"}}
		plugin, err := geoblock.New(context.TODO(), &noopHandler{}, cfg, pluginName)
		if err != nil {
			t.Fatal(err)
		}

		req := httptest.NewRequest(http.MethodGet, "/foobar", nil)
		req.Header.Set("X-Real-IP", "1.1.1.1")

		rr := httptest.NewRecorder()
		plugin.ServeHTTP(rr, req)

		if http.StatusForbidden != rr.Code {
			t.Fatalf("expected: %v is %v", http.StatusForbidden, rr.Code)
		}
	})

	t.Run("DisallowedCountries", func(t *testing.T) {
		t.Run("Pass", func(t *testing.T) {
			cfg := &geoblock.Config{Enabled: true, DatabaseFilePath: dbFilePath, DisallowedCountries: []string{"PL"}}
			plugin, err := geoblock.New(context.TODO(), &noopHandler{}, cfg, pluginName)
			if err != nil {
				t.Fatal(err)
			}

			req := httptest.NewRequest(http.MethodGet, "/foobar", nil)

			const randomCzechIP = "188.92.102.22"
			req.Header.Set("X-Real-IP", randomCzechIP)

			rr := httptest.NewRecorder()
			plugin.ServeHTTP(rr, req)

			if http.StatusTeapot != rr.Code {
				t.Fatalf("expected: %v is %v", http.StatusTeapot, rr.Code)
			}
		})

		t.Run("Forbid", func(t *testing.T) {
			cfg := &geoblock.Config{Enabled: true, DatabaseFilePath: dbFilePath, DisallowedCountries: []string{"PL"}}
			plugin, err := geoblock.New(context.TODO(), &noopHandler{}, cfg, pluginName)
			if err != nil {
				t.Fatal(err)
			}

			req := httptest.NewRequest(http.MethodGet, "/foobar", nil)
			// Define some random polish IP address.
			const ranomPolandIP = "195.136.153.130"
			req.Header.Set("X-Real-IP", ranomPolandIP)

			rr := httptest.NewRecorder()
			plugin.ServeHTTP(rr, req)

			if http.StatusForbidden != rr.Code {
				t.Fatalf("expected: %v is %v", http.StatusForbidden, rr.Code)
			}
		})
	})

	t.Run("DisallowedPrivate", func(t *testing.T) {
		cfg := &geoblock.Config{Enabled: true, DatabaseFilePath: dbFilePath, AllowedCountries: []string{}, AllowPrivate: false}
		plugin, err := geoblock.New(context.TODO(), &noopHandler{}, cfg, pluginName)
		if err != nil {
			t.Fatal(err)
		}

		req := httptest.NewRequest(http.MethodGet, "/foobar", nil)
		req.Header.Set("X-Real-IP", "192.168.178.66")

		rr := httptest.NewRecorder()
		plugin.ServeHTTP(rr, req)

		if http.StatusForbidden != rr.Code {
			t.Fatalf("expected: %v is %v", http.StatusForbidden, rr.Code)
		}
	})
}
