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
	"errors"
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

// Cidade represents supported cities.
type Cidade string

const (
	CidadeSaoPaulo      Cidade = "sp"
	CidadeBeloHorizonte Cidade = "bh"
	CidadeRecife        Cidade = "recife"
	CidadePortoAlegre   Cidade = "poa"
	CidadeFortaleza     Cidade = "fortaleza"
	CidadeCuritiba      Cidade = "curitiba"
	CidadeRioDeJaneiro  Cidade = "rj"
	CidadeBrasilia      Cidade = "brasilia"
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
	SQLBase     string  `json:"sql_base"`
	Logradouro  string  `json:"logradouro"`
	Numero      string  `json:"numero"`
	Bairro      string  `json:"bairro"`
	CEP         string  `json:"cep"`
	AreaTerreno float64 `json:"area_terreno"`
	TipoUso     string  `json:"tipo_uso"`
}

// DadosIPTU represents detailed IPTU data.
type DadosIPTU struct {
	SQL             string  `json:"sql"`
	AnoReferencia   int     `json:"ano_referencia"`
	Logradouro      string  `json:"logradouro"`
	Numero          int     `json:"numero"`
	Bairro          string  `json:"bairro"`
	CEP             string  `json:"cep"`
	AreaTerreno     float64 `json:"area_terreno"`
	AreaConstruida  float64 `json:"area_construida"`
	ValorTerreno    float64 `json:"valor_terreno"`
	ValorConstrucao float64 `json:"valor_construcao"`
	ValorVenal      float64 `json:"valor_venal"`
	Finalidade      string  `json:"finalidade"`
	TipoConstrucao  string  `json:"tipo_construcao"`
	AnoConstrucao   int     `json:"ano_construcao"`
}

// ConsultaIPTUResult represents the result from multi-city IPTU query.
type ConsultaIPTUResult struct {
	SQL             string      `json:"sql"`
	Ano             int         `json:"ano"`
	Logradouro      string      `json:"logradouro"`
	Numero          interface{} `json:"numero"` // int or string (Recife uses string)
	Complemento     *string     `json:"complemento"`
	Bairro          *string     `json:"bairro"`
	CEP             string      `json:"cep"`
	AreaTerreno     *float64    `json:"area_terreno"`
	AreaConstruida  *float64    `json:"area_construida"`
	ValorTerreno    *float64    `json:"valor_terreno"`
	ValorConstrucao *float64    `json:"valor_construcao"`
	ValorVenal      float64     `json:"valor_venal"`
	ValorImovel     *float64    `json:"valor_imovel"` // Recife: estimated total value
	ValorIPTU       *float64    `json:"valor_iptu"`   // Recife: IPTU amount
	Finalidade      *string     `json:"finalidade"`
	TipoConstrucao  *string     `json:"tipo_construcao"`
	AnoConstrucao   *int        `json:"ano_construcao"`
	Pavimentos      *int        `json:"pavimentos"`
	FracaoIdeal     *string     `json:"fracao_ideal"`
	Latitude        *float64    `json:"latitude"`  // Recife: coordinates
	Longitude       *float64    `json:"longitude"` // Recife: coordinates
	Cidade          string      `json:"cidade"`
	Fonte           string      `json:"fonte"`
}

// ConsultaEnderecoResult represents the result of an address query.
type ConsultaEnderecoResult struct {
	Success   bool                 `json:"success"`
	Data      ConsultaEnderecoData `json:"data"`
	DadosIPTU DadosIPTU            `json:"dados_iptu"`
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
	Success       bool    `json:"success"`
	ValorEstimado float64 `json:"valor_estimado"`
	ValorMinimo   float64 `json:"valor_minimo"`
	ValorMaximo   float64 `json:"valor_maximo"`
	ValorM2       float64 `json:"valor_m2"`
	Confianca     float64 `json:"confianca"`
	ModeloVersao  string  `json:"modelo_versao"`
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
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusNotFound
	}
	return false
}

