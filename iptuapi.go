// Package iptuapi provides a client for the IPTU API.
//
// SDK oficial para integração com a IPTU API.
//
// Example:
//
//	client := iptuapi.NewClient("sua_api_key")
//	resultado, err := client.ConsultaEndereco("Avenida Paulista", "1000")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(resultado)
package iptuapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultBaseURL = "https://iptuapi.com.br/api/v1"
	defaultTimeout = 30 * time.Second
)

// Client represents an IPTU API client.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// ClientOption configures the Client.
type ClientOption func(*Client)

// WithBaseURL sets a custom base URL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// WithTimeout sets a custom timeout.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

// NewClient creates a new IPTU API client.
func NewClient(apiKey string, opts ...ClientOption) *Client {
	c := &Client{
		apiKey:  apiKey,
		baseURL: defaultBaseURL,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// ConsultaEnderecoData represents the basic address data.
type ConsultaEnderecoData struct {
	SQLBase        string  `json:"sql_base"`
	Logradouro     string  `json:"logradouro"`
	Numero         string  `json:"numero"`
	Bairro         string  `json:"bairro"`
	CEP            string  `json:"cep"`
	AreaTerreno    float64 `json:"area_terreno"`
	TipoUso        string  `json:"tipo_uso"`
}

// DadosIPTU represents detailed IPTU data.
type DadosIPTU struct {
	SQL            string  `json:"sql"`
	AnoReferencia  int     `json:"ano_referencia"`
	Logradouro     string  `json:"logradouro"`
	Numero         int     `json:"numero"`
	Bairro         string  `json:"bairro"`
	CEP            string  `json:"cep"`
	AreaTerreno    float64 `json:"area_terreno"`
	AreaConstruida float64 `json:"area_construida"`
	ValorTerreno   float64 `json:"valor_terreno"`
	ValorConstrucao float64 `json:"valor_construcao"`
	ValorVenal     float64 `json:"valor_venal"`
	Finalidade     string  `json:"finalidade"`
	TipoConstrucao string  `json:"tipo_construcao"`
	AnoConstrucao  int     `json:"ano_construcao"`
}

// ConsultaEnderecoResult represents the result of an address query.
type ConsultaEnderecoResult struct {
	Success   bool                  `json:"success"`
	Data      ConsultaEnderecoData  `json:"data"`
	DadosIPTU DadosIPTU             `json:"dados_iptu"`
}

// ConsultaSQLResult represents the result of a SQL query.
type ConsultaSQLResult struct {
	SQL                  string  `json:"sql"`
	Ano                  int     `json:"ano"`
	ValorVenal           float64 `json:"valor_venal"`
	ValorVenalTerreno    float64 `json:"valor_venal_terreno"`
	ValorVenalConstrucao float64 `json:"valor_venal_construcao"`
	IPTUValor            float64 `json:"iptu_valor"`
	Logradouro           string  `json:"logradouro"`
	Numero               string  `json:"numero"`
	Bairro               string  `json:"bairro"`
	AreaTerreno          float64 `json:"area_terreno"`
	AreaConstruida       float64 `json:"area_construida"`
}

// ValuationParams contains parameters for property valuation.
type ValuationParams struct {
	AreaTerreno    float64 `json:"area_terreno"`
	AreaConstruida float64 `json:"area_construida"`
	Bairro         string  `json:"bairro"`
	Zona           string  `json:"zona"`
	TipoUso        string  `json:"tipo_uso"`
	TipoPadrao     string  `json:"tipo_padrao"`
	AnoConstrucao  int     `json:"ano_construcao,omitempty"`
}

// ValuationResult represents the result of a valuation estimate.
type ValuationResult struct {
	Success        bool    `json:"success"`
	ValorEstimado  float64 `json:"valor_estimado"`
	ValorMinimo    float64 `json:"valor_minimo"`
	ValorMaximo    float64 `json:"valor_maximo"`
	ValorM2        float64 `json:"valor_m2"`
	Confianca      float64 `json:"confianca"`
	ModeloVersao   string  `json:"modelo_versao"`
}

// APIError represents an error from the IPTU API.
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("IPTU API error (status %d): %s", e.StatusCode, e.Message)
}

// IsNotFound returns true if the error is a 404 Not Found.
func IsNotFound(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == http.StatusNotFound
	}
	return false
}

// IsRateLimit returns true if the error is a 429 Rate Limit.
func IsRateLimit(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == http.StatusTooManyRequests
	}
	return false
}

// IsAuthError returns true if the error is a 401 Authentication error.
func IsAuthError(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == http.StatusUnauthorized
	}
	return false
}

func (c *Client) doRequest(method, endpoint string, params url.Values, body interface{}, result interface{}) error {
	u, err := url.Parse(c.baseURL + endpoint)
	if err != nil {
		return err
	}

	if params != nil {
		u.RawQuery = params.Encode()
	}

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reqBody = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequest(method, u.String(), reqBody)
	if err != nil {
		return err
	}

	req.Header.Set("X-API-Key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		var errResp struct {
			Detail string `json:"detail"`
		}
		json.Unmarshal(respBody, &errResp)

		message := errResp.Detail
		if message == "" {
			switch resp.StatusCode {
			case http.StatusUnauthorized:
				message = "API Key inválida ou expirada"
			case http.StatusForbidden:
				message = "Plano não autorizado para este recurso"
			case http.StatusNotFound:
				message = "Recurso não encontrado"
			case http.StatusTooManyRequests:
				message = "Limite de requisições excedido"
			default:
				message = "Erro na API"
			}
		}

		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    message,
		}
	}

	return json.Unmarshal(respBody, result)
}

// ConsultaEndereco searches for property data by address.
func (c *Client) ConsultaEndereco(logradouro, numero string) (*ConsultaEnderecoResult, error) {
	params := url.Values{}
	params.Set("logradouro", logradouro)
	if numero != "" {
		params.Set("numero", numero)
	}

	var result ConsultaEnderecoResult
	err := c.doRequest("GET", "/consulta/endereco", params, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ConsultaSQL searches for property data by SQL number.
// Requires Starter plan or higher.
func (c *Client) ConsultaSQL(sql string) (*ConsultaSQLResult, error) {
	params := url.Values{}
	params.Set("sql", sql)

	var result ConsultaSQLResult
	err := c.doRequest("GET", "/consulta/sql", params, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ValuationEstimate estimates the market value of a property.
// Requires Pro plan or higher.
func (c *Client) ValuationEstimate(params ValuationParams) (*ValuationResult, error) {
	var result ValuationResult
	err := c.doRequest("POST", "/valuation/estimate", nil, params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
