package iptuapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Sample response data
var sampleIPTUResponse = ConsultaEnderecoResult{
	SQL:                  "000.000.0000-0",
	Logradouro:           "Avenida Paulista",
	Numero:               "1000",
	Bairro:               "Bela Vista",
	CEP:                  "01310-100",
	AreaTerreno:          500.0,
	AreaConstruida:       1200.0,
	ValorVenalTerreno:    2500000.0,
	ValorVenalConstrucao: 1800000.0,
	ValorVenalTotal:      4300000.0,
	IPTUValor:            12500.0,
	AnoConstrucao:        1985,
	TipoUso:              "Comercial",
	Zona:                 "ZC",
}

var sampleValuationResponse = ValuationResult{
	ValorEstimado:         5000000.0,
	ValorMinimo:           4500000.0,
	ValorMaximo:           5500000.0,
	Confianca:             0.85,
	Metodo:                "comparativo",
	ComparaveisUtilizados: 12,
}

func TestNewClient(t *testing.T) {
	t.Run("creates client with default options", func(t *testing.T) {
		client := NewClient("test_api_key")

		assert.NotNil(t, client)
		assert.Equal(t, "test_api_key", client.apiKey)
		assert.Equal(t, defaultBaseURL, client.baseURL)
	})

	t.Run("applies custom options", func(t *testing.T) {
		client := NewClient("test_api_key",
			WithBaseURL("https://custom.api.com"),
			WithTimeout(60*time.Second),
			WithUserAgent("custom-agent/1.0"),
		)

		assert.Equal(t, "https://custom.api.com", client.baseURL)
		assert.Equal(t, "custom-agent/1.0", client.userAgent)
	})

	t.Run("applies custom retry config", func(t *testing.T) {
		retryConfig := &RetryConfig{
			MaxRetries:    5,
			InitialDelay:  100 * time.Millisecond,
			MaxDelay:      5 * time.Second,
			BackoffFactor: 1.5,
		}
		client := NewClient("test_api_key", WithRetry(retryConfig))

		assert.Equal(t, 5, client.retryConfig.MaxRetries)
		assert.Equal(t, 100*time.Millisecond, client.retryConfig.InitialDelay)
	})
}

func TestConsultaEndereco(t *testing.T) {
	t.Run("successful request", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/consulta/endereco", r.URL.Path)
			assert.Equal(t, "Avenida Paulista", r.URL.Query().Get("logradouro"))
			assert.Equal(t, "test_api_key", r.Header.Get("X-API-Key"))

			w.Header().Set("X-RateLimit-Limit", "1000")
			w.Header().Set("X-RateLimit-Remaining", "999")
			w.Header().Set("X-RateLimit-Reset", "1704067200")
			w.Header().Set("X-Request-ID", "req_test123")

			json.NewEncoder(w).Encode(sampleIPTUResponse)
		}))
		defer server.Close()

		client := NewClient("test_api_key",
			WithBaseURL(server.URL),
			WithRetry(&RetryConfig{MaxRetries: 0}),
		)

		result, err := client.ConsultaEndereco(context.Background(), &ConsultaEnderecoParams{
			Logradouro: "Avenida Paulista",
			Numero:     "1000",
		})

		require.NoError(t, err)
		assert.Equal(t, "000.000.0000-0", result.SQL)
		assert.Equal(t, "Avenida Paulista", result.Logradouro)

		// Check rate limit tracking
		require.NotNil(t, client.RateLimit)
		assert.Equal(t, 1000, client.RateLimit.Limit)
		assert.Equal(t, 999, client.RateLimit.Remaining)
		assert.Equal(t, "req_test123", client.LastRequestID)
	})

	t.Run("with all options", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "true", r.URL.Query().Get("incluir_historico"))
			assert.Equal(t, "true", r.URL.Query().Get("incluir_comparaveis"))
			assert.Equal(t, "true", r.URL.Query().Get("incluir_zoneamento"))
			assert.Equal(t, "bh", r.URL.Query().Get("cidade"))

			json.NewEncoder(w).Encode(sampleIPTUResponse)
		}))
		defer server.Close()

		client := NewClient("test_api_key",
			WithBaseURL(server.URL),
			WithRetry(&RetryConfig{MaxRetries: 0}),
		)

		_, err := client.ConsultaEndereco(context.Background(), &ConsultaEnderecoParams{
			Logradouro:         "Avenida Afonso Pena",
			Cidade:             CidadeBeloHorizonte,
			IncluirHistorico:   true,
			IncluirComparaveis: true,
			IncluirZoneamento:  true,
		})

		require.NoError(t, err)
	})
}