// IsRateLimit returns true if the error is a 429 Rate Limit.
func IsRateLimit(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusTooManyRequests
	}
	return false
}

// IsAuthError returns true if the error is a 401 Authentication error.
func IsAuthError(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
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
		jsonBody, marshalErr := json.Marshal(body)
		if marshalErr != nil {
			return marshalErr
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
		_ = json.Unmarshal(respBody, &errResp) // Best effort to parse error detail

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

// EvaluateParams contains parameters for property evaluation.
type EvaluateParams struct {
	// SQL number of the property (alternative to address)
	SQL string `json:"sql,omitempty"`
	// Street name (alternative to SQL)
	Logradouro string `json:"logradouro,omitempty"`
	// Property number
	Numero int `json:"numero,omitempty"`
	// Unit/apartment
	Complemento string `json:"complemento,omitempty"`
	// Neighborhood
	Bairro string `json:"bairro,omitempty"`
	// City code (sp, bh)
	Cidade string `json:"cidade,omitempty"`
	// Include ITBI-based estimate (default: true)
	IncluirItbi *bool `json:"incluir_itbi,omitempty"`
	// Include comparable properties (default: true)
	IncluirComparaveis *bool `json:"incluir_comparaveis,omitempty"`
}

// AVMEstimate represents the AVM (Machine Learning) model estimate.
type AVMEstimate struct {
	ValorEstimado float64 `json:"valor_estimado"`
	ValorMinimo   float64 `json:"valor_minimo"`
	ValorMaximo   float64 `json:"valor_maximo"`
	ValorM2       float64 `json:"valor_m2"`
	Confianca     float64 `json:"confianca"`
	ModeloVersao  string  `json:"modelo_versao"`
}

// ITBIMarketEstimate represents the estimate based on real ITBI transactions.
type ITBIMarketEstimate struct {
	ValorEstimado   float64 `json:"valor_estimado"`
	FaixaMinima     float64 `json:"faixa_minima"`
	FaixaMaxima     float64 `json:"faixa_maxima"`
	ValorM2Mediana  float64 `json:"valor_m2_mediana"`
	TotalTransacoes int     `json:"total_transacoes"`
	Periodo         string  `json:"periodo"`
	Fonte           string  `json:"fonte"`
}

// FinalValuation represents the combined final valuation.
type FinalValuation struct {
	Estimado  float64 `json:"estimado"`
	Minimo    float64 `json:"minimo"`
	Maximo    float64 `json:"maximo"`
	Metodo    string  `json:"metodo"`
	PesoAvm   float64 `json:"peso_avm"`
	PesoItbi  float64 `json:"peso_itbi"`
	Confianca float64 `json:"confianca"`
	Nota      *string `json:"nota,omitempty"`
}

// PropertyEvaluationMetadata contains metadata about the evaluation.
type PropertyEvaluationMetadata struct {
	ProcessadoEm string   `json:"processado_em"`
	Fontes       []string `json:"fontes"`
	Cidade       string   `json:"cidade"`
}

// PropertyEvaluation represents the complete property evaluation result.
type PropertyEvaluation struct {
	Success       bool                       `json:"success"`
	Imovel        map[string]interface{}     `json:"imovel"`
	AvaliacaoAvm  *AVMEstimate               `json:"avaliacao_avm,omitempty"`
	AvaliacaoItbi *ITBIMarketEstimate        `json:"avaliacao_itbi,omitempty"`
	ValorFinal    FinalValuation             `json:"valor_final"`
	Comparaveis   map[string]interface{}     `json:"comparaveis,omitempty"`
	Metadata      PropertyEvaluationMetadata `json:"metadata"`
}

// ValuationEvaluate evaluates a property by address OR SQL number.
// Combines data from the AVM (ML) model with real ITBI transactions.
// Requires Pro plan or higher.
//
// Example by SQL:
//
//	result, err := client.ValuationEvaluate(iptuapi.EvaluateParams{
//	    SQL:    "123.456.0001-0",
//	    Cidade: "sp",
//	})
//
// Example by address:
//
//	result, err := client.ValuationEvaluate(iptuapi.EvaluateParams{
//	    Logradouro: "Avenida Paulista",
//	    Numero:     1000,
//	    Cidade:     "sp",
//	})
func (c *Client) ValuationEvaluate(params EvaluateParams) (*PropertyEvaluation, error) {
	// Set defaults
	if params.Cidade == "" {
		params.Cidade = "sp"
	}
	if params.IncluirItbi == nil {
		defaultTrue := true
		params.IncluirItbi = &defaultTrue
	}
	if params.IncluirComparaveis == nil {
		defaultTrue := true
		params.IncluirComparaveis = &defaultTrue
	}

	var result PropertyEvaluation
	err := c.doRequest("POST", "/valuation/evaluate", nil, params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ConsultaIPTUOptions contains optional parameters for ConsultaIPTU.
type ConsultaIPTUOptions struct {
	Numero *int
	Ano    int
	Limit  int
}

// ConsultaIPTU searches for IPTU data by address in any supported city.
// Supported cities: CidadeSaoPaulo, CidadeBeloHorizonte, CidadeRecife
// Recife includes latitude/longitude coordinates.
//
// Example:
//
//	results, err := client.ConsultaIPTU(iptuapi.CidadeRecife, "Boa Viagem", nil)
//	// or with options:
//	opts := &iptuapi.ConsultaIPTUOptions{Ano: 2025, Limit: 10}
//	results, err := client.ConsultaIPTU(iptuapi.CidadeSaoPaulo, "Paulista", opts)
func (c *Client) ConsultaIPTU(cidade Cidade, logradouro string, opts *ConsultaIPTUOptions) ([]ConsultaIPTUResult, error) {
	params := url.Values{}
	params.Set("logradouro", logradouro)

	if opts != nil {
		if opts.Numero != nil {
			params.Set("numero", fmt.Sprintf("%d", *opts.Numero))
		}
		if opts.Ano > 0 {
			params.Set("ano", fmt.Sprintf("%d", opts.Ano))
		} else {
			params.Set("ano", "2025")
		}
		if opts.Limit > 0 {
			params.Set("limit", fmt.Sprintf("%d", opts.Limit))
		} else {
			params.Set("limit", "20")
		}
	} else {
		params.Set("ano", "2025")
		params.Set("limit", "20")
	}

	var result []ConsultaIPTUResult
	err := c.doRequest("GET", fmt.Sprintf("/dados/iptu/%s/endereco", cidade), params, nil, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ConsultaIPTUSQL searches for IPTU data by property identifier in any supported city.
// For São Paulo, use the SQL number. For Belo Horizonte, use the Índice Cadastral.
// For Recife, use the Contribuinte number. Recife includes latitude/longitude.
//
// Example:
//
//	// São Paulo
//	results, err := client.ConsultaIPTUSQL(iptuapi.CidadeSaoPaulo, "00904801381", nil)
//	// Belo Horizonte
//	results, err := client.ConsultaIPTUSQL(iptuapi.CidadeBeloHorizonte, "007028 005 0086", nil)
//	// Recife
//	results, err := client.ConsultaIPTUSQL(iptuapi.CidadeRecife, "123456789", nil)
func (c *Client) ConsultaIPTUSQL(cidade Cidade, identificador string, ano *int) ([]ConsultaIPTUResult, error) {
	params := url.Values{}
	if ano != nil {
		params.Set("ano", fmt.Sprintf("%d", *ano))
	}

	var result []ConsultaIPTUResult
	err := c.doRequest("GET", fmt.Sprintf("/dados/iptu/%s/sql/%s", cidade, url.PathEscape(identificador)), params, nil, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
