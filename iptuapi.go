// Package iptuapi provides a client for the IPTU API.
//
// SDK oficial para integração com a IPTU API.
// Suporta context para cancelamento, retry automático, logging e rate limit tracking.
//
// Example:
//
//	client := iptuapi.NewClient("sua_api_key")
//	ctx := context.Background()
//	resultado, err := client.ConsultaEndereco(ctx, &iptuapi.ConsultaEnderecoParams{
//	    Logradouro: "Avenida Paulista",
//	    Numero:     "1000",
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(resultado)
package iptuapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Version is the SDK version.
const Version = "2.1.2"

const (
	defaultBaseURL = "https://iptuapi.com.br/api/v1"
	defaultTimeout = 30 * time.Second
)

// Cidade represents available cities.
type Cidade string

const (
	CidadeSaoPaulo       Cidade = "sp"
	CidadeBeloHorizonte  Cidade = "bh"
	CidadeRecife         Cidade = "recife"
	CidadePortoAlegre    Cidade = "poa"
	CidadeFortaleza      Cidade = "fortaleza"
	CidadeCuritiba       Cidade = "curitiba"
	CidadeRioDeJaneiro   Cidade = "rj"
	CidadeBrasilia       Cidade = "brasilia"
)

// Logger interface for custom logging.
type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

// DefaultLogger is a simple logger that uses the standard log package.
type DefaultLogger struct {
	Enabled bool
}

func (l *DefaultLogger) Debug(msg string, args ...interface{}) {
	if l.Enabled {
		log.Printf("[DEBUG] "+msg, args...)
	}
}

func (l *DefaultLogger) Info(msg string, args ...interface{}) {
	if l.Enabled {
		log.Printf("[INFO] "+msg, args...)
	}
}

func (l *DefaultLogger) Warn(msg string, args ...interface{}) {
	if l.Enabled {
		log.Printf("[WARN] "+msg, args...)
	}
}

func (l *DefaultLogger) Error(msg string, args ...interface{}) {
	if l.Enabled {
		log.Printf("[ERROR] "+msg, args...)
	}
}

// RetryConfig configures retry behavior.
type RetryConfig struct {
	MaxRetries      int
	InitialDelay    time.Duration
	MaxDelay        time.Duration
	BackoffFactor   float64
	RetryableStatus []int
}

// DefaultRetryConfig returns the default retry configuration.
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxRetries:      3,
		InitialDelay:    500 * time.Millisecond,
		MaxDelay:        10 * time.Second,
		BackoffFactor:   2.0,
		RetryableStatus: []int{429, 500, 502, 503, 504},
	}
}

// RateLimitInfo contains rate limit information from API response.
type RateLimitInfo struct {
	Limit     int
	Remaining int
	Reset     int64
	ResetTime time.Time
}