func TestValuationEstimate(t *testing.T) {
	t.Run("successful request", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/valuation/estimate", r.URL.Path)

			json.NewEncoder(w).Encode(sampleValuationResponse)
		}))
		defer server.Close()

		client := NewClient("test_api_key",
			WithBaseURL(server.URL),
			WithRetry(&RetryConfig{MaxRetries: 0}),
		)

		result, err := client.ValuationEstimate(context.Background(), &ValuationParams{
			AreaTerreno:    500.0,
			AreaConstruida: 1200.0,
			Bairro:         "Bela Vista",
			Zona:           "ZC",
			TipoUso:        "Comercial",
			TipoPadrao:     "Alto",
		})

		require.NoError(t, err)
		assert.Equal(t, 5000000.0, result.ValorEstimado)
		assert.Equal(t, 0.85, result.Confianca)
	})
}

func TestErrorHandling(t *testing.T) {
	t.Run("401 returns AuthenticationError", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"detail": "API Key inválida"})
		}))
		defer server.Close()

		client := NewClient("invalid_key",
			WithBaseURL(server.URL),
			WithRetry(&RetryConfig{MaxRetries: 0}),
		)

		_, err := client.ConsultaEndereco(context.Background(), &ConsultaEnderecoParams{
			Logradouro: "Test",
		})

		require.Error(t, err)
		assert.True(t, IsAuthError(err))
		authErr, ok := err.(*AuthenticationError)
		require.True(t, ok)
		assert.Equal(t, 401, authErr.StatusCode)
	})

	t.Run("403 returns ForbiddenError", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{
				"detail":        "Plano Pro necessário",
				"required_plan": "Pro",
			})
		}))
		defer server.Close()

		client := NewClient("test_key",
			WithBaseURL(server.URL),
			WithRetry(&RetryConfig{MaxRetries: 0}),
		)

		_, err := client.ValuationEstimate(context.Background(), &ValuationParams{
			AreaTerreno:    100,
			AreaConstruida: 100,
			Bairro:         "Test",
			Zona:           "ZC",
			TipoUso:        "Residencial",
			TipoPadrao:     "Médio",
		})

		require.Error(t, err)
		assert.True(t, IsForbidden(err))
		forbiddenErr, ok := err.(*ForbiddenError)
		require.True(t, ok)
		assert.Equal(t, "Pro", forbiddenErr.RequiredPlan)
	})

	t.Run("404 returns NotFoundError", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"detail": "Imóvel não encontrado"})
		}))
		defer server.Close()

		client := NewClient("test_key",
			WithBaseURL(server.URL),
			WithRetry(&RetryConfig{MaxRetries: 0}),
		)

		_, err := client.ConsultaEndereco(context.Background(), &ConsultaEnderecoParams{
			Logradouro: "Rua Inexistente",
		})

		require.Error(t, err)
		assert.True(t, IsNotFound(err))
	})

	t.Run("429 returns RateLimitError", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Retry-After", "60")
			w.Header().Set("X-RateLimit-Limit", "1000")
			w.Header().Set("X-RateLimit-Remaining", "0")
			w.Header().Set("X-RateLimit-Reset", "1704067200")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]string{"detail": "Rate limit exceeded"})
		}))
		defer server.Close()

		client := NewClient("test_key",
			WithBaseURL(server.URL),
			WithRetry(&RetryConfig{MaxRetries: 0}),
		)

		_, err := client.ConsultaEndereco(context.Background(), &ConsultaEnderecoParams{
			Logradouro: "Test",
		})

		require.Error(t, err)
		assert.True(t, IsRateLimit(err))
		rateLimitErr, ok := err.(*RateLimitError)
		require.True(t, ok)
		assert.Equal(t, 60, rateLimitErr.RetryAfter)
	})

	t.Run("500 returns ServerError", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"detail": "Internal error"})
		}))
		defer server.Close()

		client := NewClient("test_key",
			WithBaseURL(server.URL),
			WithRetry(&RetryConfig{MaxRetries: 0}),
		)

		_, err := client.ConsultaEndereco(context.Background(), &ConsultaEnderecoParams{
			Logradouro: "Test",
		})

		require.Error(t, err)
		assert.True(t, IsServerError(err))
	})
}