// Client represents an IPTU API client.
type Client struct {
	apiKey      string
	baseURL     string
	httpClient  *http.Client
	retryConfig *RetryConfig
	logger      Logger
	userAgent   string

	// Rate limit info from last request
	RateLimit     *RateLimitInfo
	LastRequestID string
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

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithRetry sets retry configuration.
func WithRetry(config *RetryConfig) ClientOption {
	return func(c *Client) {
		c.retryConfig = config
	}
}

// WithLogger sets a custom logger.
func WithLogger(logger Logger) ClientOption {
	return func(c *Client) {
		c.logger = logger
	}
}

// WithUserAgent sets a custom User-Agent header.
func WithUserAgent(userAgent string) ClientOption {
	return func(c *Client) {
		c.userAgent = userAgent
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
		retryConfig: DefaultRetryConfig(),
		logger:      &DefaultLogger{Enabled: false},
		userAgent:   "iptuapi-go/" + Version,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// =============================================================================
// Types
// =============================================================================

// ConsultaEnderecoParams contains parameters for address query.
type ConsultaEnderecoParams struct {
	Logradouro         string
	Numero             string
	Complemento        string
	Cidade             Cidade
	IncluirHistorico   bool
	IncluirComparaveis bool
	IncluirZoneamento  bool
}

// ConsultaEnderecoResult represents the result of an address query.
type ConsultaEnderecoResult struct {
	SQL                  string            `json:"sql"`
	Logradouro           string            `json:"logradouro"`
	Numero               string            `json:"numero,omitempty"`
	Complemento          string            `json:"complemento,omitempty"`
	Bairro               string            `json:"bairro,omitempty"`
	CEP                  string            `json:"cep,omitempty"`
	AreaTerreno          float64           `json:"area_terreno,omitempty"`
	AreaConstruida       float64           `json:"area_construida,omitempty"`
	ValorVenalTerreno    float64           `json:"valor_venal_terreno,omitempty"`
	ValorVenalConstrucao float64           `json:"valor_venal_construcao,omitempty"`
	ValorVenalTotal      float64           `json:"valor_venal_total,omitempty"`
	IPTUValor            float64           `json:"iptu_valor,omitempty"`
	AnoConstrucao        int               `json:"ano_construcao,omitempty"`
	TipoUso              string            `json:"tipo_uso,omitempty"`
	Zona                 string            `json:"zona,omitempty"`
	Historico            []HistoricoItem   `json:"historico,omitempty"`
	Comparaveis          []ComparavelItem  `json:"comparaveis,omitempty"`
	Zoneamento           *ZoneamentoResult `json:"zoneamento,omitempty"`
}

// ConsultaSQLResult represents the result of a SQL query.
type ConsultaSQLResult struct {
	SQL                  string  `json:"sql"`
	Ano                  int     `json:"ano,omitempty"`
	ValorVenal           float64 `json:"valor_venal,omitempty"`
	ValorVenalTerreno    float64 `json:"valor_venal_terreno,omitempty"`
	ValorVenalConstrucao float64 `json:"valor_venal_construcao,omitempty"`
	ValorVenalTotal      float64 `json:"valor_venal_total,omitempty"`
	IPTUValor            float64 `json:"iptu_valor,omitempty"`
	Logradouro           string  `json:"logradouro,omitempty"`
	Numero               string  `json:"numero,omitempty"`
	Bairro               string  `json:"bairro,omitempty"`
	AreaTerreno          float64 `json:"area_terreno,omitempty"`
	AreaConstruida       float64 `json:"area_construida,omitempty"`
}

// HistoricoItem represents a historical value entry.
type HistoricoItem struct {
	Ano                  int     `json:"ano"`
	ValorVenalTerreno    float64 `json:"valor_venal_terreno,omitempty"`
	ValorVenalConstrucao float64 `json:"valor_venal_construcao,omitempty"`
	ValorVenalTotal      float64 `json:"valor_venal_total,omitempty"`
	IPTUValor            float64 `json:"iptu_valor,omitempty"`
}

// ComparavelItem represents a comparable property.
type ComparavelItem struct {
	SQL             string  `json:"sql,omitempty"`
	Logradouro      string  `json:"logradouro,omitempty"`
	Numero          string  `json:"numero,omitempty"`
	Bairro          string  `json:"bairro,omitempty"`
	AreaTerreno     float64 `json:"area_terreno,omitempty"`
	AreaConstruida  float64 `json:"area_construida,omitempty"`
	ValorVenalTotal float64 `json:"valor_venal_total,omitempty"`
	DistanciaMetros float64 `json:"distancia_metros,omitempty"`
}

// ZoneamentoResult represents zoning data.
type ZoneamentoResult struct {
	Zona                         string  `json:"zona,omitempty"`
	ZonaDescricao                string  `json:"zona_descricao,omitempty"`
	CoeficienteAproveitamentoBasico float64 `json:"coeficiente_aproveitamento_basico,omitempty"`
	CoeficienteAproveitamentoMaximo float64 `json:"coeficiente_aproveitamento_maximo,omitempty"`
	TaxaOcupacaoMaxima            float64 `json:"taxa_ocupacao_maxima,omitempty"`
	GabaritoMaximo                int     `json:"gabarito_maximo,omitempty"`
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
	Cidade         Cidade  `json:"cidade,omitempty"`
}

// ValuationResult represents the result of a valuation estimate.
type ValuationResult struct {
	ValorEstimado        float64 `json:"valor_estimado"`
	ValorMinimo          float64 `json:"valor_minimo,omitempty"`
	ValorMaximo          float64 `json:"valor_maximo,omitempty"`
	Confianca            float64 `json:"confianca,omitempty"`
	Metodo               string  `json:"metodo,omitempty"`
	ComparaveisUtilizados int     `json:"comparaveis_utilizados,omitempty"`
	DataAvaliacao        string  `json:"data_avaliacao,omitempty"`
}

// BatchValuationResult represents batch valuation results.
type BatchValuationResult struct {
	Resultados       []ValuationResult `json:"resultados"`
	TotalProcessados int               `json:"total_processados"`
	TotalErros       int               `json:"total_erros"`
	Erros            []BatchError      `json:"erros,omitempty"`
}

// BatchError represents an error in batch processing.
type BatchError struct {
	Index int    `json:"index"`
	Error string `json:"error"`
}

// =============================================================================
// Errors
// =============================================================================

// APIError represents an error from the IPTU API.
type APIError struct {
	StatusCode   int
	Message      string
	RequestID    string
	ResponseBody map[string]interface{}
}

func (e *APIError) Error() string {
	if e.RequestID != "" {
		return fmt.Sprintf("IPTU API error (status %d, request %s): %s", e.StatusCode, e.RequestID, e.Message)
	}
	return fmt.Sprintf("IPTU API error (status %d): %s", e.StatusCode, e.Message)
}

// IsRetryable returns true if the error can be retried.
func (e *APIError) IsRetryable() bool {
	for _, status := range []int{429, 500, 502, 503, 504} {
		if e.StatusCode == status {
			return true
		}
	}
	return false
}

// AuthenticationError indicates invalid API key.
type AuthenticationError struct {
	*APIError
}

// ForbiddenError indicates plan not authorized.
type ForbiddenError struct {
	*APIError
	RequiredPlan string
}

// NotFoundError indicates resource not found.
type NotFoundError struct {
	*APIError
}

// RateLimitError indicates rate limit exceeded.
type RateLimitError struct {
	*APIError
	RetryAfter int
	Limit      int
	Remaining  int
}

// ValidationError indicates invalid parameters.
type ValidationError struct {
	*APIError
	Errors []FieldError
}

// FieldError represents a field validation error.
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ServerError indicates internal server error.
type ServerError struct {
	*APIError
}

// IsNotFound returns true if the error is a 404 Not Found.
func IsNotFound(err error) bool {
	_, ok := err.(*NotFoundError)
	return ok
}

// IsRateLimit returns true if the error is a 429 Rate Limit.
func IsRateLimit(err error) bool {
	_, ok := err.(*RateLimitError)
	return ok
}

// IsAuthError returns true if the error is a 401 Authentication error.
func IsAuthError(err error) bool {
	_, ok := err.(*AuthenticationError)
	return ok
}

// IsForbidden returns true if the error is a 403 Forbidden error.
func IsForbidden(err error) bool {
	_, ok := err.(*ForbiddenError)
	return ok
}

// IsServerError returns true if the error is a 5xx server error.
func IsServerError(err error) bool {
	_, ok := err.(*ServerError)
	return ok
}

// =============================================================================
// Internal Methods
// =============================================================================

func (c *Client) isRetryable(statusCode int) bool {
	for _, s := range c.retryConfig.RetryableStatus {
		if statusCode == s {
			return true
		}
	}
	return false
}

func (c *Client) calculateDelay(attempt int) time.Duration {
	delay := float64(c.retryConfig.InitialDelay) * math.Pow(c.retryConfig.BackoffFactor, float64(attempt))
	if delay > float64(c.retryConfig.MaxDelay) {
		delay = float64(c.retryConfig.MaxDelay)
	}
	return time.Duration(delay)
}

func (c *Client) extractRateLimit(resp *http.Response) {
	limit := resp.Header.Get("X-RateLimit-Limit")
	remaining := resp.Header.Get("X-RateLimit-Remaining")
	reset := resp.Header.Get("X-RateLimit-Reset")

	if limit != "" && remaining != "" && reset != "" {
		limitInt, _ := strconv.Atoi(limit)
		remainingInt, _ := strconv.Atoi(remaining)
		resetInt, _ := strconv.ParseInt(reset, 10, 64)

		c.RateLimit = &RateLimitInfo{
			Limit:     limitInt,
			Remaining: remainingInt,
			Reset:     resetInt,
			ResetTime: time.Unix(resetInt, 0),
		}
	}

	c.LastRequestID = resp.Header.Get("X-Request-ID")
}

func (c *Client) handleErrorResponse(resp *http.Response, body []byte) error {
	var errResp struct {
		Detail       string       `json:"detail"`
		RequiredPlan string       `json:"required_plan,omitempty"`
		Errors       []FieldError `json:"errors,omitempty"`
	}
	json.Unmarshal(body, &errResp)

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

	baseErr := &APIError{
		StatusCode: resp.StatusCode,
		Message:    message,
		RequestID:  c.LastRequestID,
	}

	switch resp.StatusCode {
	case http.StatusUnauthorized:
		return &AuthenticationError{APIError: baseErr}
	case http.StatusForbidden:
		return &ForbiddenError{APIError: baseErr, RequiredPlan: errResp.RequiredPlan}
	case http.StatusNotFound:
		return &NotFoundError{APIError: baseErr}
	case http.StatusTooManyRequests:
		retryAfter := 0
		if ra := resp.Header.Get("Retry-After"); ra != "" {
			retryAfter, _ = strconv.Atoi(ra)
		}
		return &RateLimitError{
			APIError:   baseErr,
			RetryAfter: retryAfter,
			Limit:      c.RateLimit.Limit,
			Remaining:  c.RateLimit.Remaining,
		}
	case http.StatusBadRequest, 422:
		return &ValidationError{APIError: baseErr, Errors: errResp.Errors}
	case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return &ServerError{APIError: baseErr}
	default:
		return baseErr
	}
}

func (c *Client) doRequest(ctx context.Context, method, endpoint string, params url.Values, body interface{}, result interface{}) error {
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

	var lastErr error
	for attempt := 0; attempt <= c.retryConfig.MaxRetries; attempt++ {
		if attempt > 0 {
			delay := c.calculateDelay(attempt - 1)
			c.logger.Warn("Request failed, retrying in %v (attempt %d/%d)", delay, attempt, c.retryConfig.MaxRetries)

			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		req, err := http.NewRequestWithContext(ctx, method, u.String(), reqBody)
		if err != nil {
			return err
		}

		req.Header.Set("X-API-Key", c.apiKey)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("User-Agent", c.userAgent)

		c.logger.Debug("Request: %s %s", method, u.String())

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			if attempt < c.retryConfig.MaxRetries {
				continue
			}
			return err
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = err
			continue
		}

		c.extractRateLimit(resp)
		c.logger.Debug("Response: %d %s", resp.StatusCode, u.String())

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return json.Unmarshal(respBody, result)
		}

		lastErr = c.handleErrorResponse(resp, respBody)

		// Check if retryable
		if c.isRetryable(resp.StatusCode) && attempt < c.retryConfig.MaxRetries {
			continue
		}

		return lastErr
	}

	return lastErr
}

// =============================================================================
// API Methods
// =============================================================================

// ConsultaEndereco searches for property data by address.
func (c *Client) ConsultaEndereco(ctx context.Context, p *ConsultaEnderecoParams) (*ConsultaEnderecoResult, error) {
	params := url.Values{}
	params.Set("logradouro", p.Logradouro)
	if p.Numero != "" {
		params.Set("numero", p.Numero)
	}
	if p.Complemento != "" {
		params.Set("complemento", p.Complemento)
	}
	if p.Cidade != "" {
		params.Set("cidade", string(p.Cidade))
	} else {
		params.Set("cidade", string(CidadeSaoPaulo))
	}
	if p.IncluirHistorico {
		params.Set("incluir_historico", "true")
	}
	if p.IncluirComparaveis {
		params.Set("incluir_comparaveis", "true")
	}
	if p.IncluirZoneamento {
		params.Set("incluir_zoneamento", "true")
	}

	var result ConsultaEnderecoResult
	err := c.doRequest(ctx, "GET", "/consulta/endereco", params, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ConsultaSQL searches for property data by SQL number.
func (c *Client) ConsultaSQL(ctx context.Context, sql string, cidade Cidade) (*ConsultaSQLResult, error) {
	params := url.Values{}
	if cidade != "" {
		params.Set("cidade", string(cidade))
	} else {
		params.Set("cidade", string(CidadeSaoPaulo))
	}

	var result ConsultaSQLResult
	err := c.doRequest(ctx, "GET", "/consulta/sql/"+sql, params, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ConsultaCEP searches for properties by CEP.
func (c *Client) ConsultaCEP(ctx context.Context, cep string, cidade Cidade) ([]ConsultaEnderecoResult, error) {
	params := url.Values{}
	if cidade != "" {
		params.Set("cidade", string(cidade))
	} else {
		params.Set("cidade", string(CidadeSaoPaulo))
	}

	var result []ConsultaEnderecoResult
	err := c.doRequest(ctx, "GET", "/consulta/cep/"+cep, params, nil, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ConsultaZoneamento queries zoning by coordinates.
func (c *Client) ConsultaZoneamento(ctx context.Context, latitude, longitude float64) (*ZoneamentoResult, error) {
	params := url.Values{}
	params.Set("latitude", strconv.FormatFloat(latitude, 'f', -1, 64))
	params.Set("longitude", strconv.FormatFloat(longitude, 'f', -1, 64))

	var result ZoneamentoResult
	err := c.doRequest(ctx, "GET", "/consulta/zoneamento", params, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ValuationEstimate estimates the market value of a property.
// Requires Pro plan or higher.
func (c *Client) ValuationEstimate(ctx context.Context, p *ValuationParams) (*ValuationResult, error) {
	var result ValuationResult
	err := c.doRequest(ctx, "POST", "/valuation/estimate", nil, p, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ValuationBatch estimates values for multiple properties.
// Requires Enterprise plan.
func (c *Client) ValuationBatch(ctx context.Context, imoveis []ValuationParams) (*BatchValuationResult, error) {
	body := map[string]interface{}{
		"imoveis": imoveis,
	}

	var result BatchValuationResult
	err := c.doRequest(ctx, "POST", "/valuation/estimate/batch", nil, body, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ValuationComparables finds comparable properties.
func (c *Client) ValuationComparables(ctx context.Context, bairro string, areaMin, areaMax float64, cidade Cidade, limit int) ([]ComparavelItem, error) {
	params := url.Values{}
	params.Set("bairro", bairro)
	params.Set("area_min", strconv.FormatFloat(areaMin, 'f', -1, 64))
	params.Set("area_max", strconv.FormatFloat(areaMax, 'f', -1, 64))
	if cidade != "" {
		params.Set("cidade", string(cidade))
	}
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}

	var result []ComparavelItem
	err := c.doRequest(ctx, "GET", "/valuation/comparables", params, nil, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ValuationStatisticsResult contains statistics for a neighborhood.
type ValuationStatisticsResult struct {
	Bairro      string  `json:"bairro"`
	Cidade      string  `json:"cidade"`
	TotalImoveis int    `json:"total_imoveis"`
	Media       float64 `json:"media"`
	Mediana     float64 `json:"mediana"`
	Min         float64 `json:"min"`
	Max         float64 `json:"max"`
	DesvioPadrao float64 `json:"desvio_padrao,omitempty"`
}

// ValuationStatistics gets value statistics for a neighborhood.
func (c *Client) ValuationStatistics(ctx context.Context, bairro string, cidade Cidade) (*ValuationStatisticsResult, error) {
	params := url.Values{}
	if cidade != "" {
		params.Set("cidade", string(cidade))
	}

	var result ValuationStatisticsResult
	err := c.doRequest(ctx, "GET", "/valuation/statistics/"+url.PathEscape(bairro), params, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// DadosIPTUHistorico gets IPTU value history for a property.
func (c *Client) DadosIPTUHistorico(ctx context.Context, sql string, cidade Cidade) ([]HistoricoItem, error) {
	params := url.Values{}
	if cidade != "" {
		params.Set("cidade", string(cidade))
	}

	var result []HistoricoItem
	err := c.doRequest(ctx, "GET", "/dados/iptu/historico/"+sql, params, nil, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DadosCNPJ queries company data by CNPJ.
func (c *Client) DadosCNPJ(ctx context.Context, cnpj string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := c.doRequest(ctx, "GET", "/dados/cnpj/"+cnpj, nil, nil, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// IPCAItem represents a single IPCA index entry.
type IPCAItem struct {
	Data   string  `json:"data"`
	Valor  float64 `json:"valor"`
	Acumulado12Meses float64 `json:"acumulado_12_meses,omitempty"`
}

// DadosIPCA gets historical IPCA index data.
func (c *Client) DadosIPCA(ctx context.Context, dataInicio, dataFim string) ([]IPCAItem, error) {
	params := url.Values{}
	if dataInicio != "" {
		params.Set("data_inicio", dataInicio)
	}
	if dataFim != "" {
		params.Set("data_fim", dataFim)
	}

	var result []IPCAItem
	err := c.doRequest(ctx, "GET", "/dados/ipca", params, nil, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// IPCACorrecao performs inflation adjustment using IPCA.
func (c *Client) IPCACorrecao(ctx context.Context, valor float64, dataOrigem, dataDestino string) (map[string]interface{}, error) {
	params := url.Values{}
	params.Set("valor", strconv.FormatFloat(valor, 'f', 2, 64))
	params.Set("data_origem", dataOrigem)
	if dataDestino != "" {
		params.Set("data_destino", dataDestino)
	}

	var result map[string]interface{}
	err := c.doRequest(ctx, "GET", "/dados/ipca/corrigir", params, nil, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// =============================================================================
// IPTU Tools Types (Ferramentas IPTU 2026)
// =============================================================================

// CidadeInfo represents information about a city with IPTU calendar.
type CidadeInfo struct {
	Codigo       string `json:"codigo"`
	Nome         string `json:"nome"`
	Ano          int    `json:"ano"`
	DescontoVista string `json:"desconto_vista"`
	ParcelasMax  int    `json:"parcelas_max"`
	SiteOficial  string `json:"site_oficial"`
}

// CidadesResult represents the result of listing cities.
type CidadesResult struct {
	Cidades []CidadeInfo `json:"cidades"`
	Total   int          `json:"total"`
	Nota    string       `json:"nota,omitempty"`
}

// CalendarioResult represents the IPTU calendar for a city.
type CalendarioResult struct {
	Cidade                     string   `json:"cidade"`
	Ano                        int      `json:"ano"`
	DescontoVistaPercentual    float64  `json:"desconto_vista_percentual"`
	DescontoVistaTexto         string   `json:"desconto_vista_texto"`
	ParcelasMax                int      `json:"parcelas_max"`
	ValorMinimoParcela         float64  `json:"valor_minimo_parcela"`
	IsencaoValorVenal          float64  `json:"isencao_valor_venal,omitempty"`
	IsencaoTexto               string   `json:"isencao_texto,omitempty"`
	ConsultaOnline             string   `json:"consulta_online,omitempty"`
	SiteOficial                string   `json:"site_oficial"`
	Novidades                  []string `json:"novidades,omitempty"`
	Alertas                    []string `json:"alertas,omitempty"`
	FormasPagamento            []string `json:"formas_pagamento,omitempty"`
	VencimentosCotaUnica       []string `json:"vencimentos_cota_unica"`
	VencimentosParcelado       []string `json:"vencimentos_parcelado"`
	ProximoVencimento          string   `json:"proximo_vencimento,omitempty"`
	DiasParaProximoVencimento  int      `json:"dias_para_proximo_vencimento,omitempty"`
}

// SimuladorParams contains parameters for payment simulation.
type SimuladorParams struct {
	ValorIPTU   float64 `json:"valor_iptu"`
	Cidade      string  `json:"cidade,omitempty"`
	ValorVenal  float64 `json:"valor_venal,omitempty"`
}

// SimuladorResult represents the result of payment simulation.
type SimuladorResult struct {
	ValorOriginal       float64 `json:"valor_original"`
	ValorVista          float64 `json:"valor_vista"`
	DescontoVista       float64 `json:"desconto_vista"`
	DescontoPercentual  float64 `json:"desconto_percentual"`
	Parcelas            int     `json:"parcelas"`
	ValorParcela        float64 `json:"valor_parcela"`
	ValorTotalParcelado float64 `json:"valor_total_parcelado"`
	EconomiaVista       float64 `json:"economia_vista"`
	EconomiaPercentual  float64 `json:"economia_percentual"`
	Recomendacao        string  `json:"recomendacao"`
	ElegivelIsencao     bool    `json:"elegivel_isencao"`
	IsencaoMensagem     string  `json:"isencao_mensagem,omitempty"`
	Cidade              string  `json:"cidade"`
	Ano                 int     `json:"ano"`
	ProximoVencimento   string  `json:"proximo_vencimento,omitempty"`
}

// IsencaoResult represents the result of exemption check.
type IsencaoResult struct {
	Cidade                     string   `json:"cidade"`
	ValorVenal                 float64  `json:"valor_venal"`
	LimiteIsencao              float64  `json:"limite_isencao"`
	ElegivelIsencaoTotal       bool     `json:"elegivel_isencao_total"`
	ElegivelDescontoParcial    bool     `json:"elegivel_desconto_parcial"`
	DescontoEstimadoPercentual float64  `json:"desconto_estimado_percentual,omitempty"`
	Mensagem                   string   `json:"mensagem"`
	RequisitosAdicionais       []string `json:"requisitos_adicionais"`
}

// ProximoVencimentoResult represents next due date information.
type ProximoVencimentoResult struct {
	Cidade          string  `json:"cidade"`
	DataVencimento  string  `json:"data_vencimento"`
	DiasRestantes   int     `json:"dias_restantes"`
	Status          string  `json:"status"` // em_dia, proximo, vence_hoje, vencido
	Mensagem        string  `json:"mensagem"`
	MultaEstimada   float64 `json:"multa_estimada,omitempty"`
	JurosEstimados  float64 `json:"juros_estimados,omitempty"`
}

// =============================================================================
// IPTU Tools API Methods
// =============================================================================

// IPTUToolsCidades lists all cities with available IPTU calendar.
func (c *Client) IPTUToolsCidades(ctx context.Context) (*CidadesResult, error) {
	var result CidadesResult
	err := c.doRequest(ctx, "GET", "/iptu-tools/cidades", nil, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// IPTUToolsCalendario returns the complete IPTU calendar for the specified city.
func (c *Client) IPTUToolsCalendario(ctx context.Context, cidade Cidade) (*CalendarioResult, error) {
	params := url.Values{}
	if cidade != "" {
		params.Set("cidade", string(cidade))
	} else {
		params.Set("cidade", string(CidadeSaoPaulo))
	}

	var result CalendarioResult
	err := c.doRequest(ctx, "GET", "/iptu-tools/calendario", params, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// IPTUToolsSimulador simulates IPTU payment options (lump sum vs installments).
func (c *Client) IPTUToolsSimulador(ctx context.Context, p *SimuladorParams) (*SimuladorResult, error) {
	if p.Cidade == "" {
		p.Cidade = string(CidadeSaoPaulo)
	}

	var result SimuladorResult
	err := c.doRequest(ctx, "POST", "/iptu-tools/simulador", nil, p, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// IPTUToolsIsencao checks if a property is eligible for IPTU exemption.
func (c *Client) IPTUToolsIsencao(ctx context.Context, valorVenal float64, cidade Cidade) (*IsencaoResult, error) {
	params := url.Values{}
	params.Set("valor_venal", strconv.FormatFloat(valorVenal, 'f', 2, 64))
	if cidade != "" {
		params.Set("cidade", string(cidade))
	} else {
		params.Set("cidade", string(CidadeSaoPaulo))
	}

	var result IsencaoResult
	err := c.doRequest(ctx, "GET", "/iptu-tools/isencao", params, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// IPTUToolsProximoVencimento returns information about the next IPTU due date.
func (c *Client) IPTUToolsProximoVencimento(ctx context.Context, cidade Cidade, parcela int) (*ProximoVencimentoResult, error) {
	params := url.Values{}
	if cidade != "" {
		params.Set("cidade", string(cidade))
	} else {
		params.Set("cidade", string(CidadeSaoPaulo))
	}
	if parcela > 0 {
		params.Set("parcela", strconv.Itoa(parcela))
	}

	var result ProximoVencimentoResult
	err := c.doRequest(ctx, "GET", "/iptu-tools/proximo-vencimento", params, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