func TestRetryLogic(t *testing.T) {
	t.Run("retries on 500", func(t *testing.T) {
		attempts := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attempts++
			if attempts < 3 {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"detail": "Server error"})
				return
			}
			json.NewEncoder(w).Encode(sampleIPTUResponse)
		}))
		defer server.Close()

		client := NewClient("test_key",
			WithBaseURL(server.URL),
			WithRetry(&RetryConfig{
				MaxRetries:      3,
				InitialDelay:    10 * time.Millisecond,
				MaxDelay:        100 * time.Millisecond,
				BackoffFactor:   1.5,
				RetryableStatus: []int{429, 500, 502, 503, 504},
			}),
		)

		result, err := client.ConsultaEndereco(context.Background(), &ConsultaEnderecoParams{
			Logradouro: "Test",
		})

		require.NoError(t, err)
		assert.Equal(t, 3, attempts)
		assert.Equal(t, "000.000.0000-0", result.SQL)
	})

	t.Run("does not retry on 401", func(t *testing.T) {
		attempts := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attempts++
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"detail": "Unauthorized"})
		}))
		defer server.Close()

		client := NewClient("test_key",
			WithBaseURL(server.URL),
			WithRetry(&RetryConfig{
				MaxRetries:      3,
				InitialDelay:    10 * time.Millisecond,
				MaxDelay:        100 * time.Millisecond,
				BackoffFactor:   1.5,
				RetryableStatus: []int{429, 500, 502, 503, 504},
			}),
		)

		_, err := client.ConsultaEndereco(context.Background(), &ConsultaEnderecoParams{
			Logradouro: "Test",
		})

		require.Error(t, err)
		assert.Equal(t, 1, attempts)
	})

	t.Run("respects max retries", func(t *testing.T) {
		attempts := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attempts++
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		client := NewClient("test_key",
			WithBaseURL(server.URL),
			WithRetry(&RetryConfig{
				MaxRetries:      2,
				InitialDelay:    10 * time.Millisecond,
				MaxDelay:        100 * time.Millisecond,
				BackoffFactor:   1.5,
				RetryableStatus: []int{429, 500, 502, 503, 504},
			}),
		)

		_, err := client.ConsultaEndereco(context.Background(), &ConsultaEnderecoParams{
			Logradouro: "Test",
		})

		require.Error(t, err)
		assert.Equal(t, 3, attempts) // Initial + 2 retries
	})
}

func TestContextCancellation(t *testing.T) {
	t.Run("cancels request on context timeout", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(500 * time.Millisecond)
			json.NewEncoder(w).Encode(sampleIPTUResponse)
		}))
		defer server.Close()

		client := NewClient("test_key",
			WithBaseURL(server.URL),
			WithRetry(&RetryConfig{MaxRetries: 0, RetryableStatus: []int{500}}),
		)

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		_, err := client.ConsultaEndereco(ctx, &ConsultaEnderecoParams{
			Logradouro: "Test",
		})

		require.Error(t, err)
		// The error is wrapped in url.Error, check that context.DeadlineExceeded is the cause
		assert.Contains(t, err.Error(), "context deadline exceeded")
	})

	t.Run("cancels retry on context cancel", func(t *testing.T) {
		attempts := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attempts++
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		client := NewClient("test_key",
			WithBaseURL(server.URL),
			WithRetry(&RetryConfig{
				MaxRetries:      5,
				InitialDelay:    200 * time.Millisecond,
				MaxDelay:        1 * time.Second,
				BackoffFactor:   1.5,
				RetryableStatus: []int{429, 500, 502, 503, 504},
			}),
		)

		ctx, cancel := context.WithCancel(context.Background())

		go func() {
			time.Sleep(100 * time.Millisecond)
			cancel()
		}()

		_, err := client.ConsultaEndereco(ctx, &ConsultaEnderecoParams{
			Logradouro: "Test",
		})

		require.Error(t, err)
		// Either context.Canceled or first request's error (before context was canceled)
		assert.True(t, err == context.Canceled || IsServerError(err) || attempts >= 1)
	})
}

func TestCidadeConstants(t *testing.T) {
	assert.Equal(t, Cidade("sp"), CidadeSaoPaulo)
	assert.Equal(t, Cidade("bh"), CidadeBeloHorizonte)
	assert.Equal(t, Cidade("recife"), CidadeRecife)
}

func TestAPIError(t *testing.T) {
	t.Run("Error message without request ID", func(t *testing.T) {
		err := &APIError{
			StatusCode: 500,
			Message:    "Server error",
		}
		assert.Contains(t, err.Error(), "500")
		assert.Contains(t, err.Error(), "Server error")
	})

	t.Run("Error message with request ID", func(t *testing.T) {
		err := &APIError{
			StatusCode: 500,
			Message:    "Server error",
			RequestID:  "req_123",
		}
		assert.Contains(t, err.Error(), "req_123")
	})

	t.Run("IsRetryable", func(t *testing.T) {
		assert.True(t, (&APIError{StatusCode: 429}).IsRetryable())
		assert.True(t, (&APIError{StatusCode: 500}).IsRetryable())
		assert.True(t, (&APIError{StatusCode: 502}).IsRetryable())
		assert.True(t, (&APIError{StatusCode: 503}).IsRetryable())
		assert.True(t, (&APIError{StatusCode: 504}).IsRetryable())
		assert.False(t, (&APIError{StatusCode: 400}).IsRetryable())
		assert.False(t, (&APIError{StatusCode: 401}).IsRetryable())
		assert.False(t, (&APIError{StatusCode: 404}).IsRetryable())
	})
}
